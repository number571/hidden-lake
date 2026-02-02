package utils

import "testing"

func TestNothing(t *testing.T) {
	fileInfoList, err := GetFileInfoList("./testdata", 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	list := fileInfoList.GetList()
	if len(list) != 1 {
		t.Fatal("len(list) != 1")
	}
}
