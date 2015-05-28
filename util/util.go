package util

import (
	"math"
)

func ConvertMib2Gib(size int) int {
	sizeGb := float64(size) / float64(1024)
	return int(math.Ceil(sizeGb))
}
