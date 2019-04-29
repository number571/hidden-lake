package utils

import (
	"time"
    "math/rand"
)

func Shuffle(slice []string) {
    rand.Seed(int64(time.Now().Nanosecond()))
    for i := len(slice)-1; i > 0; i-- {
        j := rand.Intn(i+1)
        slice[i], slice[j] = slice[j], slice[i]
    }
}
