package executor

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/yasamprom/balancer/internal/model"
)

const (
	routingHeader = "x-routing-key"
	servePort     = "8080"
	serveHost     = "0.0.0.0"
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

		targetHost := ex.getHost(key)
		if targetHost.Address == "" {
			w.Write([]byte("{\"response\": 400, \"status\": \"No worker for provided key.\"}"))
			return
		}

		targetReq, err := buildTargetRequest(*r, targetHost.Address, servePort)
		if err != nil {
			log.Printf("failed to construct target request: %v", err)
			w.Write([]byte("{\"response\": 500, \"status\": \"Failed to redirect request.\"}"))
			return
		}
		resp, err := http.DefaultClient.Do(targetReq)
		if err != nil {
			w.Write([]byte("{\"response\": 500, \"status\": \"Failed to execute target request.\"}"))
			return
		}
		var bytes []byte
		_, err = resp.Body.Read(bytes)
		if err != nil {
			w.Write([]byte("{\"response\": 500, \"status\": \"Failed to read response.\"}"))
			return
		}
		w.Write([]byte("{\"response\": 200, \"status\": \"done.\"}"))

	})

	log.Print("Starting serve requests...")
	http.ListenAndServe(fmt.Sprintf("%v:%v", serveHost, servePort), nil)
	return nil
}

func (ex *Executor) getHost(key uint64) model.Host {
	for r, host := range ex.mapping.mp {
		if r.From <= key && key <= r.To {
			return host
		}
	}
	return model.Host{}
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
