package transliterate_test

import (
	"sync"
	"testing"

	"github.com/neticdk/go-stdlib/xstrings/transliterate"
)

func TestCacheStats(t *testing.T) {
	// Setup
	transliterate.ClearCache()
	transliterate.ResetCacheStats()

	// Helper function to read stats
	getStats := func() (hits uint64) {
		return transliterate.GetCacheStats()
	}

	t.Run("initial stats should be zero", func(t *testing.T) {
		hits := getStats()
		if hits != 0 {
			t.Errorf("initial hits should be zero, got %d", hits)
		}
	})

	t.Run("cache misses on first access", func(t *testing.T) {
		transliterate.ResetCacheStats()

		startHits := getStats()

		// First time translation
		text := "世界"
		transliterate.String(text)

		hits := getStats()
		if hits != startHits {
			t.Errorf("shouldn't have any cache hits, got %d, expected %d", hits, startHits)
		}
	})

	t.Run("cache hits on repeated access", func(t *testing.T) {
		transliterate.ResetCacheStats()

		text := "世界"
		transliterate.String(text) // First time
		transliterate.String(text) // Second time
		transliterate.String(text) // Third time

		hits := getStats()
		expectedHits := uint64(6) // 2 chars × 3 repeated calls after the first
		if hits != expectedHits {
			t.Errorf("should have %d hits (2 chars × 3 repeated calls), got %d", expectedHits, hits)
		}
	})

	t.Run("mixed ascii and non-ascii", func(t *testing.T) {
		transliterate.ResetCacheStats()

		startHits := getStats()

		text := "Hello 世界!"
		transliterate.String(text) // First time
		transliterate.String(text) // Second time

		hits := getStats()
		// ASCII characters shouldn't affect cache stats
		expectedNewHits := uint64(4) // 2 non-ASCII chars × 2 calls after cache population
		expectedTotalHits := startHits + expectedNewHits
		if hits != expectedTotalHits {
			t.Errorf("should only count cache hits for non-ASCII chars, expected %d, got %d", expectedTotalHits, hits)
		}
	})

	t.Run("concurrent access", func(t *testing.T) {
		transliterate.ResetCacheStats()

		startHits := getStats()

		// Run concurrent translations
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			for range 100 {
				transliterate.String("世界")
			}
		}()
		go func() {
			defer wg.Done()
			for range 100 {
				transliterate.String("世界")
			}
		}()

		// Wait for both goroutines
		wg.Wait()

		hits := getStats()
		// The exact number can vary due to race conditions on the first few calls,
		// but it should be substantially greater than the start hits.
		// Expected hits: (2 chars * 199 calls) = 398 minimum if fully sequential after first call.
		// Max hits would be slightly less depending on exact timing.
		// We just check it increased significantly.
		if hits <= startHits {
			t.Errorf("should have recorded cache hits from concurrent access, got %d, expected > %d", hits, startHits)
		}
	})

	t.Run("cache eviction", func(t *testing.T) {
		transliterate.ResetCacheStats()
		transliterate.ClearCache()

		const maxEntries = 1000

		// Create first set of characters
		var text1 string
		for i := rune(0x4E00); i < rune(0x4E00+maxEntries); i++ {
			text1 += string(i)
		}

		// Create second set of different characters to fill/evict cache
		var text2 string
		for i := rune(0x3400); i < rune(0x3400+maxEntries); i++ {
			text2 += string(i)
		}

		transliterate.String(text1)
		initialHits := transliterate.GetCacheStats()
		if initialHits != 0 {
			t.Errorf("expected 0 hits after first population, got %d", initialHits)
		}

		// Run with text2 to trigger eviction
		transliterate.String(text2)

		// Run text1 again - should miss cache if evicted
		transliterate.ResetCacheStats() // Reset to clearly see misses as 0 hits
		transliterate.String(text1)
		hitsAfterEviction := transliterate.GetCacheStats()

		// If cache was evicted, hits should be 0 again for text1
		if hitsAfterEviction != 0 {
			t.Errorf("expected 0 hits after potential eviction and re-run, got %d", hitsAfterEviction)
		}
	})
}

func TestCacheEdgeCases(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		transliterate.ResetCacheStats()
		transliterate.String("")
		hits := transliterate.GetCacheStats()
		if hits != 0 {
			t.Errorf("empty string should not affect cache, got %d hits", hits)
		}
	})

	t.Run("ascii only", func(t *testing.T) {
		transliterate.ResetCacheStats()
		transliterate.String("Hello123!@#")
		hits := transliterate.GetCacheStats()
		if hits != 0 {
			t.Errorf("ASCII-only string should not affect cache, got %d hits", hits)
		}
	})

	t.Run("invalid utf8", func(t *testing.T) {
		transliterate.ResetCacheStats()
		transliterate.String("Hello\xf0\x28World")
		hits := transliterate.GetCacheStats()
		if hits != 0 {
			t.Errorf("invalid UTF-8 should not affect cache, got %d hits", hits)
		}
	})
}

