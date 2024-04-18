package main

import (
	"context"
	"log"
	"time"

	executor "github.com/yasamprom/balancer/internal/executor"
	slicer "github.com/yasamprom/balancer/internal/repo/clients/slicer"
	"github.com/yasamprom/balancer/internal/usecases/balancer"
)

func main() {
	log.Println("Start main...")
	appCtx := context.Background()
	slicerClient := slicer.New(
		slicer.Config{
			Host: "127.0.0.1", Port: "8000",
		},
	)
	uc := balancer.NewUsecases(slicerClient)
	tickers := executor.TickersConfig{
		MappingTick: 5 * time.Second,
		StateTick:   10 * time.Second,
	}
	executor := executor.New(&executor.Config{
		Usecases: *uc,
		Tickers:  tickers,
	})
	executor.RunMappingManager(appCtx)
	executor.RunStateManager(appCtx)
}
