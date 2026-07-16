# PRD - ascii-art-justify

Keep this short. The project subject + audit cases are the source of truth.

---

## 1. Problem Statement

This CLI receives text as input and prints the same text as ASCII-art using banner templates, optionally colored and aligned inside the terminal.

---

## 2. CLI Contract

Supported forms:

```
go run ./cmd "<text>"
go run ./cmd "<text>" <banner>
go run ./cmd --output=<filename> "<text>" <banner>
go run ./cmd --color=<color> "<text>"
go run ./cmd --color=<color> "<substring>" "<text>"
go run ./cmd --color=<color> "<substring>" "<text>" <banner>
go run ./cmd --color=<c1> --color=<c2> "<sub1>" "<sub2>" "<text>"
go run ./cmd --color=<c1> --color=<c2> "<sub1>" "<sub2>" "<text>" <banner>
go run ./cmd --align=<type> "<text>" <banner>
```

- `<text>`: required, must be in quotes for multi-word input
- `<banner>`: optional positional arg; defaults to `standard`; accepted values: `standard`, `shadow`, `thinkertoy`
- `--output=<filename>`: optional flag; when present, writes ASCII art to file instead of stdout; banner becomes required
- `--color=<color>`: optional flag, repeatable; colors the rendered output (see 3.5, 3.6)
- `<substring>`: positional arg, only meaningful with `--color`; the portion of `<text>` to color. With a single `--color` it is optional (omitting it colors the whole text). With multiple `--color` flags, one substring per color is required, in the same order.
- `--align=<type>`: optional flag; accepted values: `left`, `right`, `center`, `justify` (see 3.7)
- `--color`/`--output`/`--align` flags can appear in any order/position among the arguments
- Invalid `--output` format (missing `=` or empty filename): print error and exit with non-zero status
- Invalid `--color` format (missing `=`, empty value, or unknown color): print usage message and exit with non-zero status
- Invalid `--align` format (missing `=`, empty value, or unknown type): print usage message and exit with non-zero status
- Invalid chars (outside ASCII 32-126):
  - print clear error message
  - exit with non-zero status
- Banner file read failure:
  - print `Could not open file: ...`
  - exit with non-zero status

---

## 3. Rendering Functions

Documented in plain language with examples.

### 3.1 Style Selection

- Banner is passed as a positional CLI argument (default: `standard`).
- Accepted values: `standard`, `shadow`, `thinkertoy` (case-insensitive).
- Unsupported banner name: print error and exit with non-zero status.

Example:
- `go run ./cmd "Hello" shadow` -> prints `Hello` in `shadow` font.

### 3.2 Newline Handling

- Input is split on literal `\n`.
- Each segment is rendered as a separate 8-line block.
- Empty segments produce an empty line separator (except leading empty segment in current behavior).

Examples:
- `Hello\nThere` -> two rendered blocks with one blank line between.
- `Hello\n\nThere` -> two rendered blocks with two separators (one empty rendered segment in between).

### 3.3 Character-to-Glyph Mapping

- Each printable ASCII character maps to one 8-line glyph from the banner file.
- For each output row (1..8), the program concatenates the correct row from each character glyph.

Example:
- `A` -> prints 8 lines that correspond to ASCII code 65 in the selected banner.

### 3.4 File Output

- When `--output=<filename>` is provided, ASCII art is written to that file instead of stdout.
- File is created if it does not exist; overwritten if it does.
- On file creation failure: print error and exit with non-zero status.
- Banner argument is required when `--output` is used.

Example:
- `go run ./cmd --output=result.txt "Hello" standard` -> writes `Hello` in `standard` font to `result.txt`.

### 3.5 Color Highlighting

- `--color=<color>` colors the rendered ASCII art using ANSI 24-bit (`\033[38;2;R;G;Bm`) escape codes.
- Accepted color names: `red`, `green`, `blue`, `yellow`, `magenta`/`purple`, `cyan`, `orange`.
- Custom colors via `rgb(R,G,B)` (e.g. `"--color=rgb(255,128,0)"`) or hex `#RRGGBB` (e.g. `"--color=#ff8800"`); quotes required, see below.
- An optional `<substring>` positional argument selects which part of `<text>` is colored; all occurrences of the substring are colored.
- If `<substring>` is omitted (single `--color` only), the entire rendered text is colored.
- Unknown color name, malformed `rgb(...)`/`#RRGGBB`, or malformed `--color` flag (missing `=` or empty value): print usage message and exit with non-zero status.
- Shell note: `rgb(...)` contains parentheses and hex colors contain `#`, which are special characters in most shells. The whole flag must be quoted, e.g. `"--color=rgb(0,200,255)"`, `"--color=#ff8800"`.

