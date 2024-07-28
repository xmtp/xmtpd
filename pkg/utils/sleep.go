package utils

import (
	"math/rand"
	"time"
)

func RandomSleep(maxTimeMs int) {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(maxTimeMs)))
}
