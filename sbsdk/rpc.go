package sbsdk

import (
	"encoding/gob"
	"net/rpc"
)

type ProviderRPCClient struct {
	client *rpc.Client
}

func NewProviderRPCClient() Provider {
	return &ProviderRPCClient{}
}

func (p *ProviderRPCClient) Init(runnerProvider RunnerProvider) (ProviderConfig, error) {
	var result ProviderConfig
	err := p.client.Call("Plugin.Init", runnerProvider, &result)
	if err != nil {
		panic(err)
	}
	return result, nil
}

func (p *ProviderRPCClient) InitSchema() (ObjectSchema, error) {
	var result ObjectSchema
	err := p.client.Call("Plugin.InitSchema", new(interface{}), &result)

	if err != nil {
		panic(err)
	}
	return result, nil
}

func (p *ProviderRPCClient) MapPayloadToTriggerKey(data []byte) (string, error) {
	var result string
	err := p.client.Call("Plugin.MapPayloadToTriggerKey", data, &result)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (p *ProviderRPCClient) ActionNames() ([]string, error) {
	var result []string
	err := p.client.Call("Plugin.ActionNames", new(interface{}), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *ProviderRPCClient) ActionConfigurationSchema(name string) (ObjectSchema, error) {
	var result ObjectSchema
	err := p.client.Call("Plugin.ActionConfigurationSchema", name, &result)
	if err != nil {
		return ObjectSchema{}, err
	}
	return result, nil
}

func (p *ProviderRPCClient) ActionOutputType(name string) (Type, error) {
	var result Type
	err := p.client.Call("Plugin.ActionOutputType", name, &result)
	if err != nil {
		return Type{}, err
	}
	return result, nil
}

func (p *ProviderRPCClient) ActionEvaluate(contextId string, name string, input []byte) ([]byte, error) {
	var result []byte
	payload := ActionEvalData{
		ContextId: contextId,
		Name:      name,
		Input:     input,
	}
	err := p.client.Call("Plugin.ActionEvaluate", payload, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *ProviderRPCClient) TriggerKeyNames() ([]string, error) {
	var result []string
	err := p.client.Call("Plugin.TriggerKeyNames", nil, &result)
	if err != nil {
		return []string{}, err
	}
	return result, nil
}

func (p *ProviderRPCClient) TriggerConfigurationSchema() (ObjectSchema, error) {
	var result ObjectSchema
	err := p.client.Call("Plugin.TriggerConfigurationSchema", nil, &result)
	if err != nil {
		return ObjectSchema{}, err
	}
	return result, nil
}

func (p *ProviderRPCClient) TriggerOutputType(name string) (Type, error) {
	var result Type
	err := p.client.Call("Plugin.TriggerOutputType", name, &result)
	if err != nil {
		return Type{}, err
	}
	return result, nil
}

func (p *ProviderRPCClient) CreateSubscription(contextId string, input []byte) ([]byte, error) {
	var result []byte
	payload := SubscriptionData{
		ContextId: contextId,
		InputData: input,
	}
	err := p.client.Call("Plugin.CreateSubscription", payload, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *ProviderRPCClient) ReadSubscription(contextId string, subscriptionId string) ([]byte, error) {
	var result []byte
	payload := SubscriptionData{
		ContextId:      contextId,
		SubscriptionId: subscriptionId,
	}
	err := p.client.Call("Plugin.ReadSubscription", payload, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *ProviderRPCClient) UpdateSubscription(contextId string, subscriptionId string, input []byte) ([]byte, error) {
	var result []byte
	payload := SubscriptionData{
		ContextId:      contextId,
		SubscriptionId: subscriptionId,
		InputData:      input,
	}
	err := p.client.Call("Plugin.UpdateSubscription", payload, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *ProviderRPCClient) DeleteSubscription(contextId string, subscriptionId string) error {
	var result []byte
	payload := SubscriptionData{
		ContextId:      contextId,
		SubscriptionId: subscriptionId,
	}
	err := p.client.Call("Plugin.DeleteSubscription", payload, &result)
	if err != nil {
		return err
	}
	return nil
}

type ProviderRPCServer struct {
	Impl Provider
}

type InitData struct {
	GlobalConfig GlobalConfig
	UserConfig   []byte
}

func (p *ProviderRPCServer) Init(runnerProvider RunnerProvider, reply *ProviderConfig) error {
	result, err := p.Impl.Init(runnerProvider)
	if err != nil {
		return err
	}
	*reply = result
	return nil
}

func (p *ProviderRPCServer) InitSchema(_ any, reply *ObjectSchema) error {
	result, err := p.Impl.InitSchema()
	if err != nil {
		return err
	}
	*reply = result
	return nil
}

func (p *ProviderRPCServer) MapPayloadToTrigger(data []byte, reply *string) error {
	result, err := p.Impl.MapPayloadToTriggerKey(data)
	if err != nil {
		return err
	}
	*reply = result
	return nil
}

func (p *ProviderRPCServer) ActionNames(_ any, reply *[]string) error {
	result, err := p.Impl.ActionNames()
	if err != nil {
		return err
	}
	*reply = result
	return nil
}

func (p *ProviderRPCServer) ActionConfigurationSchema(name string, reply *ObjectSchema) error {
	result, err := p.Impl.ActionConfigurationSchema(name)
	if err != nil {
		return err
	}
	*reply = result
	return nil
}

func (p *ProviderRPCServer) ActionOutputType(name string, reply *Type) error {
	result, err := p.Impl.ActionOutputType(name)
	if err != nil {
		return err
	}
	*reply = result
	return nil
}

func (p *ProviderRPCServer) ActionEvaluate(payload ActionEvalData, reply *[]byte) error {
	result, err := p.Impl.ActionEvaluate(payload.ContextId, payload.Name, payload.Input)
	if err != nil {
		return err
	}
	*reply = result
	return nil
}

func (p *ProviderRPCServer) TriggerConfigurationSchema(_ any, reply *ObjectSchema) error {
	result, err := p.Impl.TriggerConfigurationSchema()
	if err != nil {
		return err
	}
	*reply = result
	return nil
}

func (p *ProviderRPCServer) TriggerOutputType(name string, reply *Type) error {
	result, err := p.Impl.TriggerOutputType(name)
	if err != nil {
		return err
	}
	*reply = result
	return nil
}

func (p *ProviderRPCServer) CreateSubscription(data SubscriptionData, reply *[]byte) error {
	result, err := p.Impl.CreateSubscription(data.ContextId, data.InputData)
	if err != nil {
		return err
	}
	*reply = result
	return nil
}

func (p *ProviderRPCServer) ReadSubscription(data SubscriptionData, reply *[]byte) error {
	result, err := p.Impl.ReadSubscription(data.ContextId, data.SubscriptionId)
	if err != nil {
		return err
	}
	*reply = result
	return nil
}

func (p *ProviderRPCServer) UpdateSubscription(data SubscriptionData, reply *[]byte) error {
	result, err := p.Impl.UpdateSubscription(data.ContextId, data.SubscriptionId, data.InputData)
	if err != nil {
		return err
	}
	*reply = result
	return nil
}

func (p *ProviderRPCServer) DeleteSubscription(data SubscriptionData) error {
	err := p.Impl.DeleteSubscription(data.ContextId, data.SubscriptionId)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	gob.Register(ActionEvalData{})
	gob.Register(InitData{})
	gob.Register(ProviderConfig{})
	gob.Register(GlobalConfig{})
	gob.Register(SubscriptionData{})
}
