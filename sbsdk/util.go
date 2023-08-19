package sbsdk

import (
	"errors"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/json"
)

func OptionalAttrSchema(name string, valType Type) *AttrSchema {
	return &AttrSchema{
		Name:     name,
		Type:     valType,
		Required: false,
	}
}

func RequiredAttrSchema(name string, valType Type) *AttrSchema {
	return &AttrSchema{
		Name:     name,
		Type:     valType,
		Required: true,
	}
}

func OptionalBlockSchema(name string, nested Schema) *BlockSchema {
	return &BlockSchema{
		Name:     name,
		Nested:   nested,
		Required: false,
	}
}

func InvokeActionEvaluation(action Action, contextId string, polymorphicKey string, input []byte) ([]byte, error) {
	var err error
	if input == nil {
		return nil, errors.New("input must not be null")
	}
	inputSchema, _ := action.ConfigurationSchema()
	inputVal, err := json.Unmarshal(input, hcldec.ImpliedType(inputSchema.Decode()))
	if err != nil {
		return nil, err
	}
	result, err := action.Evaluate(contextId, inputVal)
	if err != nil {
		return nil, err
	}

	outputType, _ := action.OutputType()
	ctyType, err := outputType.ToCty(&polymorphicKey)
	if err != nil {
		return nil, err
	}
	return json.Marshal(result, ctyType)
}

// NodeIsFormat is a simple check to make sure a particular value conforms to one of the format node structures.
// This function does not deep check any nested children, just the top level.
func NodeIsFormat(value cty.Value) bool {
	return NodeIsFormatConstraint(value) || NodeIsFormatKeyValue(value)
}

// NodeIsFormatConstraint checks to see whether a value conforms to the structure of a constraint node, which
// will/should be true for all leaf nodes in a value object
func NodeIsFormatConstraint(value cty.Value) bool {
	valType := value.Type()
	if value.IsNull() || !value.IsKnown() || !value.CanIterateElements() || !valType.IsObjectType() {
		return false
	}
	return valType.HasAttribute(FORMAT_TYPE) && valType.HasAttribute(FORMAT_REQUIRED)
}

// NodeIsFormatKeyValue checks to see whether a value conforms to the structure of a key/val node where all values
// conform to the constraint node type, where the value itself is an object with string keys and FormatConstraint looking values
func NodeIsFormatKeyValue(value cty.Value) bool {
	valType := value.Type()
	if NodeIsFormatConstraint(value) || value.IsNull() || !value.IsKnown() || !value.CanIterateElements() || !valType.IsObjectType() {
		return false
	}
	iter := value.ElementIterator()
	for iter.Next() {
		k, v := iter.Element()
		if !k.Type().Equals(cty.String) {
			return false
		}
		// all vals in key/val MUST be a normal node
		if !NodeIsFormatConstraint(v) {
			return false
		}
	}
	return true
}

// FromConstraintValueToType is a helper function that gets the type of nested values
// for a cty.Value that represents a type
func FromConstraintValueToType(val cty.Value) (*Type, error) {
	if !NodeIsFormat(val) {
		return nil, errors.New("invalid type representation")
	}

	switch {
	case NodeIsFormatKeyValue(val):
		return parseObjectType(val)

	case nodeIsPrimitiveType(val):
		return parseSpecificType(val)

	default:
		return nil, errors.New("invalid type representation")
	}
}

func parseObjectType(val cty.Value) (*Type, error) {
	outputType := &Type{TypeName: OBJECT_TYPE}
	nestedVals, err := parseNestedFormatVals(val)
	if err != nil {
		return nil, err
	}
	outputType.NestedValues = &nestedVals
	return outputType, nil
}

func parseNestedFormatVals(val cty.Value) (map[string]Type, error) {
	nestedVals := make(map[string]Type)
	iter := val.ElementIterator()

	for iter.Next() {
		k, v := iter.Element()
		res, err := FromConstraintValueToType(v)
		if err != nil {
			return nil, err
		}
		nestedVals[k.AsString()] = *res
	}

	return nestedVals, nil
}

func nodeIsPrimitiveType(val cty.Value) bool {
	typeStr := val.GetAttr(FORMAT_TYPE).AsString()
	return typeStr == NUMBER_TYPE || typeStr == BOOLEAN_TYPE || typeStr == STRING_TYPE
}

func parseSpecificType(val cty.Value) (*Type, error) {
	typeStr := val.GetAttr(FORMAT_TYPE).AsString()

	switch typeStr {
	case NUMBER_TYPE, BOOLEAN_TYPE, STRING_TYPE:
		return &Type{TypeName: typeStr}, nil

	case LIST_TYPE:
		listInnerType, err := FromConstraintValueToType(val.GetAttr(FORMAT_CHILDREN))
		if err != nil {
			return nil, err
		}
		return &Type{TypeName: LIST_TYPE, InternalType: listInnerType}, nil

	case OBJECT_TYPE:
		nestedVals, err := parseNestedFormatVals(val.GetAttr(FORMAT_CHILDREN))
		if err != nil {
			return nil, err
		}
		return &Type{TypeName: OBJECT_TYPE, NestedValues: &nestedVals}, nil

	default:
		return nil, errors.New("invalid type representation")
	}
}
