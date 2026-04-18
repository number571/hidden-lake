package limiter

import "testing"

func TestLimiter(t *testing.T) {
	t.Parallel()

	const n = 16

	// generated 1/10 token in one second
	lm := NewLimitManager(.1, n)
	l := lm.Get("id")

	for range n {
		if ok := l.Allow(); !ok {
			t.Fatal("not allowed?")
		}
	}
	if ok := l.Allow(); ok {
		t.Fatal("allowed?")
	}
}
