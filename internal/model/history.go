package model

import (
	"time"
)

// History represents queries for each range
type History struct {
	Ranges map[Range][]time.Time
}

func (h *History) Add(r Range, t time.Time) {
	h.Ranges[r] = append(h.Ranges[r], t)
}

func (h *History) CountPerPeriod(period time.Duration) map[Range]int {
	res := make(map[Range]int)

	now := time.Now()
	for r, history := range h.Ranges {
		lb := LowerBound(history, now.Add(-period))
		res[r] = len(history) - lb
		h.Ranges[r] = h.Ranges[r][lb:]
	}
	return res
}

func LowerBound(array []time.Time, target time.Time) int {
	low, high, mid := 0, len(array)-1, 0
	for low <= high {
		mid = (low + high) / 2
		if array[mid].After(target) {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return low
}
