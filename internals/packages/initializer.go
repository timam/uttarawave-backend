package packages

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/spf13/viper"
	"github.com/timam/uttaracloud-finance-backend/models"
	"os"
	"path/filepath"
	"strings"
)

var CurrentInternetPackages []models.InternetPackage //loading data from cmd/server/main.go

func InitializePackages() error {
	currentInternetPackagesFile, err := LoadCurrentInternetPackages()
	if err != nil {
		return fmt.Errorf("failed to load latest packages: %v", err)
	}

	currentInternetPackages, err := InternetPackageParser(currentInternetPackagesFile)
	if err != nil {
		return fmt.Errorf("failed to parse packages from CSV: %v", err)
	}

	CurrentInternetPackages = currentInternetPackages

	return nil
}

func LoadCurrentInternetPackages() (string, error) {
	var latestFile string

	err := filepath.Walk(viper.GetString("paths.packages.internet"), func(path string, info os.FileInfo, err error) error {
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
	return filepath.Join(viper.GetString("paths.packages.internet"), latestFile), nil
}

func InternetPackageParser(filePath string) ([]models.InternetPackage, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var packages []*models.InternetPackage
	if err := gocsv.UnmarshalFile(file, &packages); err != nil {
		return nil, err
	}

	var result []models.InternetPackage
	for _, pkg := range packages {
		result = append(result, *pkg)
	}
	return result, nil
}
