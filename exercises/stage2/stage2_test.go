package stage2

import (
	"reflect"
	"testing"
)

func TestModifyArray(t *testing.T) {
	orig := [3]int{1, 2, 3}
	got := ModifyArray(orig)
	if got != [3]int{100, 2, 3} {
		t.Errorf("ModifyArray = %v, want [100 2 3]", got)
	}
	if orig != [3]int{1, 2, 3} {
		t.Errorf("original array changed: %v", orig)
	}
}

func TestCloneAndModify(t *testing.T) {
	orig := []int{1, 2, 3}
	got := CloneAndModify(orig)
	if !reflect.DeepEqual(got, []int{99, 2, 3}) {
		t.Errorf("CloneAndModify = %v, want [99 2 3]", got)
	}
	if !reflect.DeepEqual(orig, []int{1, 2, 3}) {
		t.Errorf("original slice changed: %v", orig)
	}
	if got := CloneAndModify(nil); len(got) != 0 {
		t.Errorf("CloneAndModify(nil) = %v, want empty", got)
	}
}

func TestRemoveAt(t *testing.T) {
	tests := []struct {
		in   []int
		i    int
		want []int
	}{
		{[]int{1, 2, 3, 4}, 0, []int{2, 3, 4}},
		{[]int{1, 2, 3, 4}, 2, []int{1, 2, 4}},
		{[]int{1, 2, 3, 4}, 3, []int{1, 2, 3}},
		{[]int{1, 2, 3}, 5, []int{1, 2, 3}},
		{[]int{1, 2, 3}, -1, []int{1, 2, 3}},
	}
	for _, tc := range tests {
		in := append([]int(nil), tc.in...)
		if got := RemoveAt(in, tc.i); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("RemoveAt(%v, %d) = %v, want %v", tc.in, tc.i, got, tc.want)
		}
	}
}

func TestReverse(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	Reverse(s)
	if !reflect.DeepEqual(s, []int{5, 4, 3, 2, 1}) {
		t.Errorf("Reverse = %v", s)
	}
	e := []int{}
	Reverse(e) // must not panic
}

func TestFilter(t *testing.T) {
	even := func(n int) bool { return n%2 == 0 }
	got := Filter([]int{1, 2, 3, 4, 5, 6}, even)
	if !reflect.DeepEqual(got, []int{2, 4, 6}) {
		t.Errorf("Filter = %v, want [2 4 6]", got)
	}
	none := Filter([]int{1, 3}, even)
	if none == nil || len(none) != 0 {
		t.Errorf("Filter with no matches = %#v, want empty non-nil slice", none)
	}
}

func TestGrowthCaps(t *testing.T) {
	caps := GrowthCaps(100)
	if len(caps) != 100 {
		t.Fatalf("len(GrowthCaps(100)) = %d, want 100", len(caps))
	}
	changes := 0
	for i, c := range caps {
		if c < i+1 {
			t.Fatalf("cap %d after %d appends is smaller than length", c, i+1)
		}
		if i > 0 && c < caps[i-1] {
			t.Fatalf("capacity shrank at append %d", i+1)
		}
		if i > 0 && c != caps[i-1] {
			changes++
		}
	}
	if changes >= 50 {
		t.Errorf("capacity changed %d times; append should grow geometrically", changes)
	}
}

func TestSharedTail(t *testing.T) {
	want := []int{1, 2, 3, 100, 5}
	if got := SharedTail(); !reflect.DeepEqual(got, want) {
		t.Errorf("SharedTail() = %v, want %v", got, want)
	}
}

func TestWordFrequency(t *testing.T) {
	got := WordFrequency("Go go GO is fun  fun")
	want := map[string]int{"go": 3, "is": 1, "fun": 2}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("WordFrequency = %v, want %v", got, want)
	}
	if e := WordFrequency(""); e == nil || len(e) != 0 {
		t.Errorf("WordFrequency(\"\") = %#v, want empty non-nil map", e)
	}
}

func TestGroupByLength(t *testing.T) {
	got := GroupByLength([]string{"a", "bb", "c", "dd", "eee"})
	want := map[int][]string{1: {"a", "c"}, 2: {"bb", "dd"}, 3: {"eee"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GroupByLength = %v, want %v", got, want)
	}
}

func TestSortedKeys(t *testing.T) {
	got := SortedKeys(map[string]int{"b": 1, "c": 2, "a": 3})
	if !reflect.DeepEqual(got, []string{"a", "b", "c"}) {
		t.Errorf("SortedKeys = %v, want [a b c]", got)
	}
}

func TestLookup(t *testing.T) {
	m := map[string]int{"zero": 0, "one": 1}
	if v, ok := Lookup(m, "zero"); v != 0 || !ok {
		t.Errorf(`Lookup("zero") = (%d, %v), want (0, true)`, v, ok)
	}
	if v, ok := Lookup(m, "one"); v != 1 || !ok {
		t.Errorf(`Lookup("one") = (%d, %v), want (1, true)`, v, ok)
	}
	if v, ok := Lookup(m, "nope"); v != -1 || ok {
		t.Errorf(`Lookup("nope") = (%d, %v), want (-1, false)`, v, ok)
	}
}

func TestPersonAndEmbedding(t *testing.T) {
	e := Employee{Person: Person{Name: "Ana"}, Role: "Dev"}
	if got := e.Greet(); got != "Hi, I'm Ana" {
		t.Errorf("promoted Greet = %q", got)
	}
	if got := e.Describe(); got != "Ana works as Dev" {
		t.Errorf("Describe = %q", got)
	}
	m := Manager{Employee: e, Reports: 4}
	if got := m.Greet(); got != "Hi, I'm Ana, I manage 4 people" {
		t.Errorf("Manager.Greet = %q", got)
	}
	if got := m.Describe(); got != "Ana works as Dev" {
		t.Errorf("promoted Describe through Manager = %q", got)
	}
}

func TestStack(t *testing.T) {
	var s Stack
	if _, ok := s.Pop(); ok {
		t.Error("Pop on empty stack should return ok == false")
	}
	s.Push(1)
	s.Push(2)
	s.Push(3)
	if s.Len() != 3 {
		t.Fatalf("Len = %d, want 3", s.Len())
	}
	for _, want := range []int{3, 2, 1} {
		got, ok := s.Pop()
		if !ok || got != want {
			t.Fatalf("Pop = (%d, %v), want (%d, true)", got, ok, want)
		}
	}
	if s.Len() != 0 {
		t.Errorf("Len = %d after popping everything, want 0", s.Len())
	}
}
