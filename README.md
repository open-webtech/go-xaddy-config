# go-xaddy-config

A configuration parser and evaluator for [Caddyfile](https://caddyserver.com/docs/caddyfile)-style configuration formats, specifically implementing the subset used in [Maddy](https://github.com/foxcpp/maddy) (Caddy for Mail). This provides a familiar and proven configuration syntax for Go applications.

## Features

- **Caddyfile-style Syntax:** Uses the same intuitive configuration format as Caddy and Maddy
- **Block Structure:** Support for nested configuration blocks and directives
- **Schema-based Configuration:** Define and validate configuration structures using a flexible schema system
- **Type Safety:** Strong type checking for configuration values
- **Error Handling:** Comprehensive error reporting for configuration issues

## Configuration Format

The configuration follows the Caddyfile format style, as used in Maddy:

```caddyfile
# Basic directive
simple_directive value

# Block directive
block_name {
    nested_directive value
    another_setting value
}

# Multiple values
directive value1 value2 {
    nested_option
}
```

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

# Multiple servers can use the same macros
server backup.$(primary_domain) {
    listen 587
    tls_cert /etc/ssl/$(hostname)/cert.pem
    domains $(local_domains)
    backup_for $(hostname)
}
```

### Environment Variables

The parser supports environment variable substitution using `{env:VARIABLE}` syntax:

```caddyfile
# Environment variable expansion
server {env:SERVER_NAME} {
    listen {env:PORT}
    tls_cert {env:TLS_CERT_PATH}
    tls_key {env:TLS_KEY_PATH}
    
    # Variables work inside quotes too
    log_file "{env:LOG_DIR}/server.log"
    
    # Default fallback if variable is undefined (expands to empty string)
    debug_mode {env:DEBUG_ENABLED}
}

# Multiple environment variables in one directive
database_url {env:DB_HOST}:{env:DB_PORT}

# Combining with snippets
(ssl_config) {
    cert_file {env:SSL_CERT_FILE}
    key_file {env:SSL_KEY_FILE}
    min_version 1.2
}
```

### Snippets

Define reusable configuration blocks at the top level using parentheses:

```caddyfile
# Define a snippet for common TLS settings
(common_tls) {
    tls_min_version 1.2
    tls_max_version 1.3
    cert_file /etc/ssl/certs/server.crt
    key_file /etc/ssl/private/server.key
}

# Define a snippet for logging configuration
(debug_logging) {
    log_level debug
    log_file /var/log/app-debug.log
    verbose_errors true
}
```

### Imports

Reference snippets or include external configuration files:

```caddyfile
# Import a predefined snippet
server web {
    listen 443
    import common_tls
    
    # Other server-specific settings
    document_root /var/www/html
}

# Import from external configuration files
import /etc/app/database.conf
import tls-settings.conf

# Multiple servers can reuse the same snippet
server api {
    listen 8443  
    import common_tls
    import debug_logging
    
    api_endpoint /v1
}
```

## Installation

To install `go-xaddy-config`, use:

```bash
go get -u github.com/open-webtech/go-xaddy-config
```

## Usage

Import the library in your Go code:

```go
import "github.com/open-webtech/go-xaddy-config"
```

### Reading Configuration Files

```go
cfgNodes, err := config.ReadFile("path/to/config/file")
if err != nil {
    // handle error
}
```

### Defining Configuration Schema

The schema builder allows you to define your configuration structure using directives and blocks:

```go
// Create a new schema builder
root := schema.NewBuilder()

// Define simple directives with arguments
root.DefineDirective("log_level", args.StringArg(&cfg.LogLevel))
root.DefineDirective("max_connections", args.IntArg(&cfg.MaxConnections))

// Define directives with argument validation
root.DefineDirectiveCallback("forward", func(node parser.Node) error {
    if err := config.ExpectMinArgN(node, 1); err != nil {
        return err  // Ensures at least one target address is provided
    }
    if err := config.ExpectMaxArgN(node, 3); err != nil {
        return err  // Limits to max 3 target addresses
    }
    cfg.ForwardAddresses = node.Args
    return nil
})

// Define blocks with nested configuration
block := root.DefineBlock("server", args.StringArg(&cfg.ServerName))
block.DefineDirective("listen", args.StringArg(&cfg.Listen))
block.DefineDirective("tls", args.BoolArg(&cfg.TLS))
```

### Evaluating Configuration

Once you have defined your schema, you can evaluate configuration nodes:

```go
nodes, err := config.ReadFile("config.conf")
if err != nil {
    log.Fatal(err)
}

err = root.EvaluateTree(nodes, cfg)
if err != nil {
    log.Fatal(err)
}
```

## Development

### Generating Code Files

The project uses code generation for certain files:

- `schema/values/values_generated.go`: Value parsers and accumulators for basic types
- `schema/args/args_generated.go`: Argument type definitions and helpers

To generate all files at once, run the following command:

```bash
go generate ./...
```

To generate individual files, run one of the following commands:

```bash
go generate ./schema/values  # for values_generated.go
go generate ./schema/args    # for args_generated.go
```

The underlying generate commands are:

```bash
go run cmd/gen_values/main.go -pkg values  # for values_generated.go
go run cmd/gen_values/main.go -pkg args    # for args_generated.go
```

The generator requires:

- Go 1.16 or later
- The `goimports` tool from [golang.org/x/tools](https://pkg.go.dev/golang.org/x/tools/cmd/goimports) to add the necessary imports to the generated files

Do not modify the generated files directly as changes will be lost when regenerating.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License – see the LICENSE file for details.