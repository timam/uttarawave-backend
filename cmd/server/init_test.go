package server

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
				_ = os.MkdirAll("data/packages", os.ModePerm)
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Return latest file",
			setup: func() {
				_ = os.MkdirAll("data/packages", os.ModePerm)
				files := []string{"20230101.csv", "20230201.csv"}
				for _, file := range files {
					_ = os.WriteFile(filepath.Join("data/packages", file), []byte{}, os.ModePerm)
				}
			},
			want:    "data/packages/20230201.csv",
			wantErr: false,
		},
		{
			name: "Return File Outside the Date Scope",
			setup: func() {
				_ = os.MkdirAll("data/packages", os.ModePerm)
				files := []string{"abc.csv", "11.csv"}
				for _, file := range files {
					_ = os.WriteFile(filepath.Join("data/packages", file), []byte{}, os.ModePerm)
				}
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer os.RemoveAll("data/packages")

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

func TestInitialize(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Returns error when unable to load latest packages",
			setup: func() {
			},
			wantErr: true,
		},
		{
			name: "Returns error when unable to parse CSV",
			setup: func() {
				_ = os.MkdirAll("data/packages", os.ModePerm)
				files := []string{"20230101.csv"}
				for _, file := range files {
					_ = os.WriteFile(filepath.Join("data/packages", file), []byte{}, os.ModePerm)
				}
			},
			wantErr: true,
		},
		{
			name: "Successful initialization",
			setup: func() {
				_ = os.MkdirAll("data/packages", os.ModePerm)
				files := []string{"20230101.csv"}
				for _, file := range files {
					_ = os.WriteFile(filepath.Join("data/packages", file), []byte("Name,Bandwidth,Price,Usage,Type,RealIP\nhome10,10Mbps,500,home,shared,false\nbusiness100,100Mbps,3500,business,dedicate,true\n"), os.ModePerm)
				}
			},
			wantErr: false,
		},
		{
			name: "Initialize error - invalid packages file",
			setup: func() {
				_ = os.MkdirAll("data/packages", os.ModePerm)
				files := []string{"20230101.csv"}
				for _, file := range files {
					_ = os.WriteFile(filepath.Join("data/packages", file), []byte("Name,Bandwidth,Price,Usage,Type,RealIP\nInvalid Data Line"), os.ModePerm)
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer func() {
				err := os.RemoveAll("data/packages")
				if err != nil {

				}
			}()

			err := Initialize()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

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
			_ = os.MkdirAll("data/packages", os.ModePerm)
			tmpfile, _ := os.Create("data/packages/tmp.csv")
			tmpfile.WriteString(tt.fileSet)
			tmpfile.Close()

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

			_ = os.Remove("data/packages/tmp.csv")
		})
	}
}
