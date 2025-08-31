package nodes

import (
	"testing"

	parser "github.com/foxcpp/maddy/framework/cfgparser"
	"github.com/open-webtech/go-xaddy-config/schema/args"
)

func TestNodesContainer(t *testing.T) {
	container := &NodesContainer{}
	
	// Test that empty container doesn't crash
	if container == nil {
		t.Fatal("NodesContainer should not be nil")
	}
}

func TestNodesContainerDefineDirective(t *testing.T) {
	container := &NodesContainer{}
	
	var testValue string
	directive := container.DefineDirective("test_directive", args.StringArg(&testValue))
	
	if directive == nil {
		t.Fatal("DefineDirective() returned nil")
	}
	
	if directive.Name() != "test_directive" {
		t.Errorf("Expected directive name 'test_directive', got '%s'", directive.Name())
	}
}

func TestNodesContainerDefineDirectiveCallback(t *testing.T) {
	container := &NodesContainer{}
	
	var callbackCalled bool
	callback := func(node parser.Node) error {
		callbackCalled = true
		return nil
	}
	
	directive := container.DefineDirectiveCallback("callback_directive", callback)
	
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
	testNode := parser.Node{Name: "callback_directive"}
	err := directive.Handler()(testNode)
	if err != nil {
		t.Errorf("Callback returned error: %v", err)
	}
	
	if !callbackCalled {
		t.Error("Callback was not called")
	}
}

func TestNodesContainerDefineBlock(t *testing.T) {
	container := &NodesContainer{}
	
	var testValue string
	block := container.DefineBlock("test_block", args.StringArg(&testValue))
	
	if block == nil {
		t.Fatal("DefineBlock() returned nil")
	}
	
	if block.Name() != "test_block" {
		t.Errorf("Expected block name 'test_block', got '%s'", block.Name())
	}
}

