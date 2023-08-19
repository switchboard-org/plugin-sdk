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
	//Init is called shortly after loading a plugin, and gives the provider access to certain
	//data owned by the runner
	//
	//It hands the provider an instance of the ContextProvider interface, which will
	//be implemented by the runner, and various methods should be called to get context/config
	//information from the runner.
	Init(runnerProvider RunnerProvider) (ProviderConfig, error)

	//InitSchema is used by the CLI to validate user provided Config, and is also used in Init
	//to unmarshal the string into a cty.Value
	InitSchema() (ObjectSchema, error)

	//MapPayloadToTrigger maps an incoming event to a trigger. This is relevant
	//for providers who send many events to one endpoint, and will not be called
	//if one event is mapped to one trigger
	MapPayloadToTrigger([]byte) (string, error)

	//ActionNames returns a list of available actions from a provider
	ActionNames() ([]string, error)
	TriggerNames() ([]string, error)

	//ActionEvaluate is the implementation of a particular action as defined by the provider.
	//It is called by its name as listed in ActionNames. The input param must conform to the schema provided by ActionConfigurationSchema
	ActionEvaluate(contextId string, name string, input []byte) ([]byte, error)
	//ActionConfigurationSchema returns a data structure representing the expected schema
	//in the user-provided hcl configuration file for a particular action
	ActionConfigurationSchema(name string) (ObjectSchema, error)
	//ActionOutputType returns a data structure representing the result-type of a particular
	//action on success. This is used for both converting data to a dynamic cty.Value, and
	//for validating correct user-defined configuration
	ActionOutputType(name string) (Type, error)

	//TriggerConfigurationSchema is similar to ActionConfigurationSchema, but for triggers
	TriggerConfigurationSchema(name string) (ObjectSchema, error)
	//TriggerOutputType is similar to ActionOutputType, but for triggers
	TriggerOutputType(name string) (Type, error)

	//CreateSubscription subscribes to the provider for one or all triggers. Some providers
	//register a subscription for each trigger. Others subscribe for all in one. input param
	//may be a list of triggers, or a single trigger, as they conform to the configuration schema.
	//It returns a representation of the resulting state of the subscription
	CreateSubscription(contextId string, input []byte) ([]byte, error)
	//ReadSubscription will get the current state value of the trigger from the integration provider
	ReadSubscription(contextId string, subscriptionId string) ([]byte, error)
	//UpdateSubscription will update the trigger and return the new state value to the runner
	UpdateSubscription(contextId string, subscriptionId string, input []byte) ([]byte, error)
	//DeleteSubscription with remove the trigger and event subscription to the vendor
	DeleteSubscription(contextId string, subscriptionId string) error
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
	//Evaluate is the main function called by the runner service when a particular action is being processed. In
	// a standard integration provider, this is where the guts of integration code will be.
	Evaluate(contextId string, input cty.Value) (cty.Value, error)
}

// Trigger is an interface that maps all entry points for integrations. Triggers are registered
// with the integration all at once or individually, depending on the provider.
type Trigger interface {
	Function
}

// ProviderConfig is some static information the runner can use to decipher how to process certain types
// of provider setups.
type ProviderConfig struct {
	SubscriptionsRegisteredTogether bool
}

type ActionEvalData struct {
	ContextId string
	Name      string
	Input     []byte
}

type SubscriptionData struct {
	ContextId      string
	SubscriptionId string
	InputData      []byte
}
