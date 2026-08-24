package ast

// Node is a single parsed configuration node: either a simple directive or a
// block with child nodes.
//
//	name arg0 arg1 {
//		child0
//		child1
//	}
type Node struct {
	// Name is the first token on the node's line.
	Name string
	// Args are the remaining tokens on the node's line, before any block.
	Args []string
	// Children holds the nested nodes if the node is a block. It is non-nil
	// for a block, even an empty one, and nil for a plain directive.
	Children []Node
	// File is the name of the source file the node was read from.
	File string
	// Line is the line number the node starts on. For blocks this is the line
	// of the block header (name and arguments).
	Line int
}

// AST represents the Abstract Syntax Tree of a configuration file
type AST []Node
