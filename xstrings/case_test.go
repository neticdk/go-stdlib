package xstrings

import (
	"testing"

	"github.com/neticdk/go-stdlib/assert"
)

func TestToKebabCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "already kebab case",
			input:    "already-kebab-case",
			expected: "already-kebab-case",
		},
		{
			name:     "camel case",
			input:    "camelCase",
			expected: "camel-case",
		},
		{
			name:     "pascal case",
			input:    "PascalCase",
			expected: "pascal-case",
		},
		{
			name:     "snake case",
			input:    "snake_case",
			expected: "snake-case",
		},
		{
			name:     "screaming snake case",
			input:    "SCREAMING_SNAKE_CASE",
			expected: "screaming-snake-case",
		},
		{
			name:     "space separated",
			input:    "space separated words",
			expected: "space-separated-words",
		},
		{
			name:     "mixed separators",
			input:    "mixed.separators_with space",
			expected: "mixed-separators-with-space",
		},
		{
			name:     "consecutive uppercase",
			input:    "HTTPRequest",
			expected: "http-request",
		},
		{
			name:     "plural case consecutive uppercase",
			input:    "working with APIs",
			expected: "working-with-apis",
		},
		{
			name:     "mixed case with numbers",
			input:    "getHTTP2Data",
			expected: "get-http-2-data",
		},
		{
			name:     "mixed case with lower s and numbers",
			input:    "getHTTPs2Data",
			expected: "get-https2-data",
		},
		{
			name:     "mixed case with upper S and numbers",
			input:    "getHTTPS2Data",
			expected: "get-https-2-data",
		},
		{
			name:     "consecutive separators",
			input:    "multiple__separators",
			expected: "multiple-separators",
		},
		{
			name:     "flat case",
			input:    "flatcase",
			expected: "flatcase",
		},
		{
			name:     "single letter",
			input:    "a",
			expected: "a",
		},
		{
			name:     "acronym at start",
			input:    "APIRequest",
			expected: "api-request",
		},
		{
			name:     "acronym at end",
			input:    "requestAPI",
			expected: "request-api",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToKebabCase(tt.input)
			assert.Equal(t, result, tt.expected, "ToKebabCase/%q", tt.name)
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "already camel case",
			input:    "helloWorld",
			expected: "helloWorld",
		},

		{
			name:     "snake case",
			input:    "hello_world",
			expected: "helloWorld",
		},
		{
			name:     "kebab case",
			input:    "hello-world",
			expected: "helloWorld",
		},
		{
			name:     "dot case",
			input:    "hello.world",
			expected: "helloWorld",
		},
		{
			name:     "space separated",
			input:    "hello world",
			expected: "helloWorld",
		},
		{
			name:     "mixed delimiters",
			input:    "hello_world-example.test",
			expected: "helloWorldExampleTest",
		},
		{
			name:     "consecutive delimiters",
			input:    "hello__world--test",
			expected: "helloWorldTest",
		},
		{
			name:     "uppercase",
			input:    "HELLO_WORLD",
			expected: "helloWorld",
		},
		{
			name:     "acronyms",
			input:    "API_request",
			expected: "apiRequest",
		},
		{
			name:     "with numbers",
			input:    "user_123_name",
			expected: "user123Name",
		},
		{
			name:     "starting with number",
			input:    "123_test",
			expected: "123Test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToCamelCase(tt.input)
			assert.Equal(t, result, tt.expected, "ToCamelCase/%q", tt.name)
		})
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "already pascal case",
			input:    "HelloWorld",
			expected: "HelloWorld",
		},
		{
			name:     "camel case",
			input:    "helloWorld",
			expected: "HelloWorld",
		},
		{
			name:     "snake case",
			input:    "hello_world",
			expected: "HelloWorld",
		},
		{
			name:     "kebab case",
			input:    "hello-world",
			expected: "HelloWorld",
		},
		{
			name:     "dot case",
			input:    "hello.world",
			expected: "HelloWorld",
		},
		{
			name:     "space separated",
			input:    "hello world",
			expected: "HelloWorld",
		},
		{
			name:     "mixed delimiters",
			input:    "hello_world-example.test",
			expected: "HelloWorldExampleTest",
		},
		{
			name:     "consecutive delimiters",
			input:    "hello__world--test",
			expected: "HelloWorldTest",
		},
		{
			name:     "uppercase",
			input:    "HELLO_WORLD",
			expected: "HelloWorld",
		},
		{
			name:     "acronyms",
			input:    "API_request",
			expected: "ApiRequest",
		},
		{
			name:     "with numbers",
			input:    "user_123_name",
			expected: "User123Name",
		},
		{
			name:     "starting with number",
			input:    "123_test",
			expected: "123Test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToPascalCase(tt.input)
			assert.Equal(t, result, tt.expected, "ToPascalCase/%q", tt.name)
		})
	}
}

func TestToDelimited(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		delimiter string
		expected  string
	}{
		{
			name:      "empty delimiter",
			input:     "Hello World",
			delimiter: "",
			expected:  "helloworld",
		},
		{
			name:      "empty string",
			input:     "",
			delimiter: "_",
			expected:  "",
		},
		{
			name:      "already snake case",
			input:     "hello_world",
			delimiter: "_",
			expected:  "hello_world",
		},
		{
			name:      "already kebab case",
			input:     "hello-world",
			delimiter: "_",
			expected:  "hello_world",
		},
		{
			name:      "already dot case",
			input:     "hello.world",
			delimiter: "_",
			expected:  "hello_world",
		},
		{
			name:      "already space separated",
			input:     "hello world",
			delimiter: "_",
			expected:  "hello_world",
		},
		{
			name:      "camel case",
			input:     "helloWorld",
			delimiter: "_",
			expected:  "hello_world",
		},
		{
			name:      "pascal case",
			input:     "HelloWorld",
			delimiter: "_",
			expected:  "hello_world",
		},
		{
			name:      "uppercase",
			input:     "HELLO_WORLD",
			delimiter: "-",
			expected:  "hello-world",
		},
		{
			name:      "acronyms",
			input:     "APIRequest",
			delimiter: "-",
			expected:  "api-request",
		},
		{
			name:      "with numbers",
			input:     "user123Name",
			delimiter: "-",
			expected:  "user123-name",
		},
		{
			name:      "starting with number",
			input:     "123Test",
			delimiter: "-",
			expected:  "123-test",
		},
		{
			name:      "leading delimiter",
			input:     "-helloWorld",
			delimiter: "-",
			expected:  "hello-world",
		},
		{
			name:      "trailing delimiter",
			input:     "helloWorld-",
			delimiter: "-",
			expected:  "hello-world",
		},
		{
			name:      "multi leading delimiter",
			input:     "---helloWorld",
			delimiter: "-",
			expected:  "hello-world",
		},
		{
			name:      "multi trailing delimiter",
			input:     "helloWorld---",
			delimiter: "-",
			expected:  "hello-world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToDelimited(tt.input, tt.delimiter)
			assert.Equal(t, result, tt.expected, "ToDelimited/%q", tt.name)
		})
	}
}
