package config

import (
	"os"
	"strings"
	"testing"

	parser "github.com/foxcpp/maddy/framework/cfgparser"
)

func TestRead(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantErr  bool
		wantLen  int
		location string
	}{
		{
			name:     "simple config",
			content:  "directive value\nanother_directive value2",
			wantErr:  false,
			wantLen:  2,
			location: "test.conf",
		},
		{
			name:     "config with block",
			content:  "block_name {\n    nested value\n}",
			wantErr:  false,
			wantLen:  1,
			location: "test.conf",
		},
		{
			name:     "empty config",
			content:  "",
			wantErr:  false,
			wantLen:  0,
			location: "empty.conf",
		},
		{
			name:     "config with comments",
			content:  "# comment\ndirective value\n# another comment",
			wantErr:  false,
			wantLen:  1,
			location: "test.conf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.content)
			got, err := Read(reader, tt.location)

			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(got) != tt.wantLen {
				t.Errorf("Read() returned %d nodes, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestReadFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
		wantLen  int
	}{
		{
			name:     "simple config file",
			filename: "testdata/simple.conf",
			wantErr:  false,
			wantLen:  4, // log_level, max_connections, server web, server api
		},
		{
			name:     "config with snippets",
			filename: "testdata/with_snippets.conf",
			wantErr:  false,
			wantLen:  2, // 2 top-level blocks (main_server, backup_server) - snippets are internal
		},
		{
			name:     "config with macros",
			filename: "testdata/with_macros.conf",
			wantErr:  false,
			wantLen:  4, // 4 top-level nodes: 2 server blocks, 1 tls_files, 1 log_file - macros are internal
		},
		{
			name:     "non-existent file",
			filename: "testdata/nonexistent.conf",
			wantErr:  true,
			wantLen:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadFile(tt.filename)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(got) != tt.wantLen {
				t.Errorf("ReadFile() returned %d nodes, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestExpectMaxArgN(t *testing.T) {
	tests := []struct {
		name    string
		node    parser.Node
		maxArgs int
		wantErr bool
	}{
		{
			name: "within limit",
			node: parser.Node{
				Name: "directive",
				Args: []string{"arg1", "arg2"},
			},
			maxArgs: 3,
			wantErr: false,
		},
		{
			name: "at limit",
			node: parser.Node{
				Name: "directive",
				Args: []string{"arg1", "arg2", "arg3"},
			},
			maxArgs: 3,
			wantErr: false,
		},
		{
			name: "exceeds limit",
			node: parser.Node{
				Name: "directive",
				Args: []string{"arg1", "arg2", "arg3", "arg4"},
			},
			maxArgs: 3,
			wantErr: true,
		},
		{
			name: "no args allowed",
			node: parser.Node{
				Name: "directive",
				Args: []string{"arg1"},
			},
			maxArgs: 0,
			wantErr: true,
		},
		{
			name: "no args provided",
			node: parser.Node{
				Name: "directive",
				Args: []string{},
			},
			maxArgs: 0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ExpectMaxArgN(tt.node, tt.maxArgs)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExpectMaxArgN() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err != nil {
				// Check that error message contains expected information
				errStr := err.Error()
				if !strings.Contains(errStr, "expected at most") {
					t.Errorf("ExpectMaxArgN() error message should contain 'expected at most', got: %s", errStr)
				}
				if !strings.Contains(errStr, tt.node.Name) {
					t.Errorf("ExpectMaxArgN() error message should contain directive name '%s', got: %s", tt.node.Name, errStr)
				}
			}
		})
	}
}

func TestExpectMinArgN(t *testing.T) {
	tests := []struct {
		name    string
		node    parser.Node
		minArgs int
		wantErr bool
	}{
		{
			name: "above minimum",
			node: parser.Node{
				Name: "directive",
				Args: []string{"arg1", "arg2", "arg3"},
			},
			minArgs: 2,
			wantErr: false,
		},
		{
			name: "at minimum",
			node: parser.Node{
				Name: "directive",
				Args: []string{"arg1", "arg2"},
			},
			minArgs: 2,
			wantErr: false,
		},
		{
			name: "below minimum",
			node: parser.Node{
				Name: "directive",
				Args: []string{"arg1"},
			},
			minArgs: 2,
			wantErr: true,
		},
		{
			name: "no args when required",
			node: parser.Node{
				Name: "directive",
				Args: []string{},
			},
			minArgs: 1,
			wantErr: true,
		},
		{
			name: "no args when none required",
			node: parser.Node{
				Name: "directive",
				Args: []string{},
			},
			minArgs: 0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ExpectMinArgN(tt.node, tt.minArgs)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExpectMinArgN() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err != nil {
				// Check that error message contains expected information
				errStr := err.Error()
				if !strings.Contains(errStr, "expected at least") {
					t.Errorf("ExpectMinArgN() error message should contain 'expected at least', got: %s", errStr)
				}
				if !strings.Contains(errStr, tt.node.Name) {
					t.Errorf("ExpectMinArgN() error message should contain directive name '%s', got: %s", tt.node.Name, errStr)
				}
			}
		})
	}
}

// Test with environment variables (if supported by underlying parser)
func TestReadFileWithEnvVars(t *testing.T) {
	// Set up test environment variables
	os.Setenv("SERVER_NAME", "test-server")
	os.Setenv("PORT", "8080")
	os.Setenv("LOG_DIR", "/tmp/logs")
	defer func() {
		os.Unsetenv("SERVER_NAME")
		os.Unsetenv("PORT")
		os.Unsetenv("LOG_DIR")
	}()

	ast, err := ReadFile("testdata/with_env_vars.conf")
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if len(ast) == 0 {
		t.Fatal("ReadFile() returned empty AST")
	}

	// The AST should contain the parsed nodes
	// Note: Environment variable expansion happens at the parser level,
	// so we're mainly testing that the file can be parsed successfully
	t.Logf("Successfully parsed config with %d top-level nodes", len(ast))
}

// Test with macros (if supported by underlying parser)
func TestReadFileWithMacros(t *testing.T) {
	ast, err := ReadFile("testdata/with_macros.conf")
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if len(ast) == 0 {
		t.Fatal("ReadFile() returned empty AST")
	}

	// Check that we have the expected number of top-level nodes
	expectedNodes := 4 // 2 server blocks, 1 tls_files, 1 log_file
	if len(ast) != expectedNodes {
		t.Errorf("Expected %d top-level nodes, got %d", expectedNodes, len(ast))
	}

	// Find the first server block and verify macro expansion occurred
	var serverNode *parser.Node
	for i := range ast {
		if ast[i].Name == "server" {
			serverNode = &ast[i]
			break
		}
	}

	if serverNode == nil {
		t.Fatal("No server block found in parsed AST")
	}

	// Verify that the server block has arguments (macro should be expanded)
	if len(serverNode.Args) == 0 {
		t.Error("Server block should have arguments after macro expansion")
	}

	// The AST should contain the parsed nodes with macro expansion
	// Note: Macro expansion happens at the parser level,
	// so we're mainly testing that the file can be parsed successfully
	t.Logf("Successfully parsed config with macros - %d top-level nodes", len(ast))
}

// Test AST type conversion
func TestASTConversion(t *testing.T) {
	content := "directive value\nblock { nested value }"
	reader := strings.NewReader(content)
	
	ast, err := Read(reader, "test.conf")
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	// Test that AST can be used as []parser.Node
	nodes := []parser.Node(ast)
	if len(nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(nodes))
	}

	// Test first node
	if nodes[0].Name != "directive" {
		t.Errorf("Expected first node name 'directive', got '%s'", nodes[0].Name)
	}
	if len(nodes[0].Args) != 1 || nodes[0].Args[0] != "value" {
		t.Errorf("Expected first node args ['value'], got %v", nodes[0].Args)
	}

	// Test second node
	if nodes[1].Name != "block" {
		t.Errorf("Expected second node name 'block', got '%s'", nodes[1].Name)
	}
	if len(nodes[1].Children) != 1 {
		t.Errorf("Expected second node to have 1 child, got %d", len(nodes[1].Children))
	}
}