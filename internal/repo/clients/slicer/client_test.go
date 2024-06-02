package slicer

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yasamprom/balancer/internal/model"
)

func TestHelloName(t *testing.T) {
	body := "{\"RangeNodePairs\":[{\"Host\":\"10.244.0.29\",\"Range\":{\"From\":0,\"To\":4611686018427387903}},{\"Host\":\"10.244.0.29\",\"Range\":{\"From\":4611686018427387904,\"To\":9223372036854775807}}]}"
	res := model.RangesNodePairs{}

	_ = json.Unmarshal([]byte(body), &res)

	expected := model.RangesNodePairs{
		Pairs: []model.NodePair{
			{
				Host: "10.244.0.29",
				KeyRange: model.Range{
					From: uint64(0),
					To:   uint64(4611686018427387903),
				},
			},
			{
				Host: "10.244.0.29",
				KeyRange: model.Range{
					From: uint64(4611686018427387904),
					To:   uint64(9223372036854775807),
				},
			},
		},
	}
	assert.Equal(t, expected, res)
}
