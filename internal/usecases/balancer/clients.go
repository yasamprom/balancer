package balancer

import (
	"context"

	"github.com/yasamprom/balancer/internal/model"
)

type SlicerClient interface {
	GetMapping(ctx context.Context) (map[model.Range]model.Host, error)
}
