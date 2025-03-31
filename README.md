# Netic go standard library

[![CI](https://github.com/neticdk/go-stdlib/actions/workflows/ci.yaml/badge.svg)](https://github.com/neticdk/go-stdlib/actions/workflows/ci.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/neticdk/go-stdlib)](https://pkg.go.dev/github.com/neticdk/go-stdlib)
[![Go Report Card](https://goreportcard.com/badge/github.com/neticdk/go-stdlib)](https://goreportcard.com/report/github.com/neticdk/go-stdlib)
[![License](https://img.shields.io/github/license/neticdk/go-stdlib)](LICENSE)

The Netic go standard library is an extension to the go standard library. It
comes in the form of a collection of packages.

## Dependencies

The packages are dependency free meaning. Packages added to this module must not
use any external dependencies unless listed below.

Exceptions:

- `golang.org/x/*` - maintained by go and dependency free
- `github.com/stretchr/testify` - used for testing

CI checks the imports against regular expressions found in the
`.allowed-imports` file. To allow new imports, add them to the
`.allowed-imports` file in a separate PR.

Do *NOT* add exceptions to this list without peer review.

## Package names

- Prefix names for packages that mirror a go standard library package with `x`.
- Prefix names for packages that are likely to mirror future go standard library
  Packages with `x`.
- Use singular names for package (except in the previously mentioned cases).

## Testing

- Unit testing is mandatory.
- Go for > 95% coverage, preferably 100%.

## Documentation

- Document all exported (public) identifiers
- Maintain a `doc.go` and a `README.md` file in each package with introduction,
  installation instructions and usage examples.

## Packages

- `file` - file operations
- `unit` - unit formatting and conversion package
- `version` - version functions
- `xjson` - JSON functions
- `xset` - set data structure
- `xslices` - slice data type functions
- `xstrings` - string data type functions
- `xstructs` - struct data type functions

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
