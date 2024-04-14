package balancer

type Usecases struct {
	slicerClient SlicerClient
}

func NewUsecases(c SlicerClient) *Usecases {
	return &Usecases{slicerClient: c}
}
