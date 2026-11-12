package stage1

import (
	"reflect"
	"testing"
)

func TestSwap(t *testing.T) {
	a, b := Swap(1, 2)
	if a != 2 || b != 1 {
		t.Fatalf("Swap(1, 2) = (%d, %d), want (2, 1)", a, b)
	}
}

func TestSum(t *testing.T) {
	if got := Sum(); got != 0 {
		t.Errorf("Sum() = %d, want 0", got)
	}
	if got := Sum(1, 2, 3); got != 6 {
		t.Errorf("Sum(1, 2, 3) = %d, want 6", got)
	}
	xs := []int{10, 20, 30, 40}
	if got := Sum(xs...); got != 100 {
		t.Errorf("Sum(xs...) = %d, want 100", got)
	}
}

func TestNewCounter(t *testing.T) {
	c1, c2 := NewCounter(), NewCounter()
	for want := 1; want <= 3; want++ {
		if got := c1(); got != want {
			t.Fatalf("c1() = %d, want %d", got, want)
		}
	}
	if got := c2(); got != 1 {
		t.Errorf("c2() = %d, want 1 (counters must be independent)", got)
	}
}

func TestWeekday(t *testing.T) {
	days := []Weekday{Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday}
	names := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	for i, d := range days {
		if int(d) != i {
			t.Errorf("%s has value %d, want %d", names[i], int(d), i)
		}
		if d.String() != names[i] {
			t.Errorf("Weekday(%d).String() = %q, want %q", i, d.String(), names[i])
		}
	}
	if got := Weekday(42).String(); got != "Unknown" {
		t.Errorf("Weekday(42).String() = %q, want %q", got, "Unknown")
	}
}

func TestDescribe(t *testing.T) {
	tests := []struct {
		in   any
		want string
	}{
		{42, "int:42"},
		{"hi", "string:hi"},
		{true, "bool:true"},
		{3.14, "unknown"},
		{nil, "unknown"},
	}
	for _, tc := range tests {
		if got := Describe(tc.in); got != tc.want {
			t.Errorf("Describe(%v) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestIncrement(t *testing.T) {
	n := 41
	Increment(&n)
	if n != 42 {
		t.Errorf("n = %d after Increment, want 42", n)
	}
}

func TestSafeDeref(t *testing.T) {
	if got := SafeDeref(nil, 7); got != 7 {
		t.Errorf("SafeDeref(nil, 7) = %d, want 7", got)
	}
	v := 5
	if got := SafeDeref(&v, 7); got != 5 {
		t.Errorf("SafeDeref(&5, 7) = %d, want 5", got)
	}
}

func TestPointAdd(t *testing.T) {
	p, q := Point{1, 2}, Point{3, 4}
	got := p.Add(q)
	if got != (Point{4, 6}) {
		t.Errorf("Add = %+v, want {4 6}", got)
	}
	if p != (Point{1, 2}) {
		t.Errorf("Add mutated its receiver: %+v", p)
	}
}

func TestPointScale(t *testing.T) {
	p := Point{1, 2}
	p.Scale(3)
	if p != (Point{3, 6}) {
		t.Errorf("after Scale(3) p = %+v, want {3 6}", p)
	}
}

func TestDeferOrder(t *testing.T) {
	want := []int{3, 2, 1}
	if got := DeferOrder(); !reflect.DeepEqual(got, want) {
		t.Errorf("DeferOrder() = %v, want %v", got, want)
	}
}
