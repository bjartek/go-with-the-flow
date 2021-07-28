package gwtf

import (
	"github.com/onflow/cadence"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCadenceValueToJsonString(t *testing.T) {

	t.Run("Empty optional should be empty string", func(t *testing.T) {
		value := CadenceValueToJsonString(cadence.NewOptional(nil))
		assert.Equal(t, `""`, value)
	})
	t.Run("Unwrap optional", func(t *testing.T) {
		value := CadenceValueToJsonString(cadence.NewOptional(cadence.NewString("foo")))
		assert.Equal(t, `"foo"`, value)
	})
	t.Run("Array", func(t *testing.T) {
		value := CadenceValueToJsonString(cadence.NewArray([]cadence.Value{cadence.NewString("foo"), cadence.NewString("bar")}))
		assert.Equal(t, `[
    "foo",
    "bar"
]`, value)
	})

	t.Run("Dictionary", func(t *testing.T) {
		dict := cadence.NewDictionary([]cadence.KeyValuePair{{Key: cadence.NewString("foo"), Value: cadence.NewString("bar")}})
		value := CadenceValueToJsonString(dict)
		assert.Equal(t, `{
    "foo": "bar"
}`, value)
	})

	t.Run("Struct", func(t *testing.T) {

		s := cadence.Struct{
			Fields: []cadence.Value{cadence.NewString("bar")},
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
