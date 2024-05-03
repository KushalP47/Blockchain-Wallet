package main

import (
	"addressGenerator/api"
	"log"
	"net/http"
)

func main() {
	// Define HTTP handlers
	http.HandleFunc("/generateAddress", api.GenerateAddress)
	http.HandleFunc("/signTxn", api.SignTxn)
	// http.HandleFunc("/printAccount", api.PrintAccountHandler)

	// Start the HTTP server
	log.Println("Server started on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
