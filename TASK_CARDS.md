# Task: Empty Input Behavior

**Goal:** Keep one clear and tested behavior for empty input (`""`) so output is predictable.

**Depends on:** none

## Acceptance criteria
- [ ] `PrintAscii([]string{""}, ...)` behavior is explicitly decided and documented.
- [ ] Related unit test matches the chosen behavior.
- [ ] `go test ./...` passes.

## Tests to write first
- [ ] `TestEmptyLine`: proves empty input output is exactly as specified.
- [ ] `TestNewlineSeparator`: proves `\n` separators still work after the change.

## Notes
- Empty input and `"\\n"` are different cases and should stay distinct.

# Task: CLI Input Validation

**Goal:** Make CLI failures clear and consistent for wrong arguments and banner choices.

**Depends on:** none

## Acceptance criteria
- [ ] Missing argument shows usage guidance.
- [ ] Invalid banner name exits with clear error.
- [ ] Invalid `--output` format (missing `=` or empty filename) exits with clear error.
- [ ] Error paths return non-zero exit code.

## Tests to write first
- [ ] `TestMissingArgument`: proves app handles no argument correctly.
- [ ] `TestInvalidBannerName`: proves banner validation rejects unknown values.
- [ ] `TestInvalidOutputFlag`: proves malformed `--output` flag is rejected.

## Notes
- Keep messages short and friendly for beginner users.

# Task: Character Validation

**Goal:** Reject unsupported characters safely and consistently.

**Depends on:** none

## Acceptance criteria
- [ ] Only ASCII `32-126` is accepted.
- [ ] Accented letters and emoji are rejected with clear message.
- [ ] Program exits with non-zero status on invalid characters.

## Tests to write first
- [ ] `TestRejectAccentedCharacter`: proves non-ASCII letters are blocked.
- [ ] `TestRejectEmoji`: proves emoji input is blocked.

## Notes
- Validation should happen before rendering begins.

# Task: CI Checks

**Goal:** Automatically verify formatting and tests on each push/PR.

**Depends on:** CLI Input Validation, Character Validation

## Acceptance criteria
- [ ] CI runs `go test ./...`.
- [ ] CI fails when tests fail.
- [ ] CI checks formatting (for example `gofmt`).

## Tests to write first
- [ ] `CI: test job`: proves test command runs in pipeline.
- [ ] `CI: fmt job`: proves formatting check runs in pipeline.

## Notes
- Keep CI simple and fast.

# Task: Documentation Refresh

**Goal:** Make README clear for first-time readers and aligned with real behavior.

**Depends on:** Empty Input Behavior, CLI Input Validation

## Acceptance criteria
- [ ] README has one quick start example.
- [ ] README explains style selection and `\n` usage.
- [ ] README input rules match actual code behavior.

## Tests to write first
- [ ] `Manual doc check: quick start`: proves command works as written.
- [ ] `Manual doc check: multiline example`: proves example output flow is correct.

## Notes
- Keep wording simple enough for non-experts.

# Task: File Output Flag

**Goal:** Allow users to write ASCII art directly to a file via `--output=<filename>`.

**Depends on:** CLI Input Validation

## Acceptance criteria
- [x] `--output=<filename>` writes ASCII art to the specified file.
- [x] File is created if it does not exist; overwritten if it does.
- [x] Missing or malformed flag exits with clear error and non-zero status.
- [x] Without the flag, output still goes to stdout as before.

## Tests to write first
- [ ] `TestOutputToFile`: proves ASCII art is written correctly to a file.
- [ ] `TestOutputFileOverwrite`: proves existing file is overwritten without error.
- [ ] `TestMissingOutputFilename`: proves `--output=` with empty name is rejected.

## Notes
- Implementation lives in `internal/output.go` and `internal/config.go`.
- Banner argument is required when `--output` is used.
