package executor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/yasamprom/balancer/internal/model"
	balancer "github.com/yasamprom/balancer/internal/usecases/balancer"
)

type Config struct {
	Usecases balancer.Usecases // rename
}

type Executor struct {
	uc         balancer.Usecases // rename
	mapping    Mapping
	httpClient httpClient
	// ...
}

func New(config *Config) *Executor {
	return &Executor{
		uc: config.Usecases,
		mapping: Mapping{
			mp: make(map[model.Range]model.Host),
			mu: &sync.Mutex{},
		},
		httpClient: httpClient{},
	}
}

func (ex *Executor) RunMappingManager(ctx context.Context) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {

		ctx = context.Background()
		ticker := time.NewTicker(time.Second)
		done := make(chan bool)
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				ranges, err := ex.uc.UpdateRanges(ctx)
				if err != nil {
					fmt.Print("err")
					// log.Errorf(ctx, "mappingManager: failed to get ranges: %v", err)
				}
				ex.mapping.mu.Lock()
				ex.mapping.mp = ranges
				ex.mapping.mu.Unlock()
			}
		}
	}()
	wg.Wait()
}
