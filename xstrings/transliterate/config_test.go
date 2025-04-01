package transliterate_test

import (
	"testing"

	"github.com/neticdk/go-stdlib/xstrings/transliterate"
)

func TestConfiguration(t *testing.T) {
	// Ensure default config is restored after tests potentially modifying it
	originalMaxInputLength := 1 << 20
	originalMaxCacheSize := 1000

	t.Cleanup(func() {
		transliterate.Configure(
			transliterate.WithMaxInputLength(originalMaxInputLength),
			transliterate.WithMaxCacheSize(originalMaxCacheSize),
		)
	})

	t.Run("default configuration", func(t *testing.T) {
		// Try with string just under default 1MB limit
		text := string(make([]rune, 1<<20-1))
		_, err := transliterate.WithLimit(text)
		if err != nil {
			t.Fatalf("should accept string under default limit: unexpected error: %v", err)
		}

		// Try with string over default 1MB limit
		text = string(make([]rune, 1<<20+1))
		_, err = transliterate.WithLimit(text)
		if err == nil {
			t.Fatal("should reject string over default limit: expected an error, but got nil")
		}
	})

	t.Run("configure max input length", func(t *testing.T) {
		// Set a small limit
		transliterate.Configure(transliterate.WithMaxInputLength(15))
		t.Cleanup(func() {
			transliterate.Configure(transliterate.WithMaxInputLength(1 << 20))
		})

		// Try with string under new limit
		text := "Hello世界" // 11 bytes (5 + 3 + 3)
		result, err := transliterate.WithLimit(text)
		if err != nil {
			t.Fatalf("should accept string under configured limit: unexpected error: %v", err)
		}
		if result == "" {
			t.Fatal("expected non-empty result, but got empty string")
		}

		// Try with string over new limit
		text = "Hello世界,SubsequentText"
		_, err = transliterate.WithLimit(text)
		if err == nil {
			t.Fatal("should reject string over configured limit: expected an error, but got nil")
		}
	})

	t.Run("configure cache size", func(t *testing.T) {
		transliterate.ClearCache()
		// Set a very small cache size
		transliterate.Configure(transliterate.WithMaxCacheSize(2))
		t.Cleanup(func() {
			transliterate.Configure(transliterate.WithMaxCacheSize(1000))
			transliterate.ClearCache()
		})

		// Add three characters
		text := "世界人"
		transliterate.String(text) // Triggers cache adds internally

		// Should have cleared cache and only have last character's entry
		size := transliterate.GetCacheSize()
		if size != 1 {
			t.Fatalf("after exceeding max size, cache should be cleared and contain only last entry: expected %d, got %d", 1, size)
		}
	})

	t.Run("multiple configurations", func(t *testing.T) {
		// Configure both cache size and input length
		transliterate.Configure(
			transliterate.WithMaxCacheSize(500),
			transliterate.WithMaxInputLength(100),
		)
		t.Cleanup(func() {
			transliterate.Configure(
				transliterate.WithMaxCacheSize(1000),
				transliterate.WithMaxInputLength(1<<20),
			)
			transliterate.ClearCache()
		})

		// Test input length limit
		text := string(make([]rune, 101))
		_, err := transliterate.WithLimit(text)
		if err == nil {
			t.Fatal("should respect configured input length: expected an error, but got nil")
		}

		// Test cache size (basic check)
		transliterate.ClearCache()
		text = "世界人"
		transliterate.String(text)            // Prime cache
		transliterate.String(text)            // Hit cache
		hits := transliterate.GetCacheStats() // Assuming GetCacheStats returns only hits
		if hits <= uint64(0) {                // Assuming hits is reset or low before this test
			t.Fatalf("cache should be working with new size: expected hits (%d) to be greater than 0", hits)
		}
	})

	t.Run("invalid configurations", func(t *testing.T) {
		// Save current stats and config (assuming accessors exist or know defaults)
		currentMaxInputLength := 1 << 20 // Re-fetch or use known state
		currentMaxCacheSize := 1000      // Re-fetch or use known state

		transliterate.Configure( // Ensure a known state before testing invalid values
			transliterate.WithMaxInputLength(currentMaxInputLength),
			transliterate.WithMaxCacheSize(currentMaxCacheSize),
		)
		transliterate.ClearCache()
		oldHits := transliterate.GetCacheStats() // Get hits after reset

		// Try invalid values - should keep previous configuration
		transliterate.Configure(
			transliterate.WithMaxCacheSize(-1),
			transliterate.WithMaxInputLength(-1),
		)

		// Verify configuration did not change (would need accessors, or test behavior)

		// Test behavior implies old config: cache should still work
		text := "世界"
		transliterate.String(text) // Prime
		transliterate.String(text) // Hit
		hits := transliterate.GetCacheStats()
		if hits <= oldHits {
			t.Fatalf("cache should still work with old config: expected hits (%d) to be greater than oldHits (%d)", hits, oldHits)
		}

		// Test behavior implies old config: long text should still be rejected
		text = string(make([]rune, currentMaxInputLength+1)) // Use known limit + 1
		_, err := transliterate.WithLimit(text)
		if err == nil {
			t.Fatal("should maintain original input length limit: expected an error, but got nil")
		}
	})
}
