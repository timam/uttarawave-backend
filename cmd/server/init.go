package server

import (
	"fmt"
	"github.com/timam/uttaracloud-finance-backend/pkg/packages"
	"github.com/timam/uttaracloud-finance-backend/storage"
)

func InitializePackages() error {
	currentInternetPackagesFile, err := packages.LoadCurrentInternetPackages()
	if err != nil {
		return fmt.Errorf("failed to load latest packages: %v", err)
	}

	currentInternetPackages, err := packages.InternetPackageParser(currentInternetPackagesFile)
	if err != nil {
		return fmt.Errorf("failed to parse packages from CSV: %v", err)
	}

	storage.CurrentInternetPackages = currentInternetPackages

	return nil
}
