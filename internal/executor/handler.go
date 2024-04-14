package executor

import (
	"context"
	"net/http"
)

func (ex *Executor) StartHandle(ctx context.Context) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// The "/" matches anything not handled elsewhere. If it's not the root
		// unpack x-routing-key
		// get host by mapping
		// redurect query
	})
	return nil
}
