package main

import (
	"fmt"
	"github.com/timam/uttaracloud-finance-backend/cmd/server"
	"github.com/timam/uttaracloud-finance-backend/internals/packages"
	"log"
)

func init() {
	err := packages.InitializePackages()
	if err != nil {
		log.Fatalf("Initialization failed: %v", err)
	}

}
func main() {
	fmt.Println("Bipul")
	server.StartServer()
}
