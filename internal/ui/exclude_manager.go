package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// handleExcludeListInput 处理排除文件列表输入
func (m Model) handleExcludeListInput(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		m.State = ViewConfig
		m.ExcludeErrorMsg = ""
		return m, nil

	case "up", "k":
		if m.ExcludeListChoice > 0 {
			m.ExcludeListChoice--
		}

	case "down", "j":
		maxChoice := len(m.Cfg.Config.ExcludeFiles) + 2
		if m.ExcludeListChoice < maxChoice {
			m.ExcludeListChoice++
		}

	case "enter", " ":
		return m.handleExcludeListSelect()

	case "d", "x":
		if m.ExcludeListChoice < len(m.Cfg.Config.ExcludeFiles) {
			if err := m.Cfg.RemoveExcludePattern(m.ExcludeListChoice); err != nil {
				m.ExcludeErrorMsg = m.t("exclude.error.delete", err)
			} else {
				m.ExcludeErrorMsg = ""
				m.ExcludeDescriptions, _ = m.Cfg.GetExcludePatternDescriptions()
				if m.ExcludeListChoice >= len(m.Cfg.Config.ExcludeFiles) && m.ExcludeListChoice > 0 {
					m.ExcludeListChoice--
				}
			}
		}
	}

	return m, nil
}

// handleExcludeListSelect 处理排除列表选择
func (m Model) handleExcludeListSelect() (Model, tea.Cmd) {
	numPatterns := len(m.Cfg.Config.ExcludeFiles)

	if m.ExcludeListChoice < numPatterns {
		m.ExcludeEditIndex = m.ExcludeListChoice
		m.ExcludeEditInput = m.Cfg.Config.ExcludeFiles[m.ExcludeListChoice]
		m.State = ViewExcludeEdit
		m.ExcludeErrorMsg = ""
	} else if m.ExcludeListChoice == numPatterns {
		m.ExcludeEditInput = ""
		m.State = ViewExcludeAdd
		m.ExcludeErrorMsg = ""
	} else if m.ExcludeListChoice == numPatterns+1 {
		if err := m.Cfg.ResetExcludePatterns(); err != nil {
			m.ExcludeErrorMsg = m.t("exclude.error.reset", err)
		} else {
			m.ExcludeErrorMsg = m.t("exclude.reset.done")
			m.ExcludeDescriptions, _ = m.Cfg.GetExcludePatternDescriptions()
			m.ExcludeListChoice = 0
		}
	}

	return m, nil
}

// handleExcludeEditInput 处理排除模式编辑输入
func (m Model) handleExcludeEditInput(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.State = ViewExcludeList
		m.ExcludeErrorMsg = ""
		return m, nil

	case "enter":
		if strings.TrimSpace(m.ExcludeEditInput) == "" {
			m.ExcludeErrorMsg = m.t("exclude.error.empty")
			return m, nil
		}

		if err := m.Cfg.RemoveExcludePattern(m.ExcludeEditIndex); err != nil {
			m.ExcludeErrorMsg = m.t("exclude.error.delete", err)
			return m, nil
		}

		if err := m.Cfg.AddExcludePattern(m.ExcludeEditInput); err != nil {
			m.ExcludeErrorMsg = m.t("exclude.error.save", err)
			return m, nil
		}

		m.ExcludeDescriptions, _ = m.Cfg.GetExcludePatternDescriptions()
		m.State = ViewExcludeList
		m.ExcludeErrorMsg = ""
		return m, nil

	case "backspace":
		if len(m.ExcludeEditInput) > 0 {
			m.ExcludeEditInput = m.ExcludeEditInput[:len(m.ExcludeEditInput)-1]
		}

	default:
		if len(msg.String()) == 1 {
			m.ExcludeEditInput += msg.String()
		}
	}

	return m, nil
}

// handleExcludeAddInput 处理添加排除模式输入
func (m Model) handleExcludeAddInput(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.State = ViewExcludeList
		m.ExcludeErrorMsg = ""
		return m, nil

	case "enter":
		if strings.TrimSpace(m.ExcludeEditInput) == "" {
			m.ExcludeErrorMsg = m.t("exclude.error.empty")
			return m, nil
		}

		if err := m.Cfg.AddExcludePattern(m.ExcludeEditInput); err != nil {
			m.ExcludeErrorMsg = m.t("exclude.error.add", err)
			return m, nil
		}

		m.ExcludeDescriptions, _ = m.Cfg.GetExcludePatternDescriptions()
		m.State = ViewExcludeList
		m.ExcludeErrorMsg = ""
		m.ExcludeEditInput = ""
		return m, nil

	case "backspace":
		if len(m.ExcludeEditInput) > 0 {
			m.ExcludeEditInput = m.ExcludeEditInput[:len(m.ExcludeEditInput)-1]
		}

	default:
		if len(msg.String()) == 1 {
			m.ExcludeEditInput += msg.String()
		}
	}

	return m, nil
}

