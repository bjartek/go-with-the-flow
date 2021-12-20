package gwtf

import (
	"testing"

	"github.com/onflow/cadence"
	"github.com/stretchr/testify/assert"
)

func TestCadenceValueToJsonString(t *testing.T) {

	t.Parallel()
	t.Run("Empty value should be empty json object", func(t *testing.T) {
		value := CadenceValueToJsonString(nil)
		assert.Equal(t, "{}", value)
	})

	t.Run("Empty optional should be empty string", func(t *testing.T) {
		value := CadenceValueToJsonString(cadence.NewOptional(nil))
		assert.Equal(t, `""`, value)
	})
	t.Run("Unwrap optional", func(t *testing.T) {
		value := CadenceValueToJsonString(cadence.NewOptional(NewCadenceString("foo")))
		assert.Equal(t, `"foo"`, value)
	})
	t.Run("Array", func(t *testing.T) {
		value := CadenceValueToJsonString(cadence.NewArray([]cadence.Value{NewCadenceString("foo"), NewCadenceString("bar")}))
		assert.Equal(t, `[
    "foo",
    "bar"
]`, value)
	})

	t.Run("Dictionary", func(t *testing.T) {
		dict := cadence.NewDictionary([]cadence.KeyValuePair{{Key: NewCadenceString("foo"), Value: NewCadenceString("bar")}})
		value := CadenceValueToJsonString(dict)
		assert.Equal(t, `{
    "foo": "bar"
}`, value)
	})

	t.Run("Dictionary nested", func(t *testing.T) {
		subDict := cadence.NewDictionary([]cadence.KeyValuePair{{Key: NewCadenceString("foo"), Value: NewCadenceString("bar")}})
		dict := cadence.NewDictionary([]cadence.KeyValuePair{{Key: NewCadenceString("foo"), Value: NewCadenceString("bar")}, {Key: NewCadenceString("subdict"), Value: subDict}})
		value := CadenceValueToJsonString(dict)
		assert.Equal(t, `{
    "foo": "bar",
    "subdict": {
        "foo": "bar"
    }
}`, value)
	})

	t.Run("Dictionary", func(t *testing.T) {
		dict := cadence.NewDictionary([]cadence.KeyValuePair{{Key: cadence.NewUInt64(1), Value: cadence.NewUInt64(1)}})
		value := CadenceValueToJsonString(dict)
		assert.Equal(t, `{
    "1": "1"
}`, value)
	})

	t.Run("Struct", func(t *testing.T) {

		s := cadence.Struct{
			Fields: []cadence.Value{NewCadenceString("bar")},
			StructType: &cadence.StructType{
				Fields: []cadence.Field{{
					Identifier: "foo",
					Type:       cadence.StringType{},
				}},
			},
		}
		value := CadenceValueToJsonString(s)
		assert.Equal(t, `{
    "foo": "bar"
}`, value)
	})

}

func NewCadenceString(value string) cadence.String {
	cadenceValue, err := cadence.NewString(value)
	if err != nil {
		panic(err)
	}
	return cadenceValue
}
