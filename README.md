# Netic go standard library

[![GitHub Tag](https://img.shields.io/github/v/tag/neticdk/go-stdlib)](https://github.com/neticdk/go-stdlib/releases)
[![CI](https://github.com/neticdk/go-stdlib/actions/workflows/ci.yaml/badge.svg)](https://github.com/neticdk/go-stdlib/actions/workflows/ci.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/neticdk/go-stdlib)](https://pkg.go.dev/github.com/neticdk/go-stdlib)
[![Go Report Card](https://goreportcard.com/badge/github.com/neticdk/go-stdlib)](https://goreportcard.com/report/github.com/neticdk/go-stdlib)
[![License](https://img.shields.io/github/license/neticdk/go-stdlib)](LICENSE)

The Netic go standard library is an extension to the go standard library. It
comes in the form of a collection of packages.

## Dependencies

The packages are dependency free, meaning they must not use any external
dependencies unless explicitly listed.
Exceptions:

- `golang.org/x/*` - maintained by go and dependency free

CI checks the imports against regular expressions found in the
`.allowed-imports` file. To allow new imports, add them to the
`.allowed-imports` file in a separate PR.

Do *NOT* add exceptions to this list without peer review.

## Package names

- Prefix names for packages that mirror a go standard library package with `x`.
- Prefix names for packages that are likely to mirror future go standard library
  packages with `x`.
- Use singular names for package (except in the mentioned cases).

## Testing

- Unit testing is mandatory.
- Go for > 90% coverage, preferably 100%.

## Documentation

- Document all exported (public) identifiers
- Maintain a `doc.go` in each package with introduction, installation
  instructions and usage examples.

### doc.go minimal content

```go
// Package mypkg does ...
package mypkg
```

## Packages

- `assert` / `require` - test helpers for assertion
- `diff` / `diff/myers` / `diff/simple` - generate diffs
- `file` - file operations
- `set` - set data structure
- `unit` - unit formatting and conversion package
- `xjson` - JSON functions
- `xslices` - slice data type functions
- `xstrings` - string data type functions
- `xstructs` - struct data type functions
- `xtime` - time functions

## Installation

Install using `go get`:

```bash
go get github.com/neticdk/go-stdlib
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

Link to license or copyright notice

Copyright 2025 Netic A/S. All rights reserved.
