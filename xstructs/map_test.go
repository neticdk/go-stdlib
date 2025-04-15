package xstructs_test

import (
	"testing"

	"github.com/neticdk/go-stdlib/assert"
	"github.com/neticdk/go-stdlib/xstructs"
)

type TestingMapAny map[string]any

func TestToMap(t *testing.T) {
	tests := []struct {
		name     string
		data     any
		opts     []xstructs.ToMapOptions
		expected map[string]any
		wantErr  bool
	}{
		{
			name: "simple struct",
			data: struct {
				A int    `json:"a"`
				B string `json:"b"`
				C struct {
					D int `json:"d"`
				} `json:"c"`
				E []string `json:"e,omitempty"`
			}{
				A: 1,
				B: "test",
				C: struct {
					D int `json:"d"`
				}{
					D: 2,
				},
				E: []string{"one", "two"},
			},
			expected: map[string]any{
				"a": 1,
				"b": "test",
				"c": map[string]any{
					"d": 2,
				},
				"e": []any{"one", "two"},
			},
			wantErr: false,
		},
		{
			name: "simple struct as pointer",
			data: &struct {
				A int    `json:"a"`
				B string `json:"b"`
				C struct {
					D int `json:"d"`
				} `json:"c"`
				E []string `json:"e,omitempty"`
			}{
				A: 1,
				B: "test",
				C: struct {
					D int `json:"d"`
				}{
					D: 2,
				},
				E: []string{"one", "two"},
			},
			expected: map[string]any{
				"a": 1,
				"b": "test",
				"c": map[string]any{
					"d": 2,
				},
				"e": []any{"one", "two"},
			},
			wantErr: false,
		},
		{
			name: "with inline",
			data: &struct {
				A int    `json:"a"`
				B string `json:"b"`
				C struct {
					D int `json:"d"`
				} `json:",inline"`
			}{
				A: 1,
				B: "test",
				C: struct {
					D int `json:"d"`
				}{
					D: 2,
				},
			},
			expected: map[string]any{
				"a": 1,
				"b": "test",
				"d": 2,
			},
			wantErr: false,
		},
		{
			name: "TestingMapAny yaml inline",
			data: &struct {
				A             int    `yaml:"a"`
				B             string `yaml:"b"`
				TestingMapAny `yaml:",inline"`
			}{
				A: 1,
				B: "test",
				TestingMapAny: map[string]any{
					"d": 2,
				},
			},
			expected: map[string]any{
				"a": 1,
				"b": "test",
				"d": 2,
			},
			wantErr: false,
		},
		{
			name: "with omit '-'",
			data: &struct {
				A int    `json:"a"`
				B string `json:",-"`
				C struct {
					D int `json:"d"`
				} `json:"-"`
			}{
				A: 1,
				B: "test",
				C: struct {
					D int `json:"d"`
				}{
					D: 2,
				},
			},
			expected: map[string]any{
				"a": 1,
			},
			wantErr: false,
		},
		{
			name: "with omit '-' as name",
			data: &struct {
				A int    `json:"a"`
				B string `json:"b"`
				C struct {
					D int `json:"d"`
				} `json:"-,"`
			}{
				A: 1,
				B: "test",
				C: struct {
					D int `json:"d"`
				}{
					D: 2,
				},
			},
			expected: map[string]any{
				"a": 1,
				"b": "test",
				"-": map[string]any{
					"d": 2,
				},
			},
			wantErr: false,
		},
		{
			name: "nested struct",
			data: struct {
				A struct {
					B struct {
						C int `json:"c"`
					} `json:"b"`
				} `json:"a"`
			}{
				A: struct {
					B struct {
						C int `json:"c"`
					} `json:"b"`
				}{
					B: struct {
						C int `json:"c"`
					}{
						C: 3,
					},
				},
			},
			expected: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c": 3,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "map with string keys",
			data: map[string]any{
				"a": 1,
				"b": "test",
				"c": map[string]any{
					"d": 2,
				},
				"e": []string{"one", "two"},
				"f": struct {
					G int `json:"g"`
				}{
					G: 4,
				},
			},
			expected: map[string]any{
				"a": 1,
				"b": "test",
				"c": map[string]any{
					"d": 2,
				},
				"e": []any{"one", "two"},
				"f": map[string]any{
					"g": 4,
				},
			},
		},
		{
			name: "slice",
			data: []any{
				1,
				"test",
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "scalar",
			data:     42,
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "nil",
			data:     nil,
			expected: nil,
			wantErr:  false,
		},
		{
			name: "unexported field",
			data: struct {
				A int `json:"a"`
				b string
			}{
				A: 1,
				b: "test",
			},
			expected: map[string]any{
				"a": 1,
			},
			wantErr: false,
		},
		{
			name: "with custom tag order",
			data: struct {
				A int    `json:"-" yaml:"a"`
				B string `json:"-" yaml:"b"`
				C struct {
					D int `json:"-" yaml:"d"`
				} `json:"-" yaml:"c"`
				E []string `json:"-" yaml:"e,omitempty"`
			}{
				A: 1,
				B: "test",
				C: struct {
					D int `json:"-" yaml:"d"`
				}{
					D: 2,
				},
				E: []string{"one", "two"},
			},
			opts: []xstructs.ToMapOptions{
				xstructs.WithTags("yaml", "json"),
			},
			expected: map[string]any{
				"a": 1,
				"b": "test",
				"c": map[string]any{
					"d": 2,
				},
				"e": []any{"one", "two"},
			},
			wantErr: false,
		},
		{
			name: "with custom tags none exist",
			data: struct {
				A int    `json:"-" yaml:"a"`
				B string `json:"-" yaml:"b"`
				C struct {
					D int `json:"-" yaml:"d"`
				} `json:"-" yaml:"c"`
				E []string `json:"-" yaml:"e,omitempty"`
			}{
				A: 1,
				B: "test",
				C: struct {
					D int `json:"-" yaml:"d"`
				}{
					D: 2,
				},
				E: []string{"one", "two"},
			},
			opts: []xstructs.ToMapOptions{
				xstructs.WithTags("custom"),
			},
			expected: map[string]any{},
			wantErr:  false,
		},
		{
			name: "with allow no tags",
			data: struct {
				A int
				B string
			}{
				A: 1,
				B: "test",
			},
			opts: []xstructs.ToMapOptions{
				xstructs.WithAllowNoTags(),
			},
			expected: map[string]any{
				"A": 1,
				"B": "test",
			},
			wantErr: false,
		},
		{
			name: "with pointer value",
			data: &struct {
				A int `json:"a"`
				B *struct {
					C int `json:"c"`
				} `json:"b"`
			}{
				A: 1,
				B: &struct {
					C int `json:"c"`
				}{
					C: 2,
				},
			},
			expected: map[string]any{
				"a": 1,
				"b": map[string]any{
					"c": 2,
				},
			},
			wantErr: false,
		},
		{
			name: "with nil pointer value",
			data: &struct {
				A int `json:"a"`
				B *struct {
					C int `json:"c"`
				} `json:"b"`
			}{
				A: 1,
				B: nil,
			},
			expected: map[string]any{
				"a": 1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := xstructs.ToMap(tt.data, tt.opts...)
			if tt.wantErr {
				assert.Error(t, err)
			}
			assert.Equal(t, got, tt.expected)
		})
	}
}
