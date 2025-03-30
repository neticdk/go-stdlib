//go:build ignore
// +build ignore

// convert_tables.go reads existing generated internal/table/xNNN.go files
// and converts the transliteration data back into a tables.txt format.
// It adds comments showing the original Unicode character if printable.
//
// Usage: go run convert_tables.go [-o <output_file>] [-d <input_directory>]
// Defaults: output_file=tables.txt.converted, input_directory=internal/table

package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

const (
	defaultInputDir   = "internal/table"
	defaultOutputFile = "data/table.txt.converted" // Avoid overwriting original by default
)

func main() {
	inputDir := flag.String("d", defaultInputDir, "Directory containing xNNN.go files")
	outputFile := flag.String("o", defaultOutputFile, "Output file name for the text format")
	flag.Parse()

	log.Printf("Reading from directory: %s", *inputDir)
	log.Printf("Writing output to: %s", *outputFile)

	// Data structure to hold results before sorting and writing
	outputData := make(map[rune]string)

	// Find relevant files
	files, err := filepath.Glob(filepath.Join(*inputDir, "x*.go"))
	if err != nil {
		log.Fatalf("Error finding files in %s: %v", *inputDir, err)
	}
	if len(files) == 0 {
		log.Fatalf("No x*.go files found in %s", *inputDir)
	}

	log.Printf("Found %d potential table files.", len(files))

	// Process each file
	for _, file := range files {
		filename := filepath.Base(file)
		// Extract section hex (e.g., "00a" from "x00a.go")
		sectionHex := strings.TrimSuffix(strings.TrimPrefix(filename, "x"), ".go")
		if len(sectionHex) != 3 { // Basic validation
			log.Printf("Skipping file with unexpected name format: %s", filename)
			continue
		}

		// Parse section number
		sectionVal64, err := strconv.ParseUint(sectionHex, 16, 32)
		if err != nil {
			log.Printf("Skipping file %s, cannot parse section number '%s': %v", filename, sectionHex, err)
			continue
		}
		section := rune(sectionVal64)

		// Parse the Go file content
		fset := token.NewFileSet() // Positions are relative to fset
		node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			log.Printf("Error parsing file %s: %v", file, err)
			continue // Skip file on parse error
		}

		// Extract the string slice data
		sliceData, err := extractFirstStringSlice(node)
		if err != nil {
			log.Printf("Error extracting slice data from %s: %v", file, err)
			continue // Skip file if data not found or format is wrong
		}

		// Allow slices with 0 to 256 elements. Slices larger than 256 are problematic.
		if len(sliceData) > 256 {
			log.Printf("Warning: Expected slice length <= 256 in %s, but found %d. Skipping file.", file, len(sliceData))
			continue
		}
		// Add informational log for short slices
		if len(sliceData) < 256 {
			log.Printf("Info: Slice in %s has only %d elements (expected 256). Processing available data.", file, len(sliceData))
		}

		// Populate outputData map
		for position, translit := range sliceData { // This iterates len(sliceData) times
			if translit != "" { // Only include non-empty transliterations
				currentRune := (section << 8) | rune(position)
				outputData[currentRune] = translit
			}
		}
	}

	// Write Output File
	if len(outputData) == 0 {
		log.Println("No transliteration data found to write.")
		return
	}

	// Sort runes for deterministic output
	runes := make([]rune, 0, len(outputData))
	for r := range outputData {
		runes = append(runes, r)
	}
	slices.Sort(runes)

	// Open and buffer output file
	outFile, err := os.Create(*outputFile)
	if err != nil {
		log.Fatalf("Error creating output file %s: %v", *outputFile, err)
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	// Write header comment
	_, _ = writer.WriteString("# Converted from Go source files using convert_tables.go\n")

	// Write sorted data
	for _, r := range runes {
		translit := outputData[r]
		// Create the comment string
		var comment string
		if unicode.IsPrint(r) {
			comment = fmt.Sprintf("# U+%04X (%c)", r, r) // Add character if printable
		} else {
			comment = fmt.Sprintf("# U+%04X", r) // Only add code point if not printable
		}

		// Format using %04X for rune, %q for quoted string literal, and add the comment
		line := fmt.Sprintf("0x%04X: %q %s\n", r, translit, comment)
		_, err := writer.WriteString(line)
		if err != nil {
			log.Fatalf("Error writing to output file %s: %v", *outputFile, err)
		}
	}

	log.Printf("Successfully wrote %d entries to %s", len(outputData), *outputFile)
}

// extractFirstStringSlice walks the AST of a file and returns the content
// of the first found `var ... = []string{...}` declaration.
func extractFirstStringSlice(node ast.Node) ([]string, error) {
	var extractedSlice []string
	var found bool

	ast.Inspect(node, func(n ast.Node) bool {
		if found { // Stop searching once found
			return false
		}

		// Look for Value Specifications (var declarations)
		spec, ok := n.(*ast.ValueSpec)
		if !ok {
			return true // Continue searching
		}

		// Expecting one variable name and one value assigned
		if len(spec.Names) != 1 || len(spec.Values) != 1 {
			return true
		}

		// Check if the assigned value is a Composite Literal (`[]string{...}`)
		compLit, ok := spec.Values[0].(*ast.CompositeLit)
		if !ok {
			return true
		}

		// Iterate through elements of the literal
		results := make([]string, 0, len(compLit.Elts))
		for _, elt := range compLit.Elts {
			// Expect each element to be a basic literal (string)
			basicLit, ok := elt.(*ast.BasicLit)
			if !ok || basicLit.Kind != token.STRING {
				// If not a string literal, the structure is unexpected
				return false // Stop searching this branch
			}

			// Unquote the string value (handles escapes)
			value, err := strconv.Unquote(basicLit.Value)
			if err != nil {
				log.Printf("Warning: Could not unquote string literal %q: %v", basicLit.Value, err)
				results = append(results, "") // Append empty or handle differently
			} else {
				results = append(results, value)
			}
		}

		// Found the first valid slice literal
		extractedSlice = results
		found = true
		return false // Stop the AST walk
	})

	if !found {
		return nil, fmt.Errorf("no `var ... = []string{...}` declaration found in the file")
	}
	return extractedSlice, nil
}
