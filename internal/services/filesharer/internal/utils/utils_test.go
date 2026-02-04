package utils

import (
	"testing"
)

func TestGetFileInfoList(t *testing.T) {
	fileInfoList, err := GetFileInfoList("./testdata", 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	list := fileInfoList.GetList()
	if len(list) != 1 {
		t.Fatal("len(list) != 1")
	}
	if _, err := GetFileInfoList("./testdata/not_found", 1, 1); err != nil {
		t.Fatal(err) // return []
	}
}
