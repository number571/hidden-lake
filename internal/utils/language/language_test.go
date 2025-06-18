package language

import "testing"

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SLanguageError{str}
	if err.Error() != errPrefix+str {
		t.Fatal("incorrect err.Error()")
	}
}

func TestPanicFromLanguage(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("nothing panics")
		}
	}()

	_ = FromILanguage(111)
}

func TestToLanguage(t *testing.T) {
	lang, err := ToILanguage("ENG")
	if err != nil {
		t.Fatal(err)
	}
	if lang != CLangENG {
		t.Fatal("got invalid ENG")
	}

	lang, err = ToILanguage("RUS")
	if err != nil {
		t.Fatal(err)
	}
	if lang != CLangRUS {
		t.Fatal("got invalid RUS")
	}

	lang, err = ToILanguage("ESP")
	if err != nil {
		t.Fatal(err)
	}
	if lang != CLangESP {
		t.Fatal("got invalid ESP")
	}

	if _, err := ToILanguage("???"); err == nil {
		t.Fatal("success unknown type to language")
	}
}

func TestFromLanguage(t *testing.T) {
	if FromILanguage(CLangENG) != "ENG" {
		t.Fatal("got invalid ENG")
	}
	if FromILanguage(CLangRUS) != "RUS" {
		t.Fatal("got invalid RUS")
	}
	if FromILanguage(CLangESP) != "ESP" {
		t.Fatal("got invalid ESP")
	}
}
