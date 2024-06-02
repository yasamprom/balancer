package balancer

import (
	"context"
	"log"

	"github.com/yasamprom/balancer/internal/model"
)

func (uc *Usecases) UpdateRanges(ctx context.Context) (map[model.Range]model.Host, error) {
	ranges, err := uc.slicerClient.GetMapping(ctx)
	log.Printf("Mapping: \n%v\n", ranges)
	if err != nil {
		log.Printf("failed to get new ranges, %v", err)
		return nil, err
	}
	return ranges, nil
}