// renderExcludeList 渲染排除文件列表
func (m Model) renderExcludeList() string {
	var b strings.Builder

	b.WriteString(m.renderHeaderBlock())

	title := RenderGradientTitle("📋 " + m.t("exclude.title") + " 📋")
	b.WriteString(title + "\n\n")

	helpText := lipgloss.NewStyle().
		Foreground(m.Styles.Muted).
		Render(m.t("exclude.help"))
	b.WriteString(helpText + "\n\n")

	var listContent strings.Builder
	if len(m.Cfg.Config.ExcludeFiles) == 0 {
		listContent.WriteString(lipgloss.NewStyle().Foreground(m.Styles.Warning).Render(m.t("exclude.empty")) + "\n")
	} else {
		for i, pattern := range m.Cfg.Config.ExcludeFiles {
			var desc string
			if i < len(m.ExcludeDescriptions) {
				desc = m.ExcludeDescriptions[i]
			} else {
				desc = pattern
			}

			cursor := "  "
			style := lipgloss.NewStyle().
				Foreground(m.Styles.Foreground).
				PaddingLeft(1)

			if m.ExcludeListChoice == i {
				cursor = "› "
				style = m.Styles.SelectedMenuItem
			}

			line := fmt.Sprintf("%s%s", cursor, desc)
			listContent.WriteString(style.Render(line) + "\n")
		}
	}

	numPatterns := len(m.Cfg.Config.ExcludeFiles)

	addCursor := "  "
	addStyle := lipgloss.NewStyle().
		Foreground(m.Styles.Secondary).
		PaddingLeft(1)
	if m.ExcludeListChoice == numPatterns {
		addCursor = "› "
		addStyle = m.Styles.SelectedMenuItem
	}
	listContent.WriteString("\n" + addStyle.Render(addCursor+m.t("exclude.add")) + "\n")

	resetCursor := "  "
	resetStyle := lipgloss.NewStyle().
		Foreground(m.Styles.Warning).
		PaddingLeft(1)
	if m.ExcludeListChoice == numPatterns+1 {
		resetCursor = "› "
		resetStyle = m.Styles.SelectedMenuItem
	}
	listContent.WriteString(resetStyle.Render(resetCursor + m.t("exclude.reset")))

	b.WriteString(m.renderPanel(strings.TrimSuffix(listContent.String(), "\n"), m.Styles.Secondary) + "\n\n")

	if m.ExcludeErrorMsg != "" {
		errStyle := m.Styles.ErrorText
		if strings.Contains(m.ExcludeErrorMsg, m.t("exclude.reset.done")) {
			errStyle = m.Styles.SuccessText
		}
		b.WriteString(errStyle.Render(m.ExcludeErrorMsg) + "\n\n")
	}

	b.WriteString(m.Styles.Grid.Render(gridLine) + "\n\n")
	b.WriteString(m.Styles.Hint.Render(m.t("exclude.hint")) + "\n")

	return m.renderScreen(b.String())
}

// renderExcludeEdit 渲染编辑排除模式界面
func (m Model) renderExcludeEdit() string {
	var b strings.Builder

	b.WriteString(m.renderHeaderBlock())

	title := RenderGradientTitle("✏️ " + m.t("exclude.edit.title") + " ✏️")
	b.WriteString(title + "\n\n")

	labelStyle := m.Styles.ConfigKey
	b.WriteString(labelStyle.Render(m.t("exclude.original")) + m.Cfg.Config.ExcludeFiles[m.ExcludeEditIndex] + "\n\n")

	b.WriteString(labelStyle.Render(m.t("exclude.new")))
	inputStyle := lipgloss.NewStyle().
		Foreground(m.Styles.Primary).
		Background(m.Styles.Background).
		Border(lipgloss.NormalBorder()).
		BorderForeground(m.Styles.Secondary).
		Padding(0, 1)
	b.WriteString(inputStyle.Render(m.ExcludeEditInput+"█") + "\n\n")

	b.WriteString(m.renderPanel(m.t("exclude.examples"), m.Styles.Secondary) + "\n\n")

	if m.ExcludeErrorMsg != "" {
		b.WriteString(m.Styles.ErrorText.Render(m.ExcludeErrorMsg) + "\n\n")
	}

	b.WriteString(m.Styles.Grid.Render(gridLine) + "\n\n")
	b.WriteString(m.Styles.Hint.Render(m.t("exclude.edit.hint")) + "\n")

	return m.renderScreen(b.String())
}

// renderExcludeAdd 渲染添加排除模式界面
func (m Model) renderExcludeAdd() string {
	var b strings.Builder

	b.WriteString(m.renderHeaderBlock())

	title := RenderGradientTitle("➕ " + m.t("exclude.add.title") + " ➕")
	b.WriteString(title + "\n\n")

	labelStyle := m.Styles.ConfigKey
	b.WriteString(labelStyle.Render(m.t("exclude.new")))
	inputStyle := lipgloss.NewStyle().
		Foreground(m.Styles.Primary).
		Background(m.Styles.Background).
		Border(lipgloss.NormalBorder()).
		BorderForeground(m.Styles.Secondary).
		Padding(0, 1)
	b.WriteString(inputStyle.Render(m.ExcludeEditInput+"█") + "\n\n")

	infoStyle := m.Styles.WarningText
	b.WriteString(infoStyle.Render(m.t("exclude.add.help")) + "\n\n")

	b.WriteString(m.renderPanel(m.t("exclude.add.examples"), m.Styles.Secondary) + "\n\n")

	if m.ExcludeErrorMsg != "" {
		b.WriteString(m.Styles.ErrorText.Render(m.ExcludeErrorMsg) + "\n\n")
	}

	b.WriteString(m.Styles.Grid.Render(gridLine) + "\n\n")
	b.WriteString(m.Styles.Hint.Render(m.t("exclude.add.hint")) + "\n")

	return m.renderScreen(b.String())
}

// InitExcludeView 初始化排除文件视图
func (m *Model) InitExcludeView() {
	m.ExcludeListChoice = 0
	m.ExcludeErrorMsg = ""
	m.ExcludeEditInput = ""

	descriptions, err := m.Cfg.GetExcludePatternDescriptions()
	if err != nil {
		m.ExcludeErrorMsg = m.t("exclude.error.load", err)
		m.ExcludeDescriptions = []string{}
	} else {
		m.ExcludeDescriptions = descriptions
	}
}
