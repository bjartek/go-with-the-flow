package gwtf

import (
	"encoding/json"
	"strconv"

	"github.com/onflow/cadence"
)

//CadenceValueToJsonString converts a cadence.Value into a json pretty printed string
func CadenceValueToJsonString(value cadence.Value) string {
	if value == nil {
		return "{}"
	}
	result := CadenceValueToInterface(value)
	j, _ := json.MarshalIndent(result, "", "    ")
	return string(j)
}

//CadenceValueToInterface convert a candence.Value into interface{}
func CadenceValueToInterface(field cadence.Value) interface{} {
	if field == nil {
		return ""
	}

	switch field.(type) {
	case cadence.Optional:
		if field == nil {
			return nil
		} else {
			return CadenceValueToInterface(field.(cadence.Optional).Value)
		}
	case cadence.Dictionary:
		result := map[string]interface{}{}
		for _, item := range field.(cadence.Dictionary).Pairs {
			result[item.Key.String()] = CadenceValueToInterface(item.Value)
		}
		return result
	case cadence.Struct:
		result := map[string]interface{}{}
		subStructNames := field.(cadence.Struct).StructType.Fields
		for j, subField := range field.(cadence.Struct).Fields {
			result[subStructNames[j].Identifier] = CadenceValueToInterface(subField)
		}
		return result
	case cadence.Array:
		var result []interface{}
		for _, item := range field.(cadence.Array).Values {
			result = append(result, CadenceValueToInterface(item))
		}
		return result
	default:
		result, err := strconv.Unquote(field.String())
		if err != nil {
			return field.String()
		}
		return result
	}
}
