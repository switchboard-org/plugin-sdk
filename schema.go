package main

import (
	"encoding/gob"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/json"
)

func init() {
	gob.Register(&ObjectSchema{})
	gob.Register(&AttrSchema{})
	gob.Register(&BlockSchema{})
}

// Schema is an interface implemented by various serializable structs that are used by
// provider implementations to define the structure of hcl Config data as it is passed by the user.
// Schemas are pulled from the runner of CLI to tell the hcl decoder what the user Config
// should look like. All structs that implement this interface MUST use exportable/public properties.
// non-exportable/private properties will not work when sending data over the wire.
type Schema interface {
	//Decode transforms a Schema struct into a hcldec.Spec so that hcl Config data can be serialized
	//into a usable value and passed to the provider. This is primarily used in the runner application
	//where hcl is being parsed.
	Decode() hcldec.Spec
}

// ObjectSchema is primarily used as root of all schemas in switchboard. This enforces a key/value
// style of Config for anything using this. Many Provider interface methods will
// return this to enforce a certain style of how providers expect user Config data to look.
type ObjectSchema map[string]Schema

func (s *ObjectSchema) Decode() hcldec.Spec {
	outputSpec := hcldec.ObjectSpec{}
	for k, v := range *s {
		outputSpec[k] = v.Decode()
	}
	return outputSpec
}

// BlockSchema maps to block types from the Config. Blocks with labels are not
// supported, meaning only one block type per provided Config is permitted
type BlockSchema struct {
	Name     string
	Required bool
	Nested   Schema
}

func (b *BlockSchema) Decode() hcldec.Spec {
	return &hcldec.BlockSpec{
		TypeName: b.Name,
		Nested:   b.Nested.Decode(),
		Required: b.Required,
	}
}

// AttrSchema represents a key/val coming from the provided data structure
type AttrSchema struct {
	Name     string
	Required bool
	Type     Type
}

func (b *AttrSchema) Decode() hcldec.Spec {
	return &hcldec.AttrSpec{
		Name:     b.Name,
		Type:     b.Type.ToCty(),
		Required: b.Required,
	}
}

func UnmarshalVal(schema Schema, data []byte) (cty.Value, error) {
	t := hcldec.ImpliedType(schema.Decode())
	return json.Unmarshal(data, t)
}

func MarshalVal(schema Schema, val cty.Value) ([]byte, error) {
	t := hcldec.ImpliedType(schema.Decode())
	return json.Marshal(val, t)
}
