package utils

import (
	"testing"
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SUtilsError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestReplaceTextToEmoji(t *testing.T) {
	t.Parallel()

	s := "hello :) !"
	got := ReplaceTextToEmoji(s)
	want := "hello 🙂 !"

	if got != want {
		t.Error("got incorrect replace text to emoji")
		return
	}
}

func TestReplaceTextToURLs(t *testing.T) {
	t.Parallel()

	s := "hello https://github.com/number571/github !"
	got := ReplaceTextToURLs(s)
	want := "hello <a style='background-color:#b9cdcf;color:black;' target='_blank' href='https://github.com/number571/github'>https://github.com/number571/github</a> !"

	if got != want {
		t.Error("got incorrect replace text to urls")
		return
	}
}
