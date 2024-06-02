package executor

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/yasamprom/balancer/internal/model"
)

const (
	routingHeader = "x-routing-key"
	servePort     = "8080"
	serveHost     = "0.0.0.0"
	R             = 37
	M             = 1000000007
)

func (ex *Executor) StartHandle(ctx context.Context) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("got new request: %v, %v", r.Host, r.RequestURI)
		log.Print(r.Body)

		route := r.Header.Get(routingHeader)
		if route == "" {
			w.Write([]byte("{\"response\": 400, \"status\": \"No routing key provided.\"}"))
			return
		}
		var key uint64

		key, err := strconv.ParseUint(route, 10, 64)
		if err != nil {
			w.Write([]byte("{\"response\": 400, \"status\": \"Failed to parse routing key.\"}"))
			return
		}
		key = hashKey(key)

		// find host for redirecting query and update stats history
		keyRange, targetHost := ex.getHost(key)
		ex.history.Add(keyRange, time.Now())

		if targetHost.Address == "" {
			w.Write([]byte("{\"response\": 400, \"message\": \"No worker for provided key.\"}"))
			return
		}

		targetReq, err := buildTargetRequest(*r, targetHost.Address, servePort)
		if err != nil {
			log.Printf("failed to construct target request: %v", err)
			w.Write([]byte("{\"response\": 500, \"message\": \"Failed to redirect request.\"}"))
			return
		}
		resp, err := http.DefaultClient.Do(targetReq)
		if err != nil {
			w.Write([]byte("{\"response\": 500, \"message\": \"Failed to execute target request.\"}"))
			return
		}
		var bytes []byte
		_, err = resp.Body.Read(bytes)
		if err != nil {
			w.Write([]byte("{\"response\": 500, \"message\": \"Failed to read response.\"}"))
			return
		}
		w.Write([]byte(bytes))

	})

	log.Print("Starting serve requests...")
	http.ListenAndServe(fmt.Sprintf("%v:%v", serveHost, servePort), nil)
	return nil
}

func (ex *Executor) getHost(key uint64) (model.Range, model.Host) {
	log.Printf("Received key: %v. Current map: \n%v\n", key, ex.mapping.mp)
	for r, host := range ex.mapping.mp {
		if r.From <= key && key < r.To {
			return r, host
		}
	}
	return model.Range{}, model.Host{}
}

func buildTargetRequest(r http.Request, host, port string) (*http.Request, error) {
	fullUrl := fmt.Sprintf("http://%v:%v%v", host, port, r.RequestURI)
	req, err := http.NewRequest(r.Method, fullUrl, r.Body)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return &http.Request{}, nil
	}
	return req, nil
}

func hashKey(key uint64) uint64 {
	key = (^key) + (key << 21) // key = (key << 21) - key - 1
	key = key ^ (key >> 24)
	key = (key + (key << 3)) + (key << 8) // key * 265
	key = key ^ (key >> 14)
	key = (key + (key << 2)) + (key << 4) // key * 21
	key = key ^ (key >> 28)
	key = key + (key << 31)
	return key

}
