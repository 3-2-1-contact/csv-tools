package merge

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExecute(t *testing.T) {
	// Create temporary test files
	tempDir := t.TempDir()

	// Create test file 1
	file1Content := "header1,header2\nvalue1,value2\nvalue3,value4"
	file1Path := filepath.Join(tempDir, "test1.csv")
	if err := os.WriteFile(file1Path, []byte(file1Content), 0644); err != nil {
		t.Fatalf("Failed to create test file 1: %v", err)
	}

	// Create test file 2
	file2Content := "header1,header2\nvalue5,value6\nvalue7,value8"
	file2Path := filepath.Join(tempDir, "test2.csv")
	if err := os.WriteFile(file2Path, []byte(file2Content), 0644); err != nil {
		t.Fatalf("Failed to create test file 2: %v", err)
	}

	tests := []struct {
		name    string
		args    Flags
		wantErr bool
		want    string
	}{
		{
			name: "successful merge with new column",
			args: Flags{
				InputFiles: file1Path + "," + file2Path,
				OutputFile: filepath.Join(tempDir, "output.csv"),
				NewColumn:  "source",
			},
			wantErr: false,
			want:    "header1,header2,source\nvalue1,value2,test1\nvalue3,value4,test1\nvalue5,value6,test2\nvalue7,value8,test2\n",
		},
		{
			name: "successful merge without new column",
			args: Flags{
				InputFiles: file1Path + "," + file2Path,
				OutputFile: filepath.Join(tempDir, "output2.csv"),
			},
			wantErr: false,
			want:    "header1,header2\nvalue1,value2\nvalue3,value4\nvalue5,value6\nvalue7,value8\n",
		},
		{
			name: "error - no input files",
			args: Flags{
				InputFiles: "",
				OutputFile: filepath.Join(tempDir, "output3.csv"),
			},
			wantErr: true,
		},
		{
			name: "error - mismatched headers",
			args: Flags{
				InputFiles: file1Path + "," + filepath.Join(tempDir, "nonexistent.csv"),
				OutputFile: filepath.Join(tempDir, "output4.csv"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// If testing stdout output
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := Execute(tt.args)

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout

			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				var output string
				if tt.args.OutputFile != "" {
					// Read from output file
					content, err := os.ReadFile(tt.args.OutputFile)
					if err != nil {
						t.Fatalf("Failed to read output file: %v", err)
					}
					output = string(content)
				} else {
					// Read from stdout
					var buf bytes.Buffer
					buf.ReadFrom(r)
					output = buf.String()
				}

				// Normalize line endings
				output = strings.ReplaceAll(output, "\r\n", "\n")
				if output != tt.want {
					t.Errorf("Execute() output = %v, want %v", output, tt.want)
				}
			}
		})
	}
}

func TestCompareHeaders(t *testing.T) {
	tests := []struct {
		name string
		a    []string
		b    []string
		want bool
	}{
		{
			name: "identical headers",
			a:    []string{"header1", "header2"},
			b:    []string{"header1", "header2"},
			want: true,
		},
		{
			name: "different headers",
			a:    []string{"header1", "header2"},
			b:    []string{"header1", "header3"},
			want: false,
		},
		{
			name: "different lengths",
			a:    []string{"header1", "header2"},
			b:    []string{"header1"},
			want: false,
		},
		{
			name: "empty headers",
			a:    []string{},
			b:    []string{},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareHeaders(tt.a, tt.b); got != tt.want {
				t.Errorf("compareHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}
