package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

// TableauRequest is used to construct the request to Tableau Server.
type TableauRequest struct {
	Credentials TableauCredentials `json:"credentials"`
}

// TableauCredentials is used to store creds for the request to Tableau Server.
type TableauCredentials struct {
	Name     string      `json:"name"`
	Password string      `json:"password"`
	Site     TableauSite `json:"site"`
}

// TableauSite is used for including a content URL.
type TableauSite struct {
	ContentURL string `json:"contentUrl"`
}

func retrieveTableauToken() ([]byte, error) {
	url := os.Getenv("TABLEAU_BASE_URL") + "/api/3.1/auth/signin"
	jsonRequest, err := generateTableauRequest()

	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonRequest))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func generateTableauRequest() ([]byte, error) {
	request := TableauRequest{
		Credentials: TableauCredentials{
			Name:     os.Getenv("TABLEAU_NAME"),
			Password: os.Getenv("TABLEAU_PASSWORD"),
			Site: TableauSite{
				ContentURL: "",
			},
		},
	}

	// Format request.
	jsonRequest, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	return jsonRequest, nil
}