Examples:
- `go run ./cmd --color=red "Hello"` -> entire `Hello` rendered in red.
- `go run ./cmd --color=red "ell" "Hello"` -> only the `ell` substring rendered in red.
- `go run ./cmd "--color=rgb(0,200,255)" "Hello" shadow` -> entire `Hello` rendered in custom RGB color, `shadow` banner.

### 3.6 Multiple Colors

- `--color` may be passed more than once; each flag pairs with its own `<substring>` positional argument, in order.
- With two or more `--color` flags, every color requires its own substring (no whole-text shorthand); wrong argument count prints usage and exits non-zero.
- Overlapping substrings resolve in flag order: the first `--color` whose substring covers a character wins.

Example:
- `go run ./cmd --color=red --color=blue "He" "llo" "Hello"` -> `He` in red, `llo` in blue.

### 3.7 Alignment

- `--align=<type>` positions the rendered output inside the terminal; accepted types: `left`, `right`, `center`, `justify`.
- `left` (and no flag) leaves output unchanged; `right` pads each row so it ends at the terminal edge; `center` pads each row by half the free space.
- `justify` splits the text into words and stretches the spaces between them so every row spans exactly the terminal width. A single word (or a line wider than the terminal) is rendered unchanged.
- Terminal width detection order: `stty size` -> `tput cols` -> `COLUMNS` env var -> fallback 80.
- The flag must be written exactly `--align=<type>`; `--align <type>`, `--align=` and unknown types print usage and exit non-zero.
- Alignment composes with `--color`: substring coloring survives justify because words keep their original positions in the input line.

Examples:
- `go run ./cmd --align=right "Hello" standard` -> `Hello` pushed to the right edge.
- `go run ./cmd --align=justify "Hello There" standard` -> `Hello` at the left edge, `There` at the right edge, spaces distributed between.

---

## 4. Non-Goals

- Unicode support beyond printable ASCII
- GUI/Web interface
- Custom user-uploaded fonts
- Rich text features beyond coloring and alignment (animation, gradients)
- Word-wrapping text that is wider than the terminal

---

## 5. Acceptance Criteria

### Audit Cases

- [x] `hello` prints expected ASCII-art output in selected style.
- [x] Mixed case and spaces (for example `HeLlo HuMaN`) render correctly.
- [x] Special characters (for example `{|}~` and punctuation sets) render correctly.
- [x] Literal `\n` and `\n\n` create correct multi-line block separation.
- [x] `--align=right/center/justify` positions output correctly relative to terminal width.
- [x] `--align=justify` keeps substring coloring intact.

### Extra Golden Tests

- [x] Empty input behavior is defined and tested (`""`).
- [x] Invalid character input (for example emoji) returns clear error and non-zero exit.
- [x] Invalid style choice returns clear error and non-zero exit.
- [x] Invalid `--align` formats (`--align right`, `--align=`, unknown type) return usage error and non-zero exit.
- [x] Multiple `--color` flags map colors to substrings in order; missing substring is rejected.
- [x] Hex (`#RRGGBB`) and `orange` colors produce the expected ANSI codes; malformed hex is rejected.

---

## 6. Architecture

- We choose: Pipeline
- Because: input validation, banner selection, splitting, rendering, and alignment happen in a simple sequential flow.
- Tradeoffs we accept: less flexibility than a full parser/FSM, but easier to read and maintain.

Sketch (high-level):

`CLI args -> parse & validate -> load banner -> split by \n -> render each line (with colors) -> align rows (terminal width) -> stdout or file`

---

## 7. Milestones

1. ~~Stabilize expected behavior for empty/newline edge cases.~~ (done)
2. ~~Ensure functional output matches subject examples for key inputs.~~ (done)
3. ~~Add `--align` support (left/right/center/justify) with terminal width detection.~~ (done)
4. ~~Merge color v2 (multiple `--color` flags, hex, orange) with the justify renderer.~~ (done)
5. ~~Add/clean tests for validation, multiline, colors, and alignment.~~ (done)
6. Final docs pass (README usage + examples aligned with real behavior).

---

## 8. Risks / Open Questions

- Should final audit target exact `go run .` contract or keep `go run ./cmd`?
- Should alignment be skipped (or width taken differently) when writing to a file via `--output`?

---

## 9. Known Limitations

- Rendering is limited to printable ASCII (`32-126`).
- Input with Unicode symbols (for example Greek characters or emoji) is not supported.
- When `--output` is used, banner argument is required (no fallback to default).
- Colored substrings are matched literally (no regex/wildcards).
- ANSI color codes are written even when output is redirected to a file via `--output` (no automatic stripping for non-terminal targets).
- Alignment uses the terminal width at run time even when writing to a file via `--output`.
- Rows wider than the terminal are printed unchanged (no wrapping or truncation).
