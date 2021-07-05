package gwtf

import (
	"encoding/json"
	"strconv"

	"github.com/onflow/cadence"
)

func CadenceValueToJsonString(value cadence.Value) string {
	result := CadenceValueToInterface(value)
	j, _ := json.MarshalIndent(result, "", "    ")
	return string(j)
}

func CadenceValueToInterface(field cadence.Value) interface{} {
	switch field.(type) {
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
