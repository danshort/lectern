package ui

import tea "charm.land/bubbletea/v2"

func (m Model) updateSpec(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {

	case "q", "ctrl+c":
		return m, tea.Quit

	case "esc":
		specIdx := m.spec.Cursor
		jumpTarget := m.spec.JumpTarget
		wasFocusMode := m.spec.FocusMode
		m.enterIndex()
		if wasFocusMode && jumpTarget != "" {
			m.index.ExpandedSpecs[specIdx] = true
			m.buildIndexItems()
			for j, it := range m.index.Items {
				if it.kind == indexKindRequirement && it.idx == specIdx &&
					it.reqIdx < len(m.projectSpecs[specIdx].RequirementNames) &&
					m.projectSpecs[specIdx].RequirementNames[it.reqIdx] == jumpTarget {
					m.index.Cursor = j
					break
				}
			}
		} else {
			for i, item := range m.index.Items {
				if item.kind == indexKindSpec && item.idx == specIdx {
					m.index.Cursor = i
					break
				}
			}
		}
		m.refreshIndexViewport()

	case "j", "down":
		m.vp.ScrollDown(1)

	case "k", "up":
		m.vp.ScrollUp(1)

	case "h":
		if m.spec.FocusMode {
			ps := m.projectSpecs[m.spec.Cursor]
			if len(ps.RequirementNames) > 0 {
				m.spec.ReqCursor = (m.spec.ReqCursor - 1 + len(ps.RequirementNames)) % len(ps.RequirementNames)
				m.spec.JumpTarget = ps.RequirementNames[m.spec.ReqCursor]
				return m, m.loadViewport()
			}
		}

	case "l":
		if m.spec.FocusMode {
			ps := m.projectSpecs[m.spec.Cursor]
			if len(ps.RequirementNames) > 0 {
				m.spec.ReqCursor = (m.spec.ReqCursor + 1) % len(ps.RequirementNames)
				m.spec.JumpTarget = ps.RequirementNames[m.spec.ReqCursor]
				return m, m.loadViewport()
			}
		}

	case "e":
		if path := m.currentSpecPath(); path != "" {
			return m, m.openInEditor(path)
		}
	}
	return m, nil
}
