package sbsdk

type RunnerProvider interface {
	//UserConfig is a map of byte string arrays, where each byte string array should conform to the Schema
	//as specified by Provider.InitSchema method when marshaled into cty.JSON. Each key in the map is a context ID
	//which is essentially a unique configuration of the provider as defined in the user provided config
	UserConfig() map[string][]byte
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
