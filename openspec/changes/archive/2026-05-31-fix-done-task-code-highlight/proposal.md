## Why

In the Tasks tab, done tasks with code spans (backtick-delimited) render incorrectly: the first character of the code span appears in the dim foreground color ("shaded") while the rest appears in the terminal's default foreground ("bright white"). This happens because `underlineStyle` (only `Underline(true)`) triggers lipgloss v2.0.3's character-by-character rendering path (`useSpaceStyler`), where `\033[m` between each character resets the foreground to default. Only the first character inherits the correct foreground from the outer style.

## What Changes

- Modify `inlineMarkdown()` in `internal/ui/tasks.go` to use a combined `doneCodeStyle` (`Underline(true) + Foreground(Color("8"))`) instead of bare `underlineStyle` for done task code spans. This ensures the per-character rendering preserves the foreground color.

## Capabilities

### New Capabilities
- *(none)*

### Modified Capabilities
- *(none — this is an implementation fix with no spec-level behavior change)*

## Impact

- `internal/ui/tasks.go`: replace `underlineStyle.Render(code)` with `doneCodeStyle.Render(code)` in the done branch of `inlineMarkdown()`
- Add package-level `doneCodeStyle` variable matching `underlineStyle` plus `Foreground(Color("8"))`
- No API changes, no dependency changes, no test changes needed (existing tests verify structural output only)
