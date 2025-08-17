package appname

import (
	"testing"
)

func TestToFormatAppName(t *testing.T) {
	formatName := ToFormatAppName("hidden-lake-kernel")
	if formatName != "Hidden Lake Kernel" {
		t.Log(formatName)
		t.Fatal("invalid full name")
	}
	formatName = ToFormatAppName("hidden-lake-adapters=common")
	if formatName != "Hidden Lake Adapters = Common" {
		t.Log(formatName)
		t.Fatal("invalid full name (with subname)")
	}
}

func TestToShortAppName(t *testing.T) {
	shortName := ToShortAppName("hidden-lake-kernel")
	if shortName != "HLK" {
		t.Log(shortName)
		t.Fatal("invalid short name")
	}
	shortName = ToShortAppName("hidden-lake-adapters=common")
	if shortName != "HLA=common" {
		t.Log(shortName)
		t.Fatal("invalid short name (with subname)")
	}
}
