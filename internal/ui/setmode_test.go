package ui

import (
	"testing"

	"github.com/danshort/lectern/internal/openspec"
)

// setMode is the single chokepoint for mode transitions; these tests pin the
// state invariants it enforces (see the tui-viewer "Mode transitions enforce
// state invariants centrally" requirement).

func TestSetModeClearsFocusOnLeavingSpec(t *testing.T) {
	m := &Model{mode: ModeViewingSpec, projectSpecs: []openspec.ProjectSpec{{Name: "auth"}}}
	m.spec.FocusMode = true
	m.spec.JumpTarget = "Login"
	m.spec.ReqCursor = 3

	m.setMode(ModeIndex)

	if m.spec.FocusMode || m.spec.JumpTarget != "" || m.spec.ReqCursor != 0 {
		t.Errorf("leaving spec mode did not clear focus state: FocusMode=%v JumpTarget=%q ReqCursor=%d",
			m.spec.FocusMode, m.spec.JumpTarget, m.spec.ReqCursor)
	}
}

func TestSetModeKeepsFocusWhenNotLeavingSpec(t *testing.T) {
	// Returning to the spec view from the config overlay must not wipe focus:
	// the outgoing mode is config, not spec.
	m := &Model{mode: ModeViewingConfig, projectSpecs: []openspec.ProjectSpec{{Name: "auth"}}}
	m.spec.FocusMode = true
	m.spec.JumpTarget = "Login"
	m.spec.ReqCursor = 2

	m.setMode(ModeViewingSpec)

	if !m.spec.FocusMode || m.spec.JumpTarget != "Login" || m.spec.ReqCursor != 2 {
		t.Errorf("re-entering spec mode wrongly altered focus state: FocusMode=%v JumpTarget=%q ReqCursor=%d",
			m.spec.FocusMode, m.spec.JumpTarget, m.spec.ReqCursor)
	}
}

func TestSetModeClampsTabToAvailable(t *testing.T) {
	// Only design is present; entering a tabbed mode with proposal selected must
	// clamp to the first available tab.
	ch := openspec.Change{Design: openspec.Artifact{Present: true}}
	m := &Model{mode: ModeIndex, project: &openspec.Project{Changes: []openspec.Change{ch}}}
	m.viewer.changeIdx = 0
	m.viewer.tab = TabProposal // not available on ch

	m.setMode(ModeNormal)

	if m.viewer.tab != TabDesign {
		t.Errorf("setMode(ModeNormal) tab = %d, want TabDesign (%d)", m.viewer.tab, TabDesign)
	}
}

func TestSetModeClampsSpecCursor(t *testing.T) {
	specs := []openspec.ProjectSpec{{Name: "a"}, {Name: "b"}}
	tests := []struct {
		name   string
		cursor int
		specs  []openspec.ProjectSpec
		want   int
	}{
		{"above range clamps to last", 9, specs, 1},
		{"negative clamps to zero", -4, specs, 0},
		{"empty specs clamps to zero", 5, nil, 0},
		{"in range is preserved", 1, specs, 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := &Model{mode: ModeIndex, projectSpecs: tc.specs}
			m.spec.Cursor = tc.cursor
			m.setMode(ModeViewingSpec)
			if m.spec.Cursor != tc.want {
				t.Errorf("spec cursor = %d, want %d", m.spec.Cursor, tc.want)
			}
		})
	}
}

func TestSetModeDoesNotDropRenderCache(t *testing.T) {
	m := &Model{
		mode:        ModeNormal,
		renderCache: map[Tab]string{TabProposal: "cached"},
		project:     &openspec.Project{},
	}
	m.setMode(ModeIndex)
	if got, ok := m.renderCache[TabProposal]; !ok || got != "cached" {
		t.Errorf("setMode dropped renderCache entry: ok=%v got=%q", ok, got)
	}
}
