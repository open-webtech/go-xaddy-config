package nodes

import (
	"fmt"

	parser "github.com/foxcpp/maddy/framework/cfgparser"
	"github.com/open-webtech/go-xaddy-config/schema/args"
	"github.com/open-webtech/go-xaddy-config/schema/values"
)

// NodeDefinition defines the interface for configuration node definitions
// It provides methods to access node properties such as name, arguments, and handler
type NodeDefinition interface {
	// Name returns the name of the node
	Name() string
	// Args returns the argument definitions for the node
	Args() []*args.ArgDef
	// MinArgs returns the minimum number of arguments required
	MinArgs() int
	// MaxArgs returns the maximum number of arguments allowed
	MaxArgs() int
	// Handler returns the function that handles this node
	Handler() NodeHandler
	// Repeatable returns whether this node can appear multiple times
	Repeatable() bool
}

// NodeEvaluator defines the interface for evaluating configuration nodes
type NodeEvaluator interface {
	// Evaluate processes a configuration node and updates the config structure
	Evaluate(node parser.Node, cfg any) error
}

// NodeHandler is a function type that processes a configuration node
type NodeHandler func(node parser.Node) error

// NodeAttribute represents special attributes that can be applied to nodes
type NodeAttribute int

const (
	// Repeatable indicates that a node can appear multiple times in the configuration
	Repeatable NodeAttribute = iota
)

// CommonDef provides common functionality for node definitions
type CommonDef struct {
	name       string         // name of the node
	args       []*args.ArgDef // argument definitions
	minArgs    int            // minimum number of arguments required
	maxArgs    int            // maximum number of arguments allowed
	handler    NodeHandler    // function to handle this node
	repeatable bool           // whether this node can appear multiple times
}

func (d *CommonDef) Name() string {
	return d.name
}

func (d *CommonDef) MinArgs() int {
	return d.minArgs
}

func (d *CommonDef) MaxArgs() int {
	return d.maxArgs
}

func (d *CommonDef) Args() []*args.ArgDef {
	return d.args
}

func (d *CommonDef) addArgs(args ...*args.ArgDef) {
	beginOptional := false

	for i, arg := range args {
		if arg.Required() {
			if beginOptional {
				panic(fmt.Sprintf("node '%s': required argument %d follows optional argument(s)", d.Name(), len(d.args)+i+1))
			}
			d.minArgs++
		}
		if arg.Variadic() {
			if i != len(args)-1 {
				panic(fmt.Sprintf("node '%s': variadic argument is not the last one", d.Name()))
			}
			d.maxArgs = -1
		} else {
			d.maxArgs++
		}

		d.args = append(d.args, arg)

		if !arg.Required() {
			beginOptional = true
		}
	}
}

func (d *CommonDef) Handler() NodeHandler {
	return d.handler
}

func (d *CommonDef) Repeatable() bool {
	return d.repeatable
}

func evaluate(d NodeDefinition, node parser.Node) error {
	if node.Name != d.Name() {
		return fmt.Errorf("node '%s' is not allowed here", node.Name)
	}

	if len(node.Args) < d.MinArgs() {
		return NodeErr(node, "directive '%s' expects at least %d arguments", d.Name(), d.MinArgs())
	}
	if d.MaxArgs() != -1 && len(node.Args) > d.MaxArgs() {
		return NodeErr(node, "directive '%s' expects a maximum of %d arguments", d.Name(), d.MaxArgs())
	}

	for i, arg := range d.Args() {
		if arg.Variadic() {
			values.SetList(arg.Target(), node.Args[i:])
			break
		}
		if i < len(node.Args) {
			arg.Target().Set(node.Args[i])
		}
	}

	if d.Handler() != nil {
		return d.Handler()(node)
	}

	return nil
}
