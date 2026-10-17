package randutil

import (
	"time"

	"golang.org/x/exp/rand"
)

func ShuffleRandom[T any](slice []T) {
	rand.New(rand.NewSource(uint64(time.Now().UnixNano()))).Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}

func Shuffle[T any](seed int64, slice []T) {
	rand.New(rand.NewSource(uint64(seed))).Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}
