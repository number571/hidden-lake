package rand

import "testing"

func TestRandIntn(t *testing.T) {
	t.Parallel()

	for i := 0; i < 1000; i++ {
		if UniformIntn(10) >= 10 {
			t.Error("get invalid rand value")
			return
		}
	}
}
