//go:build ignore
// +build ignore

// make_table.go generates the transliteration table Go source files
// from a text definition file (tables.txt). It creates the necessary
// internal/table/xNNN.go files and the main decode.go initialization file.
//
// Run it with: go run make_table.go

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

const (
	inputFile          = "data/table.txt"             // Input definition file
	outputTableDir     = "internal/table"             // Directory for xNNN.go files
	outputDecoderFile  = "decode.go"                  // Output file for init function
	maxSection         = 0x1F9                        // Generate blocks up to U+FFFF
	packagePathPrefix  = "github.com/neticdk/go-stdlib/" // Adjust to your actual module path root
	internalImportPath = packagePathPrefix + "xstrings/transliterate/internal/table"
)

// Holds the parsed transliteration data
var translitData = make(map[rune]string)

// Keeps track of which sections (blocks) actually contain data
var activeSections = make(map[rune]struct{})

func main() {
	log.Println("Starting table generation...")

	// Parse Input File
	parseInputFile(inputFile)
	log.Printf("Parsed %d transliteration entries.", len(translitData))

	// Prepare Output Directory
	if err := os.MkdirAll(outputTableDir, 0755); err != nil {
		log.Fatalf("Error creating directory %s: %v", outputTableDir, err)
	}
	log.Printf("Ensured output directory exists: %s", outputTableDir)

	// Generate Table Files (xNNN.go)
	generatedSections := generateTableFiles()
	log.Printf("Generated %d table block files (xNNN.go).", len(generatedSections))

	// Generate Decoder Initialization File (decode.go)
	generateDecoderFile(generatedSections)
	log.Printf("Generated decoder file: %s", outputDecoderFile)

	log.Println("Table generation complete.")
	log.Println("Recommendation: Run 'go fmt ./...' on your package.")
}

// findQuotedStringPrefix finds the first valid Go string literal (single or double quoted)
// at the beginning of the input string. It returns the literal including the quotes.
func findQuotedStringPrefix(s string) (string, error) {
	if len(s) == 0 {
		return "", fmt.Errorf("value part is empty")
	}
	quoteChar := s[0]
	if quoteChar != '"' && quoteChar != '\'' {
		return "", fmt.Errorf("value must start with quotes (' or \")")
	}
	if len(s) < 2 {
		return "", fmt.Errorf("value too short for quotes")
	}

	escaped := false
	for i := 1; i < len(s); i++ {
		char := s[i]
		if escaped {
			escaped = false // Consume escape, next char is literal
			continue
		}
		if char == '\\' {
			escaped = true
			continue
		}
		if char == quoteChar {
			// Found the matching, unescaped closing quote
			return s[:i+1], nil // Return the substring including quotes
		}
	}
	// If loop finishes without finding closing quote
	return "", fmt.Errorf("unclosed quote %c", quoteChar)
}

// parseInputFile reads the definition file and populates translitData.
func parseInputFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening input file %s: %v", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and full-line comments
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}

		// Split by colon
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			log.Printf("Warning: Skipping malformed line %d (no colon): %q", lineNumber, line)
			continue
		}

		runeHexStr := strings.TrimSpace(parts[0])
		valuePart := strings.TrimSpace(parts[1]) // Trim space around the value/comment part

		// Find the quoted string literal at the beginning of valuePart
		translitStrLiteral, err := findQuotedStringPrefix(valuePart) // CALL happens here
		if err != nil {
			log.Printf("Warning: Skipping line %d, could not find valid string literal: %v. Line: %q", lineNumber, err, line)
			continue
		}

		// Parse rune hex value
		runeVal64, err := strconv.ParseUint(strings.TrimPrefix(runeHexStr, "0x"), 16, 32) // rune is int32
		if err != nil {
			log.Printf("Warning: Skipping line %d due to invalid rune hex %q: %v", lineNumber, runeHexStr, err)
			continue
		}
		runeVal := rune(runeVal64)

		// Unquote the extracted string literal (should always work if findQuotedStringPrefix succeeded)
		translitVal, err := strconv.Unquote(translitStrLiteral)
		if err != nil {
			// This should be rare if findQuotedStringPrefix worked, but handle defensively
			log.Printf("Error(Internal): Skipping line %d, failed to unquote presumably valid literal %q: %v", lineNumber, translitStrLiteral, err)
			continue
		}

		// Store data
		if _, exists := translitData[runeVal]; exists {
			log.Printf("Warning: Duplicate entry for rune %U (0x%X) on line %d. Overwriting.", runeVal, runeVal, lineNumber)
		}
		translitData[runeVal] = translitVal

		// Mark section as active
		section := runeVal >> 8
		activeSections[section] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input file %s: %v", filename, err)
	}
}

