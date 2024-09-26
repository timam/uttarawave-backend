package server

import (
	"os"
	"path/filepath"
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
				// setup your filesystem mock here to have no csv files
				_ = os.MkdirAll("data/packages", os.ModePerm)
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Return latest file",
			setup: func() {
				// setup your filesystem mock here to have some csv files
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
				// setup your filesystem mock here to have files outside the date scope (e.g. "abc.csv", "11.csv")
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
			assert.Equal(t, tt.want, got)
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
				// setup your LoadLatestPackages mock here to return error
			},
			wantErr: true,
		},
		{
			name: "Returns error when unable to parse CSV",
			setup: func() {
				// setup your LoadLatestPackages mock here to return a valid path
				// setup your ParseCSV mock here to return error
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
				// setup your LoadLatestPackages mock here to return a valid path
				// setup your ParseCSV mock here to return a valid list of packages
				_ = os.MkdirAll("data/packages", os.ModePerm)
				files := []string{"20230101.csv"}
				for _, file := range files {
					_ = os.WriteFile(filepath.Join("data/packages", file), []byte("Name,Bandwidth,Price,Usage,Type,RealIP\nhome10,10Mbps,500,home,shared,false\n"), os.ModePerm)
				}
			},
			wantErr: false,
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
