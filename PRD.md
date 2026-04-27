# PRD - ascii-art

Keep this short. The project subject + audit cases are the source of truth.

---

## 1. Problem Statement

This CLI receives text as input and prints the same text as ASCII-art using banner templates.

---

## 2. CLI Contract

- Command: `go run ./cmd "<text>"`
- Args:
  - exactly 1 argument is required
  - input supports letters, numbers, spaces, special chars, and literal `\n`
- Missing args:
  - print usage guidance
  - exit without rendering
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

- User selects one banner style:
  - `1` -> `standard`
  - `2` -> `shadow`
  - `3` -> `thinkertoy`
- Invalid style input is rejected.

Example:
- Input `Hello` + style `1` -> prints `Hello` in `standard` font.

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

---

## 4. Non-Goals

- Unicode support beyond printable ASCII
- GUI/Web interface
- Custom user-uploaded fonts
- Rich text features (color, alignment, animation)

---

## 5. Acceptance Criteria

### Audit Cases

- [ ] `hello` prints expected ASCII-art output in selected style.
- [ ] Mixed case and spaces (for example `HeLlo HuMaN`) render correctly.
- [ ] Special characters (for example `{|}~` and punctuation sets) render correctly.
- [ ] Literal `\n` and `\n\n` create correct multi-line block separation.

### Extra Golden Tests

- [ ] Empty input behavior is defined and tested (`""`).
- [ ] Invalid character input (for example emoji) returns clear error and non-zero exit.
- [ ] Invalid style choice (non-numeric or out of range) returns clear error and non-zero exit.

---

## 6. Architecture

- We choose: Pipeline
- Because: input validation, style selection, splitting, and rendering happen in a simple sequential flow.
- Tradeoffs we accept: less flexibility than a full parser/FSM, but easier to read and maintain.

Sketch (high-level):

`CLI arg -> validate chars -> ask style -> load banner -> split by \n -> render each line -> stdout`

---

## 7. Milestones

1. Stabilize expected behavior for empty/newline edge cases.
2. Ensure functional output matches subject examples for key inputs.
3. Add/clean tests for validation, multiline behavior, and special chars.
4. Final docs pass (README usage + examples aligned with real behavior).

---

## 8. Risks / Open Questions

- Should final audit target exact `go run .` contract or keep `go run ./cmd`?
- Should style selection remain interactive or move to a CLI flag/default style for stricter compatibility?

---

## 9. Known Limitations

- Rendering is limited to printable ASCII (`32-126`).
- Input with Unicode symbols (for example Greek characters or emoji) is not supported.
- Current UX requires interactive style selection from stdin.
