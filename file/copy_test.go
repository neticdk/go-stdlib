package file

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyFile(t *testing.T) {
	testDir := "go-stdlib-test-copy"
	err := os.MkdirAll(testDir, os.ModePerm)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(testDir)

	tests := []struct {
		name       string
		srcFile    string
		srcContent []byte
		dstFile    string
	}{
		{
			name:       "Simple Copy file",
			srcFile:    "test_src_file.txt",
			srcContent: []byte("test1"),
			dstFile:    "test_dst_file.txt",
		},
		{
			name:       "Simple Copy file with same name",
			srcFile:    "test_src_file.txt",
			srcContent: []byte("test2"),
			dstFile:    "test_dst_file.txt",
		},
		{
			name:       "Copy file with subdirectory",
			srcFile:    "subdir/test_src_file.txt",
			srcContent: []byte("test3"),
			dstFile:    "subdir/test_dst_file.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srcFile := filepath.Join(testDir, tt.srcFile)
			dstFile := filepath.Join(testDir, tt.dstFile)

			if err = os.MkdirAll(filepath.Dir(srcFile), os.ModePerm); err != nil {
				t.Fatalf("creating source file: %v", err)
			}
			if err = os.MkdirAll(filepath.Dir(dstFile), os.ModePerm); err != nil {
				t.Fatalf("creating destination file: %v", err)
			}

			// Create source file
			if err = os.WriteFile(srcFile, tt.srcContent, 0o640); err != nil {
				t.Fatalf("creating source file: %v", err)
			}

			// Copy file
			if err = Copy(srcFile, dstFile); err != nil {
				t.Fatalf("copying file: %v", err)
			}

			// Verify content
			dstContent, err := os.ReadFile(dstFile)
			if err != nil {
				t.Fatalf("reading destination file: %v", err)
			}

			if string(dstContent) != string(tt.srcContent) {
				t.Fatalf("content mismatch: expected %s, got %s", string(tt.srcContent), string(dstContent))
			}

		})
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
