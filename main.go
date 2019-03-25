package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"./tableau"
)

func main() {
	http.HandleFunc("/", handlerFunction)

	address := os.Getenv("TB_ADDRESS")
	crt := os.Getenv("TB_TLS_CRT")
	key := os.Getenv("TB_TLS_KEY")

	if len(address) <= 0 {
		address = ":8080"
	}

	if len(crt) > 0 && len(key) > 0 {
		log.Fatal(http.ListenAndServeTLS(address, crt, key, nil))
	} else {
		log.Fatal(http.ListenAndServe(address, nil))
	}
}

func handlerFunction(w http.ResponseWriter, req *http.Request) {
	// Pull token from header (request token).
	authorizeToken := req.Header.Get("Authorization-Token")

	// Authorize request.
	err := authorizeRequest(authorizeToken)

	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Tabeleau request.
	token, err := tableau.RetrieveToken()

	// Handle error on Tableau Request.
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(token)
}

func authorizeRequest(token string) error {
	// Pull token from env (valid token).
	validToken := os.Getenv("TB_API_TOKEN")

	if token != validToken {
		return errors.New("invalid token")
	}

	return nil
}
