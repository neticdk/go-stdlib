//go:build !windows

package file

import (
	"fmt"
	"os"
	"syscall"
)

func copyFileOrDir(sourcePath, destPath string, fileInfo os.FileInfo) error {
	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return fmt.Errorf("getting raw syscall.Stat_t data for %q", sourcePath)
	}

	switch fileInfo.Mode() & os.ModeType {
	case os.ModeDir:
		if err := os.MkdirAll(destPath, FileModeNewDirectory); err != nil {
			return fmt.Errorf("creating directory: %q, error: %q", destPath, err.Error())
		}
		if err := CopyDirectory(sourcePath, destPath); err != nil {
			return fmt.Errorf("copying directory: %q, error: %q", sourcePath, err.Error())
		}
	case os.ModeSymlink:
		if err := copySymLink(sourcePath, destPath); err != nil {
			return fmt.Errorf("copying symlink: %q, error: %q", sourcePath, err.Error())
		}
	default:
		if err := Copy(sourcePath, destPath); err != nil {
			return fmt.Errorf("copying file: %q, error: %q", sourcePath, err.Error())
		}
	}

	if err := os.Lchown(destPath, int(stat.Uid), int(stat.Gid)); err != nil {
		return fmt.Errorf("changing ownership of %q, error: %q", destPath, err.Error())
	}

	isSymlink := fileInfo.Mode()&os.ModeSymlink != 0
	if !isSymlink {
		if err := os.Chmod(destPath, fileInfo.Mode()); err != nil {
			return fmt.Errorf("changing mode of %q, error: %q", destPath, err.Error())
		}
	}

	return nil
}
