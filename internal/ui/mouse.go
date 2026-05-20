package ui

import (
	"fmt"
	"strings"

	"rime-wanxiang-updater/internal/types"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleMouseInput(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if msg.Action != tea.MouseActionPress || msg.Button != tea.MouseButtonLeft {
		return m, nil
	}

	switch m.State {
	case ViewWizard:
		if n, ok := choiceNumberAtY(m.renderWizard(), msg.Y, wizardChoiceCount(m.WizardStep)); ok {
			return m.handleWizardInput(numberKeyMsg(n))
		}
	case ViewMenu:
		if n, ok := choiceNumberAtY(m.renderMenu(), msg.Y, 8); ok {
			return m.handleMenuInput(numberKeyMsg(n))
		}
	case ViewCustomMenu:
		if n, ok := choiceNumberAtY(m.renderCustomMenu(), msg.Y, len(m.customMenuItems())); ok {
			return m.handleCustomMenuInput(numberKeyMsg(n))
		}
	}

	return m, nil
}

func wizardChoiceCount(step WizardStep) int {
	switch step {
	case WizardSchemeType, WizardDownloadSource:
		return 2
	case WizardSchemeVariant:
		return len(types.SchemeMap)
	default:
		return 0
	}
}

func choiceNumberAtY(rendered string, y, max int) (int, bool) {
	if max <= 0 || y < 0 {
		return 0, false
	}

	lines := strings.Split(rendered, "\n")
	for _, row := range []int{y, y - 1, y + 1} {
		if row < 0 || row >= len(lines) {
			continue
		}
		for n := 1; n <= max; n++ {
			if strings.Contains(lines[row], fmt.Sprintf("[%d]", n)) {
				return n, true
			}
		}
	}

	return 0, false
}

func numberKeyMsg(n int) tea.KeyMsg {
	return tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(fmt.Sprintf("%d", n)),
	}
}
