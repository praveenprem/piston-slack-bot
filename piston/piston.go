package piston

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type (
	Piston struct {
		*Execute
		Versions *[]Version
		Response Response
	}

	Execute struct {
		Language string   `json:"language"`
		Source   string   `json:"source"`
		Args     []string `json:"args"`
	}

	Version struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		*Aliases
	}

	Response struct {
		Ran      bool   `json:"ran"`
		Message  string `json:"message"`
		Language string `json:"language"`
		Version  string `json:"version"`
		Stdout   string `json:"stdout"`
		Stderr   string `json:"stderr"`
	}

	Aliases struct {
		Aliases []string `json:"aliases"`
	}

	Pistoner interface {
		Lang() error
		Exec() error
	}
)

func (p *Piston) Lang() error {
	log.Printf("fetching language options from Piston")
	client := new(http.Client)
	req, _ := http.NewRequest(http.MethodGet, "https://emkc.org/api/v1/piston/versions", nil)
	req.Header.Add("Content-Type", "application/json")
	resp, htErr := client.Do(req)

	if htErr != nil {
		return htErr
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("%#v", resp)
		return errors.New(resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&p.Versions); err != nil {
		return err
	}
	return nil
}

func (p *Piston) Exec() error {
	log.Printf("sending request to piston")
	client := new(http.Client)

	reqBody, err := json.Marshal(p.Execute)
	if err != nil {
		return err
	}

	req, reqErr := http.NewRequest(http.MethodPost, "https://emkc.org/api/v1/piston/execute", bytes.NewReader(reqBody))
	if reqErr != nil {
		return reqErr
	}

	req.Header.Add("Content-Type", "application/json")

	resp, htErr := client.Do(req)

	if htErr != nil {
		return htErr
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&p.Response); err != nil {
		return err
	}

	return nil
}
