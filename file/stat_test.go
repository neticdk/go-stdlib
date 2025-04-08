package file

import (
	"net"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/neticdk/go-stdlib/assert"
	"github.com/neticdk/go-stdlib/require"
)

func TestExists(t *testing.T) {
	// Create temporary test files/directories
	tmpDir := t.TempDir()
	regularFile := filepath.Join(tmpDir, "regular.txt")
	err := os.WriteFile(regularFile, []byte("test"), 0o644)
	require.NoError(t, err)
	subDir := filepath.Join(tmpDir, "subdir")
	err = os.Mkdir(subDir, 0o755)
	require.NoError(t, err)

	tests := []struct {
		name     string
		path     string
		want     bool
		wantErr  bool
		setup    func() error
		teardown func() error
	}{
		{
			name: "existing regular file",
			path: regularFile,
			want: true,
		},
		{
			name: "existing directory",
			path: subDir,
			want: true,
		},
		{
			name: "non-existent file",
			path: filepath.Join(tmpDir, "nonexistent.txt"),
			want: false,
		},
		{
			name:    "permission denied",
			path:    filepath.Join(tmpDir, "noperm", "test.txt"),
			want:    false,
			wantErr: true,
			setup: func() error {
				if err := os.Mkdir(filepath.Join(tmpDir, "noperm"), 0o755); err != nil {
					return err
				}
				if err := os.WriteFile(filepath.Join(tmpDir, "noperm", "test.txt"), []byte("test"), 0o000); err != nil {
					return err
				}
				if err := os.Chmod(filepath.Join(tmpDir, "noperm"), 0o000); err != nil {
					return err
				}
				return nil
			},
			teardown: func() error {
				if err := os.Chmod(filepath.Join(tmpDir, "noperm"), 0o755); err != nil {
					return err
				}
				return os.RemoveAll(filepath.Join(tmpDir, "noperm"))
			},
		},
		{
			name: "socket",
			path: filepath.Join(tmpDir, "socket"),
			want: true,
			setup: func() error {
				_, err := net.Listen("unix", filepath.Join(tmpDir, "socket"))
				return err
			},
		},
		{
			name: "named pipe",
			path: filepath.Join(tmpDir, "pipe"),
			want: true,
			setup: func() error {
				return syscall.Mkfifo(filepath.Join(tmpDir, "pipe"), 0o644)
			},
		},
		{
			name: "character device",
			path: filepath.Join(tmpDir, "chardev"),
			want: true,
			setup: func() error {
				return syscall.Mknod(filepath.Join(tmpDir, "chardev"), syscall.S_IFCHR|0o644, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				if err := tt.setup(); err != nil {
					t.Skipf("Failed to setup test: %s", err)
				}
			}

			got, err := Exists(tt.path)
			if !tt.wantErr {
				assert.NoError(t, err, "Exists()/%q", tt.name)
			}
			assert.Equal(t, got, tt.want, "Exists()/%q", tt.name)

			if tt.teardown != nil {
				err := tt.teardown()
				require.NoError(t, err, "teardown/%q", tt.name)
			}
		})
	}
}

func TestIsDir(t *testing.T) {
	// Create temporary test files/directories
	tmpDir := t.TempDir()
	regularFile := filepath.Join(tmpDir, "regular.txt")
	err := os.WriteFile(regularFile, []byte("test"), 0o644)
	require.NoError(t, err)
	subDir := filepath.Join(tmpDir, "subdir")
	err = os.Mkdir(subDir, 0o755)
	require.NoError(t, err)

	tests := []struct {
		name  string
		path  string
		want  bool
		setup func() error
	}{
		{
			name: "existing directory",
			path: subDir,
			want: true,
		},
		{
			name: "regular file",
			path: regularFile,
			want: false,
		},
		{
			name: "non-existent path",
			path: filepath.Join(tmpDir, "nonexistent"),
			want: false,
		},
		{
			name: "symlink to directory",
			path: filepath.Join(tmpDir, "symlink"),
			want: true,
			setup: func() error {
				return os.Symlink(subDir, filepath.Join(tmpDir, "symlink"))
			},
		},
		{
			name: "socket",
			path: filepath.Join(tmpDir, "socket"),
			want: false,
			setup: func() error {
				_, err := net.Listen("unix", filepath.Join(tmpDir, "socket"))
				return err
			},
		},
		{
			name: "named pipe",
			path: filepath.Join(tmpDir, "pipe"),
			want: false,
			setup: func() error {
				return syscall.Mkfifo(filepath.Join(tmpDir, "pipe"), 0o644)
			},
		},
		{
			name: "character device",
			path: filepath.Join(tmpDir, "chardev"),
			want: false,
			setup: func() error {
				return syscall.Mknod(filepath.Join(tmpDir, "chardev"), syscall.S_IFCHR|0o644, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				if err := tt.setup(); err != nil {
					t.Skipf("File type %s creation not supported on this platform", tt.name)
				}
			}

			got := IsDir(tt.path)
			assert.Equal(t, got, tt.want, "IsDir()/%q", tt.name)
		})
	}
}

