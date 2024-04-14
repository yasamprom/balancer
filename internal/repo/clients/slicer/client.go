package slicer

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"log"

	model "github.com/yasamprom/balancer/internal/model"
)

type Config struct {
	Host string
	Port string
}

type Client struct {
	host string
	port string
}

func New(cfg Config) *Client {
	return &Client{host: cfg.Host, port: cfg.Port}
}

func (c *Client) GetMapping(ctx context.Context) (map[model.Range]model.Host, error) {

	resp, err := http.Get(fmt.Sprintf("http://%s:%s/get_mapping", c.host, c.port))
	if err != nil {
		log.Println("http client error: %v", err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	ranges := make([]model.MapPair, 0)
	err = json.Unmarshal(body, &ranges)
	if err != nil {
		log.Println(ctx, "failed to unmarshal response: %v", err)
	}

	log.Println(ctx, "got new ranges: %v", ranges)
	res := make(map[model.Range]model.Host)

	for _, v := range ranges {
		res[model.Range{From: v.KeyRange.From, To: v.KeyRange.To}] = v.Address
	}
	return res, nil
}
