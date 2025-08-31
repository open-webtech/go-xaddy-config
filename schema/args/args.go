package args

//go:generate go run ../../cmd/gen_values/main.go -pkg args

import (
	"github.com/open-webtech/go-xaddy-config/schema/values"
)

// ValueType represents the type of a configuration argument.
type ValueType int

// ArgAttribute defines special attributes that can be applied to configuration arguments.
type ArgAttribute int

const (
	// Optional indicates that the argument is not required
	Optional ArgAttribute = iota
)

// ArgsList is a slice of argument definitions.
type ArgsList []*ArgDef

// ArgDef represents a single configuration argument definition.
type ArgDef struct {
	// name of the argument
	name string
	// target is the value that will store the argument's parsed value
	target values.Value
	// valType defines the type of the argument
	valType ValueType
	// required indicates whether the argument must be provided
	required bool
	// variadic indicates if the argument can accept multiple values
	variadic bool
}

// NewArgDef creates a new argument definition.
// target is the value that will store the parsed argument.
// t is the type of the argument.
// attributes can modify the argument's behavior (e.g., Optional).
func NewArgDef(target values.Value, t ValueType, attributes ...ArgAttribute) *ArgDef {
	arg := &ArgDef{
		target:   target,
		valType:  t,
		required: true,
	}
	for _, attribute := range attributes {
		switch attribute {
		case Optional:
			arg.required = false
		}
	}
	return arg
}

// NewVariadicArgDef creates a new variadic argument definition.
// Similar to NewArgDef, but allows multiple values.
func NewVariadicArgDef(target values.Value, t ValueType, attributes ...ArgAttribute) *ArgDef {
	arg := NewArgDef(target, t, attributes...)
	arg.variadic = true
	return arg
}

// Name returns the name of the argument.
func (d *ArgDef) Name() string {
	return d.name
}

// Type returns the value type of the argument.
func (d *ArgDef) Type() ValueType {
	return d.valType
}

// Required returns whether the argument is required.
func (d *ArgDef) Required() bool {
	return d.required
}

// Variadic returns whether the argument can accept multiple values.
func (d *ArgDef) Variadic() bool {
	return d.variadic
}

// Target returns the value that will store the parsed argument.
func (d *ArgDef) Target() values.Value {
	return d.target
}
