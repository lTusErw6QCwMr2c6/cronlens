package history

import (
	"testing"
	"time"
)

func TestAdd_NewEntry(t *testing.T) {
	h := New(10)
	h.Add("* * * * *")

	if h.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", h.Len())
	}

	list := h.List()
	if list[0].Expression != "* * * * *" {
		t.Errorf("unexpected expression: %s", list[0].Expression)
	}
	if list[0].UseCount != 1 {
		t.Errorf("expected UseCount 1, got %d", list[0].UseCount)
	}
}

func TestAdd_DuplicateIncrements(t *testing.T) {
	h := New(10)
	h.Add("0 9 * * 1")
	h.Add("0 9 * * 1")
	h.Add("0 9 * * 1")

	if h.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", h.Len())
	}
	if h.List()[0].UseCount != 3 {
		t.Errorf("expected UseCount 3, got %d", h.List()[0].UseCount)
	}
}

func TestAdd_EvictsOldestWhenFull(t *testing.T) {
	h := New(3)
	h.Add("expr1")
	h.Add("expr2")
	h.Add("expr3")
	h.Add("expr4") // should evict expr1

	if h.Len() != 3 {
		t.Fatalf("expected 3 entries, got %d", h.Len())
	}

	for _, e := range h.List() {
		if e.Expression == "expr1" {
			t.Error("expr1 should have been evicted")
		}
	}
}

func TestList_MostRecentFirst(t *testing.T) {
	h := New(10)
	h.Add("first")
	time.Sleep(time.Millisecond)
	h.Add("second")
	time.Sleep(time.Millisecond)
	h.Add("third")

	list := h.List()
	if list[0].Expression != "third" {
		t.Errorf("expected 'third' first, got %s", list[0].Expression)
	}
	if list[2].Expression != "first" {
		t.Errorf("expected 'first' last, got %s", list[2].Expression)
	}
}

func TestClear(t *testing.T) {
	h := New(10)
	h.Add("* * * * *")
	h.Add("0 0 * * *")
	h.Clear()

	if h.Len() != 0 {
		t.Errorf("expected 0 entries after clear, got %d", h.Len())
	}
}

func TestNew_DefaultMaxSize(t *testing.T) {
	h := New(0)
	for i := 0; i < 25; i++ {
		h.Add("expr")
	}
	// duplicate, so still 1 entry
	if h.Len() != 1 {
		t.Errorf("expected 1 entry, got %d", h.Len())
	}
}
