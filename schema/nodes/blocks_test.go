package nodes

import (
	"testing"

	parser "github.com/foxcpp/maddy/framework/cfgparser"
	"github.com/open-webtech/go-xaddy-config/schema/args"
)

func TestNewBlockDef(t *testing.T) {
	var testValue string
	block := NewBlockDef("test_block", args.StringArg(&testValue))
	
	if block == nil {
		t.Fatal("NewBlockDef() returned nil")
	}
	
	if block.Name() != "test_block" {
		t.Errorf("Expected name 'test_block', got '%s'", block.Name())
	}
	
	if len(block.Args()) != 1 {
		t.Errorf("Expected 1 argument, got %d", len(block.Args()))
	}
	
	if block.MinArgs() != 1 {
		t.Errorf("Expected MinArgs 1, got %d", block.MinArgs())
	}
	
	if block.MaxArgs() != 1 {
		t.Errorf("Expected MaxArgs 1, got %d", block.MaxArgs())
	}
}

func TestBlockDefSetAttrs(t *testing.T) {
	block := NewBlockDef("test_block")
	
	// Initially not repeatable
	if block.Repeatable() {
		t.Error("Expected Repeatable to be false by default")
	}
	
	// Set repeatable
	block.SetAttrs(Repeatable)
	
	if !block.Repeatable() {
		t.Error("Expected Repeatable to be true after SetAttrs")
	}
}

func TestBlockDefSetHandler(t *testing.T) {
	block := NewBlockDef("test_block")
	
	// Initially no handler
	if block.Handler() != nil {
		t.Error("Expected Handler to be nil by default")
	}
	
	var handlerCalled bool
	handler := func(node parser.Node) error {
		handlerCalled = true
		return nil
	}
	
	block.SetHandler(handler)
	
	if block.Handler() == nil {
		t.Error("Expected Handler to be set")
	}
	
	// Test handler execution
	testNode := parser.Node{Name: "test_block"}
	err := block.Handler()(testNode)
	if err != nil {
		t.Errorf("Handler returned error: %v", err)
	}
	
	if !handlerCalled {
		t.Error("Handler was not called")
	}
}

func TestBlockDefEvaluate(t *testing.T) {
	// Set up test configuration struct
	type Config struct {
		BlockName string
		NestedValue string
	}
	cfg := &Config{}
	
	// Create block definition
	block := NewBlockDef("test_block", args.StringArg(&cfg.BlockName))
	block.DefineDirective("nested_directive", args.StringArg(&cfg.NestedValue))
	
	tests := []struct {
		name    string
		node    parser.Node
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid block with children",
			node: parser.Node{
				Name: "test_block",
				Args: []string{"block_value"},
				Children: []parser.Node{
					{Name: "nested_directive", Args: []string{"nested_value"}},
				},
			},
			wantErr: false,
		},
		{
			name: "block without children",
			node: parser.Node{
				Name: "test_block",
				Args: []string{"block_value"},
				Children: []parser.Node{},
			},
			wantErr: false,
		},
		{
			name: "wrong block name",
			node: parser.Node{
				Name: "wrong_block",
				Args: []string{"block_value"},
			},
			wantErr: true,
			errMsg:  "is not allowed here",
		},
		{
			name: "missing required argument",
			node: parser.Node{
				Name: "test_block",
				Args: []string{},
			},
			wantErr: true,
			errMsg:  "expects at least",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset config
			cfg.BlockName = ""
			cfg.NestedValue = ""
			
			err := block.Evaluate(tt.node, cfg)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && tt.errMsg != "" {
				if err == nil || !contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error containing '%s', got '%v'", tt.errMsg, err)
				}
			}
			
			if !tt.wantErr && tt.node.Name == "test_block" && len(tt.node.Args) > 0 {
				if cfg.BlockName != tt.node.Args[0] {
					t.Errorf("Expected BlockName to be '%s', got '%s'", tt.node.Args[0], cfg.BlockName)
				}
				
				// Check nested directive was processed
				if len(tt.node.Children) > 0 && tt.node.Children[0].Name == "nested_directive" {
					if cfg.NestedValue != tt.node.Children[0].Args[0] {
						t.Errorf("Expected NestedValue to be '%s', got '%s'", 
							tt.node.Children[0].Args[0], cfg.NestedValue)
					}
				}
			}
		})
	}
}

