package utils

import (
	"testing"
)

func TestGetFileInfoList(t *testing.T) {
	t.Parallel()

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

func TestFileNameIsInvalid(t *testing.T) {
	t.Parallel()

	if ok := FileNameIsInvalid(""); !ok {
		t.Fatal("file name valid with void string")
	}
	if ok := FileNameIsInvalid("\x01"); !ok {
		t.Fatal("file name valid with not graphic char")
	}
	if ok := FileNameIsInvalid("/abc/file.txt"); !ok {
		t.Fatal("file name valid with path exist")
	}
	if ok := FileNameIsInvalid("file.txt"); ok {
		t.Fatal("file name is invalid but should be ok")
	}
}