// generateTableFiles creates the internal/table/xNNN.go files.
func generateTableFiles() []rune {
	var generatedSections []rune

	// Iterate through potential sections (blocks) up to maxSection
	for section := rune(0); section <= maxSection; section++ {
		// Only generate if this section has data defined in the input
		if _, ok := activeSections[section]; !ok {
			continue
		}

		filename := filepath.Join(outputTableDir, fmt.Sprintf("x%03x.go", section))
		var buf bytes.Buffer // Use buffer to write atomically later

		// Write file header
		buf.WriteString("// Code generated by make_table.go; DO NOT EDIT.\n\n")
		buf.WriteString("package table\n\n")
		buf.WriteString(fmt.Sprintf("// X%03X contains transliterations for Unicode code points U+%04X to U+%04X.\n",
			section, section<<8, (section<<8)|0xFF))
		buf.WriteString(fmt.Sprintf("var X%03X = []string{\n", section))

		// Write 256 entries for the section
		for position := range 256 {
			currentRune := (section << 8) | rune(position)
			translit, found := translitData[currentRune]

			var comment string
			if unicode.IsPrint(currentRune) {
				comment = fmt.Sprintf("// U+%04X (%c)", currentRune, currentRune) // Add character if printable
			} else {
				comment = fmt.Sprintf("// U+%04X", currentRune) // Only add code point if not printable
			}

			if found {
				// Use %q for safe Go string literal quoting and add the comment
				buf.WriteString(fmt.Sprintf("\t%q, %s\n", translit, comment))
			} else {
				// Use "" for missing entries and add the comment
				buf.WriteString(fmt.Sprintf("\t\"\", %s\n", comment))
			}
		}

		buf.WriteString("}\n")

		// Write buffer to file
		err := os.WriteFile(filename, buf.Bytes(), 0644)
		if err != nil {
			log.Fatalf("Error writing file %s: %v", filename, err)
		}
		generatedSections = append(generatedSections, section)
	}
	return generatedSections
}

// generateDecoderFile creates the decode.go file with the initialization logic.
func generateDecoderFile(sections []rune) {
	// Sort sections for deterministic output
	slices.Sort(sections)

	var buf bytes.Buffer

	// Write file header
	buf.WriteString("// Code generated by make_table.go; DO NOT EDIT.\n\n")
	buf.WriteString("package transliterate\n\n")
	buf.WriteString(fmt.Sprintf("import %q\n\n", internalImportPath)) // Use quoted import path

	buf.WriteString("// decodeTransliterations populates the global Tables map in the internal table package\n")
	buf.WriteString("// with the generated transliteration data slices (table.X000, table.X001, etc.).\n")
	buf.WriteString("// This function is intended to be called only once via decodingOnce.Do().\n")
	buf.WriteString("func decodeTransliterations() {\n")

	// Write assignments for each generated section
	for _, section := range sections {
		buf.WriteString(fmt.Sprintf("\ttable.Tables[0x%03X] = table.X%03X\n", section, section))
	}

	buf.WriteString("}\n")

	// Write buffer to file
	err := os.WriteFile(outputDecoderFile, buf.Bytes(), 0644)
	if err != nil {
		log.Fatalf("Error writing file %s: %v", outputDecoderFile, err)
	}
}
