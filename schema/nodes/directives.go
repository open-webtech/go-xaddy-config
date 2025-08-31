package nodes

import (
	parser "github.com/foxcpp/maddy/framework/cfgparser"
	"github.com/open-webtech/go-xaddy-config/schema/args"
)

// DirectiveDef represents a single configuration directive definition.
type DirectiveDef struct {
	CommonDef
}

// NewDirectiveDef creates a new directive definition with the given name and optional argument definitions.
func NewDirectiveDef(name string, args ...*args.ArgDef) *DirectiveDef {
	d := &DirectiveDef{
		CommonDef: CommonDef{name: name},
	}
	d.addArgs(args...)
	return d
}

// AddArgs adds additional argument definitions to the directive.
// Returns the directive definition for method chaining.
func (d *DirectiveDef) AddArgs(args ...*args.ArgDef) *DirectiveDef {
	d.addArgs(args...)
	return d
}

// SetAttrs sets attributes for the directive definition.
// Returns the directive definition for method chaining.
func (d *DirectiveDef) SetAttrs(attributes ...NodeAttribute) *DirectiveDef {
	for _, attribute := range attributes {
		switch attribute {
		case Repeatable:
			d.repeatable = true
		}
	}
	return d
}

// SetHandler sets the handler function for the directive definition.
// Returns the directive definition for method chaining.
func (d *DirectiveDef) SetHandler(cb NodeHandler) *DirectiveDef {
	d.handler = cb
	return d
}

// Evaluate processes a directive node, ensuring it has no child nodes.
// Returns an error if the node is a block or if evaluation fails.
func (d *DirectiveDef) Evaluate(node parser.Node, cfg any) error {
	if len(node.Children) != 0 {
		return NodeErr(node, "node '%s' may not be a block", node.Name)
	}

	return evaluate(d, node)
}
