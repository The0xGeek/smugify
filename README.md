# Smugify

A lightweight Go-based CLI utility for generating SVG containers with embedded binary payloads for controlled research, rendering analysis, and browser behavior testing.

## Features

* Minimal and fast Go implementation
* Embed arbitrary files as Base64 payloads
* Generate standalone SVG containers
* Inject payload logic into existing SVG templates
* Simple CLI interface
* No external runtime dependencies
* Open-source and easy to extend

## Installation

#### Clone Repository

```bash
git clone https://github.com/The0xGeek/smugify.git
cd smugify
```

#### Build

```bash
go build -o smugify
```

Or:

```bash
go build -o smugify cmd/smugify/main.go
```

## Usage

#### Basic Usage

```bash
./smugify -attach payload.pdf
```

This generates:

```text
output.svg
```

#### Using a Trusted SVG Template

```bash
./smugify -attach payload.pdf -trust template.svg
```

This embeds the generated payload logic into an existing SVG file.

## CLI Arguments

| Flag      | Description                                |
| --------- | ------------------------------------------ |
| `-attach` | File to encode and embed                   |
| `-trust`  | Existing SVG template file                 |
| `-out`    | Download filename presented to the browser |

## Example

```bash
./smugify \
  -attach invoice.pdf \
  -trust animation.svg \
  -out report.pdf
```

## Architecture

Smugify follows a simple processing pipeline:

```text
Input File
    ↓
Base64 Encoding
    ↓
SVG Template Rendering
    ↓
Optional Trusted SVG Injection
    ↓
Final SVG Output
```

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