func TestIsRegular(t *testing.T) {
	// Create temporary test files/directories
	tmpDir := t.TempDir()
	regularFile := filepath.Join(tmpDir, "regular.txt")
	err := os.WriteFile(regularFile, []byte("test"), 0o644)
	require.NoError(t, err)
	subDir := filepath.Join(tmpDir, "subdir")
	err = os.Mkdir(subDir, 0o755)
	require.NoError(t, err)

	tests := []struct {
		name  string
		path  string
		want  bool
		setup func() error
	}{
		{
			name: "existing directory",
			path: subDir,
			want: false,
		},
		{
			name: "regular file",
			path: regularFile,
			want: true,
		},
		{
			name: "non-existent path",
			path: filepath.Join(tmpDir, "nonexistent"),
			want: false,
		},
		{
			name: "symlink to file",
			path: filepath.Join(tmpDir, "symlink"),
			want: true,
			setup: func() error {
				return os.Symlink(regularFile, filepath.Join(tmpDir, "symlink"))
			},
		},
		{
			name: "socket",
			path: filepath.Join(tmpDir, "socket"),
			want: false,
			setup: func() error {
				_, err := net.Listen("unix", filepath.Join(tmpDir, "socket"))
				return err
			},
		},
		{
			name: "named pipe",
			path: filepath.Join(tmpDir, "pipe"),
			want: false,
			setup: func() error {
				return syscall.Mkfifo(filepath.Join(tmpDir, "pipe"), 0o644)
			},
		},
		{
			name: "character device",
			path: filepath.Join(tmpDir, "chardev"),
			want: false,
			setup: func() error {
				return syscall.Mknod(filepath.Join(tmpDir, "chardev"), syscall.S_IFCHR|0o644, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				if err := tt.setup(); err != nil {
					t.Skipf("File type %s creation not supported on this platform", tt.name)
				}
			}

			got := IsRegular(tt.path)
			assert.Equal(t, got, tt.want, "IsRegular()/%q", tt.name)
		})
	}
}

func TestIsSymlink(t *testing.T) {
	// Create temporary test files/directories
	tmpDir := t.TempDir()
	regularFile := filepath.Join(tmpDir, "regular.txt")
	err := os.WriteFile(regularFile, []byte("test"), 0o644)
	require.NoError(t, err)
	subDir := filepath.Join(tmpDir, "subdir")
	err = os.Mkdir(subDir, 0o755)
	require.NoError(t, err)

	tests := []struct {
		name  string
		path  string
		want  bool
		setup func() error
	}{
		{
			name: "existing directory",
			path: subDir,
			want: false,
		},
		{
			name: "regular file",
			path: regularFile,
			want: false,
		},
		{
			name: "non-existent path",
			path: filepath.Join(tmpDir, "nonexistent"),
			want: false,
		},
		{
			name: "symlink to file",
			path: filepath.Join(tmpDir, "symlink"),
			want: true,
			setup: func() error {
				return os.Symlink(regularFile, filepath.Join(tmpDir, "symlink"))
			},
		},
		{
			name: "socket",
			path: filepath.Join(tmpDir, "socket"),
			want: false,
			setup: func() error {
				_, err := net.Listen("unix", filepath.Join(tmpDir, "socket"))
				return err
			},
		},
		{
			name: "named pipe",
			path: filepath.Join(tmpDir, "pipe"),
			want: false,
			setup: func() error {
				return syscall.Mkfifo(filepath.Join(tmpDir, "pipe"), 0o644)
			},
		},
		{
			name: "character device",
			path: filepath.Join(tmpDir, "chardev"),
			want: false,
			setup: func() error {
				return syscall.Mknod(filepath.Join(tmpDir, "chardev"), syscall.S_IFCHR|0o644, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				if err := tt.setup(); err != nil {
					t.Skipf("File type %s creation not supported on this platform", tt.name)
				}
			}

			got := IsSymlink(tt.path)
			assert.Equal(t, got, tt.want, "IsSymlink()/%q", tt.name)
		})
	}
}

