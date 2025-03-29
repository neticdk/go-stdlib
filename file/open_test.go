package file

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSafePath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "go-common-test-")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tmpDir)

	realTmpDir, err := filepath.EvalSymlinks(tmpDir)
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		name      string
		root      string
		path      string
		setup     func() error
		teardown  func() error
		want      string
		expectErr error
	}{
		{
			name:      "Valid path within base directory",
			root:      "/safe/base/directory",
			path:      "example.txt",
			want:      "/safe/base/directory/example.txt",
			expectErr: nil,
		},
		{
			name:      "Valid absolute path within base directory",
			root:      "/safe/base/directory",
			path:      "/safe/base/directory/example.txt",
			want:      "/safe/base/directory/example.txt",
			expectErr: nil,
		},
		{
			name:      "Path traversal attack",
			root:      "/safe/base/directory",
			path:      "../example.txt",
			want:      "",
			expectErr: ErrTraversal,
		},
		{
			name:      "Path traversal attack with absolute path",
			root:      "/safe/base/directory",
			path:      "/safe/base/directory/../../example.txt",
			want:      "",
			expectErr: ErrTraversal,
		},
		{
			name:      "Empty root",
			root:      "",
			path:      "example.txt",
			want:      "",
			expectErr: ErrEmptyPath,
		},
		{
			name:      "Empty path",
			root:      "/safe/base/directory",
			path:      "",
			want:      "",
			expectErr: ErrEmptyPath,
		},
		{
			name:      "Null byte in root",
			root:      "/safe/base/directory\x00",
			path:      "example.txt",
			want:      "",
			expectErr: ErrNullByte,
		},
		{
			name:      "Null byte in path",
			root:      "/safe/base/directory",
			path:      "example.txt\x00",
			want:      "",
			expectErr: ErrNullByte,
		},
		{
			name: "Symlink resolution within base directory",
			root: tmpDir,
			path: "symlink",
			setup: func() error {
				err := os.WriteFile(filepath.Join(tmpDir, "example.txt"), []byte("example"), 0o640)
				if err != nil {
					t.Fatalf("creating example.txt: %v", err)
				}
				return os.Symlink(filepath.Join(tmpDir, "example.txt"), filepath.Join(tmpDir, "symlink"))
			},
			teardown: func() error {
				os.Remove(filepath.Join(tmpDir, "symlink"))
				return nil
			},
			want:      filepath.Join(realTmpDir, "example.txt"),
			expectErr: nil,
		},
		{
			name: "Symlink resolution outside base directory",
			root: tmpDir,
			path: "symlink",
			setup: func() error {
				return os.Symlink("/etc/passwd", filepath.Join(tmpDir, "symlink"))
			},
			teardown: func() error {
				os.Remove(filepath.Join(tmpDir, "symlink"))
				return nil
			},
			want:      "",
			expectErr: ErrTraversal,
		},
		{
			name: "Symlink to sibling directory",
			root: tmpDir,
			path: "dir1/symlink/file.txt",
			setup: func() error {
				if err := os.Mkdir(filepath.Join(tmpDir, "dir1"), 0o755); err != nil {
					return err
				}
				if err := os.Mkdir(filepath.Join(tmpDir, "dir2"), 0o755); err != nil {
					return err
				}
				if err := os.WriteFile(filepath.Join(tmpDir, "dir2", "file.txt"), []byte("test"), 0o640); err != nil {
					return err
				}
				return os.Symlink("../dir2", filepath.Join(tmpDir, "dir1", "symlink"))
			},
			want:      filepath.Join(realTmpDir, "dir2", "file.txt"),
			expectErr: nil,
		},
		{
			name: "Complex traversal attempt with symlinks",
			root: tmpDir,
			path: "dir1/dir2/../../../../etc/passwd",
			setup: func() error {
				if err := os.MkdirAll(filepath.Join(tmpDir, "dir1", "dir2"), 0o755); err != nil {
					return err
				}
				return nil
			},
			want:      "",
			expectErr: ErrTraversal,
		},
		{
			name: "Symlink to parent directory",
			root: tmpDir,
			path: "dir/symlink/file.txt",
			setup: func() error {
				if err := os.Mkdir(filepath.Join(tmpDir, "dir"), 0o755); err != nil {
					return err
				}
				return os.Symlink("../..", filepath.Join(tmpDir, "dir", "symlink"))
			},
			want:      "",
			expectErr: ErrTraversal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				if err := tt.setup(); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}

			got, err := SafePath(tt.root, tt.path)

			if err != tt.expectErr {
				t.Errorf("Expected error %v, got %v", tt.expectErr, err)
			}

			if got != tt.want {
				t.Errorf("SafePath() = %v, want %v", got, tt.want)
			}

			if tt.teardown != nil {
				if err := tt.teardown(); err != nil {
					t.Fatalf("teardown failed: %v", err)
				}
			}
		})
	}
}
