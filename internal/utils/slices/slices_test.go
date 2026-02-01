package slices

import "testing"

func TestAppendIfNotExist(t *testing.T) {
	t.Parallel()

	list := []string{"1", "2", "3"}
	list = AppendIfNotExist(list, "2")
	if len(list) != 3 {
		t.Fatal("len(list) != 3")
	}

	list = AppendIfNotExist(list, "0")
	if len(list) != 4 {
		t.Fatal("len(list) != 4")
	}
}

func TestDeleteValues(t *testing.T) {
	t.Parallel()

	list := []string{"1", "2", "3"}
	list = DeleteValues(list, "2")
	if len(list) != 2 {
		t.Fatal("len(list) != 2 (1)")
	}

	list = DeleteValues(list, "0")
	if len(list) != 2 {
		t.Fatal("len(list) != 2 (2)")
	}
}

func TestDeleteDuplicates(t *testing.T) {
	t.Parallel()

	list := []string{"1", "2", "3"}
	list = DeleteDuplicates(list)
	if len(list) != 3 {
		t.Fatal("len(list) != 3 (1)")
	}

	list = append(list, "2")
	list = DeleteDuplicates(list)
	if len(list) != 3 {
		t.Fatal("len(list) != 3 (2)")
	}
}
