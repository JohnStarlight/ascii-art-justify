TODO:
✅ args.go:39  → if len(args) < 2 || len(args) > 4
✅ args.go:50  → banner := "standard"
✅ args.go:58  → αφαίρεση case 2 banner = "standard"
✅ args.go:70  → strings.ToLower(args[1]) πριν HasPrefix
✅ args.go:91  → strings.ToLower(banner) πριν map lookup
✅ args.go:105 → for → validateText() function
✅ args.go:31-46 → αφαίρεση obvious comments
✅ args.go     → μετονομασία σε config.go
✅ output.go   → σύντομο comment, μία γραμμή
✅ printascii.go:52-65 → Normalize πριν Validate
✅ main.go:12-20 → αφαίρεση obvious comments
✅ main.go:62-65 → αφαίρεση obvious comments
□ README.md   → ανανέωση
□ PRD.md      → έλεγχος & ανανέωση
□ TASK_CARDS.md → έλεγχος & ανανέωση

1. line 39
- if len(args) != 2 && len(args) != 3 && len(args) != 4 {
+ if len(args) < 2 || len(args) > 4 {  - more go idiomatic

2. 31-46 (31-42 New)
-34 // Valid command formats:
-35	//
-36	// go run ./cmd [STRING]
-37	// go run ./cmd [STRING] [BANNER]
-38	// go run ./cmd --output=<fileName.txt> [STRING] [BANNER] 
+36 // Valid command formats:

3. line 50 (46 new)
50 -var banner string
46 +banner := "standard" // Default banner    - declare default

4. line 58 (53 new)
54 -banner = "standard"  - no need because of default

5. line 70 (64 New)
64 +args[1] = strings.ToLower(args[1]) - to avoid caps in args

6. line 91 (88 New)
88. +banner = strings.ToLower(banner)  - to avoid caps in args

7. line: 105-137 (97-100 && 110-148 New)

105 - // Validate input characters.
106 -	//
107 -	// Allowed:
108 -	// - printable ASCII characters (32–126)
109 -	// - the escape sequence "\n"
110 - 137 Neo func

97  + // Validate the input text.
98  + if err := validateText(text); err != nil {
99	+ return Config{}, err
100	+ }

110 + // Validate input characters and escape sequences.
111 + func validateText(text string) error {
112 +  for i := 0; i < len(text); i++ {
113 +  char := text[i]
114 +
115 +   // Handle escape sequences beginning with '\'.
116     if char == '\\' {
117 
118     // '\' cannot appear as the final character.
119       if i+1 >= len(text) {
113 -     return Config{}, fmt.Errorf(
120 +     return fmt.Errorf(
121       "invalid escape sequence: trailing backslash",
122       )
123     }
124 
125     // Only "\n" is supported.
126     if text[i+1] != 'n' {
120 -		  return Config{}, fmt.Errorf(
127 +     return fmt.Errorf(
128       "invalid escape sequence \\%c",
129       text[i+1],
130       )
131     }
132 
133       // Skip the validated 'n'.
134       i++
135       continue
136     }
137 
138       // Reject non-printable ASCII characters.
139       if char < 32 || char > 126 {
133 -			return Config{}, fmt.Errorf(
140 +     return fmt.Errorf(
141       "invalid character %q (only printable ASCII is supported)",
142       char,
143       )
144     }
145   }
146 + return nil
147 + }

8. file args.go -> config.go

/internal/output.go

9. lines: 8-12 (8 New)

8  - // PrepareOutputFile creates or truncates
9  - // the specified output file.
10 - //
11 - // If the file already exists,
12 - // its contents are erased.

8 + // Creates or truncates output file and returns a handle to it.

/internal/printascii.go

10. Lines: 52-65 - Move Validate after Normalize to ensure accurate newline count

52 - Validate
52  + // Normalize Windows line endings (\r\n)
53	+ // into Unix format (\n).
54	+ normalized := strings.ReplaceAll(
55	+ 	string(data),
56	+ 	"\r\n",
57	+ 	"\n",
58	+ )
59
60  - Normilize
60  +	// Validate banner structure before rendering.
61  +	if strings.Count(normalized, "\n") != expectedNewlines {
62  +		return fmt.Errorf(
63  +			"banner file is corrupt or invalid",
64  +		)
65  +	}

- Validate banner structure after normalizing line endings
  to ensure accurate newline count regardless of OS format.
  Previously, validation ran on raw data before normalization,
  causing false negatives on Windows-formatted files.

/cmd/main.go

11. remove obvious comments - need to comment Why not what!

13-20
- //
- // It:
- // - parses and validates CLI arguments
- // - prepares the output destination
- // - splits logical lines
- // - renders ASCII art
- //
- // Any error is returned to main for centralized handling.

14
- // Parse and validate command-line arguments.

19
- // By default, write ASCII art to the terminal.


22-23
-	// If an output file was requested,
-	// create/truncate the file and redirect output there.

36-37
- // Convert literal "\n" sequences
-	// into separate logical lines.

36
+ // Split the input text into separate logical lines.

57
- // Print errors to stderr instead of stdout.

62 -65 
-	// main is responsible only for:
-	// - starting the application
-	// - handling the final error
-	// - returning the correct exit status

61 
+ // Any error is returned to main for centralized handling.

70-71
-		// Exit with non-zero status code
-		// to indicate program failure.

12. README.md  - all what comments is going here