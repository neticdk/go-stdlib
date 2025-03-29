package file

import (
	"net"
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

func TestExists(t *testing.T) {
	// Create temporary test files/directories
	tmpDir := t.TempDir()
	regularFile := filepath.Join(tmpDir, "regular.txt")
	if err := os.WriteFile(regularFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}

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
				if err := os.Mkdir(filepath.Join(tmpDir, "noperm"), 0755); err != nil {
					return err
				}
				if err := os.WriteFile(filepath.Join(tmpDir, "noperm", "test.txt"), []byte("test"), 0000); err != nil {
					return err
				}
				if err := os.Chmod(filepath.Join(tmpDir, "noperm"), 0000); err != nil {
					return err
				}
				return nil
			},
			teardown: func() error {
				if err := os.Chmod(filepath.Join(tmpDir, "noperm"), 0755); err != nil {
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
				return syscall.Mkfifo(filepath.Join(tmpDir, "pipe"), 0644)
			},
		},
		{
			name: "character device",
			path: filepath.Join(tmpDir, "chardev"),
			want: true,
			setup: func() error {
				return syscall.Mknod(filepath.Join(tmpDir, "chardev"), syscall.S_IFCHR|0644, 0)
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
			if (err != nil) != tt.wantErr {
				t.Errorf("Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}

			if tt.teardown != nil {
				if err := tt.teardown(); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestIsDir(t *testing.T) {
	// Create temporary test files/directories
	tmpDir := t.TempDir()
	regularFile := filepath.Join(tmpDir, "regular.txt")
	if err := os.WriteFile(regularFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}

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
				return syscall.Mkfifo(filepath.Join(tmpDir, "pipe"), 0644)
			},
		},
		{
			name: "character device",
			path: filepath.Join(tmpDir, "chardev"),
			want: false,
			setup: func() error {
				return syscall.Mknod(filepath.Join(tmpDir, "chardev"), syscall.S_IFCHR|0644, 0)
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

			if got := IsDir(tt.path); got != tt.want {
				t.Errorf("IsDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsRegular(t *testing.T) {
	// Create temporary test files/directories
	tmpDir := t.TempDir()
	regularFile := filepath.Join(tmpDir, "regular.txt")
	if err := os.WriteFile(regularFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}

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
				return syscall.Mkfifo(filepath.Join(tmpDir, "pipe"), 0644)
			},
		},
		{
			name: "character device",
			path: filepath.Join(tmpDir, "chardev"),
			want: false,
			setup: func() error {
				return syscall.Mknod(filepath.Join(tmpDir, "chardev"), syscall.S_IFCHR|0644, 0)
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

			if got := IsRegular(tt.path); got != tt.want {
				t.Errorf("IsRegular() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSymlink(t *testing.T) {
	// Create temporary test files/directories
	tmpDir := t.TempDir()
	regularFile := filepath.Join(tmpDir, "regular.txt")
	if err := os.WriteFile(regularFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}

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
				return syscall.Mkfifo(filepath.Join(tmpDir, "pipe"), 0644)
			},
		},
		{
			name: "character device",
			path: filepath.Join(tmpDir, "chardev"),
			want: false,
			setup: func() error {
				return syscall.Mknod(filepath.Join(tmpDir, "chardev"), syscall.S_IFCHR|0644, 0)
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

			if got := IsSymlink(tt.path); got != tt.want {
				t.Errorf("IsSymlink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSocket(t *testing.T) {
	// Create temporary test files/directories
	tmpDir := t.TempDir()
	regularFile := filepath.Join(tmpDir, "regular.txt")
	if err := os.WriteFile(regularFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}

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
				return syscall.Mkfifo(filepath.Join(tmpDir, "pipe"), 0644)
			},
		},
		{
			name: "character device",
			path: filepath.Join(tmpDir, "chardev"),
			want: false,
			setup: func() error {
				return syscall.Mknod(filepath.Join(tmpDir, "chardev"), syscall.S_IFCHR|0644, 0)
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

			if got := IsSocket(tt.path); got != tt.want {
				t.Errorf("IsFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNamedPipe(t *testing.T) {
	// Create temporary test files/directories
	tmpDir := t.TempDir()
	regularFile := filepath.Join(tmpDir, "regular.txt")
	if err := os.WriteFile(regularFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}

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
				return syscall.Mkfifo(filepath.Join(tmpDir, "pipe"), 0644)
			},
		},
		{
			name: "character device",
			path: filepath.Join(tmpDir, "chardev"),
			want: false,
			setup: func() error {
				return syscall.Mknod(filepath.Join(tmpDir, "chardev"), syscall.S_IFCHR|0644, 0)
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

			if got := IsNamedPipe(tt.path); got != tt.want {
				t.Errorf("IsFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsDevice(t *testing.T) {
	// Create temporary test files/directories
	tmpDir := t.TempDir()
	regularFile := filepath.Join(tmpDir, "regular.txt")
	if err := os.WriteFile(regularFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}

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
				return syscall.Mkfifo(filepath.Join(tmpDir, "pipe"), 0644)
			},
		},
		{
			name: "character device",
			path: filepath.Join(tmpDir, "chardev"),
			want: true,
			setup: func() error {
				return syscall.Mknod(filepath.Join(tmpDir, "chardev"), syscall.S_IFCHR|0644, 0)
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

			if got := IsDevice(tt.path); got != tt.want {
				t.Errorf("IsFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsFile(t *testing.T) {
	// Create temporary test files/directories
	tmpDir := t.TempDir()
	regularFile := filepath.Join(tmpDir, "regular.txt")
	if err := os.WriteFile(regularFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}

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
				return syscall.Mkfifo(filepath.Join(tmpDir, "pipe"), 0644)
			},
		},
		{
			name: "character device",
			path: filepath.Join(tmpDir, "chardev"),
			want: true,
			setup: func() error {
				return syscall.Mknod(filepath.Join(tmpDir, "chardev"), syscall.S_IFCHR|0644, 0)
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

			if got := IsFile(tt.path); got != tt.want {
				t.Errorf("IsFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
