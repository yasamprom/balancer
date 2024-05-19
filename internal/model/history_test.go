package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_History(t *testing.T) {
	t.Run("Add few times", func(t *testing.T) {
		h := History{
			Ranges: make(map[Range][]time.Time),
		}
		for i := 0; i < 10; i++ {
			h.Add(Range{From: 1, To: 2}, time.Now())
		}
	})

	t.Run("Add few times and count", func(t *testing.T) {
		h := History{
			Ranges: make(map[Range][]time.Time),
		}
		times := make([]time.Time, 0)
		cur := time.Now()
		for i := 0; i < 10; i++ {
			times = append(times, cur.Add(-time.Duration(10*time.Minute-time.Duration(i)*time.Minute)))
		}
		for i := 0; i < len(times); i++ {
			h.Add(Range{From: 1, To: 2}, times[i])
		}
		res := h.CountPerPeriod(time.Minute + time.Second)

		expected := map[Range]int{
			{1, 2}: 1,
		}
		assert.Equal(t, expected, res, "count per period failed")
	})

	t.Run("Add same times and count", func(t *testing.T) {
		h := History{
			Ranges: make(map[Range][]time.Time),
		}
		times := make([]time.Time, 0)
		cur := time.Now()
		for i := 0; i < 10; i++ {
			times = append(times, cur.Add(-time.Duration(5*time.Minute)))
		}
		for i := 0; i < len(times); i++ {
			h.Add(Range{From: 1, To: 2}, times[i])
		}
		res := h.CountPerPeriod(time.Minute + time.Second)

		expected := map[Range]int{
			{1, 2}: 0,
		}
		assert.Equal(t, expected, res, "count per period failed")
	})

	t.Run("Add same times and count", func(t *testing.T) {
		h := History{
			Ranges: make(map[Range][]time.Time),
		}
		times := make([]time.Time, 0)
		cur := time.Now()
		for i := 0; i < 10; i++ {
			times = append(times, cur.Add(-time.Duration(5*time.Minute)))
		}
		for i := 0; i < len(times); i++ {
			h.Add(Range{From: 1, To: 2}, times[i])
		}
		res := h.CountPerPeriod(5*time.Minute + time.Second)

		expected := map[Range]int{
			{1, 2}: len(times),
		}
		assert.Equal(t, expected, res, "count per period failed")
	})

	t.Run("Empty events list", func(t *testing.T) {
		h := History{
			Ranges: map[Range][]time.Time{
				{1, 2}: nil,
			},
		}

		res := h.CountPerPeriod(time.Minute)

		expected := map[Range]int{
			{1, 2}: 0,
		}
		assert.Equal(t, expected, res, "count per period failed")
	})

}
