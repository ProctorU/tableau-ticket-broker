package main

import (
	"errors"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", handlerFunction)
	log.Fatal(http.ListenAndServe(":8080", nil))
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
	token, err := retrieveTableauToken()

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
	validToken := os.Getenv("TABLEAU_SERVER_TOKEN")

	if token != validToken {
		return errors.New("invalid token")
	}

	return nil
}
