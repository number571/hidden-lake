package dto

import "testing"

func TestFileInfoList(t *testing.T) {
	t.Parallel()

	f1, err := NewFileInfo("./testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	f2, err := NewFileInfo("./testdata/file.txt")
	if err != nil {
		t.Fatal(err)
	}

	f3, err := NewFileInfo("./testdata/something.txt")
	if err != nil {
		t.Fatal(err)
	}

	origFileInfoList := []IFileInfo{f1, f2, f3}
	fileInfoList, err := LoadFileInfoList(origFileInfoList)
	if err != nil {
		t.Fatal(err)
	}

	gotFileInfoList := fileInfoList.GetList()
	for i := range gotFileInfoList {
		if gotFileInfoList[i].GetHash() != origFileInfoList[i].GetHash() {
			t.Fatal("gotFileInfoList[i].GetHash() != origFileInfoList[i].GetHash()")
		}
	}

	if _, err := LoadFileInfoList(f1); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadFileInfoList(fileInfoList.ToBytes()); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadFileInfoList(fileInfoList.ToString()); err != nil {
		t.Fatal(err)
	}
}
