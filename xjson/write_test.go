package xjson

import (
	"bytes"
	"testing"

	"github.com/neticdk/go-stdlib/assert"
)

func TestPrettyPrintJSON(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		input := []byte(`{"name": "John", "age": 30}`)
		expected := "{\n  \"name\": \"John\",\n  \"age\": 30\n}\n"
		var output bytes.Buffer
		err := PrettyPrintJSON(input, &output)
		assert.NoError(t, err)
		assert.Equal(t, output.String(), expected)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		input := []byte(`{"name": "John", "age": 30`) // missing closing brace
		var output bytes.Buffer
		err := PrettyPrintJSON(input, &output)
		assert.Error(t, err)
	})
}
