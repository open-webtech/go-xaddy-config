# go-xaddy-config - Information for Coding Agents

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is **go-xaddy-config**, a Go library that implements a configuration parser and evaluator for Caddyfile-style configuration formats, specifically implementing the subset used in [Maddy](https://github.com/foxcpp/maddy) (Caddy for Mail). The library provides a familiar and proven configuration syntax for Go applications using the same format as Caddy and Maddy servers.

## Key Dependencies

- `github.com/foxcpp/maddy v0.7.1` - Core configuration parser from Maddy project
- `golang.org/x/text v0.14.0` - Text processing utilities for code generation
- Go 1.23.2+ required

## Architecture Overview

The library is structured around these core components:

### Configuration Parsing (`config.go`)
- **`config.Read()`** and **`config.ReadFile()`**: Parse configuration files into AST (Abstract Syntax Tree)
- **`config.ExpectMinArgN()`** and **`config.ExpectMaxArgN()`**: Utility functions for argument validation

### Schema System (`schema/`)
- **`schema.Builder`**: Main interface for defining configuration schemas
- **`nodes/`**: Node definitions, evaluators, and error handling for configuration blocks and directives
- **`args/`**: Argument definitions and type handling (includes generated code)
- **`values/`**: Value parsers and accumulators for basic types (includes generated code)

### Code Generation (`cmd/gen_values/`)
- Generates type-safe value parsers and argument helpers from `schema/values/values.json`
- Two generated files: `schema/values/values_generated.go` and `schema/args/args_generated.go`

## Common Development Commands

### Code Generation
Generated files are critical to the library's type system. Always regenerate after modifying `schema/values/values.json`:

```bash
# Generate all files
go generate ./...

# Generate individual packages
go generate ./schema/values  # for values_generated.go
go generate ./schema/args    # for args_generated.go

# Direct generator commands
go run cmd/gen_values/main.go -pkg values  # for values_generated.go
go run cmd/gen_values/main.go -pkg args    # for args_generated.go
```

### Testing and Building
```bash
# Run tests (standard Go testing)
go test ./...

# Build the library
go build ./...

# Install dependencies
go mod tidy
```

## Configuration Format

The library parses Caddyfile-style configuration with support for:
- Simple directives: `directive value`
- Block structures: `block_name { nested_directive value }`
- Multiple arguments: `directive arg1 arg2 { nested_option }`

### Snippets and Imports

The parser supports Maddy's snippet and import functionality for configuration reuse:

**Snippets** - Reusable configuration blocks defined at the top level:
```caddyfile
(snippet_name) {
    directive value
    nested_block {
        option setting
    }
}
```

**Imports** - Reference snippets or include external files:
```caddyfile
# Import a predefined snippet
import snippet_name

# Import from external configuration file
import tls.conf
import /path/to/config.conf
```

**Import behavior**:
- Imported content is directly expanded into the configuration
- File paths can be relative or absolute
- If both a snippet and file have the same name, the snippet takes precedence
- Enables configuration modularization and reuse across multiple files

### Environment Variables

The parser supports environment variable substitution using `{env:VARIABLENAME}` syntax:

```caddyfile
# Environment variable expansion
server {env:SERVER_NAME} {
    listen {env:PORT}
    tls_cert {env:TLS_CERT_PATH}
    log_file "{env:LOG_DIR}/server.log"
    debug_mode {env:DEBUG_ENABLED}
}

# Multiple variables in one directive
database_url {env:DB_HOST}:{env:DB_PORT}
```

**Environment variable behavior**:
- **Syntax**: `{env:VARIABLENAME}` references environment variables
- **Undefined variables**: Expand to empty strings but remain in argument list
- **Quote expansion**: Variables expand inside quoted strings
- **Partial syntax**: Incomplete placeholders like `{env:VAR` are left unchanged
- **Case sensitivity**: Variable names are case-sensitive
- **Security**: Keeps sensitive values out of configuration files

### Macros

The parser supports macros with `$(...)` syntax for defining reusable configuration variables:

```caddyfile
# Macro definitions at top level
$(hostname) = mail.example.org
$(primary_domain) = example.org
$(local_domains) = $(primary_domain)

# Using macros throughout configuration
server $(hostname) {
    listen 25
    tls_cert /etc/ssl/$(hostname)/cert.pem
    domains $(local_domains)
}
```

**Macro behavior**:
- **Definition**: `$(name) = value1 value2...` defines a macro at top-level only
- **Reference**: `$(name)` expands to the defined values anywhere in the configuration
- **Substitution**: Macros are expanded before configuration parsing
- **Nesting**: Macros can reference other macros (like `$(local_domains)` referencing `$(primary_domain)`)
- **Use cases**: Centralize domain names, hostnames, and other repeated configuration values

## Code Generation Requirements

The code generator requires:
- **`goimports`** tool from golang.org/x/tools for import management
- The generator reads `schema/values/values.json` to create type-safe parsers
- Generated files should NEVER be modified directly as changes will be lost

## Important Implementation Notes

- The library uses Maddy's configuration parser as the underlying parsing engine
- Schema builders define the structure and validation rules for configurations  
- Node definitions handle both directives (simple config lines) and blocks (nested config sections)
- Value types are defined in JSON and generated into Go code for type safety
- Error reporting includes file location information from the parsed configuration

## File Organization

- `config.go` - Main parsing interface and utilities
- `schema/builder.go` - Schema definition entry point  
- `schema/nodes/` - Node definition system (directives, blocks, containers)
- `schema/values/` - Value parsing and accumulation (includes generated code)
- `schema/args/` - Argument definitions (includes generated code)
- `cmd/gen_values/` - Code generator for type-safe parsers