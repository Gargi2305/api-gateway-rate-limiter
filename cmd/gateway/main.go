package main

import (
	"log"

	"api-gateway/internal/server"
)

func main() {
	srv := server.New()

	log.Println("API Gateway starting on port 8080...")
	log.Fatal(srv.ListenAndServe())
}
