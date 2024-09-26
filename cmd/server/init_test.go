package server

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Test for LoadLatestPackages function
func TestLoadLatestPackages(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		want    string
		wantErr bool
	}{
		{
			name: "Return error when no csv files",
			setup: func() {
				createDir(t, "data/packages")
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Return latest file",
			setup: func() {
				createDir(t, "data/packages")
				files := []string{"20230101.csv", "20230201.csv"}
				for _, file := range files {
					writeFile(t, filepath.Join("data/packages", file), []byte{})
				}
			},
			want:    "data/packages/20230201.csv",
			wantErr: false,
		},
		{
			name: "Return File Outside the Date Scope",
			setup: func() {
				createDir(t, "data/packages")
				files := []string{"abc.csv", "11.csv"}
				for _, file := range files {
					writeFile(t, filepath.Join("data/packages", file), []byte{})
				}
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer cleanupDir(t, "data")

			got, err := LoadLatestPackages()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if (got != "" && filepath.Base(got) != filepath.Base(tt.want)) || (got == "" && tt.want != "") {
				t.Errorf("got = %v, want = %v", got, tt.want)
			}
		})
	}
}

// Test for Initialize function
func TestInitialize(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T)
		wantErr bool
	}{
		{
			name: "Returns error when unable to load latest packages",
			setup: func(t *testing.T) {
				// Simulate inability to load latest packages (e.g., by not creating any file)
			},
			wantErr: true,
		},
		{
			name: "Returns error when unable to parse CSV",
			setup: func(t *testing.T) {
				dirPath := "data/packages"
				createDir(t, dirPath)
				filePath := filepath.Join(dirPath, "20230101.csv")
				writeFile(t, filePath, []byte{})
			},
			wantErr: true,
		},
		{
			name: "Successful initialization",
			setup: func(t *testing.T) {
				dirPath := "data/packages"
				createDir(t, dirPath)
				fileContent := `Name,Bandwidth,Price,Usage,Type,RealIP
home10,10Mbps,500,home,shared,false
business100,100Mbps,3500,business,dedicate,true
`
				filePath := filepath.Join(dirPath, "20230101.csv")
				writeFile(t, filePath, []byte(fileContent))
			},
			wantErr: false,
		},
		{
			name: "Initialize error - invalid packages file",
			setup: func(t *testing.T) {
				dirPath := "data/packages"
				createDir(t, dirPath)
				fileContent := `Name,Bandwidth,Price,Usage,Type,RealIP
Invalid Data Line`
				filePath := filepath.Join(dirPath, "20230101.csv")
				writeFile(t, filePath, []byte(fileContent))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)

			defer cleanupDir(t, "data")

			err := Initialize()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Helper function to create directories
func createDir(t *testing.T, dirPath string) {
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	// Ensure the created directory is cleaned up after the test
	t.Cleanup(func() {
		cleanupDir(t, dirPath)
	})
}

// Helper function to write content to files
func writeFile(t *testing.T, filePath string, content []byte) {
	err := os.WriteFile(filePath, content, os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
}

// Helper function to clean up directories
func cleanupDir(t *testing.T, dirPath string) {
	err := os.RemoveAll(dirPath)
	if err != nil {
		return
	}
}

// Test for ParseCSV function
func TestParseCSV(t *testing.T) {
	tests := []struct {
		name    string
		fileSet string
		wantErr bool
	}{
		{
			name:    "Valid CSV file",
			fileSet: "Name,Bandwidth,Price,Usage,Type,RealIP\nhome10,10Mbps,500,home,shared,false\nbusiness100,100Mbps,3500,business,dedicate,true\n",
			wantErr: false,
		},
		{
			name:    "Invalid CSV file",
			fileSet: "Name,Bandwidth,Price,Usage,Type,RealIP\nInvalid CSV Content\n",
			wantErr: true,
		},
		{
			name:    "Empty CSV file",
			fileSet: "", // Empty CSV file should still have headers
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createDir(t, "data/packages")
			tmpfile, _ := os.Create("data/packages/tmp.csv")
			tmpfile.WriteString(tt.fileSet)
			tmpfile.Close()

			defer cleanupDir(t, "data")

			if strings.TrimSpace(tt.fileSet) == "" {
				_, err := ParseCSV("data/packages/tmp.csv")
				assert.Error(t, err)
			} else {
				_, err := ParseCSV("data/packages/tmp.csv")
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}
