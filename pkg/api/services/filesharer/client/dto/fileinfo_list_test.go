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
	fileInfoList, err := LoadFileInfoList(FileInfoListToString(origFileInfoList))
	if err != nil {
		t.Fatal(err)
	}

	for i := range fileInfoList {
		if fileInfoList[i].GetHash() != origFileInfoList[i].GetHash() {
			t.Fatal("gotFileInfoList[i].GetHash() != origFileInfoList[i].GetHash()")
		}
	}

	if _, err := LoadFileInfoList(111); err == nil {
		t.Fatal("success load file info with invalid type")
	}
	if _, err := LoadFileInfoList(`["name":"xxx.txt","size":10,"hash":"123"]`); err == nil {
		t.Fatal("success load file info with invalid hash")
	}
}
