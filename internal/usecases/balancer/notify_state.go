package balancer

import (
	"context"
	"log"

	"github.com/yasamprom/balancer/internal/model"
)

func (uc *Usecases) SendState(ctx context.Context) error {
	// to be collected from metrics
	states := model.HostState{}

	err := uc.slicerClient.NotifyState(ctx, states)
	if err != nil {
		log.Println("failed to notify states, %v", err)
		return err
	}
	return nil
}
