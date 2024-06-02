package slicer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

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

	resp, err := http.Get(fmt.Sprintf("http://%s:%s/api/v1/get_mapping", c.host, c.port))
	if err != nil {
		log.Printf("http client error: %v", err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	ranges := model.RangesNodePairs{}
	err = json.Unmarshal(body, &ranges)
	if err != nil {
		log.Println(ctx, "failed to unmarshal response: %v", err)
	}

	log.Printf("got new ranges: %v", ranges)
	res := make(map[model.Range]model.Host)

	for _, v := range ranges.Pairs {
		res[model.Range{From: v.KeyRange.From, To: v.KeyRange.To}] = model.Host{Address: v.Host}
	}
	log.Printf("Got %v ranges\n", len(res))
	log.Printf("Unmarshaled ranges:\n\n %v\n\n", res)
	return res, nil
}

func (c *Client) NotifyState(ctx context.Context, state model.HostState) error {
	jsonBody, err := json.Marshal(state)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(jsonBody)

	requestURL := fmt.Sprintf("http://%s:%s/update_state", c.host, c.port)
	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		log.Printf("http client error: %v", err)
		return err
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	_, err = client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return err
	}
	return nil
}

func (c *Client) SendStats(ctx context.Context, stats map[model.Range]int) error {
	jsonBody, err := json.Marshal(stats)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(jsonBody)

	requestURL := fmt.Sprintf("http://%s:%s/send_stats", c.host, c.port)
	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		log.Printf("http client error: %v", err)
		return err
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	_, err = client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return err
	}
	return nil
}
