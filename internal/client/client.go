package client

import (
	"mangahub/pkg/utils/config"
	"net/http"
)

type Client struct {
	Config     *config.ClientConfig
	HttpClient *http.Client
}

func NewClient(config *config.ClientConfig) *Client {
	return &Client{
		Config:     config,
		HttpClient: http.DefaultClient,
	}
}
