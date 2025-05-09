package xstrings

import (
	"testing"

	"github.com/neticdk/go-stdlib/assert"
)

func TestSlugifyDefaultOptions(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expected      string
		lowercase     bool
		decamelize    bool
		transliterate bool
	}{
		{
			name:          "empty all true",
			input:         "",
			expected:      "",
			lowercase:     true,
			decamelize:    true,
			transliterate: true,
		},
		{
			name:          "empty all true except transliterate",
			input:         "",
			expected:      "",
			lowercase:     true,
			decamelize:    true,
			transliterate: false,
		},
		{
			name:          "empty all false except lowercase",
			input:         "",
			expected:      "",
			lowercase:     true,
			decamelize:    false,
			transliterate: false,
		},
		{
			name:          "empty all false except transliterate",
			input:         "",
			expected:      "",
			lowercase:     true,
			decamelize:    false,
			transliterate: true,
		},
		{
			name:          "empty all false",
			input:         "",
			expected:      "",
			lowercase:     false,
			decamelize:    false,
			transliterate: false,
		},
		{
			name:          "empty all false except transliterate",
			input:         "",
			expected:      "",
			lowercase:     false,
			decamelize:    false,
			transliterate: true,
		},
		{
			name:          "empty all false except decamelize and transliterate",
			input:         "",
			expected:      "",
			lowercase:     false,
			decamelize:    true,
			transliterate: true,
		},
		{
			name:          "empty all false except decamelize",
			input:         "",
			expected:      "",
			lowercase:     false,
			decamelize:    true,
			transliterate: false,
		},
		{
			name:          "hello world all true",
			input:         "hello world",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    true,
			transliterate: true,
		},
		{
			name:          "hello world all true except transliterate",
			input:         "hello world",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    true,
			transliterate: false,
		},
		{
			name:          "hello world all false except lowercase",
			input:         "hello world",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    false,
			transliterate: false,
		},
		{
			name:          "hello world all false except transliterate",
			input:         "hello world",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    false,
			transliterate: true,
		},
		{
			name:          "hello world all false",
			input:         "hello world",
			expected:      "hello-world",
			lowercase:     false,
			decamelize:    false,
			transliterate: false,
		},
		{
			name:          "hello world all false except transliterate",
			input:         "hello world",
			expected:      "hello-world",
			lowercase:     false,
			decamelize:    false,
			transliterate: true,
		},
		{
			name:          "hello world all false except decamelize and transliterate",
			input:         "hello world",
			expected:      "hello-world",
			lowercase:     false,
			decamelize:    true,
			transliterate: true,
		},
		{
			name:          "hello world all false except decamelize",
			input:         "hello world",
			expected:      "hello-world",
			lowercase:     false,
			decamelize:    true,
			transliterate: false,
		},
		{
			name:          "hello double space all true",
			input:         "hello  world",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    true,
			transliterate: true,
		},
		{
			name:          "hello double space all true except transliterate",
			input:         "hello  world",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    true,
			transliterate: false,
		},
		{
			name:          "hello double space all false except lowercase",
			input:         "hello  world",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    false,
			transliterate: false,
		},
		{
			name:          "hello double space all false except transliterate",
			input:         "hello  world",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    false,
			transliterate: true,
		},
		{
			name:          "hello double space all false",
			input:         "hello  world",
			expected:      "hello-world",
			lowercase:     false,
			decamelize:    false,
			transliterate: false,
		},
		{
			name:          "hello double space all false except transliterate",
			input:         "hello  world",
			expected:      "hello-world",
			lowercase:     false,
			decamelize:    false,
			transliterate: true,
		},
		{
			name:          "hello double space all false except decamelize and transliterate",
			input:         "hello  world",
			expected:      "hello-world",
			lowercase:     false,
			decamelize:    true,
			transliterate: true,
		},
		{
			name:          "hello double space all false except decamelize",
			input:         "hello  world",
			expected:      "hello-world",
			lowercase:     false,
			decamelize:    true,
			transliterate: false,
		},
		{
			name:          "hello triple dash all true",
			input:         "hello---world",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    true,
			transliterate: true,
		},
		{
			name:          "hello triple dash all true except transliterate",
			input:         "hello---world",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    true,
			transliterate: false,
		},
		{
			name:          "hello triple dash all false except lowercase",
			input:         "hello---world",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    false,
			transliterate: false,
		},
		{
			name:          "hello triple dash all false except transliterate",
			input:         "hello---world",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    false,
			transliterate: true,
		},
		{
			name:          "hello triple dash all false",
			input:         "hello---world",
			expected:      "hello-world",
			lowercase:     false,
			decamelize:    false,
			transliterate: false,
		},
		{
			name:          "hello triple dash all false except transliterate",
			input:         "hello---world",
			expected:      "hello-world",
			lowercase:     false,
			decamelize:    false,
			transliterate: true,
		},
		{
			name:          "hello triple dash all false except decamelize and transliterate",
			input:         "hello---world",
			expected:      "hello-world",
			lowercase:     false,
			decamelize:    true,
			transliterate: true,
		},
		{
			name:          "hello triple dash all false except decamelize",
			input:         "hello---world",
			expected:      "hello-world",
			lowercase:     false,
			decamelize:    true,
			transliterate: false,
		},
		{
			name:          "hello space dash space all true",
			input:         "hello- -world",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    true,
			transliterate: true,
		},
		{
			name:          "hello space dash space all true except transliterate",
			input:         "hello- -world",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    true,
			transliterate: false,
		},
		{
			name:          "hello space dash space all false except lowercase",
			input:         "hello- -world",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    false,
			transliterate: false,
		},
		{
			name:          "hello space dash space all false except transliterate",
			input:         "hello- -world",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    false,
			transliterate: true,
		},
		{
			name:          "hello space dash space all false",
			input:         "hello- -world",
			expected:      "hello-world",
			lowercase:     false,
			decamelize:    false,
			transliterate: false,
		},
		{
			name:          "hello space dash space all false except transliterate",
			input:         "hello- -world",
			expected:      "hello-world",
			lowercase:     false,
			decamelize:    false,
			transliterate: true,
		},
		{
			name:          "hello space dash space all false except decamelize and transliterate",
			input:         "hello- -world",
			expected:      "hello-world",
			lowercase:     false,
			decamelize:    true,
			transliterate: true,
		},
		{
			name:          "hello space dash space all false except decamelize",
			input:         "hello- -world",
			expected:      "hello-world",
			lowercase:     false,
			decamelize:    true,
			transliterate: false,
		},
		{
			name:          "hello space dash space all true",
			input:         "hello - - world",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    true,
			transliterate: true,
		},
		{
			name:          "hello space dash space all true except transliterate",
			input:         "hello - - world",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    true,
			transliterate: false,
		},
		{
			name:          "hello space dash space all false except lowercase",
			input:         "hello - - world",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    false,
			transliterate: false,
		},
		{
			name:          "hello space dash space all false except transliterate",
			input:         "hello - - world",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    false,
			transliterate: true,
		},
		{
			name:          "hello space dash space all false",
			input:         "hello - - world",
			expected:      "hello-world",
			lowercase:     false,
			decamelize:    false,
			transliterate: false,
		},
		{
			name:          "hello space dash space all false except transliterate",
			input:         "hello - - world",
			expected:      "hello-world",
			lowercase:     false,
			decamelize:    false,
			transliterate: true,
		},
		{
			name:          "hello space dash space all false except decamelize and transliterate",
			input:         "hello - - world",
			expected:      "hello-world",
			lowercase:     false,
			decamelize:    true,
			transliterate: true,
		},
		{
			name:          "hello space dash space all false except decamelize",
			input:         "hello - - world",
			expected:      "hello-world",
			lowercase:     false,
			decamelize:    true,
			transliterate: false,
		},
		{
			name:          "dashes all true",
			input:         "-------",
			expected:      "",
			lowercase:     true,
			decamelize:    true,
			transliterate: true,
		},
		{
			name:          "dashes all true except transliterate",
			input:         "-------",
			expected:      "",
			lowercase:     true,
			decamelize:    true,
			transliterate: false,
		},
		{
			name:          "dashes all false except lowercase",
			input:         "-------",
			expected:      "",
			lowercase:     true,
			decamelize:    false,
			transliterate: false,
		},
		{
			name:          "dashes all false except transliterate",
			input:         "-------",
			expected:      "",
			lowercase:     true,
			decamelize:    false,
			transliterate: true,
		},
		{
			name:          "dashes all false",
			input:         "-------",
			expected:      "",
			lowercase:     false,
			decamelize:    false,
			transliterate: false,
		},
		{
			name:          "dashes all false except transliterate",
			input:         "-------",
			expected:      "",
			lowercase:     false,
			decamelize:    false,
			transliterate: true,
		},
		{
			name:          "dashes all false except decamelize and transliterate",
			input:         "-------",
			expected:      "",
			lowercase:     false,
			decamelize:    true,
			transliterate: true,
		},
		{
			name:          "dashes all false except decamelize",
			input:         "-------",
			expected:      "",
			lowercase:     false,
			decamelize:    true,
			transliterate: false,
		},
		{
			name:          "dash hello world all true",
			input:         "-hello world",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    true,
			transliterate: true,
		},
		{
			name:          "dash hello world all true except transliterate",
			input:         "-hello world",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    true,
			transliterate: false,
		},
		{
			name:          "dash hello world all false except lowercase",
			input:         "-hello world",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    false,
			transliterate: false,
		},
		{
			name:          "dash hello world all false except transliterate",
			input:         "-hello world",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    false,
			transliterate: true,
		},
		{
			name:          "dash hello world all false",
			input:         "-hello world",
			expected:      "hello-world",
			lowercase:     false,
			decamelize:    false,
			transliterate: false,
		},
		{
			name:          "dash hello world all false except transliterate",
			input:         "-hello world",
			expected:      "hello-world",
			lowercase:     false,
			decamelize:    false,
			transliterate: true,
		},
		{
			name:          "dash hello world all false except decamelize and transliterate",
			input:         "-hello world",
			expected:      "hello-world",
			lowercase:     false,
			decamelize:    true,
			transliterate: true,
		},
		{
			name:          "dash hello world all false except decamelize",
			input:         "-hello world",
			expected:      "hello-world",
			lowercase:     false,
			decamelize:    true,
			transliterate: false,
		},
		{
			name:          "space hello world space all true",
			input:         " hello world ",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    true,
			transliterate: true,
		},
		{
			name:          "space hello world space all true except transliterate",
			input:         " hello world ",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    true,
			transliterate: false,
		},
		{
			name:          "space hello world space all false except lowercase",
			input:         " hello world ",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    false,
			transliterate: false,
		},
		{
			name:          "space hello world space all false except transliterate",
			input:         " hello world ",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    false,
			transliterate: true,
		},
		{
			name:          "space hello world space all false",
			input:         " hello world ",
			expected:      "hello-world",
			lowercase:     false,
			decamelize:    false,
			transliterate: false,
		},
		{
			name:          "space hello world space all false except transliterate",
			input:         " hello world ",
			expected:      "hello-world",
			lowercase:     false,
			decamelize:    false,
			transliterate: true,
		},
		{
			name:          "space hello world space all false except decamelize and transliterate",
			input:         " hello world ",
			expected:      "hello-world",
			lowercase:     false,
			decamelize:    true,
			transliterate: true,
		},
		{
			name:          "space hello world space all false except decamelize",
			input:         " hello world ",
			expected:      "hello-world",
			lowercase:     false,
			decamelize:    true,
			transliterate: false,
		},
		{
			name:          "special characters all true",
			input:         "!@#$%^&*()_+=-`~[]\\{}|;':\",./<>?",
			expected:      "",
			lowercase:     true,
			decamelize:    true,
			transliterate: true,
		},
		{
			name:          "special characters all true except transliterate",
			input:         "!@#$%^&*()_+=-`~[]\\{}|;':\",./<>?",
			expected:      "",
			lowercase:     true,
			decamelize:    true,
			transliterate: false,
		},
		{
			name:          "special characters all false except lowercase",
			input:         "!@#$%^&*()_+=-`~[]\\{}|;':\",./<>?",
			expected:      "",
			lowercase:     true,
			decamelize:    false,
			transliterate: false,
		},
		{
			name:          "special characters all false except transliterate",
			input:         "!@#$%^&*()_+=-`~[]\\{}|;':\",./<>?",
			expected:      "",
			lowercase:     true,
			decamelize:    false,
			transliterate: true,
		},
		{
			name:          "special characters all false",
			input:         "!@#$%^&*()_+=-`~[]\\{}|;':\",./<>?",
			expected:      "",
			lowercase:     false,
			decamelize:    false,
			transliterate: false,
		},
		{
			name:          "special characters all false except transliterate",
			input:         "!@#$%^&*()_+=-`~[]\\{}|;':\",./<>?",
			expected:      "",
			lowercase:     false,
			decamelize:    false,
			transliterate: true,
		},
		{
			name:          "special characters all false except decamelize and transliterate",
			input:         "!@#$%^&*()_+=-`~[]\\{}|;':\",./<>?",
			expected:      "",
			lowercase:     false,
			decamelize:    true,
			transliterate: true,
		},
		{
			name:          "special characters all false except decamelize",
			input:         "!@#$%^&*()_+=-`~[]\\{}|;':\",./<>?",
			expected:      "",
			lowercase:     false,
			decamelize:    true,
			transliterate: false,
		},
		{
			name:          "unicode all true",
			input:         "HélLø Wörld",
			expected:      "hel-lo-world",
			lowercase:     true,
			decamelize:    true,
			transliterate: true,
		},
		{
			name:          "unicode all true except transliterate",
			input:         "HélLø Wörld",
			expected:      "h-l-l-w-rld",
			lowercase:     true,
			decamelize:    true,
			transliterate: false,
		},
		{
			name:          "unicode all false except lowercase",
			input:         "HélLø Wörld",
			expected:      "h-ll-w-rld",
			lowercase:     true,
			decamelize:    false,
			transliterate: false,
		},
		{
			name:          "unicode all false except transliterate",
			input:         "HélLø Wörld",
			expected:      "hello-world",
			lowercase:     true,
			decamelize:    false,
			transliterate: true,
		},
		{
			name:          "unicode all false",
			input:         "HélLø Wörld",
			expected:      "H-lL-W-rld",
			lowercase:     false,
			decamelize:    false,
			transliterate: false,
		},
		{
			name:          "unicode all false except transliterate",
			input:         "HélLø Wörld",
			expected:      "HelLo-World",
			lowercase:     false,
			decamelize:    false,
			transliterate: true,
		},
		{
			name:          "unicode all false except decamelize and transliterate",
			input:         "HélLø Wörld",
			expected:      "hel-lo-world",
			lowercase:     false,
			decamelize:    true,
			transliterate: true,
		},
		{
			name:          "unicode all false except decamelize",
			input:         "HélLø Wörld",
			expected:      "h-l-l-w-rld",
			lowercase:     false,
			decamelize:    true,
			transliterate: false,
		},
		{
			name:          "dots all true",
			input:         "with.dots.everywhere",
			expected:      "with-dots-everywhere",
			lowercase:     true,
			decamelize:    true,
			transliterate: true,
		},
		{
			name:          "dots all true except transliterate",
			input:         "with.dots.everywhere",
			expected:      "with-dots-everywhere",
			lowercase:     true,
			decamelize:    true,
			transliterate: false,
		},
		{
			name:          "dots all false except lowercase",
			input:         "with.dots.everywhere",
			expected:      "with-dots-everywhere",
			lowercase:     true,
			decamelize:    false,
			transliterate: false,
		},
		{
			name:          "dots all false except transliterate",
			input:         "with.dots.everywhere",
			expected:      "with-dots-everywhere",
			lowercase:     true,
			decamelize:    false,
			transliterate: true,
		},
		{
			name:          "dots all false",
			input:         "with.dots.everywhere",
			expected:      "with-dots-everywhere",
			lowercase:     false,
			decamelize:    false,
			transliterate: false,
		},
		{
			name:          "dots all false except transliterate",
			input:         "with.dots.everywhere",
			expected:      "with-dots-everywhere",
			lowercase:     false,
			decamelize:    false,
			transliterate: true,
		},
		{
			name:          "dots all false except decamelize and transliterate",
			input:         "with.dots.everywhere",
			expected:      "with-dots-everywhere",
			lowercase:     false,
			decamelize:    true,
			transliterate: true,
		},
		{
			name:          "dots all false except decamelize",
			input:         "with.dots.everywhere",
			expected:      "with-dots-everywhere",
			lowercase:     false,
			decamelize:    true,
			transliterate: false,
		},
		{
			name:          "numbers all true",
			input:         "12345",
			expected:      "12345",
			lowercase:     true,
			decamelize:    true,
			transliterate: true,
		},
		{
			name:          "numbers all true except transliterate",
			input:         "12345",
			expected:      "12345",
			lowercase:     true,
			decamelize:    true,
			transliterate: false,
		},
		{
			name:          "numbers all false except lowercase",
			input:         "12345",
			expected:      "12345",
			lowercase:     true,
			decamelize:    false,
			transliterate: false,
		},
		{
			name:          "numbers all false except transliterate",
			input:         "12345",
			expected:      "12345",
			lowercase:     true,
			decamelize:    false,
			transliterate: true,
		},
		{
			name:          "numbers all false",
			input:         "12345",
			expected:      "12345",
			lowercase:     false,
			decamelize:    false,
			transliterate: false,
		},
		{
			name:          "numbers all false except transliterate",
			input:         "12345",
			expected:      "12345",
			lowercase:     false,
			decamelize:    false,
			transliterate: true,
		},
		{
			name:          "numbers all false except decamelize and transliterate",
			input:         "12345",
			expected:      "12345",
			lowercase:     false,
			decamelize:    true,
			transliterate: true,
		},
		{
			name:          "numbers all false except decamelize",
			input:         "12345",
			expected:      "12345",
			lowercase:     false,
			decamelize:    true,
			transliterate: false,
		},
		{
			name:          "mixed case all true",
			input:         "MixedCase123",
			expected:      "mixed-case123",
			lowercase:     true,
			decamelize:    true,
			transliterate: true,
		},
		{
			name:          "mixed case all true except transliterate",
			input:         "MixedCase123",
			expected:      "mixed-case123",
			lowercase:     true,
			decamelize:    true,
			transliterate: false,
		},
		{
			name:          "mixed case all false except lowercase",
			input:         "MixedCase123",
			expected:      "mixedcase123",
			lowercase:     true,
			decamelize:    false,
			transliterate: false,
		},
		{
			name:          "mixed case all false except transliterate",
			input:         "MixedCase123",
			expected:      "mixedcase123",
			lowercase:     true,
			decamelize:    false,
			transliterate: true,
		},
		{
			name:          "mixed case all false",
			input:         "MixedCase123",
			expected:      "MixedCase123",
			lowercase:     false,
			decamelize:    false,
			transliterate: false,
		},
		{
			name:          "mixed case all false except transliterate",
			input:         "MixedCase123",
			expected:      "MixedCase123",
			lowercase:     false,
			decamelize:    false,
			transliterate: true,
		},
		{
			name:          "mixed case all false except decamelize and transliterate",
			input:         "MixedCase123",
			expected:      "mixed-case123",
			lowercase:     false,
			decamelize:    true,
			transliterate: true,
		},
		{
			name:          "mixed case all false except decamelize",
			input:         "MixedCase123",
			expected:      "mixed-case123",
			lowercase:     false,
			decamelize:    true,
			transliterate: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			actual := Slugify(
				tt.input,
				WithDecamelize(tt.decamelize),
				WithLowercase(tt.lowercase),
				WithTransliterate(tt.transliterate),
			)
			assert.Equal(t, actual, tt.expected, "Slugify/%q", tt.name)
		})
	}
}
