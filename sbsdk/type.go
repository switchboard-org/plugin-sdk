package sbsdk

import (
	"encoding/gob"
	"errors"
	"github.com/zclconf/go-cty/cty"
)

const (
	NUMBER_TYPE            = "number"
	BOOLEAN_TYPE           = "bool"
	STRING_TYPE            = "string"
	OBJECT_TYPE            = "object"
	MAP_TYPE               = "map"
	LIST_TYPE              = "list"
	DYNAMIC_TYPE           = "dynamic"
	INVALID_TYPE           = "invalid"
	TYPE_NAME_KEY          = "type_name"
	TYPE_NESTED_VALUES_KEY = "nested_values"
	TYPE_INTERNAL_TYPE_KEY = "internal_type"
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
	//TypeName is a string representation of the underlying type
	TypeName string
	//NestedValues is used exclusively for an "object" type
	NestedValues *map[string]Type

	//InternalType is used to represent what type the value of a list or map value is.
	InternalType *Type
}

// ValConformsToTypeStructure is a recursive function that checks whether a cty.Value
// object comforms to the appropriate structure expected to generate a Type object
func valConformsToTypeStructure(val cty.Value, isRoot bool) error {

	if !val.Type().HasAttribute(TYPE_NAME_KEY) || !val.GetAttr(TYPE_NAME_KEY).Type().Equals(cty.String) {
		return errors.New("invalid value to convert to type")
	}
	typeNameVal := val.GetAttr(TYPE_NAME_KEY).AsString()

	if isRoot && typeNameVal != OBJECT_TYPE {
		return errors.New("root type must be an object")
	}
	if typeNameVal == OBJECT_TYPE {
		if !(val.Type().HasAttribute(TYPE_NESTED_VALUES_KEY) && val.GetAttr(TYPE_NESTED_VALUES_KEY).Type().IsObjectType()) {
			return errors.New("objects must have a nested_values attribute set to an object value")
		}
		iter := val.GetAttr(TYPE_NESTED_VALUES_KEY).ElementIterator()
		for iter.Next() {
			_, el := iter.Element()
			err := valConformsToTypeStructure(el, false)
			if err != nil {
				return err
			}
		}
	}

	if typeNameVal == MAP_TYPE || typeNameVal == LIST_TYPE {
		if !(val.Type().HasAttribute(TYPE_INTERNAL_TYPE_KEY) && val.GetAttr(TYPE_INTERNAL_TYPE_KEY).Type().IsObjectType()) {
			return errors.New("maps must have an internal_type attribute set to an object value")
		}
		err := valConformsToTypeStructure(val.GetAttr(TYPE_INTERNAL_TYPE_KEY), false)
		if err != nil {
			return err
		}
	}
	return nil
}

func FromCtyToType(val cty.Value, isRoot bool) (*Type, error) {
	if isRoot {
		err := valConformsToTypeStructure(val, true)
		if err != nil {
			return nil, err
		}
	}
	switch val.GetAttr(TYPE_NAME_KEY).AsString() {
	case OBJECT_TYPE:
		iter := val.GetAttr(TYPE_NESTED_VALUES_KEY).ElementIterator()
		nestedOut := make(map[string]Type)
		for iter.Next() {
			k, v := iter.Element()
			ty, err := FromCtyToType(v, false)
			if err != nil {
				return nil, err
			}
			nestedOut[k.AsString()] = *ty
		}
		out := Type{
			TypeName:     OBJECT_TYPE,
			NestedValues: &nestedOut,
			InternalType: nil,
		}
		return &out, nil
	case LIST_TYPE:
	case MAP_TYPE:
		internalType, err := FromCtyToType(val.GetAttr(TYPE_INTERNAL_TYPE_KEY), false)
		if err != nil {
			return nil, err
		}
		out := Type{
			TypeName:     val.GetAttr(TYPE_NAME_KEY).AsString(),
			NestedValues: nil,
			InternalType: internalType,
		}
		return &out, nil
	default:
		return &Type{
			TypeName: val.GetAttr(TYPE_NAME_KEY).AsString(),
		}, nil
	}
	return nil, nil
}

// ToCty is used for converting types into the cty.Type. The cty library is a hard dep in
// hcl library for decoding user config values, but it is not serializable, which is why
// we have our own type representation here.
func (t *Type) ToCty() cty.Type {
	switch t.TypeName {
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

// Dynamic represents a type that can be anything. This is not recommended
// as these values cannot be validated ahead of time in user configuration.
// Only use this value when the underlying type is uncertain and variable.
var Dynamic Type

// Invalid is used as an error value for types
var Invalid Type

// Object creates structured map of passes key/values as a Type object.
// Use this when your map keys are known and static
func Object(data map[string]Type) Type {
	return Type{
		TypeName:     OBJECT_TYPE,
		NestedValues: &data,
	}
}

// Map creates a map with unknown keys but known values as a Type object.
// Use this when the map keys are unknown/variable, but values are known
func Map(valType Type) Type {
	return Type{
		TypeName:     MAP_TYPE,
		InternalType: &valType,
	}
}

// List creates a list with known types as list values.
func List(valType Type) Type {
	return Type{
		TypeName:     LIST_TYPE,
		InternalType: &valType,
	}
}

func init() {
	String = Type{
		TypeName: STRING_TYPE,
	}
	Number = Type{
		TypeName: NUMBER_TYPE,
	}
	Bool = Type{
		TypeName: BOOLEAN_TYPE,
	}
	Dynamic = Type{
		TypeName: DYNAMIC_TYPE,
	}
	Invalid = Type{
		TypeName: INVALID_TYPE,
	}

	gob.Register(&Type{})
}
