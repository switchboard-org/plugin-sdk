package sbsdk

import (
	"errors"
	"github.com/hashicorp/hcl/v2/hcldec"
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

func InvokeActionEvaluation(contextId string, action Action, input []byte) ([]byte, error) {
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
	ctyType := outputType.ToCty()
	if err != nil {
		return nil, err
	}
	return json.Marshal(result, ctyType)
}
