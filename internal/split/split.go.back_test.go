package split

//
// func TestExecute(t *testing.T) {
//	// Create a temporary test directory
//	tmpDir := t.TempDir()
//	originalWd, _ := os.Getwd()
//	defer os.Chdir(originalWd)
//	os.Chdir(tmpDir)
//
//	// Create a test CSV file
//	testData := `Header1,Header2,Header3
// value1,group1,data1
// value2,group1,data2
// value3,group2,data3
// value4,group2,data4
// value5,group3,data5`
//
//	err := os.WriteFile("test.csv", []byte(testData), 0644)
//	if err != nil {
//		t.Fatalf("Failed to create test file: %v", err)
//	}
//
//	tests := []struct {
//		name    string
//		args    Flags
//		wantErr bool
//		errMsg  string
//		files   []string // Expected output files
//	}{
//		{
//			name: "successful split",
//			args: Flags{
//				InputFile: "test.csv",
//				Filter:    "Header1=value1,value2",
//				SplitCol:  "Header2",
//			},
//			wantErr: false,
//			files:   []string{"group1-value1-value2.csv"},
//		},
//		{
//			name: "missing filter",
//			args: Flags{
//				InputFile: "test.csv",
//				SplitCol:  "Header2",
//			},
//			wantErr: true,
//			errMsg:  "filter and split column must be specified",
//		},
//		{
//			name: "missing split column",
//			args: Flags{
//				InputFile: "test.csv",
//				Filter:    "Header1=value1",
//			},
//			wantErr: true,
//			errMsg:  "filter and split column must be specified",
//		},
//		{
//			name: "invalid filter format",
//			args: Flags{
//				InputFile: "test.csv",
//				Filter:    "Header1value1",
//				SplitCol:  "Header2",
//			},
//			wantErr: true,
//			errMsg:  "filter must be in format column=value1,value2,...",
//		},
//		{
//			name: "non-existent filter column",
//			args: Flags{
//				InputFile: "test.csv",
//				Filter:    "NonExistentHeader=value1",
//				SplitCol:  "Header2",
//			},
//			wantErr: true,
//			errMsg:  "filter column 'NonExistentHeader' not found in headers",
//		},
//		{
//			name: "non-existent split column",
//			args: Flags{
//				InputFile: "test.csv",
//				Filter:    "Header1=value1",
//				SplitCol:  "NonExistentHeader",
//			},
//			wantErr: true,
//			errMsg:  "split column 'NonExistentHeader' not found in headers",
//		},
//		{
//			name: "non-existent input file",
//			args: Flags{
//				InputFile: "nonexistent.csv",
//				Filter:    "Header1=value1",
//				SplitCol:  "Header2",
//			},
//			wantErr: true,
//			errMsg:  "error opening input file",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			// Clean up any existing output files
//			files, _ := filepath.Glob("*.csv")
//			for _, f := range files {
//				if f != "test.csv" {
//					os.Remove(f)
//				}
//			}
//
//			// Execute the function
//			err := Execute(tt.args)
//
//			// Check error conditions
//			if tt.wantErr {
//				if err == nil {
//					t.Errorf("Execute() expected error but got none")
//					return
//				}
//				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
//					t.Errorf("Execute() error = %v, want error containing %v", err, tt.errMsg)
//				}
//				return
//			}
//
//			if err != nil {
//				t.Errorf("Execute() unexpected error: %v", err)
//				return
//			}
//
//			// Check if expected files were created
//			for _, expectedFile := range tt.files {
//				if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
//					t.Errorf("Expected file %s was not created", expectedFile)
//				}
//			}
//		})
//	}
//}
//
//// Helper function to check if a string contains another string
// func contains(s, substr string) bool {
//	return strings.Contains(s, substr)
//}
//
//// TestExecuteWithStdin tests the function with stdin input
// func TestExecuteWithStdin(t *testing.T) {
//	// Save original stdin
//	oldStdin := os.Stdin
//	defer func() { os.Stdin = oldStdin }()
//
//	// Create a pipe
//	r, w, err := os.Pipe()
//	if err != nil {
//		t.Fatalf("Failed to create pipe: %v", err)
//	}
//
//	// Set stdin to our pipe
//	os.Stdin = r
//
//	// Write test data to pipe
//	testData := `Header1,Header2,Header3
// value1,group1,data1
// value2,group1,data2`
//
//	go func() {
//		w.Write([]byte(testData))
//		w.Close()
//	}()
//
//	// Create temporary directory for output files
//	tmpDir := t.TempDir()
//	originalWd, _ := os.Getwd()
//	defer os.Chdir(originalWd)
//	os.Chdir(tmpDir)
//
//	// Execute function
//	err = Execute(Flags{
//		Filter:   "Header1=value1,value2",
//		SplitCol: "Header2",
//	})
//
//	if err != nil {
//		t.Errorf("Execute() with stdin failed: %v", err)
//	}
//
//	// Check if output file was created
//	expectedFile := "group1-value1-value2.csv"
//	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
//		t.Errorf("Expected file %s was not created", expectedFile)
//	}
//}
