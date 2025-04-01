package xstrings

import (
	"testing"
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
			if result != tt.expected {
				t.Errorf("ToKebabCase(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"already camel case", "helloWorld", "helloWorld"},
		{"snake case", "hello_world", "helloWorld"},
		{"kebab case", "hello-world", "helloWorld"},
		{"dot case", "hello.world", "helloWorld"},
		{"space separated", "hello world", "helloWorld"},
		{"mixed delimiters", "hello_world-example.test", "helloWorldExampleTest"},
		{"consecutive delimiters", "hello__world--test", "helloWorldTest"},
		{"uppercase", "HELLO_WORLD", "helloWorld"},
		{"acronyms", "API_request", "apiRequest"},
		{"with numbers", "user_123_name", "user123Name"},
		{"starting with number", "123_test", "123Test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToCamelCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToCamelCase(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"already pascal case", "HelloWorld", "HelloWorld"},
		{"camel case", "helloWorld", "HelloWorld"},
		{"snake case", "hello_world", "HelloWorld"},
		{"kebab case", "hello-world", "HelloWorld"},
		{"dot case", "hello.world", "HelloWorld"},
		{"space separated", "hello world", "HelloWorld"},
		{"mixed delimiters", "hello_world-example.test", "HelloWorldExampleTest"},
		{"consecutive delimiters", "hello__world--test", "HelloWorldTest"},
		{"uppercase", "HELLO_WORLD", "HelloWorld"},
		{"acronyms", "API_request", "ApiRequest"},
		{"with numbers", "user_123_name", "User123Name"},
		{"starting with number", "123_test", "123Test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToPascalCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToPascalCase(%q) = %q; want %q", tt.input, result, tt.expected)
			}
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
		{"empty delimiter", "Hello World", "", "helloworld"},
		{"empty string", "", "_", ""},
		{"already snake case", "hello_world", "_", "hello_world"},
		{"already kebab case", "hello-world", "_", "hello_world"},
		{"already dot case", "hello.world", "_", "hello_world"},
		{"already space separated", "hello world", "_", "hello_world"},
		{"camel case", "helloWorld", "_", "hello_world"},
		{"pascal case", "HelloWorld", "_", "hello_world"},
		{"uppercase", "HELLO_WORLD", "-", "hello-world"},
		{"acronyms", "APIRequest", "-", "api-request"},
		{"with numbers", "user123Name", "-", "user-123-name"},
		{"starting with number", "123Test", "-", "123-test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToDelimited(tt.input, tt.delimiter)
			if result != tt.expected {
				t.Errorf("ToDelimited(%q, %q) = %q; want %q", tt.input, tt.delimiter, result, tt.expected)
			}
		})
	}
}
