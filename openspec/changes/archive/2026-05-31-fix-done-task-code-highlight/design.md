## Context

The Tasks tab renders done-task code spans using `underlineStyle.Render(code) + doneRestore` in `inlineMarkdown()`. Lipgloss v2.0.3, when a style has `Underline(true)` and no foreground set, activates the `useSpaceStyler` code path (style.go:318). This renders each character individually:

```go
for _, r := range line {
    b.WriteString(te.Styled(string(r)))
}
```

Each `te.Styled(char)` wraps the character in `\033[4;4m<char>\033[m`. The `\033[m` resets all terminal attributes — including the foreground set by the outer `taskDoneStyle`. Only the first character of the span happens to retain the outer foreground (inherited from the text immediately preceding the span). Subsequent characters render with the terminal's default foreground, creating a visible color mismatch.

The bug was introduced in `replace-ansi-with-lipgloss` (archive 2026-05-30) when `\033[4m...\033[0m` + restore was replaced with `underlineStyle.Render()`. The raw-ANSI version used `\033[0m` only once (at the end of the span), so the foreground was preserved throughout.

## Goals / Non-Goals

**Goals:**
- Fix the first-character color mismatch for code spans in done tasks
- Preserve the underline visual distinction for done-task code spans

**Non-Goals:**
- Changing how pending-task code spans render (they use `cyanStyle` with foreground, no underline — unaffected)
- Changing non-code rendering of done tasks
- Refactoring inlineMarkdown logic or regex patterns

## Decisions

**Combine underline with foreground color in a single style.**
- Create `doneCodeStyle = lipgloss.NewStyle().Underline(true).Foreground(lipgloss.Color("8"))` at package level
- Use it in place of `underlineStyle` in the `done == true` branch of `inlineMarkdown()`
- This ensures the foreground color "8" (dark gray) is part of the style's `te` object, so character-by-character rendering preserves it between `\033[m` resets

**Alternatives considered:**
1. Raw ANSI `\033[4m` / `\033[24m` (underline on/off without full reset) — rejected because `replace-ansi-with-lipgloss` intentionally migrated away from raw ANSI
2. Setting `underlineSpaces(true)` on underlineStyle — doesn't prevent `useSpaceStyler`; the space-styler path would still be active
3. Disabling Width on taskDoneStyle — would break proper line padding in the task list

## Risks / Trade-offs

- **Duplicate "4" in ANSI sequence** → `\033[4;90;4m` instead of `\033[4;90m`. Harmless (duplicate parameters are ignored by terminals), but slightly longer output. No visual impact.
- **Existing tests may not catch regressions** → Current tests check for presence of "▶" cursor marker but not precise ANSI output. Should add a visual-regression test if practical.
