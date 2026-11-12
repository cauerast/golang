package stage3

import (
	"errors"
	"math"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"testing/iotest"
)

// Compile-time checks that the types satisfy the interfaces.
var (
	_ Shape = Rectangle{}
	_ Shape = Circle{}
	_ error = (*ValidationError)(nil)
)

func near(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestRectangle(t *testing.T) {
	r := Rectangle{Width: 3, Height: 4}
	if !near(r.Area(), 12) {
		t.Errorf("Area = %v, want 12", r.Area())
	}
	if !near(r.Perimeter(), 14) {
		t.Errorf("Perimeter = %v, want 14", r.Perimeter())
	}
}

func TestCircle(t *testing.T) {
	c := Circle{Radius: 2}
	if !near(c.Area(), math.Pi*4) {
		t.Errorf("Area = %v, want %v", c.Area(), math.Pi*4)
	}
	if !near(c.Perimeter(), 4*math.Pi) {
		t.Errorf("Perimeter = %v, want %v", c.Perimeter(), 4*math.Pi)
	}
}

func TestTotalArea(t *testing.T) {
	shapes := []Shape{Rectangle{2, 3}, Circle{1}}
	if got, want := TotalArea(shapes), 6+math.Pi; !near(got, want) {
		t.Errorf("TotalArea = %v, want %v", got, want)
	}
	if got := TotalArea(nil); got != 0 {
		t.Errorf("TotalArea(nil) = %v, want 0", got)
	}
}

func TestLargest(t *testing.T) {
	if got := Largest(nil); got != nil {
		t.Errorf("Largest(nil) = %v, want nil", got)
	}
	big := Rectangle{10, 10}
	got := Largest([]Shape{Circle{1}, big, Rectangle{1, 2}})
	if got != Shape(big) {
		t.Errorf("Largest = %v, want %v", got, big)
	}
}

func TestCountWords(t *testing.T) {
	n, err := CountWords(strings.NewReader("the quick  brown\nfox\tjumps"))
	if err != nil || n != 5 {
		t.Errorf("CountWords = (%d, %v), want (5, nil)", n, err)
	}
	n, err = CountWords(strings.NewReader(""))
	if err != nil || n != 0 {
		t.Errorf("CountWords(empty) = (%d, %v), want (0, nil)", n, err)
	}
	boom := errors.New("boom")
	if _, err := CountWords(iotest.ErrReader(boom)); !errors.Is(err, boom) {
		t.Errorf("CountWords must return the reader's error, got %v", err)
	}
}

func TestSafeSqrt(t *testing.T) {
	got, err := SafeSqrt(9)
	if err != nil || !near(got, 3) {
		t.Errorf("SafeSqrt(9) = (%v, %v), want (3, nil)", got, err)
	}
	if _, err := SafeSqrt(-1); !errors.Is(err, ErrNegativeInput) {
		t.Errorf("SafeSqrt(-1) error = %v, want ErrNegativeInput", err)
	}
}

func TestLoadConfig(t *testing.T) {
	data := map[string]string{"host": "localhost"}
	v, err := LoadConfig(data, "host")
	if err != nil || v != "localhost" {
		t.Errorf("LoadConfig(host) = (%q, %v)", v, err)
	}
	_, err = LoadConfig(data, "port")
	if !errors.Is(err, ErrKeyNotFound) {
		t.Fatalf("error %v does not wrap ErrKeyNotFound", err)
	}
	if err == ErrKeyNotFound {
		t.Error("error must be wrapped with context, not the bare sentinel")
	}
	if !strings.Contains(err.Error(), "port") {
		t.Errorf("error message %q should mention the key", err.Error())
	}
}

func TestValidateAge(t *testing.T) {
	for _, ok := range []int{0, 30, 150} {
		if err := ValidateAge(ok); err != nil {
			t.Errorf("ValidateAge(%d) = %v, want nil", ok, err)
		}
	}
	for _, bad := range []int{-1, 151} {
		err := ValidateAge(bad)
		var ve *ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("ValidateAge(%d) = %v, want *ValidationError", bad, err)
		}
		if ve.Field != "age" {
			t.Errorf("Field = %q, want %q", ve.Field, "age")
		}
		if ve.Error() != ve.Field+": "+ve.Message || ve.Message == "" {
			t.Errorf("Error() = %q, want %q", ve.Error(), ve.Field+": "+ve.Message)
		}
	}
}

func TestParseAll(t *testing.T) {
	got, err := ParseAll([]string{"1", "2", "3"})
	if err != nil || !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Errorf("ParseAll(valid) = (%v, %v)", got, err)
	}
	got, err = ParseAll([]string{"1", "x", "3", "y"})
	if !reflect.DeepEqual(got, []int{1, 3}) {
		t.Errorf("ParseAll valid subset = %v, want [1 3]", got)
	}
	if !errors.Is(err, strconv.ErrSyntax) {
		t.Fatalf("error %v should satisfy errors.Is(err, strconv.ErrSyntax)", err)
	}
	if !strings.Contains(err.Error(), "x") || !strings.Contains(err.Error(), "y") {
		t.Errorf("joined error %q should report both bad inputs", err.Error())
	}
}
