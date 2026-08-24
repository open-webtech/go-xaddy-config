package parser

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/open-webtech/go-xaddy-config/ast"
)

// maxImportDepth bounds recursive snippet and file imports, so that import
// cycles fail with a clear error instead of endless recursion.
const maxImportDepth = 100

var (
	snippetNameRe = regexp.MustCompile(`^\(([^()\s]+)\)$`)
	macroDefRe    = regexp.MustCompile(`^\$\(([^)\s]+)\)$`)
	macroRefRe    = regexp.MustCompile(`\$\(([^)\s]+)\)`)
	envRefRe      = regexp.MustCompile(`\{env:([^{}]+)\}`)
)

// expander tracks the snippets available while resolving import directives.
type expander struct {
	snippets map[string][]ast.Node
}

// expandTop processes the top level of a file: snippet definitions are
// collected and removed, then import directives in the remaining nodes are
// replaced by the content they refer to.
func (e *expander) expandTop(nodes []ast.Node, depth int) ([]ast.Node, error) {
	rest := make([]ast.Node, 0, len(nodes))
	for _, n := range nodes {
		m := snippetNameRe.FindStringSubmatch(n.Name)
		if m == nil {
			rest = append(rest, n)
			continue
		}
		if len(n.Args) != 0 {
			return nil, nodeErrf(n, "snippet '%s' cannot have arguments", n.Name)
		}
		e.snippets[m[1]] = n.Children
	}
	return e.expandImports(rest, depth)
}

// expandImports replaces every import directive in the tree with the body of
// the snippet or file it names. Snippets take precedence over files.
func (e *expander) expandImports(nodes []ast.Node, depth int) ([]ast.Node, error) {
	if depth > maxImportDepth {
		return nil, errors.New("too many nested imports (import cycle?)")
	}
	out := make([]ast.Node, 0, len(nodes))
	for _, n := range nodes {
		if n.Name == "import" {
			if len(n.Args) != 1 || n.Children != nil {
				return nil, nodeErrf(n, "import expects exactly one argument: a snippet name or file path")
			}
			imported, err := e.resolveImport(n, depth)
			if err != nil {
				return nil, err
			}
			out = append(out, imported...)
			continue
		}
		if len(n.Children) > 0 {
			children, err := e.expandImports(n.Children, depth)
			if err != nil {
				return nil, err
			}
			n.Children = children
		}
		out = append(out, n)
	}
	return out, nil
}

func (e *expander) resolveImport(n ast.Node, depth int) ([]ast.Node, error) {
	target := n.Args[0]
	if body, ok := e.snippets[target]; ok {
		return e.expandImports(copyNodes(body), depth+1)
	}
	f, err := os.Open(target)
	if err != nil {
		return nil, nodeErrf(n, "import '%s' does not name a snippet or a readable file: %v", target, err)
	}
	defer f.Close()
	tree, err := parseTree(f, target)
	if err != nil {
		return nil, err
	}
	return e.expandTop(tree, depth+1)
}

func copyNodes(nodes []ast.Node) []ast.Node {
	out := make([]ast.Node, len(nodes))
	for i, n := range nodes {
		out[i] = n
		out[i].Args = append([]string(nil), n.Args...)
		if n.Children != nil {
			out[i].Children = copyNodes(n.Children)
		}
	}
	return out
}

// expandMacros collects top-level macro definitions of the form
// "$(name) = value..." and substitutes references to them throughout the
// tree. A definition may reference macros defined before it. A reference
// that stands alone as an argument splices in all of the macro's values;
// a reference inside a longer string is replaced textually. References to
// undefined macros are left unchanged.
func expandMacros(nodes []ast.Node) ([]ast.Node, error) {
	macros := map[string][]string{}
	rest := make([]ast.Node, 0, len(nodes))
	for _, n := range nodes {
		m := macroDefRe.FindStringSubmatch(n.Name)
		if m == nil || len(n.Args) == 0 || n.Args[0] != "=" {
			rest = append(rest, n)
			continue
		}
		if n.Children != nil {
			return nil, nodeErrf(n, "macro definition '%s' cannot be a block", n.Name)
		}
		macros[m[1]] = expandMacroStrings(n.Args[1:], macros)
	}
	if len(macros) == 0 {
		return rest, nil
	}
	expandMacroNodes(rest, macros)
	return rest, nil
}

func expandMacroNodes(nodes []ast.Node, macros map[string][]string) {
	for i := range nodes {
		nodes[i].Name = expandMacroString(nodes[i].Name, macros)
		nodes[i].Args = expandMacroStrings(nodes[i].Args, macros)
		expandMacroNodes(nodes[i].Children, macros)
	}
}

func expandMacroStrings(args []string, macros map[string][]string) []string {
	out := make([]string, 0, len(args))
	for _, arg := range args {
		if m := macroDefRe.FindStringSubmatch(arg); m != nil {
			if values, ok := macros[m[1]]; ok {
				out = append(out, values...)
				continue
			}
		}
		out = append(out, expandMacroString(arg, macros))
	}
	return out
}

func expandMacroString(s string, macros map[string][]string) string {
	return macroRefRe.ReplaceAllStringFunc(s, func(ref string) string {
		name := macroRefRe.FindStringSubmatch(ref)[1]
		if values, ok := macros[name]; ok {
			return strings.Join(values, " ")
		}
		return ref
	})
}

// expandEnv replaces {env:NAME} placeholders in node names and arguments.
// Undefined variables expand to an empty string; malformed placeholders are
// left untouched.
func expandEnv(nodes []ast.Node) []ast.Node {
	for i := range nodes {
		nodes[i].Name = expandEnvString(nodes[i].Name)
		for j := range nodes[i].Args {
			nodes[i].Args[j] = expandEnvString(nodes[i].Args[j])
		}
		expandEnv(nodes[i].Children)
	}
	return nodes
}

func expandEnvString(s string) string {
	return envRefRe.ReplaceAllStringFunc(s, func(ref string) string {
		return os.Getenv(ref[len("{env:") : len(ref)-1])
	})
}

// nodeErrf formats an error message, prefixed with the node's source
// location when one is known.
func nodeErrf(node ast.Node, format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)
	if node.File == "" {
		return errors.New(msg)
	}
	return fmt.Errorf("%s:%d: %s", node.File, node.Line, msg)
}
