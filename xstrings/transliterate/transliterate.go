package transliterate

import (
	"fmt"
	"strings"
	"sync"
	"unicode"

	"github.com/neticdk/stdlib/xstrings/transliterate/internal/table"
)

//go:generate go run tools/make_tables/main.go

var decodingOnce sync.Once

// bufferPool maintains a pool of strings.Builder instances for reuse
// The buffer is reused to avoid unnecessary allocations and is useful in
// high-concurrency scenarios.
var bufferPool = sync.Pool{
	New: func() any {
		return &strings.Builder{}
	},
}

// String transliterates a Unicode string into its closest ASCII representation.
// For example, "é" becomes "e". Characters without a known approximation are omitted.
// Invalid UTF-8 sequences are also omitted.
func String(s string) string {
	// Lazily initialize the transliteration tables exactly once.
	decodingOnce.Do(decodeTransliterations)

	// Get buffer from pool
	// We can safely assume this type assertion as we fully control the pool's
	// New function and Put operations
	sb := bufferPool.Get().(*strings.Builder) //nolint
	sb.Reset()
	defer bufferPool.Put(sb)

	// Pre-allocate builder capacity using a heuristic to reduce re-allocations.
	// The factor 1.5 assumes some expansion but might need tuning for specific data.
	sb.Grow(len(s) + len(s)/2)

	for _, r := range s {
		sb.WriteString(getTransliteration(r))
	}
	return sb.String()
}

// WithLimit transliterates a Unicode string into its closest ASCII
// representation, but limits the input string length to prevent excessive
// memory usage.
// For example, "é" becomes "e". Characters without a known approximation are
// omitted.
// Invalid UTF-8 sequences are also omitted.
func WithLimit(s string) (string, error) {
	if len(s) > defaultConfig.maxInputLen {
		return "", fmt.Errorf("input string too long: %d > %d",
			len(s), defaultConfig.maxInputLen)
	}
	return String(s), nil
}

// getTransliteration handles the actual character lookup/translation
func getTransliteration(r rune) string {
	// Fast path for ASCII
	if r > 0 && r <= unicode.MaxASCII {
		return string(r)
	}

	// Check cache first
	if cached, ok := defaultCache.get(r); ok {
		return cached
	}

	// Calculate table lookup
	section := r >> 8
	position := r & 0xFF //revive:disable-line:add-constant

	var result string
	if tb, ok := table.Tables[section]; ok {
		// Technically redundant check since position is always within the
		// bounds of the table. But we include it for clarity and to avoid
		// potential future changes that might alter this behavior.
		if len(tb) > int(position) {
			result = tb[position]
			// If tb[position] is "", the character is effectively skipped.
		}
	}

	// Cache the result
	defaultCache.set(r, result)
	return result
}
