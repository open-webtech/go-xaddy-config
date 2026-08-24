package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/open-webtech/go-xaddy-config/ast"
)

func read(t *testing.T, content string) []ast.Node {
	t.Helper()
	nodes, err := Read(strings.NewReader(content), "test.conf")
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	return nodes
}

func TestReadSyntax(t *testing.T) {
	t.Run("directive with args", func(t *testing.T) {
		nodes := read(t, "directive arg1 arg2")
		if len(nodes) != 1 || nodes[0].Name != "directive" {
			t.Fatalf("unexpected nodes: %+v", nodes)
		}
		if len(nodes[0].Args) != 2 || nodes[0].Args[0] != "arg1" || nodes[0].Args[1] != "arg2" {
			t.Errorf("unexpected args: %v", nodes[0].Args)
		}
	})

	t.Run("multi-line block", func(t *testing.T) {
		nodes := read(t, "block name {\n\tchild value\n\tother {\n\t\tinner 1\n\t}\n}")
		if len(nodes) != 1 || len(nodes[0].Children) != 2 {
			t.Fatalf("unexpected nodes: %+v", nodes)
		}
		inner := nodes[0].Children[1]
		if inner.Name != "other" || len(inner.Children) != 1 || inner.Children[0].Name != "inner" {
			t.Errorf("unexpected nested block: %+v", inner)
		}
	})

	t.Run("single-line block", func(t *testing.T) {
		nodes := read(t, "block { child value }")
		if len(nodes) != 1 || len(nodes[0].Children) != 1 {
			t.Fatalf("unexpected nodes: %+v", nodes)
		}
		child := nodes[0].Children[0]
		if child.Name != "child" || len(child.Args) != 1 || child.Args[0] != "value" {
			t.Errorf("unexpected child: %+v", child)
		}
	})

	t.Run("quoted arguments", func(t *testing.T) {
		nodes := read(t, `directive "two words" plain`)
		if len(nodes[0].Args) != 2 || nodes[0].Args[0] != "two words" {
			t.Errorf("unexpected args: %v", nodes[0].Args)
		}
	})

	t.Run("comments", func(t *testing.T) {
		nodes := read(t, "# comment\ndirective value # trailing\n# another")
		if len(nodes) != 1 || nodes[0].Name != "directive" {
			t.Fatalf("unexpected nodes: %+v", nodes)
		}
	})

	t.Run("file and line positions", func(t *testing.T) {
		nodes := read(t, "first\n\nblock {\n\tchild\n}")
		if nodes[0].File != "test.conf" || nodes[0].Line != 1 {
			t.Errorf("first node position = %s:%d", nodes[0].File, nodes[0].Line)
		}
		if nodes[1].Line != 3 || nodes[1].Children[0].Line != 4 {
			t.Errorf("block positions = %d, child %d", nodes[1].Line, nodes[1].Children[0].Line)
		}
	})

	t.Run("unclosed block", func(t *testing.T) {
		if _, err := Read(strings.NewReader("block {\n\tchild"), "test.conf"); err == nil {
			t.Error("expected error for unclosed block")
		}
	})

	t.Run("stray closing brace", func(t *testing.T) {
		if _, err := Read(strings.NewReader("directive value\n}"), "test.conf"); err == nil {
			t.Error("expected error for stray '}'")
		}
	})

	t.Run("brace instead of name", func(t *testing.T) {
		if _, err := Read(strings.NewReader("directive\n{\n}"), "test.conf"); err == nil {
			t.Error("expected error for '{' in name position")
		}
	})
}