func TestBlockDefNestedStructure(t *testing.T) {
	// Set up nested configuration struct
	type NestedConfig struct {
		ServerName string
		Listen     string
		TLS        bool
		LogLevel   string
	}
	cfg := &NestedConfig{}
	
	// Create server block with nested directives
	serverBlock := NewBlockDef("server", args.StringArg(&cfg.ServerName))
	serverBlock.DefineDirective("listen", args.StringArg(&cfg.Listen))
	serverBlock.DefineDirective("tls", args.BoolArg(&cfg.TLS))
	serverBlock.DefineDirective("log_level", args.StringArg(&cfg.LogLevel))
	
	node := parser.Node{
		Name: "server",
		Args: []string{"web"},
		Children: []parser.Node{
			{Name: "listen", Args: []string{"80"}},
			{Name: "tls", Args: []string{"false"}},
			{Name: "log_level", Args: []string{"info"}},
		},
	}
	
	err := serverBlock.Evaluate(node, cfg)
	if err != nil {
		t.Fatalf("Evaluate() failed: %v", err)
	}
	
	// Verify all values were set correctly
	if cfg.ServerName != "web" {
		t.Errorf("Expected ServerName 'web', got '%s'", cfg.ServerName)
	}
	if cfg.Listen != "80" {
		t.Errorf("Expected Listen '80', got '%s'", cfg.Listen)
	}
	if cfg.TLS != false {
		t.Errorf("Expected TLS false, got %t", cfg.TLS)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("Expected LogLevel 'info', got '%s'", cfg.LogLevel)
	}
}

func TestNewModuleBlockDef(t *testing.T) {
	moduleBlock := NewModuleBlockDef("module_block")
	
	if moduleBlock == nil {
		t.Fatal("NewModuleBlockDef() returned nil")
	}
	
	if moduleBlock.Name() != "module_block" {
		t.Errorf("Expected name 'module_block', got '%s'", moduleBlock.Name())
	}
	
	// Should have one argument for module name
	if len(moduleBlock.Args()) != 1 {
		t.Errorf("Expected 1 argument for module name, got %d", len(moduleBlock.Args()))
	}
}

func TestModuleBlockDefWithArgs(t *testing.T) {
	var extraArg string
	moduleBlock := NewModuleBlockDef("module_block").WithArgs(args.StringArg(&extraArg))
	
	// Should have two arguments: module name + extra arg
	if len(moduleBlock.Args()) != 2 {
		t.Errorf("Expected 2 arguments, got %d", len(moduleBlock.Args()))
	}
	
	if moduleBlock.MinArgs() != 2 {
		t.Errorf("Expected MinArgs 2, got %d", moduleBlock.MinArgs())
	}
}

func TestModuleBlockDefEvaluate(t *testing.T) {
	type Config struct {
		Value      string
		ModuleName string
	}
	cfg := &Config{}
	
	// Create module block def with proper module name target
	moduleBlock := &ModuleBlockDef{
		modules:    make(map[string]*NodesContainer),
		moduleName: &cfg.ModuleName,
		CommonDef:  CommonDef{name: "auth"},
	}
	moduleBlock.addArgs(args.StringArg(&cfg.ModuleName))
	
	// Define a module type
	testModule := &NodesContainer{}
	testModule.DefineDirective("setting", args.StringArg(&cfg.Value))
	moduleBlock.modules["test_module"] = testModule
	
	tests := []struct {
		name       string
		node       parser.Node
		wantErr    bool
		errMsg     string
		expectValue string
	}{
		{
			name: "valid module",
			node: parser.Node{
				Name: "auth",
				Args: []string{"test_module"},
				Children: []parser.Node{
					{Name: "setting", Args: []string{"test_value"}},
				},
			},
			wantErr:     false,
			expectValue: "test_value",
		},
		{
			name: "unknown module",
			node: parser.Node{
				Name: "auth",
				Args: []string{"unknown_module"},
			},
			wantErr: true,
			errMsg:  "unknown module",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg.Value = ""
			
			err := moduleBlock.Evaluate(tt.node, cfg)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && tt.errMsg != "" {
				if err == nil || !contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error containing '%s', got '%v'", tt.errMsg, err)
				}
			}
			
			if !tt.wantErr && tt.expectValue != "" {
				if cfg.Value != tt.expectValue {
					t.Errorf("Expected Value to be '%s', got '%s'", tt.expectValue, cfg.Value)
				}
			}
		})
	}
}