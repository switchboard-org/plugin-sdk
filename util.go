package main

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

func InvokeActionEvaluation(action Action, p Provider, name string, config []byte, input []byte) ([]byte, error) {

	var contextVal cty.Value
	var err error
	if config != nil {
		configSchema, _ := p.InitSchema()
		contextVal, err = json.Unmarshal(config, hcldec.ImpliedType(configSchema.Decode()))
		if err != nil {
			return nil, err
		}
	} else {
		contextVal = cty.NilVal
	}
	if input == nil {
		return nil, errors.New("input must not be null")
	}
	inputSchema, _ := p.ActionConfigurationSchema(name)
	inputVal, err := json.Unmarshal(input, hcldec.ImpliedType(inputSchema.Decode()))
	if err != nil {
		return nil, err
	}
	result, err := action.Evaluate(contextVal, inputVal)
	if err != nil {
		return nil, err
	}

	outputType, _ := p.ActionOutputType(name)
	return json.Marshal(result, outputType.ToCty())
}
