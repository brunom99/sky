package utils

import (
	"math/rand"
)

func RandInt(min int, max int, r ...*rand.Rand) int {
	if min == max {
		return min
	}
	if max-min <= 0 {
		return 0
	}
	if len(r) > 0 {
		return r[0].Intn(max-min) + min
	}
	return rand.Intn(max-min) + min
}

func RandFloat(min float64, max float64, r ...*rand.Rand) float64 {
	if len(r) > 0 {
		return min + r[0].Float64()*(max-min)
	}
	return min + rand.Float64()*(max-min)
}

func RandLatitudeLongitude(excludePoles bool, r ...*rand.Rand) (float64, float64) {
	if excludePoles {
		return RandFloat(-70, 70, r...), RandFloat(-180, 180, r...)
	}
	return RandFloat(-90, 90, r...), RandFloat(-180, 180, r...)
}
