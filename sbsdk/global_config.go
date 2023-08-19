package sbsdk

type RunnerProvider interface {
	//UserConfig is a byte string that should conform to the Schema
	//as specified by Provider.InitSchema method. The contextId parameter
	//is a specific namespace for a provider (you can have multiple provider configurations
	//in switchboard). This also provides a future-compatible interface for embedded
	//integration setups where actions in one workflow may have a dynamic
	//configuration context
	UserConfig(contextId string) []byte
	GlobalConfig() GlobalConfig
}

// GlobalConfig contains details related to the core switchboard instance
// that each provider may use as they see fit in their implementations
type GlobalConfig struct {
	//PublicUri tells each provider where they should register
	//triggers that are http based
	PublicIngestUri string
	//PrivateUri is used by providers that have event subscriptions that are received directly by the provider,
	//such as event-queue triggers or polling-style event listeners. This
	PrivateIngestUri string
}