// Test cache size limits
func TestCacheSizeLimits(t *testing.T) {
	t.Run("at_capacity", func(t *testing.T) {
		transliterate.ResetCacheStats()
		transliterate.ClearCache()

		maxCacheSize := 1000
		// Make sure the cache size is set
		transliterate.Configure(transliterate.WithMaxCacheSize(maxCacheSize))

		// Fill cache to exactly maxCacheSize
		var text string
		for i := rune(0x4E00); i < rune(0x4E00+maxCacheSize); i++ {
			text += string(i)
		}

		transliterate.String(text)
		size := transliterate.GetCacheSize()
		if size != maxCacheSize {
			t.Errorf("cache size should be %d, got %d", maxCacheSize, size)
		}

		// Add one more character + run again
		transliterate.String(text + "世") // Should trigger hits for original `text` chars
		newSize := transliterate.GetCacheSize()

		expectedSize := maxCacheSize
		if newSize != expectedSize {
			t.Errorf("cache size should be %d after adding one more character, got %d", expectedSize, newSize)
		}
	})

	t.Run("repeated characters", func(t *testing.T) {
		transliterate.ResetCacheStats()
		transliterate.ClearCache()

		// Test with repeated characters
		// First pass - caches '世', hits on 2nd and 3rd '世'
		text := "世世世"
		transliterate.String(text) // First pass - should cache but no hits
		hits1 := transliterate.GetCacheStats()

		// Second pass - should hit on all three '世'
		transliterate.String(text)
		hits2 := transliterate.GetCacheStats()

		expectedHits1 := uint64(2) // Hits for 2nd and 3rd char
		if hits1 != expectedHits1 {
			t.Errorf("first pass should record %d hits, got %d", expectedHits1, hits1)
		}

		expectedHits2 := uint64(5) // 2 from first pass + 3 from second pass
		if hits2 != expectedHits2 {
			t.Errorf("second pass should record %d total hits, got %d", expectedHits2, hits2)
		}
	})

	t.Run("repeated character patterns", func(t *testing.T) {
		patterns := []struct {
			name     string
			text     string
			expected uint64
		}{
			{"single", "世", 0},         // One unique char: no hits
			{"double", "世世", 1},        // One unique char, second occurrence hits
			{"different", "世界", 0},     // Two unique chars: no hits
			{"mixed_repeat", "世界世", 1}, // Two unique chars, one repeat: one hit
			{"alternating", "世界世界", 2}, // Two unique chars, each repeated once: two hits
		}

		for _, p := range patterns {
			t.Run(p.name, func(t *testing.T) {
				transliterate.ResetCacheStats()
				transliterate.ClearCache()

				transliterate.String(p.text)
				hits := transliterate.GetCacheStats()

				if hits != p.expected {
					t.Errorf("unexpected number of hits for pattern '%s': expected %d, got %d", p.text, p.expected, hits)
				}
			})
		}
	})

	t.Run("repeated access patterns", func(t *testing.T) {
		transliterate.ResetCacheStats()
		transliterate.ClearCache()

		text := "世界"

		// First pass - caching, no hits
		transliterate.String(text)
		hits1 := transliterate.GetCacheStats()
		if hits1 != 0 {
			t.Errorf("first pass should have 0 hits, got %d", hits1)
		}

		// Second pass - hits for both chars
		transliterate.String(text)
		hits2 := transliterate.GetCacheStats()
		expectedHits2 := uint64(2) // 0 + 2
		if hits2 != expectedHits2 {
			t.Errorf("second pass should have %d total hits, got %d", expectedHits2, hits2)
		}

		// Third pass - more hits for both chars
		transliterate.String(text)
		hits3 := transliterate.GetCacheStats()
		expectedHits3 := uint64(4) // 2 + 2
		if hits3 != expectedHits3 {
			t.Errorf("third pass should have %d total hits, got %d", expectedHits3, hits3)
		}
	})
}

func TestCacheConcurrency(t *testing.T) {
	t.Run("high contention", func(t *testing.T) {
		transliterate.ResetCacheStats()
		transliterate.ClearCache()

		const goroutines = 100
		const iterations = 1000

		var wg sync.WaitGroup
		wg.Add(goroutines)

		// Mix of reading and writing to cache
		for i := range goroutines {
			go func(n int) {
				defer wg.Done()
				for range iterations {
					// Use different characters based on goroutine number
					char := rune(0x4E00 + n%10)
					transliterate.String(string(char))
				}
			}(i)
		}

		wg.Wait()
		hits := transliterate.GetCacheStats()

		minExpectedHits := uint64(1000) // A reasonable lower bound check
		if hits <= minExpectedHits {    // Use a reasonable lower bound instead of just > 0
			t.Errorf("should record significant hits under high concurrency, got %d, expected > %d", hits, minExpectedHits)
		}
	})
}

// Benchmark cache performance
func BenchmarkCachePerformance(b *testing.B) {
	texts := []struct {
		name string
		text string
	}{
		{"cached", "世界"},       // Should be cached after first run
		{"uncached", "未知文字"},   // New characters each time
		{"mixed", "Hello 世界!"}, // Mix of ASCII and non-ASCII
	}

	for _, tt := range texts {
		b.Run(tt.name, func(b *testing.B) {
			// Reset timer to exclude first run that populates cache
			transliterate.String(tt.text)
			b.ResetTimer()

			for b.Loop() {
				transliterate.String(tt.text)
			}
		})
	}
}
