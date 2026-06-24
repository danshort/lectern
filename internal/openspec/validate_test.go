package openspec

import "testing"

func TestValidateSpec(t *testing.T) {
	valid := `# Foo Specification

## Purpose
Does a thing.

## Requirements

### Requirement: Does the thing
The system SHALL do the thing.

#### Scenario: It does it
- **WHEN** asked
- **THEN** it does
`

	tests := []struct {
		name      string
		content   string
		wantValid bool
	}{
		{"valid spec", valid, true},
		{
			"missing purpose",
			"# Foo\n\n## Requirements\n\n### Requirement: X\n\n#### Scenario: Y\n- **WHEN** a\n- **THEN** b\n",
			false,
		},
		{
			"missing requirements",
			"# Foo\n\n## Purpose\nThing.\n",
			false,
		},
		{
			"requirement without scenario",
			"# Foo\n\n## Purpose\nThing.\n\n## Requirements\n\n### Requirement: Lonely\nNo scenario here.\n",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateSpec(tt.content)
			if tt.wantValid && len(errs) != 0 {
				t.Errorf("expected valid, got errors: %v", errs)
			}
			if !tt.wantValid && len(errs) == 0 {
				t.Errorf("expected errors, got none")
			}
		})
	}
}

func TestValidateSpecMultipleRequirements(t *testing.T) {
	// First requirement has a scenario, second does not.
	content := `## Purpose
P.

## Requirements

### Requirement: First
#### Scenario: ok
- **WHEN** a
- **THEN** b

### Requirement: Second
no scenario
`
	errs := ValidateSpec(content)
	if len(errs) != 1 {
		t.Fatalf("expected exactly 1 error, got %d: %v", len(errs), errs)
	}
}

func TestValidateChange(t *testing.T) {
	goodDelta := NamedSpec{
		Name:    "cap",
		Content: "## ADDED Requirements\n\n### Requirement: X\n#### Scenario: Y\n- **WHEN** a\n- **THEN** b\n",
	}

	tests := []struct {
		name      string
		change    Change
		wantValid bool
	}{
		{
			"valid change",
			Change{Proposal: Artifact{Present: true}, SpecFiles: []NamedSpec{goodDelta}},
			true,
		},
		{
			"valid change with no delta specs",
			Change{Proposal: Artifact{Present: true}},
			true,
		},
		{
			"missing proposal",
			Change{Proposal: Artifact{Present: false}},
			false,
		},
		{
			"delta without header",
			Change{Proposal: Artifact{Present: true}, SpecFiles: []NamedSpec{
				{Name: "cap", Content: "### Requirement: X\n#### Scenario: Y\n- **WHEN** a\n- **THEN** b\n"},
			}},
			false,
		},
		{
			"added requirement without scenario",
			Change{Proposal: Artifact{Present: true}, SpecFiles: []NamedSpec{
				{Name: "cap", Content: "## ADDED Requirements\n\n### Requirement: X\nno scenario\n"},
			}},
			false,
		},
		{
			"removed requirement without scenario is exempt",
			Change{Proposal: Artifact{Present: true}, SpecFiles: []NamedSpec{
				{Name: "cap", Content: "## REMOVED Requirements\n\n### Requirement: Old\nreason for removal\n"},
			}},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateChange(tt.change)
			if tt.wantValid && len(errs) != 0 {
				t.Errorf("expected valid, got errors: %v", errs)
			}
			if !tt.wantValid && len(errs) == 0 {
				t.Errorf("expected errors, got none")
			}
		})
	}
}
