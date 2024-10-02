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

var CurrentInternetPackages []models.InternetPackage
var CurrentCableTVPackages []models.CableTVPackage

var LatestInternetPackagesFile string
var LatestCableTVPackagesFile string

func InitializePackages() error {
	// Initialize Internet packages
	currentInternetPackagesFile, err := LoadCurrentInternetPackages()
	if err != nil {
		return fmt.Errorf("failed to load latest internet packages: %v", err)
	}

	currentInternetPackages, err := InternetPackageParser(currentInternetPackagesFile)
	if err != nil {
		return fmt.Errorf("failed to parse internet packages from CSV: %v", err)
	}

	CurrentInternetPackages = currentInternetPackages
	LatestInternetPackagesFile = getFileNameWithoutExtension(currentInternetPackagesFile)

	// Initialize Cable TV packages
	currentCableTVPackagesFile, err := LoadCurrentCableTVPackages()
	if err != nil {
		return fmt.Errorf("failed to load latest cable TV packages: %v", err)
	}

	currentCableTVPackages, err := CableTVPackageParser(currentCableTVPackagesFile)
	if err != nil {
		return fmt.Errorf("failed to parse cable TV packages from CSV: %v", err)
	}

	CurrentCableTVPackages = currentCableTVPackages
	LatestCableTVPackagesFile = getFileNameWithoutExtension(currentCableTVPackagesFile)

	return nil
}

func LoadCurrentInternetPackages() (string, error) {
	return loadLatestPackageFile(viper.GetString("paths.packages.internet"))
}

func LoadCurrentCableTVPackages() (string, error) {
	return loadLatestPackageFile(viper.GetString("paths.packages.cabletv"))
}

func loadLatestPackageFile(directoryPath string) (string, error) {
	var latestFile string

	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
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
		return "", fmt.Errorf("no CSV files found in '%s' directory", directoryPath)
	}
	return filepath.Join(directoryPath, latestFile), nil
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

func CableTVPackageParser(filePath string) ([]models.CableTVPackage, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var packages []*models.CableTVPackage
	if err := gocsv.UnmarshalFile(file, &packages); err != nil {
		return nil, err
	}

	var result []models.CableTVPackage
	for _, pkg := range packages {
		result = append(result, *pkg)
	}
	return result, nil
}

func getFileNameWithoutExtension(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}
