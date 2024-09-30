package server

import (
	"fmt"
	"github.com/timam/uttaracloud-finance-backend/pkg/packages"
	"github.com/timam/uttaracloud-finance-backend/storage"
)

func Initialize() error {
	latestPackagesFile, err := packages.LoadCurrentPackages()
	if err != nil {
		return fmt.Errorf("failed to load latest packages: %v", err)
	}

	pack, err := packages.ParseCSV(latestPackagesFile)
	if err != nil {
		return fmt.Errorf("failed to parse packages from CSV: %v", err)
	}

	storage.LoadedPackages = pack

	return nil
}
