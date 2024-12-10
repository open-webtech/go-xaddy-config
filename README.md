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
