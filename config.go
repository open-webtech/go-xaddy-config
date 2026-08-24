package config

import (
	"io"
	"os"

	"github.com/open-webtech/go-xaddy-config/ast"
	"github.com/open-webtech/go-xaddy-config/internal/parser"
	"github.com/open-webtech/go-xaddy-config/schema/nodes"
)

// Read parses configuration from an io.Reader and returns the AST
func Read(r io.Reader, location string) (ast.AST, error) {
	nl, err := parser.Read(r, location)
	return ast.AST(nl), err
}

// ReadFile reads and parses configuration from a file and returns the AST
func ReadFile(filename string) (ast.AST, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Read(f, filename)
}

// ExpectMaxArgN checks if a configuration node has at most the specified number of arguments
func ExpectMaxArgN(node ast.Node, num int) error {
	if len(node.Args) > num {
		return nodes.NodeErr(node, "expected at most %d arguments to %s, got %d", num, node.Name, len(node.Args))
	}
	return nil
}

// ExpectMinArgN checks if a configuration node has at least the specified number of arguments
func ExpectMinArgN(node ast.Node, num int) error {
	if len(node.Args) < num {
		return nodes.NodeErr(node, "expected at least %d arguments to %s, got %d", num, node.Name, len(node.Args))
	}
	return nil
}