func TestIsSocket(t *testing.T) {
	// Create temporary test files/directories
	tmpDir := t.TempDir()
	regularFile := filepath.Join(tmpDir, "regular.txt")
	err := os.WriteFile(regularFile, []byte("test"), 0o644)
	require.NoError(t, err)
	subDir := filepath.Join(tmpDir, "subdir")
	err = os.Mkdir(subDir, 0o755)
	require.NoError(t, err)

	tests := []struct {
		name  string
		path  string
		want  bool
		setup func() error
	}{
		{
			name: "existing directory",
			path: subDir,
			want: false,
		},
		{
			name: "regular file",
			path: regularFile,
			want: false,
		},
		{
			name: "non-existent path",
			path: filepath.Join(tmpDir, "nonexistent"),
			want: false,
		},
		{
			name: "symlink to file",
			path: filepath.Join(tmpDir, "symlink"),
			want: false,
			setup: func() error {
				return os.Symlink(regularFile, filepath.Join(tmpDir, "symlink"))
			},
		},
		{
			name: "socket",
			path: filepath.Join(tmpDir, "socket"),
			want: true,
			setup: func() error {
				_, err := net.Listen("unix", filepath.Join(tmpDir, "socket"))
				return err
			},
		},
		{
			name: "named pipe",
			path: filepath.Join(tmpDir, "pipe"),
			want: false,
			setup: func() error {
				return syscall.Mkfifo(filepath.Join(tmpDir, "pipe"), 0o644)
			},
		},
		{
			name: "character device",
			path: filepath.Join(tmpDir, "chardev"),
			want: false,
			setup: func() error {
				return syscall.Mknod(filepath.Join(tmpDir, "chardev"), syscall.S_IFCHR|0o644, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				if err := tt.setup(); err != nil {
					t.Skipf("File type %s creation not supported on this platform", tt.name)
				}
			}

			got := IsSocket(tt.path)
			assert.Equal(t, got, tt.want, "IsSocket()/%q", tt.name)
		})
	}
}

func TestIsNamedPipe(t *testing.T) {
	// Create temporary test files/directories
	tmpDir := t.TempDir()
	regularFile := filepath.Join(tmpDir, "regular.txt")
	err := os.WriteFile(regularFile, []byte("test"), 0o644)
	require.NoError(t, err)
	subDir := filepath.Join(tmpDir, "subdir")
	err = os.Mkdir(subDir, 0o755)
	require.NoError(t, err)

	tests := []struct {
		name  string
		path  string
		want  bool
		setup func() error
	}{
		{
			name: "existing directory",
			path: subDir,
			want: false,
		},
		{
			name: "regular file",
			path: regularFile,
			want: false,
		},
		{
			name: "non-existent path",
			path: filepath.Join(tmpDir, "nonexistent"),
			want: false,
		},
		{
			name: "symlink to file",
			path: filepath.Join(tmpDir, "symlink"),
			want: false,
			setup: func() error {
				return os.Symlink(regularFile, filepath.Join(tmpDir, "symlink"))
			},
		},
		{
			name: "socket",
			path: filepath.Join(tmpDir, "socket"),
			want: false,
			setup: func() error {
				_, err := net.Listen("unix", filepath.Join(tmpDir, "socket"))
				return err
			},
		},
		{
			name: "named pipe",
			path: filepath.Join(tmpDir, "pipe"),
			want: true,
			setup: func() error {
				return syscall.Mkfifo(filepath.Join(tmpDir, "pipe"), 0o644)
			},
		},
		{
			name: "character device",
			path: filepath.Join(tmpDir, "chardev"),
			want: false,
			setup: func() error {
				return syscall.Mknod(filepath.Join(tmpDir, "chardev"), syscall.S_IFCHR|0o644, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				if err := tt.setup(); err != nil {
					t.Skipf("File type %s creation not supported on this platform", tt.name)
				}
			}

			got := IsNamedPipe(tt.path)
			assert.Equal(t, got, tt.want, "IsNamedPipe()/%q", tt.name)
		})
	}
}

