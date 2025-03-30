# Table Converter

Converts generated Go table files back to text format for inspection or modification.

## Usage

From the package root directory:

```bash
go run tools/convert_tables/main.go [-o output] [-d input_dir]
```

Options:
- `-o`: Output file (default: data/table.txt.converted)
- `-d`: Input directory (default: internal/table)

## Example

```bash
# Convert current tables
go run tools/convert_tables/main.go

# Convert to specific file
go run tools/convert_tables/main.go -o new_tables.txt
```

## Output Format

Generates a text file with format:
```
0x0041: "A" # U+0041 (A)
0x00E9: "e" # U+00E9 (é)
```
