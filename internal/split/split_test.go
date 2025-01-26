package split

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name        string
	flags       Flags
	expectFiles []string // Expected output files relative to testdata directory
	wantErr     bool
	errMsg      string
}

func TestExecute(t *testing.T) {
	// Get absolute path to testdata directory before changing directories
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	testDataDir := filepath.Join(originalWd, "test-data")
	tmpDir := t.TempDir()

	tests := []testCase{
		{
			name: "basic split",
			flags: Flags{
				InputFile: "input/merged.csv",
				Filter:    "Tag=Aircraft",
				SplitCol:  "County",
			},
			expectFiles: []string{
				// "internal/split/test-data/expected/bexar-Aircraft.csv",
				"expected/bexar-Aircraft.csv",
				"expected/bexar-Aircraft.csv",
				"expected/comal-Aircraft.csv",
				"expected/guadalupe-Aircraft.csv",
				"expected/hays-Aircraft.csv",
			},
			wantErr: false,
		},

		{
			name: "missing column",
			flags: Flags{
				InputFile: "input/missing_column.csv",
				Filter:    "type=car,truck",
				SplitCol:  "nonexistent",
			},
			wantErr: true,
			errMsg:  "split column 'nonexistent' not found in headers",
		},
		{
			name: "empty input",
			flags: Flags{
				InputFile: "input/empty.csv",
				Filter:    "type=car,truck",
				SplitCol:  "color",
			},
			expectFiles: []string{},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up tmp dir between tests
			cleanTmpDir(t, tmpDir)

			// Copy input file to tmp dir if specified
			if tt.flags.InputFile != "" {
				srcPath := filepath.Join(testDataDir, tt.flags.InputFile)
				destPath := filepath.Join(tmpDir, filepath.Base(tt.flags.InputFile))
				if err := copyFile(srcPath, destPath); err != nil {
					t.Fatalf("Failed to copy input file: %v", err)
				}
				// Update flags to use tmp dir path
				tt.flags.InputFile = destPath
			}

			// Change working directory to tmp dir for output files
			if err := os.Chdir(tmpDir); err != nil {
				t.Fatalf("Failed to change working directory: %v", err)
			}
			defer func() {
				if err := os.Chdir(originalWd); err != nil {
					t.Fatalf("Failed to restore working directory: %v", err)
				}
			}()

			// Execute the function
			err = Execute(tt.flags)

			// Check error conditions
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("Execute() error = %v, want error %v", err, tt.errMsg)
				return
			}

			// Verify output files
			for _, expectFile := range tt.expectFiles {
				expectedPath := filepath.Join(testDataDir, expectFile)
				actualPath := filepath.Join(tmpDir, filepath.Base(expectFile))

				expectedContent, err := os.ReadFile(expectedPath)
				if err != nil {
					t.Fatalf("Failed to read expected file: %v", err)
				}

				actualContent, err := os.ReadFile(actualPath)
				if err != nil {
					t.Fatalf("Failed to read actual file: %v", err)
				}

				normalizedExpected := strings.ReplaceAll(strings.ReplaceAll(string(expectedContent), "\r\n", "\n"), "\r", "\n")
				normalizedActual := strings.ReplaceAll(strings.ReplaceAll(string(actualContent), "\r\n", "\n"), "\r", "\n")

				assert.Equal(t, normalizedExpected, normalizedActual)
			}
		})
	}
}

// Helper functions.
func cleanTmpDir(t *testing.T, dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("Failed to read tmp dir: %v", err)
	}
	for _, entry := range entries {
		if err := os.Remove(filepath.Join(dir, entry.Name())); err != nil {
			t.Fatalf("Failed to clean tmp dir: %v", err)
		}
	}
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0o644)
}
