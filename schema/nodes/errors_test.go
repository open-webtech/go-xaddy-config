package nodes

import (
	"strings"
	"testing"

	parser "github.com/foxcpp/maddy/framework/cfgparser"
)

func TestNodeErr(t *testing.T) {
	node := parser.Node{
		Name: "test_directive",
		Args: []string{"arg1", "arg2"},
		File: "test.conf",
		Line: 42,
	}

	tests := []struct {
		name        string
		format      string
		args        []interface{}
		expectFile  bool
		expectLine  bool
		expectMsg   bool
	}{
		{
			name:       "simple error message",
			format:     "test error",
			args:       []interface{}{},
			expectFile: true,
			expectLine: true,
			expectMsg:  true,
		},
		{
			name:       "formatted error message",
			format:     "expected %d arguments, got %d",
			args:       []interface{}{2, 1},
			expectFile: true,
			expectLine: true,
			expectMsg:  true,
		},
		{
			name:       "error with directive name",
			format:     "directive '%s' is not allowed here",
			args:       []interface{}{"test_directive"},
			expectFile: true,
			expectLine: true,
			expectMsg:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NodeErr(node, tt.format, tt.args...)
			
			if err == nil {
				t.Fatal("NodeErr() returned nil")
			}

			errStr := err.Error()
			
			// Check that error contains file information
			if tt.expectFile && !strings.Contains(errStr, node.File) {
				t.Errorf("Error should contain file name '%s', got: %s", node.File, errStr)
			}

			// Check that error contains line information
			if tt.expectLine && !strings.Contains(errStr, "42") {
				t.Errorf("Error should contain line number '42', got: %s", errStr)
			}

			// Check that error contains the formatted message
			if tt.expectMsg {
				expectedMsg := tt.format
				if len(tt.args) > 0 {
					// For formatted strings, just check that some formatting occurred
					// by ensuring the error is not exactly the format string
					if errStr == tt.format {
						t.Errorf("Error should be formatted, but got exact format string: %s", errStr)
					}
				} else {
					// For simple strings, check that it contains the message
					if !strings.Contains(errStr, expectedMsg) {
						t.Errorf("Error should contain message '%s', got: %s", expectedMsg, errStr)
					}
				}
			}
		})
	}
}

func TestNodeErrWithEmptyNode(t *testing.T) {
	// Test with node that has minimal information
	node := parser.Node{
		Name: "minimal_node",
	}

	err := NodeErr(node, "test error message")
	if err == nil {
		t.Fatal("NodeErr() returned nil")
	}

	errStr := err.Error()
	if !strings.Contains(errStr, "test error message") {
		t.Errorf("Error should contain message 'test error message', got: %s", errStr)
	}

	// Should still work even without file/line info
	if errStr == "" {
		t.Error("Error string should not be empty")
	}
}

func TestNodeErrWithVariousArguments(t *testing.T) {
	node := parser.Node{
		Name: "test_node",
		File: "config.conf",
		Line: 10,
	}

	// Test with string argument
	err1 := NodeErr(node, "invalid value '%s'", "bad_value")
	if !strings.Contains(err1.Error(), "bad_value") {
		t.Errorf("Error should contain 'bad_value', got: %s", err1.Error())
	}

	// Test with integer argument
	err2 := NodeErr(node, "expected %d arguments", 3)
	if !strings.Contains(err2.Error(), "3") {
		t.Errorf("Error should contain '3', got: %s", err2.Error())
	}

	// Test with multiple arguments
	err3 := NodeErr(node, "expected between %d and %d arguments, got %d", 1, 3, 5)
	errStr := err3.Error()
	if !strings.Contains(errStr, "1") || !strings.Contains(errStr, "3") || !strings.Contains(errStr, "5") {
		t.Errorf("Error should contain '1', '3', and '5', got: %s", errStr)
	}
}

func TestNodeErrFormatting(t *testing.T) {
	node := parser.Node{
		Name: "format_test",
		File: "test.conf", 
		Line: 123,
	}

	err := NodeErr(node, "test %s with %d values", "formatting", 42)
	errStr := err.Error()

	// The exact format may vary, but should include key information
	expectedParts := []string{
		"test.conf",  // file name
		"123",        // line number
		"formatting", // string argument
		"42",         // integer argument
	}

	for _, part := range expectedParts {
		if !strings.Contains(errStr, part) {
			t.Errorf("Error should contain '%s', got: %s", part, errStr)
		}
	}
}

// Test that NodeErr works with nodes that have different field combinations
func TestNodeErrWithDifferentNodeTypes(t *testing.T) {
	tests := []struct {
		name string
		node parser.Node
	}{
		{
			name: "node with all fields",
			node: parser.Node{
				Name: "full_node",
				Args: []string{"arg1", "arg2"},
				File: "full.conf",
				Line: 100,
				Children: []parser.Node{
					{Name: "child", Args: []string{"child_arg"}},
				},
			},
		},
		{
			name: "node with just name",
			node: parser.Node{
				Name: "simple_node",
			},
		},
		{
			name: "node with name and args",
			node: parser.Node{
				Name: "node_with_args",
				Args: []string{"value1", "value2", "value3"},
			},
		},
		{
			name: "node with file but no line",
			node: parser.Node{
				Name: "file_only",
				File: "file_only.conf",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NodeErr(tt.node, "test error for %s", tt.node.Name)
			
			if err == nil {
				t.Fatal("NodeErr() returned nil")
			}

			errStr := err.Error()
			
			// Should always contain the node name
			if !strings.Contains(errStr, tt.node.Name) {
				t.Errorf("Error should contain node name '%s', got: %s", tt.node.Name, errStr)
			}

			// Should contain file if present
			if tt.node.File != "" && !strings.Contains(errStr, tt.node.File) {
				t.Errorf("Error should contain file '%s', got: %s", tt.node.File, errStr)
			}

			// Error should not be empty
			if errStr == "" {
				t.Error("Error string should not be empty")
			}
		})
	}
}