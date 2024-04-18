package executor

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/yasamprom/balancer/internal/model"
	balancer "github.com/yasamprom/balancer/internal/usecases/balancer"
)

// Config for initializing Executor
type Config struct {
	Usecases balancer.Usecases
	Tickers  TickersConfig
}

// TickersConfig for flexible tickers settings
type TickersConfig struct {
	MappingTick time.Duration
	StateTick   time.Duration
}

// Executor is main struct for implementing balancer life cycle
type Executor struct {
	uc         balancer.Usecases
	mapping    Mapping
	httpClient httpClient
	tickers    TickersConfig
}

func New(config *Config) *Executor {
	return &Executor{
		uc: config.Usecases,
		mapping: Mapping{
			mp: make(map[model.Range]model.Host),
			mu: &sync.Mutex{},
		},
		httpClient: httpClient{},
		tickers:    config.Tickers,
	}
}

// RunMappingManager starts manager for updating ranges
func (ex Executor) RunMappingManager(ctx context.Context) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {

		ctx = context.Background()
		ticker := time.NewTicker(ex.tickers.MappingTick)
		done := make(chan bool)
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				ranges, err := ex.uc.UpdateRanges(ctx)
				if err != nil {
					log.Println(ctx, "failed to update key ranges: ", err)
				}
				ex.mapping.mu.Lock()
				ex.mapping.mp = ranges
				ex.mapping.mu.Unlock()
			}
		}
	}()
	wg.Wait()
}

func (ex Executor) RunStateManager(ctx context.Context) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {

		ctx = context.Background()
		ticker := time.NewTicker(ex.tickers.StateTick)
		done := make(chan bool)
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				err := ex.uc.SendState(ctx)
				if err != nil {
					log.Println(ctx, "failed to notify state: ", err)
				}

			}
		}
	}()
	wg.Wait()
}
