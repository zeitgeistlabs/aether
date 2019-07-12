package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
)

var urlBase = "aether"

func readToken() string {
	file, err := os.Open("/var/run/secrets/tokens/aether-token")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	return string(b)
}

type AetherAuthResponse struct {
	Token string
}

func main() {
	fmt.Println("Running Aether in-pod tests")

	serviceToken := readToken()


	body := url.Values{
		"token": []string{serviceToken},
	}
	encodedBody := body.Encode()

	req, err := http.NewRequest("POST", "http://" + path.Join(urlBase, "auth"), strings.NewReader(encodedBody))
	fmt.Println(req)
	if err != nil {
		fmt.Println("Error building request")
		panic(err)
	}
	req.Header.Add("Host", "aether")
	req.Header.Add("Token", serviceToken)
	fmt.Println(serviceToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(encodedBody)))

	client := &http.Client{}
	resp, err := client.Do(req)
	fmt.Println(resp)
	if err != nil {
		fmt.Println("Error submitting request")
		panic(err)
	}
	if resp.StatusCode != 200 {
		panic("http code " + resp.Status)
	}

	authResponse := &AetherAuthResponse{}
	err = json.NewDecoder(resp.Body).Decode(authResponse)
	if err != nil {
		fmt.Println("Error decoding response")
		panic(err)
	}

	if len(authResponse.Token) == 0 {
		panic(errors.New("received no aether token"))
	}
}
