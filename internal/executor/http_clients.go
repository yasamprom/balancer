package executor

import (
	"net/http"

	"github.com/yasamprom/balancer/internal/model"
)

type httpClientsConfig struct {
	hosts []model.Host
}

type httpClient struct {
	hosts  []model.Host
	client http.Client
}

func (c *httpClient) New(config httpClientsConfig) *httpClient {
	client := httpClient{
		hosts:  config.hosts,
		client: http.Client{},
	}
	return &client
}
