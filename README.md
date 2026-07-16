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
- Optionally colors the output using one or more `--color=<color>` flags (whole text, one substring, or several substrings each in its own color)
- Optionally aligns the output inside the terminal using `--align=<type>` (`left`, `center`, `right`, `justify`)

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

```bash
go run ./cmd --color=red "Hello"
```

```bash
go run ./cmd --color=red "ell" "Hello"
```

```bash
go run ./cmd --align=center "Hello" standard
```

```bash
go run ./cmd --align=justify "Hello There" standard
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

## Color Support

Use `--color=<color>` to colorize the rendered output.

```bash
go run ./cmd --color=red "Hello"
```

Colors the entire `Hello` in red.

```bash
go run ./cmd --color=red "ell" "Hello"
```

Colors only the `ell` substring (every occurrence) in red; the rest stays uncolored.

```bash
go run ./cmd "--color=rgb(0,200,255)" "Hello" shadow
```

Colors `Hello` using a custom RGB color, rendered in the `shadow` banner.

```bash
go run ./cmd "--color=#ff8800" "Hello"
```

Colors `Hello` using a custom hex color (`#RRGGBB`).

Supported color names: `red`, `green`, `blue`, `yellow`, `magenta`/`purple`, `cyan`, `orange`.
Custom colors: `rgb(R,G,B)` (e.g. `rgb(255,128,0)`) or hex `#RRGGBB` (e.g. `#ff8800`).

Important: when using `rgb(...)` or `#RRGGBB`, wrap the whole flag in quotes (`"--color=rgb(R,G,B)"`, `"--color=#ff8800"`).
Parentheses and `#` are special characters in the shell and will break the command otherwise.

### Multiple Colors

You can pass `--color` more than once. Each flag pairs with its own substring, in the same order:

```bash
go run ./cmd --color=red --color=blue "He" "llo" "Hello"
```

Colors `He` in red and `llo` in blue. With two or more `--color` flags, every color requires its own substring. If two substrings overlap, the first `--color` flag wins.

`--color`, `--output` and `--align` can be combined and used in any order.

## Alignment Support

Use `--align=<type>` to position the rendered output inside the terminal:

```bash
go run ./cmd --align=right "Hello" standard
go run ./cmd --align=center "Hello" standard
go run ./cmd --align=justify "Hello There friend" standard
```

- `left` — default behavior, output starts at the left edge
- `center` — output is centered in the terminal
- `right` — output is pushed to the right edge
- `justify` — spaces between words are stretched so the text spans the full terminal width

The flag must be written exactly as `--align=<type>` (`--align center` is rejected).
The terminal width is detected via `stty size`, then `tput cols`, then the `COLUMNS` environment variable, with a final fallback of 80 columns.

`justify` needs at least two words; a single word (or a line wider than the terminal) is printed unchanged. Alignment works together with `--color` — colored substrings stay correctly colored when words are spread apart.

## Input Rules

- You must pass exactly one argument
- Only printable ASCII characters are accepted (`32` to `126`)
- Non-ASCII characters (for example `é` or emoji) are rejected
- The sequence `\n` (backslash + n) is used to create new lines in the output

Example:

```bash
go run ./cmd "Hello\nThere"
```

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
- `internal/align.go` - alignment and justify logic
- `internal/terminal.go` - terminal width detection
- `internal/color.go` - color name/RGB/hex parsing for `--color`
- `internal/output.go` - output file creation
- `banners/*.txt` - banner templates
- `test/printascii_test.go` - core unit tests
- `test/audit_examples_test.go` - audit/instruction sample tests
- `test/color_test.go` - color flag and rendering tests
- `test/align_test.go` - alignment and justify tests

## Known Limitations

- Supports printable ASCII characters only (`32` to `126`)
- Unicode characters (for example Greek letters or emoji) are rejected
- When using `--output`, the banner argument is required (no default when the flag is present)
- Colored substrings must match literally (no regex/wildcards)
- ANSI color codes are written even when output goes to a file via `--output` (no automatic stripping)
- Alignment uses the width of the terminal at run time, even when writing to a file with `--output`
- If the terminal is narrower than the rendered text, alignment leaves the output unchanged

## License

MIT. See `LICENSE`.
