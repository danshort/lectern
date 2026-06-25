package ui

import (
	"errors"
	"strings"
	"testing"

	"github.com/danshort/lectern/internal/openspec"
)

func newIndexModel() Model {
	m := Model{mode: ModeIndex, width: 80, project: &openspec.Project{}}
	m.index.ExpandedSpecs = map[int]bool{}
	m.index.ExpandedArchives = map[int]bool{}
	return m
}

func TestUnreadableSpecMarker(t *testing.T) {
	m := newIndexModel()
	m.projectSpecs = []openspec.ProjectSpec{
		{Name: "unreadable-spec", ReadErr: errors.New("permission denied"), Content: "⚠ couldn't read .../spec.md: permission denied"},
		{Name: "valid-spec", Content: "## Purpose\nP\n\n## Requirements\n\n### Requirement: R\n#### Scenario: S\n- **WHEN** a\n- **THEN** b\n"},
	}
	m.buildIndexItems()

	out, _ := m.renderIndexContent()
	if !strings.Contains(out, "⚠") {
		t.Error("expected ⚠ marker for the unreadable spec")
	}
	// The unreadable spec must NOT also be flagged ✗ (read failure ≠ invalid),
	// and the valid spec is fine — so no ✗ anywhere.
	if strings.Contains(out, "✗") {
		t.Error("unreadable spec should show ⚠ in place of ✗, and the valid spec none")
	}
}

func TestUnreadableChangeMarker(t *testing.T) {
	t.Run("unreadable artifact shows warn marker", func(t *testing.T) {
		m := newIndexModel()
		m.project.Changes = []openspec.Change{
			{Name: "feat", Proposal: openspec.Artifact{Present: true, ReadErr: errors.New("EIO")}},
		}
		m.buildIndexItems()
		out, _ := m.renderIndexContent()
		if !strings.Contains(out, "⚠") {
			t.Error("expected ⚠ for a change with an unreadable artifact")
		}
		if strings.Contains(out, "✗") {
			t.Error("unreadable artifact should not also produce a ✗")
		}
	})

	t.Run("genuinely missing proposal still shows validation cross", func(t *testing.T) {
		m := newIndexModel()
		m.project.Changes = []openspec.Change{{Name: "feat"}} // no proposal present
		m.buildIndexItems()
		out, _ := m.renderIndexContent()
		if !strings.Contains(out, "✗") {
			t.Error("a change missing its proposal should still show ✗")
		}
	})
}
