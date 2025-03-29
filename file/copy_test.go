package file

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyFile(t *testing.T) {
	srcFile := "test_src_file.txt"
	dstFile := "test_dst_file.txt"
	content := []byte("Hello, World!")

	// Create source file
	if err := os.WriteFile(srcFile, content, 0o640); err != nil {
		t.Fatalf("creating source file: %v", err)
	}
	defer os.Remove(srcFile)
	defer os.Remove(dstFile)

	// Copy file
	if err := Copy(srcFile, dstFile); err != nil {
		t.Fatalf("copying file: %v", err)
	}

	// Verify content
	dstContent, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("reading destination file: %v", err)
	}
	if string(dstContent) != string(content) {
		t.Fatalf("content mismatch: expected %s, got %s", string(content), string(dstContent))
	}
}

func TestCopyDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "go-common-test-")
	assert.NoError(t, err)
	srcDir := filepath.Join(tmpDir, "src")
	dstDir := filepath.Join(tmpDir, "dst")
	fileContent := []byte("Hello, Directory!")
	defer os.RemoveAll(tmpDir)

	// Create source directory and file
	if err := os.Mkdir(srcDir, 0o750); err != nil {
		t.Fatalf("creating source directory: %v", err)
	}

	srcFile := filepath.Join(srcDir, "file.txt")
	if err := os.WriteFile(srcFile, fileContent, 0o640); err != nil {
		t.Fatalf("creating source file: %v", err)
	}

	// Copy directory
	if err := CopyDirectory(srcDir, dstDir); err != nil {
		t.Fatalf("copying directory: %v", err)
	}

	// Verify content
	dstFile := filepath.Join(dstDir, "file.txt")
	dstContent, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("reading destination file: %v", err)
	}
	if string(dstContent) != string(fileContent) {
		t.Fatalf("content mismatch: expected %s, got %s", string(fileContent), string(dstContent))
	}
}
