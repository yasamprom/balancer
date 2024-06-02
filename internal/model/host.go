package model

// Host ...
type Host struct {
	Address string `json:"host"`
}

// MapPair contains map: server -> key range
type RangesNodePairs struct {
	Pairs []NodePair `json:"RangeNodePairs"`
}

type NodePair struct {
	Host     string `json:"Host"`
	KeyRange Range  `json:"Range"`
}

// HostState contains base information about server
type HostState struct {
	Host         Host
	Availability bool
	CPU          float64
	GPU          float64
	Network      float64
}
