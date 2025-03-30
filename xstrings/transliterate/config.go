package transliterate

// Default configuration values
const (
	DefaultCacheSize   = 1000
	DefaultMaxInputLen = 1 << 20 // 1MB
)

// CacheOption allows configuration of the transliteration cache
type CacheOption func(*config)

// config holds all configurable parameters for the transliteration package.
// It is not exported as configuration should be done through the Configure
// function.
type config struct {
	maxCacheSize int
	maxInputLen  int
}

// Default configuration values
var defaultConfig = config{
	maxCacheSize: DefaultCacheSize,
	maxInputLen:  DefaultMaxInputLen,
}

// WithMaxCacheSize sets the maximum size of the cache
func WithMaxCacheSize(size int) CacheOption {
	return func(c *config) {
		if size > 0 {
			c.maxCacheSize = size
		}
	}
}

// WithMaxInputLength sets the maximum input string length
func WithMaxInputLength(length int) CacheOption {
	return func(c *config) {
		if length > 0 {
			c.maxInputLen = length
		}
	}
}

// Configure applies the given options to the configuration
func Configure(opts ...CacheOption) {
	for _, opt := range opts {
		opt(&defaultConfig)
	}

	// Update cache size if it changed
	defaultCache.Lock()
	defaultCache.maxSize = defaultConfig.maxCacheSize
	defaultCache.Unlock()
}
