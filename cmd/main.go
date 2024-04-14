package main

import (
	"context"
	"log"

	executor "github.com/yasamprom/balancer/internal/executor"
	slicer "github.com/yasamprom/balancer/internal/repo/clients/slicer"
	"github.com/yasamprom/balancer/internal/usecases/balancer"
)

func main() {
	log.Println("Start main...")

	slicerClient := slicer.New(
		slicer.Config{
			Host: "127.0.0.1", Port: "8000",
		},
	)
	uc := balancer.NewUsecases(slicerClient)
	executor := executor.New(&executor.Config{Usecases: *uc})
	executor.RunMappingManager(context.Background())
}
