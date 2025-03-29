package file

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const OSWindows = "windows"

var (
	ErrEmptyPath = errors.New("path cannot be empty")
	ErrNullByte  = errors.New("path contains null byte")
	ErrTraversal = errors.New("path traversal attempt detected")
	ErrSymlink   = errors.New("symlink traversal attempt detected")
)

// SafeOpenFile opens a file with the specified path, base directory, flags, and mode.
// It ensures the file operation is secure by validating the mode and path.
// Returns a file handle and any error encountered.
func SafeOpenFile(root, path string, flag int, mode int64) (*os.File, error) {
	fileMode, err := ValidMode(mode)
	if err != nil {
		return nil, err
	}

	path, err = SafePath(root, path)
	if err != nil {
		return nil, err
	}

	// Open the file with the specified mode
	outFile, err := os.OpenFile(path, flag, fileMode) // #nosec
	if err != nil {
		return nil, err
	}

	return outFile, nil
}

// SafeOpen opens a file for read-only access in a secure manner.
// It uses SafeOpenFile with read-only flag and default permissions.
func SafeOpen(root, path string) (*os.File, error) {
	return SafeOpenFile(root, path, os.O_RDONLY, 0)
}

// SafeCreate creates or truncates a file with the specified mode.
// It uses SafeOpenFile with write-only, create and truncate flags.
func SafeCreate(root, path string, mode int64) (*os.File, error) {
	return SafeOpenFile(root, path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
}

// SafeReadFile reads the entire contents of a file securely.
// It ensures the file path is safe before reading.
// Returns the file contents and any error encountered.
func SafeReadFile(root, path string) ([]byte, error) {
	path, err := SafePath(root, path)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(path) // #nosec
}

// ValidMode checks if the given mode is a valid file mode
func ValidMode(mode int64) (os.FileMode, error) {
	if mode < FileModeMinValid || mode > FileModeMaxValid {
		return 0, fmt.Errorf("invalid file mode: %o", mode)
	}
	return os.FileMode(mode), nil
}

// SafePath ensures the given path is safe to use within the specified root directory.
// It returns the cleaned absolute path and an error if the path is unsafe.
// The function performs the following checks and operations:
// 1. Validates that the path is not empty and does not contain null bytes.
// 2. Cleans and converts both the root and input paths to absolute paths.
// 3. Resolves any symlinks in the input path, even if the path does not exist.
// 4. Ensures the resolved path is within the root directory to prevent path traversal attacks.
// 5. Works on both Windows and Unix-like systems, handling platform-specific path separators and case sensitivity.
//
// Parameters:
// - root: The root directory within which the path must be contained.
// - path: The input path to be validated and resolved. It can be either an absolute or relative path.
func SafePath(root, path string) (string, error) {
	// Check for empty path
	if root == "" {
		return "", ErrEmptyPath
	}

	if path == "" {
		return "", ErrEmptyPath
	}

	// Check for null bytes
	if strings.Contains(root, "\x00") {
		return "", ErrNullByte
	}
	if strings.Contains(path, "\x00") {
		return "", ErrNullByte
	}

	// Clean and get absolute paths for both root and input path
	cleanRoot := filepath.Clean(root)
	absRoot, err := filepath.Abs(cleanRoot)
	if err != nil {
		return "", err
	}

	// Resolve any symlinks in the root path if it exists
	realRoot := absRoot
	if _, err := os.Stat(absRoot); err == nil {
		realRoot, err = filepath.EvalSymlinks(absRoot)
		if err != nil {
			return "", err
		}
	}

	// Clean and get absolute path for the input path
	cleanPath := filepath.Clean(path)
	var absPath string
	if filepath.IsAbs(cleanPath) {
		absPath = cleanPath
	} else {
		absPath = filepath.Join(realRoot, cleanPath)
	}

	// Resolve symlinks in the input path
	resolvedPath, err := resolveSymlinks(absPath)
	if err != nil {
		return "", err
	}

	// Check if the resolved path is within the root directory
	if !isWithinRoot(realRoot, resolvedPath) {
		return "", ErrTraversal
	}

	return resolvedPath, nil
}

// The `resolveSymlinks` function is designed to resolve any symbolic links
// (symlinks) in a given path. This is important for security purposes, as
// symlinks can potentially point to locations outside the intended directory,
// leading to path traversal attacks. The function ensures that the final
// resolved path is safe to use.
//
// It will resolve each component of the path, ensuring that the final path is a
// valid and safe location.
//
// If `file.txt` is a symlink pointing to another file, `resolveSymlinks` will
// resolve the symlink and return the actual path to the target file. If any
// part of the path contains a symlink, it will be resolved to its target
// location, so it works even though `file.txt` itself does not exist.
func resolveSymlinks(path string) (string, error) {
	var currentPath string
	if runtime.GOOS == OSWindows {
		currentPath = filepath.VolumeName(path) + string(filepath.Separator)
	} else if filepath.IsAbs(path) {
		currentPath = string(filepath.Separator)
	}

	components := strings.Split(filepath.ToSlash(path), "/")
	for _, component := range components {
		if component == "" || component == "." {
			continue
		}

		nextPath := filepath.Join(currentPath, component)
		if runtime.GOOS != OSWindows && filepath.IsAbs(path) && !strings.HasPrefix(nextPath, "/") {
			nextPath = "/" + nextPath
		}

		fileInfo, err := os.Lstat(nextPath)
		if err != nil {
			if os.IsNotExist(err) {
				currentPath = nextPath
				continue
			}
			return "", err
		}

		if fileInfo.Mode()&os.ModeSymlink != 0 {
			resolvedPath, err := filepath.EvalSymlinks(nextPath)
			if err != nil {
				return "", err
			}
			currentPath = resolvedPath
		} else {
			currentPath = nextPath
		}
	}

	if runtime.GOOS != OSWindows && filepath.IsAbs(path) && !strings.HasPrefix(currentPath, "/") {
		currentPath = "/" + currentPath
	}

	return currentPath, nil
}

func isWithinRoot(root, path string) bool {
	cleanRoot := filepath.Clean(root)
	cleanPath := filepath.Clean(path)

	if runtime.GOOS == OSWindows {
		cleanRoot = strings.ToLower(cleanRoot)
		cleanPath = strings.ToLower(cleanPath)
	}

	if cleanPath == cleanRoot {
		return true
	}

	rootWithSep := cleanRoot + string(filepath.Separator)
	return strings.HasPrefix(cleanPath, rootWithSep)
}
