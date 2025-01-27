package rand

import (
	"math"

	"github.com/number571/go-peer/pkg/crypto/random"
)

// Uniform random uint64 in [0;n)
func UniformIntn(n uint64) uint64 {
	random := random.NewRandom()
	pow2 := uint64(math.Pow(2, math.Ceil(math.Log2(float64(n)))))
	for {
		u64 := random.GetUint64() % pow2
		if u64 >= n {
			continue
		}
		return u64
	}
}
