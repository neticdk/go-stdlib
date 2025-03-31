// Package transliterate provides functionality to convert Unicode text
// into plain ASCII equivalents.
//
// It takes Unicode characters and replaces them with their closest ASCII
// representation (e.g., 'é' becomes 'e', 'ü' becomes 'u'). Characters
// without a known ASCII approximation are generally omitted.
//
// Emoji and other pictographic symbols are omitted from the output
// as they have no standardized ASCII representation. If you need
// emoji-to-text conversion, consider using a dedicated emoji
// processing library.
//
// The transliteration tables are based on the tables used here
// https://github.com/mozillazg/go-unidecode but they have been modified
// to include additional characters and improve accuracy.
//
// Currently supports Unicode BMP (U+0000-U+FFFF) and some supplementary planes:
//   - x1d4: Mathematical Alphanumeric Symbols
//   - x1d5: Mathematical Alphanumeric Symbols
//   - x1d6: Mathematical Alphanumeric Symbols
//   - x1d7: Mathematical Alphanumeric Symbols
//   - x1f1: Enclosed Alphanumeric Supplement
//   - x1f6: Transport and Map Symbols + Emoji Symbols
//
// # Thread Safety
//
// String() and WithLimit() are safe for concurrent use. The package maintains
// an internal cache that is thread-safe and uses a buffer pool for improved
// performance under high concurrency.
//
// # Cache Behavior
//
// The package maintains a cache of up to 1000 character translations.
// When the cache becomes full, it is completely cleared. This approach
// favors simplicity over granular eviction but may affect performance
// for workloads with highly varying character sets.
//
// # Configuration
//
// The package can be configured using the Configure function with various options:
//
// Cache size (default 1000 entries):
//
//	transliterate.Configure(transliterate.WithMaxCacheSize(5000))
//
// Maximum input length (default 1MB):
//
//	transliterate.Configure(transliterate.WithMaxInputLength(1 << 24)) // 16MB
//
// Multiple options can be combined:
//
//	transliterate.Configure(
//	    transliterate.WithMaxCacheSize(5000),
//	    transliterate.WithMaxInputLength(1 << 24),
//	)
//
// Configuration should be done early in your application lifecycle,
// preferably before any calls to String() or WithLimit().
//
// # Table Generation
//
// The transliteration tables are generated from a text definition file.
// See tools/make_tables/README.md and tools/convert_tables/README.md
// for details on maintaining the tables.
//
// Examples:
//
//	ascii := transliterate.String("これはひらがなです") // Output: "korehahiraganadesu"
//	ascii := transliterate.String("你好，世界") // Output: "Ni Hao, Shi Jie" (Depends on table)
package transliterate

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o README.md
