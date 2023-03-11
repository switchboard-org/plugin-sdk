package sbsdk

import (
	"encoding/gob"
	"github.com/zclconf/go-cty/cty"
)

const (
	NUMBER_TYPE  = "number"
	BOOLEAN_TYPE = "bool"
	STRING_TYPE  = "string"
	OBJECT_TYPE  = "object"
	MAP_TYPE     = "map"
	LIST_TYPE    = "list"
	INVALID_TYPE = "invalid"
)

// typeImpl is an interface implemented by the Type struct that can be
// sent over the wire and converted to cty.Type when necessary for serializing/deserializing
// payloads sent to/from the runner.
type typeImpl interface {
	//ToCty converts this type to a comparable cty.Type for the purpose of encoding a cty/json value
	//into a cty.Value
	ToCty() cty.Type
}

// Type is a serializable data structure that helps the CLI and runner understand what data structures
// in configuration should look like.
type Type struct {
	//Name is a string representation of the underlying type
	Name string `json:"Name"`
	//NestedValues is used exclusively for an "object" type
	NestedValues *map[string]Type `json:"nestedValues,omitempty"`

	//InternalType is used to represent what type the value of a list or map value is.
	InternalType *Type `json:"internalType,omitempty"`
}

// ToCty is used for converting types into the cty.Type. The cty library is a hard dep in
// hcl library for decoding user config values, but it is not serializable, which is why
// we have our own type representation here.
func (t *Type) ToCty() cty.Type {
	switch t.Name {
	case NUMBER_TYPE:
		return cty.Number
	case BOOLEAN_TYPE:
		return cty.Bool
	case STRING_TYPE:
		return cty.String
	case OBJECT_TYPE:
		mappedTypes := make(map[string]cty.Type)
		if t.NestedValues == nil {
			return cty.Object(map[string]cty.Type{})
		}
		for k, v := range *t.NestedValues {
			mappedTypes[k] = v.ToCty()
		}
		return cty.Object(mappedTypes)
	case MAP_TYPE:
		if t.InternalType == nil {
			return cty.Map(cty.String)
		}
		return cty.Map(t.InternalType.ToCty())
	case LIST_TYPE:
		if t.InternalType == nil {
			return cty.List(cty.String)
		}
		return cty.List(t.InternalType.ToCty())
	default:
		return cty.NilType
	}
}

// String represents a primitive string type
var String Type

// Number represents a primitive numeric type
var Number Type

// Bool represents a primitive boolean type
var Bool Type

// Invalid is used as an error value for types
var Invalid Type

// Object creates structured map of passes key/values as a Type object.
// Use this when your map keys are known and static
func Object(data map[string]Type) Type {
	return Type{
		Name:         OBJECT_TYPE,
		NestedValues: &data,
	}
}

// Map creates an unstructured map of passed key/values as a Type object.
// Use this when the map keys are unknown or variable
func Map(valType Type) Type {
	return Type{
		Name:         MAP_TYPE,
		InternalType: &valType,
	}
}

func List(valType Type) Type {
	return Type{
		Name:         LIST_TYPE,
		InternalType: &valType,
	}
}

func init() {
	String = Type{
		Name: STRING_TYPE,
	}
	Number = Type{
		Name: NUMBER_TYPE,
	}
	Bool = Type{
		Name: BOOLEAN_TYPE,
	}
	Invalid = Type{
		Name: INVALID_TYPE,
	}

	gob.Register(&Type{})
}
