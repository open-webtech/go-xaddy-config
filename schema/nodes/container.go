package nodes

import (
	parser "github.com/foxcpp/maddy/framework/cfgparser"
	"github.com/open-webtech/go-xaddy-config/schema/args"
)

// NodesContainer represents a collection of configuration directives and blocks
type NodesContainer struct {
	// Directives is a list of defined directive configurations
	Directives []*DirectiveDef
	// Blocks is a list of defined block configurations
	Blocks []*BlockDef
}

// AddBlock adds one or more block definitions to the container.
// Returns the container for method chaining.
func (nc *NodesContainer) AddBlock(blocks ...*BlockDef) *NodesContainer {
	nc.Blocks = append(nc.Blocks, blocks...)
	return nc
}

// DefineBlock creates a new block definition and adds it to the container.
// Returns the newly created block definition.
func (nc *NodesContainer) DefineBlock(name string, args ...*args.ArgDef) *BlockDef {
	b := NewBlockDef(name, args...)
	nc.AddBlock(b)
	return b
}

// DefineBlockCallback creates a block definition with a custom handler and adds it to the container.
// Returns the newly created block definition.
func (nc *NodesContainer) DefineBlockCallback(name string, cb NodeHandler) *BlockDef {
	b := NewBlockDef(name)
	b.SetHandler(cb)
	b.maxArgs = -1
	nc.AddBlock(b)
	return b
}

// AddDirective adds one or more directive definitions to the container.
// Returns the container for method chaining.
func (nc *NodesContainer) AddDirective(directives ...*DirectiveDef) *NodesContainer {
	nc.Directives = append(nc.Directives, directives...)
	return nc
}

// DefineDirective creates a new directive definition and adds it to the container.
// Returns the newly created directive definition.
func (nc *NodesContainer) DefineDirective(name string, args ...*args.ArgDef) *DirectiveDef {
	d := NewDirectiveDef(name, args...)
	nc.AddDirective(d)
	return d
}

// DefineDirectiveCallback creates a directive definition with a custom handler and adds it to the container.
// Returns the newly created directive definition.
func (nc *NodesContainer) DefineDirectiveCallback(name string, cb NodeHandler) *DirectiveDef {
	d := NewDirectiveDef(name)
	d.SetHandler(cb)
	d.maxArgs = -1
	nc.AddDirective(d)
	return d
}

// EvaluateTree evaluates and validates the configuration tree.
// Returns an error if the evaluation fails.
func (nc *NodesContainer) EvaluateTree(nodes []parser.Node, cfg any) error {
	var usedDirectives = make(map[string]bool)
	var usedBlocks = make(map[string]bool)

	for _, node := range nodes {
		for _, def := range nc.Directives {
			if node.Name == def.Name() {
				if !def.Repeatable() && usedDirectives[node.Name] {
					return NodeErr(node, "directive '%s' may not be repeated", node.Name)
				}

				usedDirectives[node.Name] = true
				if err := def.Evaluate(node, cfg); err != nil {
					return err
				}
			}
		}
		for _, def := range nc.Blocks {
			if node.Name == def.Name() {
				if !def.Repeatable() && usedBlocks[node.Name] {
					return NodeErr(node, "block '%s' may not be repeated", node.Name)
				}

				usedBlocks[node.Name] = true
				if err := def.Evaluate(node, cfg); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
