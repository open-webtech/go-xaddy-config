package nodes

import (
	"fmt"

	parser "github.com/foxcpp/maddy/framework/cfgparser"
)

// NodeErr creates a formatted error message for configuration nodes.
// If a file location is available, it prepends the file and line number to the error message.
// If no file location is available, it returns a standard formatted error.
func NodeErr(node parser.Node, errMsg string, args ...interface{}) error {
	if node.File == "" {
		return fmt.Errorf(errMsg, args...)
	}
	return fmt.Errorf("%s:%d: %s", node.File, node.Line, fmt.Sprintf(errMsg, args...))
}
