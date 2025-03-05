package server

import (
	"math"
)

func meanTime(times []int64) float64 {
	if len(times) == 0 {
		return 0
	}

	sum := float64(0)
	for _, t := range times {
		sum += float64(t)
	}

	return sum / float64(len(times))
}

func stdDev(times []int64, meanInterval float64) float64 {
	if len(times) <= 1 {
		return 0
	}

	numIntervals := float64(len(times))
	sum := float64(0)

	for _, t := range times {
		sum += math.Pow(float64(t)-meanInterval, 2)
	}

	return math.Sqrt(sum / (numIntervals - 1))
}
