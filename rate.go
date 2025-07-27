package main

import (
	"math"
	"math/rand/v2"
	"strconv"
	"strings"
)

func getMultiplierByFloat(rate float64) int {
	numStr := strconv.FormatFloat(rate, 'f', -1, 64)
	if i := strings.Index(numStr, "."); i >= 0 {
		return int(math.Pow(10, float64(len(numStr)-i-1)))
	}
	return 0
}

func tryLuck(rate float64, multiplier int) bool {
	intRate := int(float64(multiplier) * rate)
	randChance := rand.IntN(multiplier)
	return randChance <= intRate
}
