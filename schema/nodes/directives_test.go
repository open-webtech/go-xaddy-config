package nodes

import (
	"testing"

	parser "github.com/foxcpp/maddy/framework/cfgparser"
	"github.com/open-webtech/go-xaddy-config/schema/args"
)

func TestNewDirectiveDef(t *testing.T) {
	var testValue string
	directive := NewDirectiveDef("test_directive", args.StringArg(&testValue))
	
	if directive == nil {
		t.Fatal("NewDirectiveDef() returned nil")
	}
	
	if directive.Name() != "test_directive" {
		t.Errorf("Expected name 'test_directive', got '%s'", directive.Name())
	}
	
	if len(directive.Args()) != 1 {
		t.Errorf("Expected 1 argument, got %d", len(directive.Args()))
	}
	
	if directive.MinArgs() != 1 {
		t.Errorf("Expected MinArgs 1, got %d", directive.MinArgs())
	}
	
	if directive.MaxArgs() != 1 {
		t.Errorf("Expected MaxArgs 1, got %d", directive.MaxArgs())
	}
	
	if directive.Repeatable() {
		t.Error("Expected Repeatable to be false by default")
	}
}

func TestDirectiveDefAddArgs(t *testing.T) {
	var testValue1, testValue2 string
	directive := NewDirectiveDef("test_directive")
	
	// Initially no args
	if len(directive.Args()) != 0 {
		t.Errorf("Expected 0 arguments initially, got %d", len(directive.Args()))
	}
	
	// Add args
	directive.AddArgs(args.StringArg(&testValue1), args.StringArg(&testValue2))
	
	if len(directive.Args()) != 2 {
		t.Errorf("Expected 2 arguments after AddArgs, got %d", len(directive.Args()))
	}
	
	if directive.MinArgs() != 2 {
		t.Errorf("Expected MinArgs 2, got %d", directive.MinArgs())
	}
	
	if directive.MaxArgs() != 2 {
		t.Errorf("Expected MaxArgs 2, got %d", directive.MaxArgs())
	}
}

func TestDirectiveDefSetAttrs(t *testing.T) {
	directive := NewDirectiveDef("test_directive")
	
	// Initially not repeatable
	if directive.Repeatable() {
		t.Error("Expected Repeatable to be false by default")
	}
	
	// Set repeatable
	directive.SetAttrs(Repeatable)
	
	if !directive.Repeatable() {
		t.Error("Expected Repeatable to be true after SetAttrs")
	}
}

func TestDirectiveDefSetHandler(t *testing.T) {
	directive := NewDirectiveDef("test_directive")
	
	// Initially no handler
	if directive.Handler() != nil {
		t.Error("Expected Handler to be nil by default")
	}
	
	var handlerCalled bool
	handler := func(node parser.Node) error {
		handlerCalled = true
		return nil
	}
	
	directive.SetHandler(handler)
	
	if directive.Handler() == nil {
		t.Error("Expected Handler to be set")
	}
	
	// Test handler execution
	testNode := parser.Node{Name: "test_directive"}
	err := directive.Handler()(testNode)
	if err != nil {
		t.Errorf("Handler returned error: %v", err)
	}
	
	if !handlerCalled {
		t.Error("Handler was not called")
	}
}

func TestDirectiveDefEvaluate(t *testing.T) {
	tests := []struct {
		name    string
		node    parser.Node
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid directive",
			node: parser.Node{
				Name: "test_directive",
				Args: []string{"value"},
			},
			wantErr: false,
		},
		{
			name: "directive with block should fail",
			node: parser.Node{
				Name: "test_directive",
				Args: []string{"value"},
				Children: []parser.Node{
					{Name: "nested", Args: []string{"value"}},
				},
			},
			wantErr: true,
			errMsg:  "may not be a block",
		},
		{
			name: "wrong directive name",
			node: parser.Node{
				Name: "wrong_name",
				Args: []string{"value"},
			},
			wantErr: true,
			errMsg:  "is not allowed here",
		},
		{
			name: "too few arguments",
			node: parser.Node{
				Name: "test_directive",
				Args: []string{},
			},
			wantErr: true,
			errMsg:  "expects at least",
		},
		{
			name: "too many arguments",
			node: parser.Node{
				Name: "test_directive",
				Args: []string{"arg1", "arg2"},
			},
			wantErr: true,
			errMsg:  "expects a maximum of",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testValue string
			directive := NewDirectiveDef("test_directive", args.StringArg(&testValue))
			
			err := directive.Evaluate(tt.node, nil)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && tt.errMsg != "" {
				if err == nil || !contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error containing '%s', got '%v'", tt.errMsg, err)
				}
			}
			
			if !tt.wantErr && tt.node.Name == "test_directive" && len(tt.node.Args) == 1 {
				if testValue != tt.node.Args[0] {
					t.Errorf("Expected testValue to be '%s', got '%s'", tt.node.Args[0], testValue)
				}
			}
		})
	}
}

func TestDirectiveDefWithOptionalArgs(t *testing.T) {
	var required, optional string
	directive := NewDirectiveDef("test_directive",
		args.StringArg(&required),
		args.StringArg(&optional, args.Optional),
	)
	
	// Test with required arg only
	node1 := parser.Node{
		Name: "test_directive",
		Args: []string{"req_value"},
	}
	
	err := directive.Evaluate(node1, nil)
	if err != nil {
		t.Errorf("Evaluate() with required arg only failed: %v", err)
	}
	
	if required != "req_value" {
		t.Errorf("Expected required to be 'req_value', got '%s'", required)
	}
	
	// Test with both args
	node2 := parser.Node{
		Name: "test_directive",
		Args: []string{"req_value2", "opt_value"},
	}
	
	err = directive.Evaluate(node2, nil)
	if err != nil {
		t.Errorf("Evaluate() with both args failed: %v", err)
	}
	
	if required != "req_value2" {
		t.Errorf("Expected required to be 'req_value2', got '%s'", required)
	}
	if optional != "opt_value" {
		t.Errorf("Expected optional to be 'opt_value', got '%s'", optional)
	}
}

func TestDirectiveDefWithVariadicArgs(t *testing.T) {
	var values []string
	directive := NewDirectiveDef("test_directive", args.VariadicStringArg(&values))
	
	node := parser.Node{
		Name: "test_directive",
		Args: []string{"value1", "value2", "value3"},
	}
	
	err := directive.Evaluate(node, nil)
	if err != nil {
		t.Errorf("Evaluate() with variadic args failed: %v", err)
	}
	
	expectedValues := []string{"value1", "value2", "value3"}
	if len(values) != len(expectedValues) {
		t.Errorf("Expected %d values, got %d", len(expectedValues), len(values))
	}
	
	for i, expected := range expectedValues {
		if i >= len(values) || values[i] != expected {
			t.Errorf("Expected values[%d] to be '%s', got '%s'", i, expected, values[i])
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && 
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		 func() bool {
			for i := 1; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}()))
}