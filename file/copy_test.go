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
