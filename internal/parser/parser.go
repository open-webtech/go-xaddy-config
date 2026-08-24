// Package parser reads Caddyfile-style configuration into an ast.Node tree.
//
// Tokenization (quoting, comments, source positions) is provided by the
// Apache-2.0 licensed caddyfile package maintained by CoreDNS. The dialect
// extensions layered on top — snippets, imports, macros and {env:...}
// substitution — are implemented in this package.
package parser

import (
	"io"

	"github.com/coredns/caddy/caddyfile"
	"github.com/open-webtech/go-xaddy-config/ast"
)

// Read parses a configuration stream into a node tree and applies the
// dialect expansions in order: snippets and imports, then macros, then
// environment variable placeholders.
func Read(r io.Reader, location string) ([]ast.Node, error) {
	tree, err := parseTree(r, location)
	if err != nil {
		return nil, err
	}
	e := &expander{snippets: map[string][]ast.Node{}}
	tree, err = e.expandTop(tree, 0)
	if err != nil {
		return nil, err
	}
	tree, err = expandMacros(tree)
	if err != nil {
		return nil, err
	}
	return expandEnv(tree), nil
}

// parseTree reads the raw node tree from a stream, with no dialect
// expansions applied.
func parseTree(r io.Reader, location string) ([]ast.Node, error) {
	d := caddyfile.NewDispenser(location, r)
	var nodes []ast.Node
	for d.Next() {
		if d.Val() == "}" {
			return nil, d.Err("unexpected '}'")
		}
		node, closed, err := readNode(&d)
		if err != nil {
			return nil, err
		}
		if closed {
			return nil, d.Err("unexpected '}'")
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

// readNode reads one node starting at the token under the cursor, which is
// the node's name. Arguments are the following tokens on the same line; a
// '{' among them opens a child block. closed reports that a '}' terminating
// the enclosing block was consumed while collecting arguments, as happens
// with single-line blocks such as "block { child value }".
func readNode(d *caddyfile.Dispenser) (node ast.Node, closed bool, err error) {
	if d.Val() == "{" {
		return node, false, d.SyntaxErr("directive or block name")
	}
	node = ast.Node{Name: d.Val(), File: d.File(), Line: d.Line()}
	for d.NextArg() {
		switch d.Val() {
		case "{":
			node.Children, err = readBlock(d)
			return node, false, err
		case "}":
			return node, true, nil
		default:
			node.Args = append(node.Args, d.Val())
		}
	}
	return node, false, nil
}

// readBlock reads the children of a block whose opening '{' was just
// consumed, up to and including the matching '}'.
func readBlock(d *caddyfile.Dispenser) ([]ast.Node, error) {
	var children []ast.Node
	for {
		if !d.Next() {
			return nil, d.EOFErr()
		}
		if d.Val() == "}" {
			return children, nil
		}
		child, closed, err := readNode(d)
		if err != nil {
			return nil, err
		}
		children = append(children, child)
		if closed {
			return children, nil
		}
	}
}
