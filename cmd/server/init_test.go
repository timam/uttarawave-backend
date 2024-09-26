package server

import (
	"os"
	"path/filepath"
	"testing"
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

			if (err != nil) != tt.wantErr {
				t.Errorf("LoadLatestPackages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LoadLatestPackages() = %v, want %v", got, tt.want)
			}
		})
	}
}
