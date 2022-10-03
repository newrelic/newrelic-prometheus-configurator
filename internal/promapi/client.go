package promapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

const (
	metadataPath = "/api/v1/metadata"
	reloadPath   = "/-/reload"
)

var ErrMetadataEndpointFail = errors.New("checking metadata: status not successful")

type PromAPI interface {
	Metadata() (Metadata, error)
}

type Client struct {
	prometheusURL url.URL
}

func New(prometheusURL url.URL) Client {
	return Client{prometheusURL: prometheusURL}
}

func (c Client) Metadata() (Metadata, error) {
	// TODO
	resp, err := http.Get(c.prometheusURL.JoinPath(metadataPath).String())
	if err != nil {
		// TODO
		return Metadata{}, err
	}
	defer resp.Body.Close()

	metadata := Metadata{}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&metadata); err != nil {
		return Metadata{}, err
	}

	if metadata.Status != "success" {
		return Metadata{}, ErrMetadataEndpointFail
	}

	return metadata, nil
}

func (c Client) Reload() error {
	resp, err := http.Post(c.prometheusURL.JoinPath(reloadPath).String(), "", nil)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
