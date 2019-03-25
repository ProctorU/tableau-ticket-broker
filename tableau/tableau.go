package tableau

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

// Payload is used to construct the request to Tableau Server.
type Payload struct {
	Credentials Credentials `json:"credentials"`
}

// Credentials is used to store creds for the request to Tableau Server.
type Credentials struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Site     Site   `json:"site"`
}

// Site is used for including a content URL.
type Site struct {
	ContentURL string `json:"contentUrl"`
}

// RetrieveToken is used to retrieve the TableauToken
func RetrieveToken() ([]byte, error) {
	url := os.Getenv("TB_TABLEAU_BASE_URL") + "/api/3.1/auth/signin"
	p := generatePayload()
	j, err := json.Marshal(p)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

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

func generatePayload() Payload {
	return Payload{
		Credentials: Credentials{
			Name:     os.Getenv("TB_TABLEAU_USERNAME"),
			Password: os.Getenv("TB_TABLEAU_PASSWORD"),
			Site: Site{
				ContentURL: "",
			},
		},
	}
}
