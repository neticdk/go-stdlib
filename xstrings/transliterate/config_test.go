package transliterate_test

import (
	"testing"

	"github.com/neticdk/go-stdlib/xstrings/transliterate"
	"github.com/stretchr/testify/assert"
)

func TestConfiguration(t *testing.T) {
	t.Run("default configuration", func(t *testing.T) {
		// Try with string just under default 1MB limit
		text := string(make([]rune, 1<<20-1))
		_, err := transliterate.WithLimit(text)
		assert.NoError(t, err, "should accept string under default limit")

		// Try with string over default 1MB limit
		text = string(make([]rune, 1<<20+1))
		_, err = transliterate.WithLimit(text)
		assert.Error(t, err, "should reject string over default limit")
	})

	t.Run("configure max input length", func(t *testing.T) {
		// Set a small limit
		transliterate.Configure(transliterate.WithMaxInputLength(15))

		// Try with string under new limit
		text := "Hello世界" // 11 bytes (5 + 3 + 3)
		result, err := transliterate.WithLimit(text)
		assert.NoError(t, err, "should accept string under configured limit")
		assert.NotEmpty(t, result)

		// Try with string over new limit
		text = "Hello世界,SubsequentText"
		_, err = transliterate.WithLimit(text)
		assert.Error(t, err, "should reject string over configured limit")

		// Reset to default for other tests
		transliterate.Configure(transliterate.WithMaxInputLength(1 << 20))
	})

	t.Run("configure cache size", func(t *testing.T) {
		transliterate.ClearCache()

		// Set a very small cache size
		transliterate.Configure(transliterate.WithMaxCacheSize(2))

		// Add three characters
		text := "世界人"
		transliterate.String(text)

		// Should have cleared cache and only have last character
		size := transliterate.GetCacheSize()
		assert.Equal(t, 1, size, "after exceeding max size, cache should be cleared and contain only last entry")

		// Reset to default for other tests
		transliterate.Configure(transliterate.WithMaxCacheSize(1000))
	})

	t.Run("multiple configurations", func(t *testing.T) {
		// Configure both cache size and input length
		transliterate.Configure(
			transliterate.WithMaxCacheSize(500),
			transliterate.WithMaxInputLength(100),
		)

		// Test input length limit
		text := string(make([]rune, 101))
		_, err := transliterate.WithLimit(text)
		assert.Error(t, err, "should respect configured input length")

		// Test cache size (basic check)
		transliterate.ClearCache()
		text = "世界人"
		transliterate.String(text)
		transliterate.String(text)
		hits := transliterate.GetCacheStats()
		assert.Greater(t, hits, uint64(0), "cache should be working with new size")

		// Reset to defaults
		transliterate.Configure(
			transliterate.WithMaxCacheSize(1000),
			transliterate.WithMaxInputLength(1<<20),
		)
	})

	t.Run("invalid configurations", func(t *testing.T) {
		// Save current stats
		oldHits := transliterate.GetCacheStats()

		// Try invalid values - should keep previous configuration
		transliterate.Configure(
			transliterate.WithMaxCacheSize(-1),
			transliterate.WithMaxInputLength(-1),
		)

		text := "世界"
		transliterate.String(text)
		transliterate.String(text)

		// Cache should still work with old configuration
		hits := transliterate.GetCacheStats()
		assert.Greater(t, hits, oldHits, "cache should still work with old config")

		// Long text should still be rejected with original limit
		text = string(make([]rune, 1<<20+1))
		_, err := transliterate.WithLimit(text)
		assert.Error(t, err, "should maintain original input length limit")
	})
}
