package file

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/neticdk/go-stdlib/assert"
	"github.com/neticdk/go-stdlib/require"
)

func TestCopyFile(t *testing.T) {
	srcFile := "test_src_file.txt"
	dstFile := "test_dst_file.txt"
	content := []byte("Hello, World!")

	// Create source file
	err := os.WriteFile(srcFile, content, 0o640)
	require.NoError(t, err, "creating source file")
	defer os.Remove(srcFile)
	defer os.Remove(dstFile)

	// Copy file
	err = Copy(srcFile, dstFile)
	require.NoError(t, err)

	// Verify content
	dstContent, err := os.ReadFile(dstFile)
	require.NoError(t, err, "reading destination file")
	assert.Equal(t, string(content), string(dstContent), "content mismatch")

	testDir := "go-stdlib-test-copy"
	err = os.MkdirAll(testDir, os.ModePerm)
  require.NoError(t, err, "creating test directory")
	defer func() { _ = os.RemoveAll(testDir) }()

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

			if err := os.MkdirAll(filepath.Dir(srcFile), os.ModePerm); err != nil {
				t.Fatalf("creating source file: %v", err)
			}
			if err := os.MkdirAll(filepath.Dir(dstFile), os.ModePerm); err != nil {
				t.Fatalf("creating destination file: %v", err)
			}

			// Create source file
			if err := os.WriteFile(srcFile, tt.srcContent, 0o640); err != nil {
				t.Fatalf("creating source file: %v", err)
			}

			// Copy file
			if err := Copy(srcFile, dstFile); err != nil {
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
	require.NoError(t, err)
	srcDir := filepath.Join(tmpDir, "src")
	dstDir := filepath.Join(tmpDir, "dst")
	fileContent := []byte("Hello, Directory!")
	defer os.RemoveAll(tmpDir)

	// Create source directory and file
	err = os.Mkdir(srcDir, 0o750)
	require.NoError(t, err, "creating source directory")

	srcFile := filepath.Join(srcDir, "file.txt")
	err = os.WriteFile(srcFile, fileContent, 0o640)
	require.NoError(t, err, "creating source file")

	// Copy directory
	err = CopyDirectory(srcDir, dstDir)
	require.NoError(t, err)

	// Verify content
	dstFile := filepath.Join(dstDir, "file.txt")
	dstContent, err := os.ReadFile(dstFile)
	require.NoError(t, err, "reading destination file")
	assert.Equal(t, string(dstContent), string(fileContent), "content mismatch")
}
