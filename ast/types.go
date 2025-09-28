package ast

import (
	parser "github.com/foxcpp/maddy/framework/cfgparser"
)

// Node wraps the cfgparser.Node type for convenience
type Node parser.Node

// AST represents the Abstract Syntax Tree of a configuration file
type AST []parser.Node
