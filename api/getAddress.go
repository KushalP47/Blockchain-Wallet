package api

import "net/http"

func GetAddress(w http.ResponseWriter, r *http.Request) {
	// Get address
	w.Write([]byte(myCredentials.Address.Hex()))
}