func TestIsDevice(t *testing.T) {
	// Create temporary test files/directories
	tmpDir := t.TempDir()
	regularFile := filepath.Join(tmpDir, "regular.txt")
	err := os.WriteFile(regularFile, []byte("test"), 0o644)
	require.NoError(t, err)
	subDir := filepath.Join(tmpDir, "subdir")
	err = os.Mkdir(subDir, 0o755)
	require.NoError(t, err)

	tests := []struct {
		name  string
		path  string
		want  bool
		setup func() error
	}{
		{
			name: "existing directory",
			path: subDir,
			want: false,
		},
		{
			name: "regular file",
			path: regularFile,
			want: false,
		},
		{
			name: "non-existent path",
			path: filepath.Join(tmpDir, "nonexistent"),
			want: false,
		},
		{
			name: "symlink to file",
			path: filepath.Join(tmpDir, "symlink"),
			want: false,
			setup: func() error {
				return os.Symlink(regularFile, filepath.Join(tmpDir, "symlink"))
			},
		},
		{
			name: "socket",
			path: filepath.Join(tmpDir, "socket"),
			want: false,
			setup: func() error {
				_, err := net.Listen("unix", filepath.Join(tmpDir, "socket"))
				return err
			},
		},
		{
			name: "named pipe",
			path: filepath.Join(tmpDir, "pipe"),
			want: false,
			setup: func() error {
				return syscall.Mkfifo(filepath.Join(tmpDir, "pipe"), 0o644)
			},
		},
		{
			name: "character device",
			path: filepath.Join(tmpDir, "chardev"),
			want: true,
			setup: func() error {
				return syscall.Mknod(filepath.Join(tmpDir, "chardev"), syscall.S_IFCHR|0o644, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				if err := tt.setup(); err != nil {
					t.Skipf("File type %s creation not supported on this platform", tt.name)
				}
			}

			got := IsDevice(tt.path)
			assert.Equal(t, got, tt.want, "IsDevice()/%q", tt.name)
		})
	}
}

func TestIsFile(t *testing.T) {
	// Create temporary test files/directories
	tmpDir := t.TempDir()
	regularFile := filepath.Join(tmpDir, "regular.txt")
	err := os.WriteFile(regularFile, []byte("test"), 0o644)
	require.NoError(t, err)
	subDir := filepath.Join(tmpDir, "subdir")
	err = os.Mkdir(subDir, 0o755)
	require.NoError(t, err)

	tests := []struct {
		name  string
		path  string
		want  bool
		setup func() error
	}{
		{
			name: "existing directory",
			path: subDir,
			want: false,
		},
		{
			name: "regular file",
			path: regularFile,
			want: true,
		},
		{
			name: "non-existent path",
			path: filepath.Join(tmpDir, "nonexistent"),
			want: false,
		},
		{
			name: "symlink to file",
			path: filepath.Join(tmpDir, "symlink"),
			want: true,
			setup: func() error {
				return os.Symlink(regularFile, filepath.Join(tmpDir, "symlink"))
			},
		},
		{
			name: "socket",
			path: filepath.Join(tmpDir, "socket"),
			want: true,
			setup: func() error {
				_, err := net.Listen("unix", filepath.Join(tmpDir, "socket"))
				return err
			},
		},
		{
			name: "named pipe",
			path: filepath.Join(tmpDir, "pipe"),
			want: true,
			setup: func() error {
				return syscall.Mkfifo(filepath.Join(tmpDir, "pipe"), 0o644)
			},
		},
		{
			name: "character device",
			path: filepath.Join(tmpDir, "chardev"),
			want: true,
			setup: func() error {
				return syscall.Mknod(filepath.Join(tmpDir, "chardev"), syscall.S_IFCHR|0o644, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				if err := tt.setup(); err != nil {
					t.Skipf("File type %s creation not supported on this platform", tt.name)
				}
			}

			got := IsFile(tt.path)
			assert.Equal(t, got, tt.want, "IsFile()/%q", tt.name)
		})
	}
}
