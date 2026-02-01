package dto

import (
	"path/filepath"
	"testing"
)

const (
	testFileNotFound = "./testdata/random-932080293i029if94u0g3/example.txt"
	testFilePath     = "./testdata/example.txt"
	testFileHash     = "33d3864ac917b56d4d26f2c453f28db3fcbdaff5f56ac84cea971eb38e84931e9d9a6479d636cc19c04a0cb8c15ca2f9"
	testFileSize     = 26
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SError{str}
	if err.Error() != errPrefix+str {
		t.Fatal("incorrect err.Error()")
	}
}

func TestFileInfo(t *testing.T) {
	t.Parallel()

	fileInfo, err := NewFileInfo(testFilePath)
	if err != nil {
		t.Fatal(err)
	}

	if fileInfo.GetName() != filepath.Base(testFilePath) {
		t.Fatal("invalid file name")
	}
	if fileInfo.GetSize() != testFileSize {
		t.Fatal("invalid file size")
	}
	if fileInfo.GetHash() != testFileHash {
		t.Fatal("invalid file hash")
	}

	fileInfo2, err := LoadFileInfo(fileInfo.ToBytes())
	if err != nil {
		t.Fatal(err)
	}

	if fileInfo2.GetName() != filepath.Base(testFilePath) {
		t.Fatal("invalid file name (2)")
	}
	if fileInfo2.GetSize() != testFileSize {
		t.Fatal("invalid file size (2)")
	}
	if fileInfo2.GetHash() != testFileHash {
		t.Fatal("invalid file hash (2)")
	}

	fileInfo3, err := LoadFileInfo(fileInfo.ToString())
	if err != nil {
		t.Fatal(err)
	}

	if fileInfo3.GetName() != filepath.Base(testFilePath) {
		t.Fatal("invalid file name (3)")
	}
	if fileInfo3.GetSize() != testFileSize {
		t.Fatal("invalid file size (3)")
	}
	if fileInfo3.GetHash() != testFileHash {
		t.Fatal("invalid file hash (3)")
	}

	if _, err := LoadFileInfo(1); err == nil {
		t.Fatal("success load invalid type")
	}
	if _, err := LoadFileInfo([]byte{1}); err == nil {
		t.Fatal("success load invalid json")
	}

	fileInfoX, err := LoadFileInfo(fileInfo.ToString())
	if err != nil {
		t.Fatal(err)
	}
	fileInfoXRaw := fileInfoX.(*sFileInfo)

	fileInfoXRaw.FHash = "_"
	if _, err := LoadFileInfo(fileInfoX.ToBytes()); err == nil {
		t.Fatal("success load invalid hash value")
	}

	fileInfoXRaw.FHash = "FF"
	if _, err := LoadFileInfo(fileInfoX.ToBytes()); err == nil {
		t.Fatal("success load invalid hash size")
	}

	if _, err := NewFileInfo(testFileNotFound); err == nil {
		t.Fatal("success new not found file")
	}
}
