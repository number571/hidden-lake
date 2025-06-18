package chars

import "testing"

func TestHasNotGraphicCharacters(t *testing.T) {
	t.Parallel()

	if HasNotGraphicCharacters("hello, world!") {
		t.Fatal("message contains only graphic chars")
	}

	if !HasNotGraphicCharacters("hello,\nworld!") {
		t.Fatal("message contains not graphic chars")
	}
}
