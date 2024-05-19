package balancer

import (
	"context"

	"github.com/yasamprom/balancer/internal/model"
)

type SlicerClient interface {
	GetMapping(ctx context.Context) (map[model.Range]model.Host, error)
	NotifyState(ctx context.Context, state model.HostState) error
	SendStats(ctx context.Context, stats map[model.Range]int) error
}
