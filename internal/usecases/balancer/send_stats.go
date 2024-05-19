package balancer

import (
	"context"
	"log"

	"github.com/yasamprom/balancer/internal/model"
)

func (uc *Usecases) SendStats(ctx context.Context, stats map[model.Range]int) error {
	err := uc.slicerClient.SendStats(ctx, stats)
	if err != nil {
		log.Printf("failed to send stats to Slicer, %v", err)
		return err
	}
	return nil
}
