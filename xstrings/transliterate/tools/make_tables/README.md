# Table Generator

Generates Go source files for Unicode transliteration tables from a text definition.

## Usage

From the package root directory:

```bash
go run tools/make_tables/main.go
```

## Input

Reads from `data/table.txt` with format:
```
0x0041: "A" # U+0041 (A)
0x00E9: "e" # U+00E9 (é)
```

## Output

Generates:
- `internal/table/x*.go` files containing transliteration tables
- `decode.go` with initialization code

## File Format

Each line in `table.txt` should be:
- Hex code point: `0xNNNN`
- Colon separator: `:`
- Quoted string: `"ascii"`
- Optional comment: `# description`

## Example
```bash
# Generate tables
go run tools/make_table/main.go

# Format generated code
go fmt ./...
```
