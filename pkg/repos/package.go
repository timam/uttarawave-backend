package repos

import (
	"encoding/csv"
	"fmt"
	"github.com/timam/uttaracloud-finance-backend/pkg/models"
	"os"
	"strconv"
)

var CurrentPackages []models.Package

func LoadLatestPackages(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("could not read CSV file: %v", err)
	}

	for i, record := range records[1:] { // skip header
		if len(record) != 6 {
			fmt.Printf("Skipping invalid record at line %d: %+v\n", i+2, record)
			continue
		}

		name := record[0]
		speedStr := record[1]
		priceStr := record[2]
		connection := models.PackageType(record[3])
		typ := models.ConnectionType(record[4])
		realIPBool, parseErr := strconv.ParseBool(record[5])
		if parseErr != nil {
			fmt.Printf("Invalid RealIP value at line %d: %+v\n", i+2, record[5])
			realIPBool = false
		}
		realIP := strconv.FormatBool(realIPBool)

		pkg := models.Package{
			Name:       name,
			Speed:      speedStr,
			Price:      priceStr,
			Connection: connection,
			Type:       typ,
			RealIP:     realIP,
		}

		CurrentPackages = append(CurrentPackages, pkg)

		// Example log for each package loaded
		fmt.Printf("Loaded package #%d: %+v\n", i+1, pkg)
	}

	return nil
}

//TODO: Define LoadPackagesFromCSV : Should load packages from any given CSV file
