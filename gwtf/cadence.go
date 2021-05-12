package gwtf

import (
	"encoding/json"
	"strconv"

	"github.com/onflow/cadence"
)

func CadenceValueToJsonString(value cadence.Value) string {
	result := CadenceValueToInterface(value)
	json, _ := json.MarshalIndent(result, "", "    ")
	return string(json)
}

func CadenceValueToInterface(field cadence.Value) interface{} {
	dictionaryValue, isDictionary := field.(cadence.Dictionary)
	structValue, isStruct := field.(cadence.Struct)
	arrayValue, isArray := field.(cadence.Array)
	if isStruct {
		subStructNames := structValue.StructType.Fields
		result := map[string]interface{}{}
		for j, subField := range structValue.Fields {
			result[subStructNames[j].Identifier] = CadenceValueToInterface(subField)
		}
		return result
	} else if isDictionary {
		result := map[string]interface{}{}
		for _, item := range dictionaryValue.Pairs {
			result[item.Key.String()] = CadenceValueToInterface(item.Value)
		}
		return result
	} else if isArray {
		result := []interface{}{}
		for _, item := range arrayValue.Values {
			result = append(result, CadenceValueToInterface(item))
		}
		return result
	}
	result, err := strconv.Unquote(field.String())
	if err != nil {
		return field.String()
	}
	return result
}
