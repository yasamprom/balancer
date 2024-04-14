package model

type Host struct {
	Address string `json:"host"`
}

type MapPair struct {
	Address  Host
	KeyRange Range
}