func TestSnippetsAndImports(t *testing.T) {
	t.Run("snippet import", func(t *testing.T) {
		nodes := read(t, "(common) {\n\tlog_level debug\n\ttimeout 30\n}\nserver {\n\timport common\n\tlisten 443\n}")
		if len(nodes) != 1 {
			t.Fatalf("snippet definition should not appear in output, got %d nodes", len(nodes))
		}
		children := nodes[0].Children
		if len(children) != 3 || children[0].Name != "log_level" || children[2].Name != "listen" {
			t.Errorf("unexpected expansion: %+v", children)
		}
	})

	t.Run("snippet used twice is independent", func(t *testing.T) {
		nodes := read(t, "(c) {\n\topt x\n}\na {\n\timport c\n}\nb {\n\timport c\n}")
		nodes[0].Children[0].Args[0] = "changed"
		if nodes[1].Children[0].Args[0] != "x" {
			t.Error("snippet expansions share state")
		}
	})

	t.Run("file import", func(t *testing.T) {
		dir := t.TempDir()
		imported := filepath.Join(dir, "base.conf")
		if err := os.WriteFile(imported, []byte("base_setting value1\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		nodes := read(t, "import "+imported+"\nmain {\n\tlisten 80\n}")
		if len(nodes) != 2 || nodes[0].Name != "base_setting" {
			t.Fatalf("unexpected nodes: %+v", nodes)
		}
		if nodes[0].File != imported {
			t.Errorf("imported node should keep its source file, got %s", nodes[0].File)
		}
	})

	t.Run("snippet takes precedence over file", func(t *testing.T) {
		dir := t.TempDir()
		name := filepath.Join(dir, "shared.conf")
		if err := os.WriteFile(name, []byte("from_file\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		nodes := read(t, "("+name+") {\n\tfrom_snippet\n}\nimport "+name)
		if len(nodes) != 1 || nodes[0].Name != "from_snippet" {
			t.Errorf("expected snippet to win over file, got %+v", nodes)
		}
	})

	t.Run("import cycle", func(t *testing.T) {
		dir := t.TempDir()
		a := filepath.Join(dir, "a.conf")
		b := filepath.Join(dir, "b.conf")
		if err := os.WriteFile(a, []byte("import "+b+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(b, []byte("import "+a+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := Read(strings.NewReader("import "+a), "test.conf"); err == nil {
			t.Error("expected error for import cycle")
		}
	})

	t.Run("unknown import target", func(t *testing.T) {
		_, err := Read(strings.NewReader("import no_such_snippet"), "test.conf")
		if err == nil {
			t.Error("expected error for unknown import target")
		}
	})
}

func TestMacros(t *testing.T) {
	t.Run("basic and nested definitions", func(t *testing.T) {
		nodes := read(t, "$(host) = mail.example.org\n$(domains) = $(host) example.org\nserver $(host) {\n\tdomains $(domains)\n}")
		if len(nodes) != 1 {
			t.Fatalf("macro definitions should not appear in output, got %d nodes", len(nodes))
		}
		if nodes[0].Args[0] != "mail.example.org" {
			t.Errorf("macro not expanded in args: %v", nodes[0].Args)
		}
		domains := nodes[0].Children[0]
		if len(domains.Args) != 2 || domains.Args[0] != "mail.example.org" || domains.Args[1] != "example.org" {
			t.Errorf("multi-value macro not spliced: %v", domains.Args)
		}
	})

	t.Run("substring expansion", func(t *testing.T) {
		nodes := read(t, "$(host) = example.org\ntls_cert /etc/ssl/$(host)/cert.pem")
		if nodes[0].Args[0] != "/etc/ssl/example.org/cert.pem" {
			t.Errorf("substring macro not expanded: %v", nodes[0].Args)
		}
	})

	t.Run("undefined reference left unchanged", func(t *testing.T) {
		nodes := read(t, "directive $(missing)")
		if nodes[0].Args[0] != "$(missing)" {
			t.Errorf("undefined macro should stay as-is: %v", nodes[0].Args)
		}
	})
}

func TestEnvExpansion(t *testing.T) {
	t.Setenv("XADDY_TEST_HOST", "test-host")

	t.Run("defined variable", func(t *testing.T) {
		nodes := read(t, "server {env:XADDY_TEST_HOST}")
		if nodes[0].Args[0] != "test-host" {
			t.Errorf("env var not expanded: %v", nodes[0].Args)
		}
	})

	t.Run("inside longer string", func(t *testing.T) {
		nodes := read(t, `log_file "{env:XADDY_TEST_HOST}/server.log"`)
		if nodes[0].Args[0] != "test-host/server.log" {
			t.Errorf("env var not expanded in string: %v", nodes[0].Args)
		}
	})

	t.Run("undefined variable becomes empty but stays", func(t *testing.T) {
		nodes := read(t, "listen {env:XADDY_TEST_UNDEFINED}")
		if len(nodes[0].Args) != 1 || nodes[0].Args[0] != "" {
			t.Errorf("undefined env var should stay as empty arg: %q", nodes[0].Args)
		}
	})

	t.Run("incomplete placeholder untouched", func(t *testing.T) {
		nodes := read(t, "directive {env:INCOMPLETE")
		if nodes[0].Args[0] != "{env:INCOMPLETE" {
			t.Errorf("incomplete placeholder should stay as-is: %v", nodes[0].Args)
		}
	})

	t.Run("expansion applies to macro values", func(t *testing.T) {
		nodes := read(t, "$(base) = {env:XADDY_TEST_HOST}\nlog_dir /var/log/$(base)")
		if nodes[0].Args[0] != "/var/log/test-host" {
			t.Errorf("env var in macro value not expanded: %v", nodes[0].Args)
		}
	})
}
