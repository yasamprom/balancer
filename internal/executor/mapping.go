package executor

import (
	"sync"

	"github.com/yasamprom/balancer/internal/model"
)

type Mapping struct {
	mp map[model.Range]model.Host
	mu *sync.Mutex
}
