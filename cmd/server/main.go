package server

import (
	"github.com/timam/uttaracloud-finance-backend/pkg/routers"
	"log"
)

func startServer(initFunc func() error) error {
	err := initFunc()
	if err != nil {
		return err
	}

	router := routers.InitRouter()
	log.Println("Starting server...")
	err = router.Run(":8080")
	if err != nil {
		return err
	}
	return nil
}

func StartServer() {
	if err := startServer(Initialize); err != nil {
		log.Fatalf("Initialization failed: %v", err)
	}
}
