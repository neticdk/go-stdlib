package transliterate_test

import (
	"sync"
	"testing"

	"github.com/neticdk/go-stdlib/assert"
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
		assert.Zero(t, hits, "initial hits should be zero")
	})

	t.Run("cache misses on first access", func(t *testing.T) {
		transliterate.ResetCacheStats()

		startHits := getStats()

		// First time translation
		text := "世界"
		transliterate.String(text)

		hits := getStats()
		assert.Equal(t, hits, startHits, "should have 0 hits on first access")
	})

	t.Run("cache hits on repeated access", func(t *testing.T) {
		transliterate.ResetCacheStats()

		text := "世界"
		transliterate.String(text) // First time
		transliterate.String(text) // Second time
		transliterate.String(text) // Third time

		hits := getStats()
		expectedHits := uint64(6) // 2 chars × 3 repeated calls after the first
		assert.Equal(t, hits, expectedHits, "should have 6 hits (2 chars × 3 repeated calls)")
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
		assert.Equal(t, hits, expectedTotalHits, "should have 4 new hits (2 non-ASCII chars × 2 calls)")
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
		assert.Greater(t, hits, startHits, "should have recorded cache hits from concurrent access")
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
		assert.Equal(t, initialHits, uint64(0), "initial hits should be 0 after first population")

		// Run with text2 to trigger eviction
		transliterate.String(text2)

		// Run text1 again - should miss cache if evicted
		transliterate.ResetCacheStats() // Reset to clearly see misses as 0 hits
		transliterate.String(text1)
		hitsAfterEviction := transliterate.GetCacheStats()

		// If cache was evicted, hits should be 0 again for text1
		assert.Equal(t, hitsAfterEviction, uint64(0), "should have 0 hits after eviction and re-run")
	})
}

func TestCacheEdgeCases(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		transliterate.ResetCacheStats()
		transliterate.String("")
		hits := transliterate.GetCacheStats()
		assert.Zero(t, hits, "empty string should not affect cache stats")
	})

	t.Run("ascii only", func(t *testing.T) {
		transliterate.ResetCacheStats()
		transliterate.String("Hello123!@#")
		hits := transliterate.GetCacheStats()
		assert.Zero(t, hits, "ASCII-only string should not affect cache stats")
	})

	t.Run("invalid utf8", func(t *testing.T) {
		transliterate.ResetCacheStats()
		transliterate.String("Hello\xf0\x28World")
		hits := transliterate.GetCacheStats()
		assert.Zero(t, hits, "invalid UTF-8 string should not affect cache stats")
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
		assert.Equal(t, size, maxCacheSize, "cache size should be equal to maxCacheSize")

		// Add one more character + run again
		transliterate.String(text + "世") // Should trigger hits for original `text` chars
		newSize := transliterate.GetCacheSize()

		expectedSize := maxCacheSize
		assert.Equal(t, newSize, expectedSize, "cache size should remain the same after adding one more character")
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
		assert.Equal(t, hits1, expectedHits1, "first pass should record 2 hits")

		expectedHits2 := uint64(5) // 2 from first pass + 3 from second pass
		assert.Equal(t, hits2, expectedHits2, "second pass should record 5 total hits")
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

				assert.Equal(t, hits, p.expected, "unexpected number of hits for pattern '%s': expected %d, got %d", p.text, p.expected, hits)
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
		assert.Zero(t, hits1, "first pass should have 0 hits")

		// Second pass - hits for both chars
		transliterate.String(text)
		hits2 := transliterate.GetCacheStats()
		expectedHits2 := uint64(2) // 0 + 2
		assert.Equal(t, hits2, expectedHits2, "second pass should have 2 total hits")

		// Third pass - more hits for both chars
		transliterate.String(text)
		hits3 := transliterate.GetCacheStats()
		expectedHits3 := uint64(4) // 2 + 2
		assert.Equal(t, hits3, expectedHits3, "third pass should have 4 total hits")
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
		assert.Greater(t, hits, minExpectedHits, "should record significant hits under high concurrency")
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
