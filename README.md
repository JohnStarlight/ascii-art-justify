# ASCII Art CLI

A simple Go command-line tool that turns text into ASCII art.

## What This Project Does

- Takes one text argument from the terminal
- Converts each character to ASCII art (8 lines tall)
- Lets you choose one of 3 styles:
  - `standard`
  - `shadow`
  - `thinkertoy`
- Optionally writes the output to a file using `--output=<filename>`

## Quick Start

From project root:

```bash
go run ./cmd "Hello"
```

```bash
go run ./cmd "Hello" shadow
```

```bash
go run ./cmd --output=result.txt "Hello" shadow
```

Important:
Always pass your full input inside quotes (`"..."`) for correct results.

## Usage Example

```bash
go run ./cmd "Hello There" thinkertoy
```

Prints the ASCII art to the terminal.

```bash
go run ./cmd --output=out.txt "Hello There" standard
```

Writes the ASCII art to `out.txt` instead of printing to the terminal. The file is created (or overwritten) automatically.

## Input Rules

- You must pass exactly one argument
- Only printable ASCII characters are accepted (`32` to `126`)
- Non-ASCII characters (for example `é` or emoji) are rejected
- The sequence `\n` (backslash + n) is used to create new lines in the output

Example:

```bash
go run ./cmd "Hello\nThere"

## Run Tests

1. Run all tests:

```bash
go test ./...
```

Runs all test files in all project packages.

2. Run all tests with detailed output:

```bash
go test ./... -v
```

Shows `PASS/FAIL` for each test and subtest.

3. Force fresh run (no cache):

```bash
go test ./... -v -count=1
```

Forces tests to run again even if previous results were cached.

## Project Structure

- `cmd/main.go` - CLI entrypoint and argument handling
- `internal/config.go` - argument parsing and validation
- `internal/printascii.go` - ASCII rendering logic
- `internal/output.go` - output file creation
- `banners/*.txt` - banner templates
- `test/printascii_test.go` - core unit tests
- `test/audit_examples_test.go` - audit/instruction sample tests

## Known Limitations

- Supports printable ASCII characters only (`32` to `126`)
- Unicode characters (for example Greek letters or emoji) are rejected
- When using `--output`, the banner argument is required (no default when the flag is present)

## License

MIT. See `LICENSE`.
