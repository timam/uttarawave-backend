package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/timam/uttaracloud-finance-backend/pkg/routers"
	"github.com/timam/uttaracloud-finance-backend/pkg/storage"
	"sync"
	"testing"
)

var (
	initRouterMutex sync.Mutex
	oldRouter       func() *gin.Engine
	origInitialize  func() error
)

func init() {
	oldRouter = routers.InitRouter
	origInitialize = InitializePackages
}

func InitializePackages() error {
	latestPackagesFile, err := LoadLatestPackages()
	if err != nil {
		return fmt.Errorf("failed to load latest packages: %v", err)
	}

	packages, err := ParseCSV(latestPackagesFile)
	if err != nil {
		return fmt.Errorf("failed to parse packages from CSV: %v", err)
	}

	storage.LoadedPackages = packages

	return nil
}

func TestStartServer(t *testing.T) {
	// Mock the router
	initRouterMutex.Lock()
	defer initRouterMutex.Unlock()

	defer func() {}()

	// Test cases
	testCases := []struct {
		name    string
		setUp   func()
		cleanup func()
		check   func(*testing.T)
	}{
		{
			name: "Server Initialization Fails",
			setUp: func() {

			},
			cleanup: func() {

			},
			check: func(t *testing.T) {
				errorChannel := startServerAsync()
				assert.Error(t, <-errorChannel)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setUp()
			defer tc.cleanup()
			tc.check(t)
		})
	}
}

func startServerAsync() chan error {
	errorChannel := make(chan error, 1)
	go func() {
		defer close(errorChannel)
		defer func() {
			if r := recover(); r != nil {
				errorChannel <- fmt.Errorf("recovered from panic: %v", r)
			}
		}()
		StartServer()
	}()
	return errorChannel
}
