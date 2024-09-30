package packages

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/timam/uttaracloud-finance-backend/models"
	"os"
	"path/filepath"
	"strings"
)

func LoadLatestPackages() (string, error) {
	var latestFile string

	err := filepath.Walk("data/packages", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".csv" {
			filename := filepath.Base(path)
			if len(filename) == len("YYYYMMDD.csv") && strings.HasSuffix(filename, ".csv") {
				if filename > latestFile {
					latestFile = filename
				}
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
	return filepath.Join("data/packages", latestFile), nil
}

func ParseCSV(filePath string) ([]models.Package, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var packages []*models.Package
	if err := gocsv.UnmarshalFile(file, &packages); err != nil {
		return nil, err
	}

	var result []models.Package
	for _, pkg := range packages {
		result = append(result, *pkg)
	}
	return result, nil
}
