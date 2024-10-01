package server

import (
	"github.com/timam/uttaracloud-finance-backend/routers"
	"log"
)

func StartServer() {

	router := routers.InitRouter()
	log.Println("Starting server...")
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
