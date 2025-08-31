package nodes

import (
	parser "github.com/foxcpp/maddy/framework/cfgparser"
	"github.com/open-webtech/go-xaddy-config/schema/args"
)

// BlockDef represents a block-style configuration node that can contain child nodes.
type BlockDef struct {
	CommonDef
	NodesContainer
}

// NewBlockDef creates a new block definition with the given name and optional argument definitions.
func NewBlockDef(name string, args ...*args.ArgDef) *BlockDef {
	d := &BlockDef{
		CommonDef: CommonDef{name: name},
	}
	d.addArgs(args...)
	return d
}

// SetAttrs sets attributes for the block definition.
// It returns the block definition for method chaining.
func (d *BlockDef) SetAttrs(attributes ...NodeAttribute) *BlockDef {
	for _, attribute := range attributes {
		switch attribute {
		case Repeatable:
			d.repeatable = true
		}
	}
	return d
}

// SetHandler sets the handler function for the block definition.
// It returns the block definition for method chaining.
func (d *BlockDef) SetHandler(cb NodeHandler) *BlockDef {
	d.handler = cb
	return d
}

// Evaluate processes a block node and its children, updating the configuration.
func (d *BlockDef) Evaluate(node parser.Node, cfg any) error {
	if err := evaluate(d, node); err != nil {
		return err
	}

	return d.EvaluateTree(node.Children, cfg)
}

// ModuleBlockDef represents a block definition that can handle different module types.
type ModuleBlockDef struct {
	modules    map[string]*NodesContainer // map of module names to their node containers
	moduleName *string                    // name of the current module
	CommonDef
}

// NewModuleBlockDef creates a new module block definition with the given name.
func NewModuleBlockDef(name string) *ModuleBlockDef {
	d := &ModuleBlockDef{
		modules:   make(map[string]*NodesContainer),
		CommonDef: CommonDef{name: name},
	}
	d.addArgs(args.StringArg(d.moduleName))
	return d
}

// WithArgs adds additional argument definitions to the module block definition.
// It returns the module block definition for method chaining.
func (d *ModuleBlockDef) WithArgs(args ...*args.ArgDef) *ModuleBlockDef {
	d.addArgs(args...)
	return d
}

// Evaluate processes a module block node and its children, updating the configuration.
func (d *ModuleBlockDef) Evaluate(node parser.Node, cfg any) error {
	if err := evaluate(d, node); err != nil {
		return err
	}

	m, ok := d.modules[*d.moduleName]
	if !ok {
		return NodeErr(node, "unknown module '%s'", *d.moduleName)
	}

	return m.EvaluateTree(node.Children, cfg)
}
