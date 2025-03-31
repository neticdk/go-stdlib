package transliterate_test

import (
	"sync"
	"testing"

	"github.com/neticdk/go-stdlib/xstrings/transliterate"
	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, startHits, hits, "shouldn't have any cache hits")
	})

	t.Run("cache hits on repeated access", func(t *testing.T) {
		transliterate.ResetCacheStats()

		text := "世界"
		transliterate.String(text) // First time
		transliterate.String(text) // Second time
		transliterate.String(text) // Third time

		hits := getStats()
		assert.Equal(t, uint64(6), hits, "should have 6 hits (2 chars × 3 repeated calls)")
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

		assert.Equal(t, startHits+expectedNewHits, hits,
			"should only count cache hits for non-ASCII chars")
	})

	t.Run("concurrent access", func(t *testing.T) {
		transliterate.ResetCacheStats()

		startHits := getStats()

		// Run concurrent translations
		done := make(chan struct{})
		go func() {
			for range 100 {
				transliterate.String("世界")
			}
			done <- struct{}{}
		}()
		go func() {
			for range 100 {
				transliterate.String("世界")
			}
			done <- struct{}{}
		}()

		// Wait for both goroutines
		<-done
		<-done

		hits := getStats()
		assert.Greater(t, hits, startHits,
			"should have recorded cache hits from concurrent access")
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

		// Create second set of different characters
		var text2 string
		for i := rune(0x3400); i < rune(0x3400+maxEntries); i++ {
			text2 += string(i)
		}

		// First pass with text1
		transliterate.String(text1)
		hits1 := transliterate.GetCacheStats()

		// Run with text2 to trigger eviction
		transliterate.String(text2)

		// Run text1 again - should get new hits if cache was evicted
		transliterate.ResetCacheStats() // Reset to clearly see new hits
		transliterate.String(text1)
		hits2 := transliterate.GetCacheStats()

		// If cache was evicted, we should see fresh hits for text1
		assert.Equal(t, hits1, hits2,
			"should see same number of hits after eviction and reload")
	})
}

func TestCacheEdgeCases(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		transliterate.ResetCacheStats()
		transliterate.String("")
		hits := transliterate.GetCacheStats()
		assert.Zero(t, hits, "empty string should not affect cache")
	})

	t.Run("ascii only", func(t *testing.T) {
		transliterate.ResetCacheStats()
		transliterate.String("Hello123!@#")
		hits := transliterate.GetCacheStats()
		assert.Zero(t, hits, "ASCII-only string should not affect cache")
	})

	t.Run("invalid utf8", func(t *testing.T) {
		transliterate.ResetCacheStats()
		transliterate.String("Hello\xf0\x28World")
		hits := transliterate.GetCacheStats()
		assert.Zero(t, hits, "invalid UTF-8 should not affect cache")
	})
}

// Test cache size limits
func TestCacheSizeLimits(t *testing.T) {
	t.Run("at_capacity", func(t *testing.T) {
		transliterate.ResetCacheStats()
		transliterate.ClearCache()

		// Fill cache to exactly maxCacheSize
		var text string
		for i := rune(0x4E00); i < rune(0x4E00+1000); i++ {
			text += string(i)
		}

		transliterate.String(text)
		hits1 := transliterate.GetCacheStats()

		// Add one more character
		transliterate.String(text + "世")
		hits2 := transliterate.GetCacheStats()

		assert.Greater(t, hits2, hits1, "cache should handle boundary conditions")
	})

	t.Run("repeated characters", func(t *testing.T) {
		transliterate.ResetCacheStats()
		transliterate.ClearCache()

		// Test with repeated characters
		text := "世世世"
		transliterate.String(text) // First pass - should cache but no hits
		hits1 := transliterate.GetCacheStats()

		// Second pass - should get hits
		transliterate.String(text)
		hits2 := transliterate.GetCacheStats()

		assert.Equal(t, uint64(2), hits1,
			"first pass should record 2 hits")
		assert.Equal(t, uint64(5), hits2,
			"second pass should record 5 hits (one per character)")
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

				assert.Equal(t, p.expected, hits,
					"unexpected number of hits for pattern: %s", p.text)
			})
		}
	})

	t.Run("repeated access patterns", func(t *testing.T) {
		transliterate.ResetCacheStats()
		transliterate.ClearCache()

		text := "世界"

		// First pass - caching
		transliterate.String(text)
		hits1 := transliterate.GetCacheStats()

		// Second pass - hits
		transliterate.String(text)
		hits2 := transliterate.GetCacheStats()

		// Third pass - more hits
		transliterate.String(text)
		hits3 := transliterate.GetCacheStats()

		assert.Equal(t, uint64(0), hits1, "first pass should have no hits")
		assert.Equal(t, uint64(2), hits2, "second pass should have 2 hits")
		assert.Equal(t, uint64(4), hits3, "third pass should have 4 hits")
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
		assert.Greater(t, hits, uint64(0), "should record hits under high concurrency")
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
