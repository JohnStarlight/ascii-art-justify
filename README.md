# ASCII Art CLI

A simple Go command-line tool that turns text into ASCII art.

## What This Project Does

- Takes the text to render from the terminal (plus optional substrings when coloring)
- Converts each character to ASCII art (8 lines tall)
- Lets you choose one of 3 styles (the banner name is optional and defaults to `standard`):
  - `standard`
  - `shadow`
  - `thinkertoy`
- Optionally writes the output to a file using `--output=<filename>`
- Optionally colors the output using one or more `--color=<color>` flags (whole text, one substring, or several substrings each in its own color)
- Optionally aligns the output inside the terminal using `--align=<type>` (`left`, `center`, `right`, `justify`)

## Command Shape

```
go run ./cmd [--color=<color>]... [--output=<file>] [--align=<type>] [SUBSTRING]... [STRING] [BANNER]
```

- `[STRING]` is required; `[BANNER]` is optional and defaults to `standard`
- Flags may appear anywhere in the argument list — before, between, or after the positional arguments
- Flag names, the `--align` value and the banner name are case-insensitive (`--ALIGN=CENTER "Hi" STANDARD` works)
- Positional order matters: substrings come first, then the string, then the banner
- Errors are printed to `stderr` as `Error: ...` and the program exits with status `1`

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

`--color`, `--output` and `--align` can be combined and used in any order — see [Combining Flags](#combining-flags) for an example of every combination.

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

## Combining Flags

`--color`, `--output` and `--align` are independent and can be mixed freely, with or without a banner. The table below lists every possible combination, followed by a runnable example of each.

| # | `--color` | `--output` | `--align` | Banner |
|---|---|---|---|---|
| 1 | – | – | – | – |
| 2 | – | – | – | yes |
| 3 | – | yes | – | – |
| 4 | – | yes | – | yes |
| 5 | – | – | yes | – |
| 6 | – | – | yes | yes |
| 7 | – | yes | yes | – |
| 8 | – | yes | yes | yes |
| 9 | one (whole text) | – | – | – |
| 10 | one (substring) | – | – | – |
| 11 | one (substring) | – | – | yes |
| 12 | one | yes | – | – |
| 13 | one | yes | – | yes |
| 14 | one | – | yes | – |
| 15 | one | – | yes | yes |
| 16 | one | yes | yes | – |
| 17 | one | yes | yes | yes |
| 18 | many | – | – | – |
| 19 | many | – | – | yes |
| 20 | many | yes | – | – |
| 21 | many | yes | – | yes |
| 22 | many | – | yes | – |
| 23 | many | – | yes | yes |
| 24 | many | yes | yes | – |
| 25 | many | yes | yes | yes |

### No color

```bash
# 1 — plain text, default standard banner
go run ./cmd "Hello There"

# 2 — plain text with an explicit banner
go run ./cmd "Hello There" shadow

# 3 — write to a file
go run ./cmd --output=result.txt "Hello There"

# 4 — write to a file with a banner
go run ./cmd --output=result.txt "Hello There" thinkertoy

# 5 — align only
go run ./cmd --align=center "Hello There"

# 6 — align with a banner
go run ./cmd --align=right "Hello There" shadow

# 7 — align and write to a file
go run ./cmd --output=result.txt --align=justify "Hello There"

# 8 — align, file and banner
go run ./cmd --output=result.txt --align=center "Hello There" thinkertoy
```

### One `--color`

With a single `--color`, the substring is optional: omit it and the whole string is colored.

```bash
# 9 — color the whole string
go run ./cmd --color=red "Hello There"

# 10 — color one substring only
go run ./cmd --color=red "ell" "Hello There"

# 11 — color one substring, with a banner
go run ./cmd --color=green "There" "Hello There" shadow

# 12 — color plus output file
go run ./cmd --color=blue --output=result.txt "Hello There"

# 13 — color, substring, output file and banner
go run ./cmd --color=blue --output=result.txt "Hello" "Hello There" thinkertoy

# 14 — color plus alignment
go run ./cmd --color=cyan --align=center "Hello There"

# 15 — color, alignment and banner
go run ./cmd --color=cyan --align=right "Hello There" shadow

# 16 — color, substring, justify and output file
go run ./cmd --color=yellow --align=justify --output=result.txt "There" "Hello There"

# 17 — everything with a single color
go run ./cmd "--color=#ff8800" --align=justify --output=result.txt "Hello" "Hello There" standard
```

### Multiple `--color`

With two or more `--color` flags, **every** color needs its own substring, in the same order as the flags.

```bash
# 18 — two colors, two substrings
go run ./cmd --color=red --color=blue "He" "llo" "Hello There"

# 19 — two colors with a banner
go run ./cmd --color=red --color=blue "Hello" "There" "Hello There" shadow

# 20 — two colors written to a file
go run ./cmd --color=red --color=green --output=result.txt "Hello" "There" "Hello There"

# 21 — two colors, file and banner
go run ./cmd --color=red --color=green --output=result.txt "Hello" "There" "Hello There" thinkertoy

# 22 — two colors plus alignment
go run ./cmd --color=magenta --color=cyan --align=center "Hello" "There" "Hello There"

# 23 — two colors, justify and banner
go run ./cmd --color=magenta --color=cyan --align=justify "Hello" "There" "Hello There" shadow

# 24 — two colors, justify and output file
go run ./cmd --color=orange --color=blue --align=justify --output=result.txt "Hello" "There" "Hello There"

# 25 — everything at once, including custom rgb/hex colors
go run ./cmd "--color=rgb(255,128,0)" "--color=#00c8ff" --align=justify --output=result.txt "Hello" "There" "Hello There" shadow
```

### Multi-line combinations

`\n` works with every flag; each line is aligned independently.

```bash
go run ./cmd --align=center "Hello\nThere"
go run ./cmd --color=red --align=justify "Hello There\nGeneral Kenobi"
go run ./cmd --color=red --color=blue --align=right --output=result.txt "Hello" "There" "Hello There\nAgain" shadow
```

## Input Rules

- You must pass exactly one string to render (plus one substring per `--color` flag when using two or more colors)
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

- `cmd/main.go` - CLI entrypoint; orchestrates parsing, output setup, and rendering
- `internal/config.go` - CLI argument parsing and validation (`ParseArgs`)
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

- The banner files are loaded from the relative path `banners/`, so the program must be run from the project root
- Supports printable ASCII characters only (`32` to `126`)
- Unicode characters (for example Greek letters or emoji) are rejected
- Colored substrings must match literally (no regex/wildcards)
- ANSI color codes are written even when output goes to a file via `--output` (no automatic stripping)
- Alignment uses the width of the terminal at run time, even when writing to a file with `--output`
- If the terminal is narrower than the rendered text, alignment leaves the output unchanged

## License

MIT. See `LICENSE`.
