package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var VERSION string = ""
var BUILD string = ""
var RUNNING_SINCE = time.Now().UTC()

func main() {

	address := os.Getenv("TB_ADDRESS")
	crt := os.Getenv("TB_TLS_CRT")
	key := os.Getenv("TB_TLS_KEY")

	if len(address) <= 0 {
		address = ":8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlerFunction)
	mux.HandleFunc("/healthz", healthHandler)

	cfg := &tls.Config{
		MinVersion: tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
		},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		},
	}
	srv := &http.Server{
		Addr:         address,
		Handler:      mux,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	if len(crt) > 0 && len(key) > 0 {
		log.Fatal(srv.ListenAndServeTLS(crt, key))
	} else {
		log.Fatal(http.ListenAndServe(address, nil))
	}
}

func healthHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	body := []byte(fmt.Sprintf(`{"status":"ok","running_since":"%v"}`, RUNNING_SINCE))
	w.Write(body)
}

func handlerFunction(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Pull token from header (request token).
	authorizeToken := req.Header.Get("Authorization-Token")
	usernames, ok := req.URL.Query()["username"]

	if !ok || len(usernames[0]) <= 0 {
		http.Error(w, `{"error":"mising param username"}`, http.StatusBadRequest)
		return
	}
	username := usernames[0]

	// Authorize request.
	err := authorizeRequest(authorizeToken)

	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusUnauthorized)
		return
	}

	// Tabeleau request.
	token, err := getTableauToken(username)

	// Handle error on Tableau Request.
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	if string(token) == "-1" {
		http.Error(w, `{"error":"could not find username in tableau"}`, http.StatusNotFound)
		return
	}

	// Return response.
	w.WriteHeader(http.StatusOK)
	body := []byte(fmt.Sprintf(`{"token":"%s"}`, token))
	w.Write(body)
}

func authorizeRequest(token string) error {
	// Pull token from env (valid token).
	validToken := os.Getenv("TB_API_TOKEN")

	if token != validToken {
		return errors.New("invalid authorization token")
	}

	return nil
}

// getTableauToken is used to retrieve the TableauToken
func getTableauToken(username string) ([]byte, error) {
	url := os.Getenv("TB_TABLEAU_BASE_URL") + "/trusted?username=" + username

	req, err := http.NewRequest("POST", url, nil)

	c := &http.Client{}
	res, err := c.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return b, nil
}
