package sbsdk

import (
	"github.com/zclconf/go-cty/cty"
)

// Provider is the main network interface that must be implemented by every integration provider
// in order to work with the Switchboard runner service.
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

	//ActionNames returns a list of available actions from a provider
	ActionNames() ([]string, error)

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

	//TriggerKeyNames returns a list of keys that an action can filter by. This method is
	//used during the validation process to make sure action filters are valid
	TriggerKeyNames() ([]string, error)
	//TriggerConfigurationSchema is similar to ActionConfigurationSchema, but for triggers
	TriggerConfigurationSchema() (ObjectSchema, error)
	//MapPayloadToTriggerKey maps an incoming event to a trigger key. This is helpful
	//for actions that want to fire on isolated event types and for fetching type hints.
	//This is particularly useful for triggers that have 2+ event types that the subscription listens to
	MapPayloadToTriggerKey([]byte) (string, error)
	//TriggerOutputType provides the runner with a type hint based on the value that gets returned
	//from MapPayloadToTriggerKey.
	TriggerOutputType(string) (Type, error)

	//CreateSubscription subscribes to the provider for one trigger.
	//It returns a representation of the resulting state of the subscription in raw JSON byte string
	CreateSubscription(contextId string, input []byte) ([]byte, error)
	//ReadSubscription will get the current state value of the trigger from the integration provider
	// in raw JSON byte string
	ReadSubscription(contextId string, subscriptionId string) ([]byte, error)
	//UpdateSubscription will update the trigger and return the new state value to the runner
	// as a raw JSON byte string
	UpdateSubscription(contextId string, subscriptionId string, input []byte) ([]byte, error)
	//DeleteSubscription with remove the trigger and event subscription to the vendor
	DeleteSubscription(contextId string, subscriptionId string) error
}

// Action to be implemented by providers for each individual method/endpoint tied the provider.
// Each action should be inclusive of all possible ways of calling a method/endpoint for the integration.
type Action interface {
	// ConfigurationSchema returns an ObjectSpec that returns all required/optional blocks and attributes
	// in an Action or Trigger. This should include all general configuration settings, as well as all details
	// pertinent to an individual interaction (api call, event publish, etc.) with the integration
	ConfigurationSchema() (ObjectSchema, error)
	// OutputType provides a schema structure for the result of an Action or Trigger. This is an essential component
	// of using output from one action in a workflow as the Input of another, as well as pre-deployment
	// configuration validation. Note, this does not return an ObjectSchema because that type is primarily
	// used for helping the calling application know how the hcl Config data should look.
	OutputType() (Type, error)
	//Evaluate is the main function called by the runner service when a particular action is being processed. In
	// a standard integration provider, this is where the guts of integration code will be.
	Evaluate(contextId string, input cty.Value) (cty.Value, error)
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
