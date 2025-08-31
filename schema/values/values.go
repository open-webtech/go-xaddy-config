package values

import (
	"reflect"
	"strings"
)

//go:generate go run ../../cmd/gen_values/main.go -pkg values

// Value is the interface to the dynamic value stored in an argument.
// It allows the contents of a Value to be set and retrieved.
// (The default value is represented as a string.)
type Value interface {
	// Set updates the value from a string representation
	Set(string) error
	// String returns a string representation of the value
	String() string
	// Get returns the underlying value
	Get() any
}

// Accumulator is a generic value collector that uses reflection to accumulate values into a slice.
type Accumulator struct {
	// element is a function that creates a Value for each element in the slice
	element func(value interface{}) Value
	// typ is the type of elements in the slice
	typ reflect.Type
	// slice is the reflect.Value of the target slice
	slice reflect.Value
}

// NewAccumulator creates a new Accumulator for collecting values into a slice.
// slice must be a pointer to a slice.
// element is a function that creates a Value for each element.
// Panics if slice is not a pointer to a slice.
func NewAccumulator(slice interface{}, element func(value interface{}) Value) *Accumulator {
	typ := reflect.TypeOf(slice)
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Slice {
		panic("expected a pointer to a slice")
	}
	return &Accumulator{
		element: element,
		typ:     typ.Elem().Elem(),
		slice:   reflect.ValueOf(slice),
	}
}

// String returns a comma-separated string of all accumulated values.
func (a *Accumulator) String() string {
	out := []string{}
	s := a.slice.Elem()
	for i := 0; i < s.Len(); i++ {
		out = append(out, a.element(s.Index(i).Addr().Interface()).String())
	}
	return strings.Join(out, ",")
}

// Set adds a single value to the accumulated slice.
func (a *Accumulator) Set(value string) error {
	e := reflect.New(a.typ)
	if err := a.element(e.Interface()).Set(value); err != nil {
		return err
	}
	slice := reflect.Append(a.slice.Elem(), e.Elem())
	a.slice.Elem().Set(slice)
	return nil
}

// Get returns the current slice of accumulated values.
func (a *Accumulator) Get() interface{} {
	return a.slice.Interface()
}

// IsCumulative indicates that this value type can accumulate multiple values.
func (a *Accumulator) IsCumulative() bool {
	return true
}

// SetList adds multiple values to the accumulated slice.
func SetList(a Value, values []string) error {
	for _, value := range values {
		if err := a.Set(value); err != nil {
			return err
		}
	}
	return nil
}
