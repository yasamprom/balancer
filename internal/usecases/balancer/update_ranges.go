package balancer

import (
	"context"
	"log"

	"github.com/yasamprom/balancer/internal/model"
)

func (uc *Usecases) UpdateRanges(ctx context.Context) (map[model.Range]model.Host, error) {
	ranges, err := uc.slicerClient.GetMapping(ctx)
	if err != nil {
		log.Println("failed to get new ranges, %v", err)
		return nil, err
	}
	return ranges, nil
}
