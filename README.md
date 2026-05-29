# Smugify

A lightweight Go-based CLI tool for generating SVG containers with embedded JavaScript payloads and optional code transformation pipeline for controlled research, rendering analysis, and browser behavior testing.

## Features

* Simple CLI interface
* Minimal and fast Go implementation
* Embed arbitrary files as Base64 payloads
* Generate standalone SVG containers
* Inject payload logic into existing SVG templates
* JavaScript transformation pipeline (obfuscation layer)
* Comment stripping, string encoding, hex escaping
* No external runtime dependencies
* Open-source and easy to extend

## Installation

### Option 1 — Install via Go (Recommended)

If you have Go installed (1.20+ recommended):

```bash
go install github.com/The0xGeek/smugify@latest
```

After installation, the binary will be available in:
```bash
$(go env GOPATH)/bin
```

Make sure this line is in your PATH:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Option 2 — Build from source

```bash
git clone https://github.com/The0xGeek/smugify.git
cd smugify
go build -o smugify
```
or:
```bash
go build -o smugify cmd/smugify/main.go
```

## Usage

You can use `payload.pdf` and `template.svg` files just for testing, or use your custom payloads.

#### Basic Usage

```bash
smugify -attach payload.pdf
```

This generates:

```text
output.svg
```

#### Advanced Usage:

```bash
./smugify -attach payload.pdf \
    -trust template.svg \
    - obfuscate \
    - hex-escape
    - advanced
```

This embeds the generated payload logic into an existing SVG file with an obfuscation layer on javascript payload.

## CLI Flags

#### Core Flags:

| Flag      | Description                                |
| --------- | ------------------------------------------ |
| `-attach` | File to encode and embed                   |
| `-trust`  | Existing SVG template file                 |
| `-out`    | Download filename presented to the browser |

#### JavaScript Transformation Pipeline:

| Flag      | Description                                |
| --------- | ------------------------------------------ |
| `-obfuscate` | Enable JS transformation pipeline                   |
| `-hex-escape`  | Convert string literals to hex escape sequences                 |
| `-remove-comments`    | Strip single-line and multi-line comments |
| `-minify`    | Minify JavaScript output |
| `-advanced`    | Enable advanced multi-layer encoding (eval-based loader) |

## Security Notice

This project is intended for:

* Browser rendering research
* Defensive testing
* SVG parsing analysis
* Controlled security labs
* Educational purposes

Do not use this project against systems or environments without explicit authorization.

## Contributing

Pull requests, issue reports, and architecture suggestions are welcome.

Before contributing:

1. Format code with `go fmt`
2. Keep functions small and testable
3. Prefer idiomatic Go patterns
4. Avoid unnecessary dependencies

<br>
Developed by The0xGeek.