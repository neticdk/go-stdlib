package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyDirectory copies a directory from src to dest
func CopyDirectory(srcDir, dest string) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("reading directory: %q, error: %q", srcDir, err.Error())
	}
	if !IsDir(dest) {
		if err := os.MkdirAll(dest, FileModeNewDirectory); err != nil {
			return fmt.Errorf("creating directory: %q, error: %q", dest, err.Error())
		}
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		fileInfo, err := entry.Info()
		if err != nil {
			return fmt.Errorf("getting file info for %q", sourcePath)
		}

		if err := copyFileOrDir(sourcePath, destPath, fileInfo); err != nil {
			return err
		}
	}
	return nil
}

// Copy copies a file from src to dest
func Copy(srcFile, dstFile string) error {
	in, err := SafeOpen(filepath.Dir(srcFile), srcFile)
	if err != nil {
		return fmt.Errorf("opening file: %q, error: %q", srcFile, err.Error())
	}
	defer func() {
		if cerr := in.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "closing file: %q, error: %q\n", srcFile, cerr.Error())
		}
	}()

	stat, err := in.Stat()
	var mode int64 = 0o640
	if err == nil {
		mode = int64(stat.Mode().Perm())
	}
	out, err := SafeCreate(filepath.Dir(dstFile), dstFile, mode)
	if err != nil {
		return fmt.Errorf("creating file: %q, error: %q", dstFile, err.Error())
	}
	defer func() {
		if cerr := out.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "closing file: %q, error: %q\n", dstFile, cerr.Error())
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return fmt.Errorf("copying file: %q, error: %q", srcFile, err.Error())
	}

	return nil
}

func copySymLink(source, dest string) error {
	link, err := os.Readlink(source)
	if err != nil {
		return fmt.Errorf("reading symlink: %q, error: %q", source, err.Error())
	}
	return os.Symlink(link, dest)
}
