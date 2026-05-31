## 1. Core Implementation

- [x] 1.1 Add `doneCodeStyle` package-level variable combining `Underline(true)` + `Foreground(Color("8"))` in `internal/ui/tasks.go`
- [x] 1.2 Replace `underlineStyle.Render(code)` with `doneCodeStyle.Render(code)` in the `done == true` branch of `inlineMarkdown()` in `internal/ui/tasks.go`
- [x] 1.3 Run `go build ./cmd/dossier/` to verify compilation

## 2. Verification

- [x] 2.1 Run `go test ./internal/ui/` to confirm existing tests pass
- [x] 2.2 Launch app and visually confirm done-task code spans render uniformly (no first-character color mismatch) in kitty and gnome-terminal
