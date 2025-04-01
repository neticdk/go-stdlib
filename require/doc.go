// Package require provides assertion functions that wrap the functions
// from the `is` package. Unlike `is`, functions in `require` call
// t.FailNow() upon failure, immediately stopping the current test.
//
// This is useful when a test cannot proceed meaningfully after a
// specific assertion fails.
package require

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o README.md
