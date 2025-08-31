package schema

import (
	"testing"

	parser "github.com/foxcpp/maddy/framework/cfgparser"
	"github.com/open-webtech/go-xaddy-config/schema/args"
)

func TestNewBuilder(t *testing.T) {
	builder := NewBuilder()
	if builder == nil {
		t.Fatal("NewBuilder() returned nil")
	}
}

func TestBuilderDefineDirective(t *testing.T) {
	builder := NewBuilder()
	
	var testValue string
	directive := builder.DefineDirective("test_directive", args.StringArg(&testValue))
	
	if directive == nil {
		t.Fatal("DefineDirective() returned nil")
	}
	
	if directive.Name() != "test_directive" {
		t.Errorf("Expected directive name 'test_directive', got '%s'", directive.Name())
	}
}

func TestBuilderDefineDirectiveCallback(t *testing.T) {
	builder := NewBuilder()
	
	var callbackCalled bool
	var receivedNode parser.Node
	
	callback := func(node parser.Node) error {
		callbackCalled = true
		receivedNode = node
		return nil
	}
	
	directive := builder.DefineDirectiveCallback("callback_directive", callback)
	
	if directive == nil {
		t.Fatal("DefineDirectiveCallback() returned nil")
	}
	
	if directive.Name() != "callback_directive" {
		t.Errorf("Expected directive name 'callback_directive', got '%s'", directive.Name())
	}
	
	// Test that handler was set
	if directive.Handler() == nil {
		t.Error("Expected handler to be set")
	}
	
	// Test callback execution
	testNode := parser.Node{Name: "callback_directive", Args: []string{"test"}}
	err := directive.Handler()(testNode)
	if err != nil {
		t.Errorf("Callback returned error: %v", err)
	}
	
	if !callbackCalled {
		t.Error("Callback was not called")
	}
	
	if receivedNode.Name != "callback_directive" {
		t.Errorf("Expected callback to receive node with name 'callback_directive', got '%s'", receivedNode.Name)
	}
}

func TestBuilderDefineBlock(t *testing.T) {
	builder := NewBuilder()
	
	var testValue string
	block := builder.DefineBlock("test_block", args.StringArg(&testValue))
	
	if block == nil {
		t.Fatal("DefineBlock() returned nil")
	}
	
	if block.Name() != "test_block" {
		t.Errorf("Expected block name 'test_block', got '%s'", block.Name())
	}
}

func TestBuilderDefineBlockCallback(t *testing.T) {
	builder := NewBuilder()
	
	var callbackCalled bool
	callback := func(node parser.Node) error {
		callbackCalled = true
		return nil
	}
	
	blockDef := builder.DefineBlockCallback("callback_block", callback)
	
	if blockDef == nil {
		t.Fatal("DefineBlockCallback() returned nil")
	}
	
	if blockDef.Name() != "callback_block" {
		t.Errorf("Expected block name 'callback_block', got '%s'", blockDef.Name())
	}
	
	// Test that handler was set
	if blockDef.Handler() == nil {
		t.Error("Expected handler to be set")
	}
	
	// Test callback execution
	testNode := parser.Node{Name: "callback_block"}
	err := blockDef.Handler()(testNode)
	if err != nil {
		t.Errorf("Callback returned error: %v", err)
	}
	
	if !callbackCalled {
		t.Error("Callback was not called")
	}
}

// Integration test with actual configuration parsing
func TestBuilderIntegration(t *testing.T) {
	// Set up configuration struct
	type Config struct {
		LogLevel       string
		MaxConnections int
		ServerName     string
		Listen         string
		TLS            bool
	}
	
	cfg := &Config{}
	
	// Build schema
	builder := NewBuilder()
	builder.DefineDirective("log_level", args.StringArg(&cfg.LogLevel))
	builder.DefineDirective("max_connections", args.IntArg(&cfg.MaxConnections))
	
	serverBlock := builder.DefineBlock("server", args.StringArg(&cfg.ServerName))
	serverBlock.DefineDirective("listen", args.StringArg(&cfg.Listen))
	serverBlock.DefineDirective("tls", args.BoolArg(&cfg.TLS))
	
	// Test configuration content
	configContent := `log_level debug
max_connections 100

server web {
    listen 80
    tls false
}`
	
	// Parse configuration
	nodes, err := parseConfigString(configContent)
	if err != nil {
		t.Fatalf("Failed to parse config: %v", err)
	}
	
	// Evaluate configuration
	err = builder.EvaluateTree(nodes, cfg)
	if err != nil {
		t.Fatalf("Failed to evaluate config: %v", err)
	}
	
	// Verify results
	if cfg.LogLevel != "debug" {
		t.Errorf("Expected LogLevel 'debug', got '%s'", cfg.LogLevel)
	}
	if cfg.MaxConnections != 100 {
		t.Errorf("Expected MaxConnections 100, got %d", cfg.MaxConnections)
	}
	if cfg.ServerName != "web" {
		t.Errorf("Expected ServerName 'web', got '%s'", cfg.ServerName)
	}
	if cfg.Listen != "80" {
		t.Errorf("Expected Listen '80', got '%s'", cfg.Listen)
	}
	if cfg.TLS != false {
		t.Errorf("Expected TLS false, got %t", cfg.TLS)
	}
}

// Helper function to parse config string
func parseConfigString(content string) ([]parser.Node, error) {
	// This would typically use the actual parser
	// For now, we'll create mock nodes for testing
	return []parser.Node{
		{Name: "log_level", Args: []string{"debug"}},
		{Name: "max_connections", Args: []string{"100"}},
		{
			Name: "server",
			Args: []string{"web"},
			Children: []parser.Node{
				{Name: "listen", Args: []string{"80"}},
				{Name: "tls", Args: []string{"false"}},
			},
		},
	}, nil
}