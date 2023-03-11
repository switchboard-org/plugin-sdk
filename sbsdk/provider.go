package sbsdk

import (
	"github.com/zclconf/go-cty/cty"
)

// Provider is the main interface that must be implemented by every integration provider
// in order to work with the Switchboard runner service. In addition to provider specific
// methods, the provider also includes proxy methods to call on specific Action or Trigger
// implementations (this way, we don't have to register every Action and Trigger as a plugin)
//
// Every method returns an error as the last return type so that we can gracefully deal with
// any RPC related errors in the go-plugin client implementation of this interface.
type Provider interface {
	//Init is called shortly after loading a plugin, and provides some stateful values
	//to the provider, such as auth keys.
	//
	//The byte array param must be the result of json.Marshal
	//from the github.com/zclconf/go-cty/cty/json package and should immediately be
	//unmarshalled into a cty.Value based on the type as inferred from the hcldec.ObjectSpec from InitSchema
	Init([]byte) error

	//InitSchema is used by the CLI to validate user provided Config, and is also used in Init
	//to unmarshal the string into a cty.Value
	InitSchema() (ObjectSchema, error)
	ActionNames() ([]string, error)
	ActionEvaluate(name string, config []byte, input []byte) ([]byte, error)
	ActionConfigurationSchema(name string) (ObjectSchema, error)
	ActionOutputType(name string) (Type, error)
	//TriggerNames() ([]string, error)
	//TriggerConfigurationSchema(Name string) (hcldec.ObjectSpec, error)
	//TriggerOutputType(Name string) (cty.typeImpl, error)
	//TriggerRegister(Name string, Config cty.Value, Input cty.Value) (string, error)
	//TriggerUnregister(Name string, Config cty.Value, Input cty.Value) (string, error)
	//TriggerShouldListen(Name string) (bool, error)
	//TriggerListen(Name string, cb func(output cty.Value)) error
	//TriggerIsRegistered(Name string, Config cty.Value, Input cty.Value) (bool, error)

	// ConfigurationSchema is used by the runner to parse provider configuration set by the user. It will be parsed and
	// saved in memory and the resulting value will be used as the first param of Action.Evaluate, Trigger.Register,
	// and Trigger.Unregister. Most often, the parsed value will include global authentication settings.
	//ConfigurationSchema() (hcldec.ObjectSpec, error)
}

type Function interface {
	// ConfigurationSchema returns an ObjectSpec that returns all required/optional blocks and attributes
	// in an Action or Trigger. This should include all general configuration settings, as well as all details
	// pertinent to an individual interaction (api call, event publish, etc.) with the integration
	ConfigurationSchema() (ObjectSchema, error)
	// OutputType provides a schema structure for the result of an Action or Trigger. This is an essential component
	// of using output from one action in a workflow as the Input of another, as well as pre-deployment
	// configuration validation. Note, this does not return an ObjectSchema because that type is primarily
	// used for helping the calling application know how the hcl Config data should look.
	OutputType() (Type, error)
}

type Action interface {
	Function
	// Evaluate accepts two value parameters. The first contains provider Config details. The second
	//is a value that should map to the Function.ConfigurationSchema output. This
	// is the main function called by the runner service when a particular action is being processed. In
	// a standard integration provider, this is where the guts of integration code will be.
	Evaluate(config cty.Value, input cty.Value) (cty.Value, error)
}

type Trigger interface {
	Function
	IsRegistered(config cty.Value, input cty.Value) (bool, error)
	// Register accepts two value parameters. The first contains provider Config details. The second is the
	// Input data for the integration. Register is responsible for working with the respective integration to
	//register the trigger with switchboard. This may be an event-subscription, polling api, or webhook registration.
	Register(config cty.Value, input cty.Value) (string, error)
	Unregister(config cty.Value, input cty.Value) (string, error)
	// ShouldListen tells the caller if it should listen for messages directly from the plugin.
	// Only applicable for unique cases like event queues where messages are received at the subscription point
	ShouldListen() (bool, error)
	// Listen is passed a callback that will be fired any time a trigger where ShouldListen() = true receives
	// an event from the subscription point
	Listen(cb func(output cty.Value)) error
}
