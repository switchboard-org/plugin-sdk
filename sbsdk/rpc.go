package sbsdk

import (
	"encoding/gob"
	"log"
	"net/rpc"
)

type ProviderRPCClient struct {
	client *rpc.Client
}

func NewProviderRPCClient() Provider {
	return &ProviderRPCClient{}
}

func (p *ProviderRPCClient) Init(val []byte) error {
	var result interface{}
	err := p.client.Call("Plugin.Init", val, &result)
	if err != nil {
		panic(err)
	}
	return err
}

func (p *ProviderRPCClient) InitSchema() (ObjectSchema, error) {
	var result ObjectSchema
	err := p.client.Call("Plugin.InitSchema", new(interface{}), &result)

	if err != nil {
		panic(err)
	}
	return result, nil
}

func (p *ProviderRPCClient) ActionNames() ([]string, error) {
	var result []string
	err := p.client.Call("Plugin.ActionNames", new(interface{}), &result)
	if err != nil {
		return []string{}, err
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
		return Invalid, err
	}
	return result, nil
}

func (p *ProviderRPCClient) ActionEvaluate(name string, config []byte, input []byte) ([]byte, error) {
	var result []byte
	payload := EvaluateServerPayload{
		Name:   name,
		Config: config,
		Input:  input,
	}
	err := p.client.Call("Plugin.ActionEvaluate", payload, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type ProviderRPCServer struct {
	Impl Provider
}

func (p *ProviderRPCServer) Init(val []byte, _ *any) error {
	return p.Impl.Init(val)
}

func (p *ProviderRPCServer) InitSchema(_ any, reply *ObjectSchema) error {
	result, err := p.Impl.InitSchema()
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

type EvaluateServerPayload struct {
	Name   string
	Config []byte
	Input  []byte
}

func (p *ProviderRPCServer) ActionEvaluate(payload EvaluateServerPayload, reply *[]byte) error {
	result, err := p.Impl.ActionEvaluate(payload.Name, payload.Config, payload.Input)
	if err != nil {
		log.Println(err)
		return err
	}
	*reply = result
	return nil
}

func init() {
	gob.Register(EvaluateServerPayload{})
}
