package sbsdk

import (
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

func MapInputToCtyValue(input []byte, schema ObjectSchema) (cty.Value, error) {
	inputVal, err := json.Unmarshal(input, hcldec.ImpliedType(schema.Decode()))
	if err != nil {
		return cty.NilVal, err
	}
	return inputVal, nil
}

func MapCtyValueToByteString(val cty.Value, outputType Type) ([]byte, error) {
	ctyType := outputType.ToCty()
	return json.Marshal(val, ctyType)
}
