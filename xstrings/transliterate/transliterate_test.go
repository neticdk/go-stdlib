package transliterate

import (
	"strings"
	"testing"

	"github.com/neticdk/go-stdlib/assert"
)

func TestTransliterate(t *testing.T) {
	// Note: Expected values depend heavily on the specific transliteration tables used.
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Basic & Existing
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "ascii only",
			input:    "Hello World 123!@#",
			expected: "Hello World 123!@#",
		},
		{
			name:     "basic accented characters",
			input:    "Héllò Wórld",
			expected: "Hello World",
		},
		{
			name:     "mixed basic",
			input:    "Café 123",
			expected: "Cafe 123",
		},
		{
			name:     "basic symbols",
			input:    "©®™",
			expected: "(c)(r)(tm)",
		},

		// More Latin Variants
		{
			name:     "german umlauts",
			input:    "Königstraße",
			expected: "Konigstrasse",
		},
		{
			name:     "french cedilla accent",
			input:    "français garçon",
			expected: "francais garcon",
		},
		{
			name:     "spanish tilde",
			input:    "español",
			expected: "espanol",
		},
		{
			name:     "scandinavian",
			input:    "Smørrebrød Århus Æ Ø",
			expected: "Smorrebrod Aarhus AE O", // Or AE OE
		},
		{
			name:     "eastern european",
			input:    "Česká Šibenik Žilina Łódź",
			expected: "Ceska Sibenik Zilina Lodz",
		},
		{
			name:     "vietnamese",
			input:    "Tiếng Việt",
			expected: "Tieng Viet",
		},

		// Ligatures
		{
			name:     "latin fi fl ligatures",
			input:    "ﬁ ﬂ",
			expected: "fi fl",
		},
		{
			name:     "latin oe ligature",
			input:    "Œdipe œuvre",
			expected: "OEdipe oeuvre", // Or Oedipe oeuvre
		},

		// Other Scripts
		{
			name:     "cyrillic",
			input:    "Привет мир",
			expected: "Privet mir",
		},
		{
			name:     "greek",
			input:    "Αλφαβητικός κατάλογος",
			expected: "Alphabetikos katalogos", // Approximation may vary
		},
		{
			name:     "arabic",
			input:    "مرحبا بالعالم",
			expected: "mrHb bl`lm", // Output heavily depends on table style
		},
		{
			name:     "hebrew",
			input:    "שלום עולם",
			expected: "SHlvm `vlm", // Output heavily depends on table style
		},
		{
			name:     "japanese hiragana katakana",
			input:    "こんにちは、カタカナ",
			expected: "konnichiha, katakana",
		},
		{
			name:     "korean hangul",
			input:    "안녕하세요",
			expected: "annyeonghaseyo",
		},
		{
			name:     "thai",
			input:    "สวัสดี ประเทศไทย",
			expected: "swasdii praethsaithy", // Approximation may vary
		},
		{
			name:     "hindi devanagari",
			input:    "परीक्षा",
			expected: "priikssaa",
		},
		{
			name:     "chinese hanzi",
			input:    "你好，世界",
			expected: "Ni Hao ,Shi Jie ",
		},
		{
			name:     "mixed complex",
			input:    "Café परीक्षा 测试 résumé 123",
			expected: "Cafe priikssaa Ce Shi  resume 123",
		},

		// More Symbols
		{
			name:     "currency",
			input:    "€ £ ¥ ₹ Ƀ",
			expected: "EUR PS Y= Rs B",
		},
		{
			name:     "math symbols",
			input:    "≠ ≤ ≥ × ÷ ∑ ∫ ∞",
			expected: "!= <= >= x /   ",
		},
		{
			name:     "more punctuation",
			input:    "… „“ «» — –",
			expected: "... \"\" <<>> -- -", // Common fallback
		},

		// Edge Cases
		{
			name:     "invalid utf8 start",
			input:    "\xf0\x90\x80test", // Incomplete 4-byte sequence
			expected: "test",             // Assuming replacement char (from range loop) is omitted
		},
		{
			name:     "invalid utf8 middle",
			input:    "valid\xe2\x28\xa1invalid", // Malformed sequence E2 28 A1
			expected: "valid(invalid",            // Correct: range finds '(' amidst invalid bytes
		},
		{
			name:     "invalid utf8 end",
			input:    "test\xf0\x90\x80\x80", // Valid sequence (U+10000)
			expected: "test",                 // Or "" or " " depending on table for U+10000
		},
		{
			name: "valid but overlong encoding",
			// C0 80 is an overlong encoding of U+0000 (NUL) - range should treat as invalid
			input:    "test\xc0\x80test",
			expected: "testtest", // Assuming replacement char is omitted
		},
		{
			name:     "explicit replacement char",
			input:    "Hello\uFFFDWorld",
			expected: "HelloWorld", // Assuming U+FFFD has "" mapping in tables
		},
		{
			name:     "null character",
			input:    "Hello\x00World",
			expected: "HelloWorld", // Assuming U+0000 has "" mapping
		},
		{
			name:     "supplementary plane emoji",
			input:    "Test 👍 Test", // U+1F44D
			expected: "Test  Test",  // Assuming no specific mapping, results in ""
		},
		{
			name:     "supplementary plane cjk",
			input:    "Test \U00020000 Test", // U+20000 (CJK Ext B)
			expected: "Test  Test",           // Assuming no specific mapping, results in ""
		},
		{
			name:     "high private use area",
			input:    "Test \U0010FFFD Test", // U+10FFFD (Max PUA)
			expected: "Test  Test",           // Assuming no specific mapping
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := String(tt.input)
			assert.Equal(t, result, tt.expected, "Transliteration mismatch/%q", tt.name)
		})
	}
}

func TestWithLimit(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:        "empty string",
			input:       "",
			expected:    "",
			expectError: false,
		},
		{
			name:        "normal string",
			input:       "Héllò Wórld",
			expected:    "Hello World",
			expectError: false,
		},
		{
			name:        "string at limit",
			input:       strings.Repeat("a", 1<<20),
			expected:    strings.Repeat("a", 1<<20),
			expectError: false,
		},
		{
			name:        "string exceeding limit",
			input:       strings.Repeat("a", (1<<20)+1),
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := WithLimit(tt.input)
			if tt.expectError {
				assert.NotNil(t, err, "WithLimit/%q", tt.name)
				assert.Empty(t, result, "WithLimit/%q", tt.name)
			} else {
				assert.NoError(t, err, "WithLimit/%q", tt.name)
				assert.Equal(t, result, tt.expected, "WithLimit/%q", tt.name)
			}
		})
	}
}

func BenchmarkString(b *testing.B) {
	benchmarks := []struct {
		name  string
		input string
	}{
		{
			name:  "ascii_only",
			input: strings.Repeat("Hello World! ", 100),
		},
		{
			name:  "latin_accents",
			input: strings.Repeat("Héllò Wórld! ", 100),
		},
		{
			name:  "mixed_scripts",
			input: strings.Repeat("Hello 你好 こんにちは Café", 50),
		},
		{
			name:  "cyrillic",
			input: strings.Repeat("Привет мир! ", 100),
		},
		{
			name:  "empty_string",
			input: "",
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				_ = String(bm.input)
			}
		})
	}
}
