package stage3

import "errors"

// ErrNegativeInput is a sentinel error: callers compare with errors.Is.
var ErrNegativeInput = errors.New("negative input")

// ErrKeyNotFound is a sentinel error returned (wrapped) by LoadConfig.
var ErrKeyNotFound = errors.New("key not found")

// SafeSqrt returns the square root of x. For negative x it returns
// (0, ErrNegativeInput).
//
// TODO: Return the sentinel error directly; no panics, no exceptions.
func SafeSqrt(x float64) (float64, error) {
	return 0, nil
}

// LoadConfig returns data[key]. When the key is missing it must return an
// error that wraps ErrKeyNotFound AND mentions the key in its message, e.g.
// `loadconfig "port": key not found`.
//
// TODO: Wrap with fmt.Errorf and the %w verb so errors.Is still works.
func LoadConfig(data map[string]string, key string) (string, error) {
	return "", nil
}

// ValidationError is a custom error type carrying structured data.
type ValidationError struct {
	Field   string
	Message string
}

// Error returns "<Field>: <Message>".
//
// TODO: Implement the error interface (Error() string).
func (e *ValidationError) Error() string {
	return ""
}

// ValidateAge returns nil when 0 <= age <= 150. Otherwise it returns a
// *ValidationError with Field "age".
//
// TODO: Return a *ValidationError. Beware: returning a nil *ValidationError
// through the error interface is NOT a nil error. Return a literal nil instead.
func ValidateAge(age int) error {
	return nil
}

// ParseAll converts every string to an int. Valid numbers are collected in the
// result. Invalid ones do not stop the loop: all failures are combined into a
// single returned error (nil when everything parsed). The returned error must
// satisfy errors.Is(err, strconv.ErrSyntax) when any input is malformed.
//
// TODO: Use strconv.Atoi and errors.Join.
func ParseAll(inputs []string) ([]int, error) {
	return nil, nil
}