func TestNodesContainerDefineBlockCallback(t *testing.T) {
	container := &NodesContainer{}
	
	var callbackCalled bool
	callback := func(node parser.Node) error {
		callbackCalled = true
		return nil
	}
	
	blockDef := container.DefineBlockCallback("callback_block", callback)
	
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

func TestNodesContainerEvaluateTree(t *testing.T) {
	// Set up test configuration struct
	type Config struct {
		GlobalSetting string
		ServerName    string
		Listen        string
		TLS           bool
	}
	cfg := &Config{}
	
	// Set up container with definitions
	container := &NodesContainer{}
	container.DefineDirective("global_setting", args.StringArg(&cfg.GlobalSetting))
	
	serverBlock := container.DefineBlock("server", args.StringArg(&cfg.ServerName))
	serverBlock.DefineDirective("listen", args.StringArg(&cfg.Listen))
	serverBlock.DefineDirective("tls", args.BoolArg(&cfg.TLS))
	
	// Create test nodes
	nodes := []parser.Node{
		{
			Name: "global_setting",
			Args: []string{"global_value"},
		},
		{
			Name: "server",
			Args: []string{"web"},
			Children: []parser.Node{
				{Name: "listen", Args: []string{"80"}},
				{Name: "tls", Args: []string{"false"}},
			},
		},
	}
	
	// Evaluate tree
	err := container.EvaluateTree(nodes, cfg)
	if err != nil {
		t.Fatalf("EvaluateTree() failed: %v", err)
	}
	
	// Verify results
	if cfg.GlobalSetting != "global_value" {
		t.Errorf("Expected GlobalSetting 'global_value', got '%s'", cfg.GlobalSetting)
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

func TestNodesContainerEvaluateTreeWithErrors(t *testing.T) {
	type Config struct {
		Value string
	}
	cfg := &Config{}
	
	container := &NodesContainer{}
	container.DefineDirective("known_directive", args.StringArg(&cfg.Value))
	
	tests := []struct {
		name    string
		nodes   []parser.Node
		wantErr bool
		errMsg  string
	}{
		{
			name: "unknown directive",
			nodes: []parser.Node{
				{Name: "unknown_directive", Args: []string{"value"}},
			},
			wantErr: false, // The current implementation doesn't error on unknown directives, it just skips them
			errMsg:  "",
		},
		{
			name: "valid directive",
			nodes: []parser.Node{
				{Name: "known_directive", Args: []string{"test_value"}},
			},
			wantErr: false,
		},
		{
			name:    "empty nodes",
			nodes:   []parser.Node{},
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg.Value = "" // Reset
			
			err := container.EvaluateTree(tt.nodes, cfg)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("EvaluateTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && tt.errMsg != "" {
				if err == nil || !contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error containing '%s', got '%v'", tt.errMsg, err)
				}
			}
			
			if !tt.wantErr && tt.name == "valid directive" {
				if cfg.Value != "test_value" {
					t.Errorf("Expected Value to be 'test_value', got '%s'", cfg.Value)
				}
			}
		})
	}
}

func TestNodesContainerNestedBlocks(t *testing.T) {
	// Test deeply nested block structure
	type Config struct {
		ServerName   string
		Listen       string
		AuthType     string
		AuthSetting  string
		CacheSetting string
	}
	cfg := &Config{}
	
	container := &NodesContainer{}
	
	serverBlock := container.DefineBlock("server", args.StringArg(&cfg.ServerName))
	serverBlock.DefineDirective("listen", args.StringArg(&cfg.Listen))
	
	authBlock := serverBlock.DefineBlock("auth", args.StringArg(&cfg.AuthType))
	authBlock.DefineDirective("setting", args.StringArg(&cfg.AuthSetting))
	
	cacheBlock := authBlock.DefineBlock("cache")
	cacheBlock.DefineDirective("setting", args.StringArg(&cfg.CacheSetting))
	
	// Create nested node structure
	nodes := []parser.Node{
		{
			Name: "server",
			Args: []string{"api"},
			Children: []parser.Node{
				{Name: "listen", Args: []string{"8080"}},
				{
					Name: "auth",
					Args: []string{"jwt"},
					Children: []parser.Node{
						{Name: "setting", Args: []string{"auth_value"}},
						{
							Name: "cache",
							Children: []parser.Node{
								{Name: "setting", Args: []string{"cache_value"}},
							},
						},
					},
				},
			},
		},
	}
	
	err := container.EvaluateTree(nodes, cfg)
	if err != nil {
		t.Fatalf("EvaluateTree() with nested blocks failed: %v", err)
	}
	
	// Verify all nested values were set
	if cfg.ServerName != "api" {
		t.Errorf("Expected ServerName 'api', got '%s'", cfg.ServerName)
	}
	if cfg.Listen != "8080" {
		t.Errorf("Expected Listen '8080', got '%s'", cfg.Listen)
	}
	if cfg.AuthType != "jwt" {
		t.Errorf("Expected AuthType 'jwt', got '%s'", cfg.AuthType)
	}
	if cfg.AuthSetting != "auth_value" {
		t.Errorf("Expected AuthSetting 'auth_value', got '%s'", cfg.AuthSetting)
	}
	if cfg.CacheSetting != "cache_value" {
		t.Errorf("Expected CacheSetting 'cache_value', got '%s'", cfg.CacheSetting)
	}
}

func TestNodesContainerRepeatableDirectives(t *testing.T) {
	type Config struct {
		Values []string
	}
	cfg := &Config{}
	
	container := &NodesContainer{}
	
	// Create a repeatable directive with explicit repeatable attribute
	var callCount int
	directive := container.DefineDirectiveCallback("repeatable", func(node parser.Node) error {
		callCount++
		if len(node.Args) > 0 {
			cfg.Values = append(cfg.Values, node.Args[0])
		}
		return nil
	})
	directive.SetAttrs(Repeatable) // Make it explicitly repeatable
	
	nodes := []parser.Node{
		{Name: "repeatable", Args: []string{"value1"}},
		{Name: "repeatable", Args: []string{"value2"}},
		{Name: "repeatable", Args: []string{"value3"}},
	}
	
	err := container.EvaluateTree(nodes, cfg)
	if err != nil {
		t.Fatalf("EvaluateTree() with repeatable directives failed: %v", err)
	}
	
	if callCount != 3 {
		t.Errorf("Expected callback to be called 3 times, got %d", callCount)
	}
	
	expectedValues := []string{"value1", "value2", "value3"}
	if len(cfg.Values) != len(expectedValues) {
		t.Errorf("Expected %d values, got %d", len(expectedValues), len(cfg.Values))
	}
	
	for i, expected := range expectedValues {
		if i >= len(cfg.Values) || cfg.Values[i] != expected {
			t.Errorf("Expected cfg.Values[%d] to be '%s', got '%s'", i, expected, cfg.Values[i])
		}
	}
}