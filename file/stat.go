package file

import (
	"errors"
	"io/fs"
	"os"
)

// Exists returns true if the given path exists
//
// It returns false and an error on any error, e.g. on insufficient permissions
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}

// IsDir returns true if the given path is a directory
//
// It returns false on any error, e.g. if the file does not exist or on insufficient permissions
func IsDir(path string) bool {
	return isFileMode(path, os.ModeDir)
}

// IsRegular returns true if the given path is a regular file
//
// It resolves all symbolic links
// It returns false on any error, e.g. if the file does not exist or on insufficient permissions
func IsRegular(path string) bool {
	if s, err := os.Stat(path); err == nil {
		return s.Mode().IsRegular()
	}
	return false
}

// IsSymlink returns true if the given path is a symlink
//
// It does not resolve any symbolic links
// It returns false on any error, e.g. if the file does not exist or on insufficient permissions
func IsSymlink(path string) bool {
	return isFileModeL(path, os.ModeSymlink)
}

// IsNamedPipe returns true if the given path is a named pipe
//
// It resolves all symbolic links
// It returns false on any error, e.g. if the file does not exist or on insufficient permissions
func IsNamedPipe(path string) bool {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return fileInfo.Mode()&os.ModeNamedPipe == os.ModeNamedPipe
}

// IsSocket returns true if the given path is a socket
//
// It resolves all symbolic links
// It returns false on any error, e.g. if the file does not exist or on insufficient permissions
func IsSocket(path string) bool {
	return isFileMode(path, os.ModeSocket)
}

// IsDevice returns true if the given path is a device
//
// It resolves all symbolic links
// It returns false on any error, e.g. if the file does not exist or on insufficient permissions
func IsDevice(path string) bool {
	return isFileMode(path, os.ModeDevice)
}

// IsFile returns true if the given path is a regular file, symlink, socket, or device
//
// It returns false on any error, e.g. if the file does not exist or on insufficient permissions
func IsFile(path string) bool {
	return IsRegular(path) ||
		IsSymlink(path) ||
		IsNamedPipe(path) ||
		IsSocket(path) ||
		IsDevice(path)
}

// isFileMode returns true if the given path has the given file mode, resolving any symbolic links
//
// It returns false on any error, e.g. if the file does not exist or on insufficient permissions
func isFileMode(path string, fileMode fs.FileMode) bool {
	if s, err := os.Stat(path); err == nil {
		return s.Mode()&fileMode != 0
	}
	return false
}

// isFileModeL returns true if the given path has the given file mode, not resolving any symbolic links
//
// It returns false on any error, e.g. if the file does not exist or on insufficient permissions
func isFileModeL(path string, fileMode fs.FileMode) bool {
	if s, err := os.Lstat(path); err == nil {
		return s.Mode()&fileMode != 0
	}
	return false
}
