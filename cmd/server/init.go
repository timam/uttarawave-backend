package server

import (
	"fmt"
	"github.com/timam/uttaracloud-finance-backend/pkg/repos"
	"os"
	"path/filepath"
	"time"
)

func LoadLatestPackages() (string, error) {
	var latestFile string
	var latestTime time.Time

	err := filepath.Walk("data/packages", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".csv" {
			fileTime := info.ModTime()
			if fileTime.After(latestTime) {
				latestTime = fileTime
				latestFile = path
			}
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	if latestFile == "" {
		return "", fmt.Errorf("no CSV files found in 'data/packages' directory")
	}

	return latestFile, nil
}

func Initialize() error {
	latestPackages, err := LoadLatestPackages()
	if err != nil {
		return fmt.Errorf("failed to load latest packages: %v", err)
	}

	err = repos.LoadLatestPackages(latestPackages)
	if err != nil {
		return fmt.Errorf("failed to load latest packages: %v", err)
	}
	return nil
}
