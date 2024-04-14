package pb

import "google.golang.org/grpc"

type SlicerClient struct {
	cc grpc.ClientConnInterface
}

func NewSlicerClient(c grpc.ClientConnInterface) SlicerClient {
	return SlicerClient{
		cc: c,
	}
}

func (c *SlicerClient) GetMapping() interface{} {
	return nil
}
