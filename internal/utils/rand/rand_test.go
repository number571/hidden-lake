package rand

import "testing"

func TestRandIntn(t *testing.T) {
	t.Parallel()

	for i := 0; i < 1000; i++ {
		if UniformUint64n(10) >= 10 {
			t.Fatal("get invalid rand value")
		}
	}
}
