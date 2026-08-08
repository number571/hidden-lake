package alg

import "testing"

func TestNothing(t *testing.T) {
	t.Parallel()

	if GetAnonymityType("111") != CUnknownAlg {
		t.Fatal("got != unknown_alg")
	}
	if GetAnonymityType("") != CQBPAlg || GetAnonymityType("qbp") != CQBPAlg {
		t.Fatal("got != qbp")
	}
	if GetAnonymityType("silent") != CSilentAlg {
		t.Fatal("got != silent")
	}

	if CUnknownAlg.String() != "<nil>" { // nolint: goconst
		t.Fatal("unknown.string() != <nil>")
	}
	if CQBPAlg.String() != "qbp" { // nolint: goconst
		t.Fatal("qbp.string() != qbp")
	}
	if CSilentAlg.String() != "silent" { // nolint: goconst
		t.Fatal("silent.string() != silent")
	}
}
