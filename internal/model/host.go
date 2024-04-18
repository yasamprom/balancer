package model

// Host ...
type Host struct {
	Address string `json:"host"`
}

// MapPair contains map: server -> key range
type MapPair struct {
	Address  Host
	KeyRange Range
}

// HostState contains base information about server
type HostState struct {
	Host         Host
	Availability bool
	CPU          float64
	GPU          float64
	Network      float64
}
