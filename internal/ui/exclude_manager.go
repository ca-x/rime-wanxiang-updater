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
				m.ExcludeErrorMsg = fmt.Sprintf("删除失败: %v", err)
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
			m.ExcludeErrorMsg = fmt.Sprintf("重置失败: %v", err)
		} else {
			m.ExcludeErrorMsg = "已重置为默认排除模式"
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
			m.ExcludeErrorMsg = "模式不能为空"
			return m, nil
		}

		if err := m.Cfg.RemoveExcludePattern(m.ExcludeEditIndex); err != nil {
			m.ExcludeErrorMsg = fmt.Sprintf("删除失败: %v", err)
			return m, nil
		}

		if err := m.Cfg.AddExcludePattern(m.ExcludeEditInput); err != nil {
			m.ExcludeErrorMsg = fmt.Sprintf("保存失败: %v", err)
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
			m.ExcludeErrorMsg = "模式不能为空"
			return m, nil
		}

		if err := m.Cfg.AddExcludePattern(m.ExcludeEditInput); err != nil {
			m.ExcludeErrorMsg = fmt.Sprintf("添加失败: %v", err)
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

	b.WriteString(logoStyle.Render(asciiLogo) + "\n")
	b.WriteString(scanLineStyle.Render(scanLine) + "\n\n")

	title := RenderGradientTitle("📋 排除文件管理 📋")
	b.WriteString(title + "\n\n")

	helpText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FFFF")).
		Render("支持三种模式: 通配符(*.yaml) | 正则(^sync/.*$) | 精确(user.yaml)")
	b.WriteString(helpText + "\n\n")

	if len(m.Cfg.Config.ExcludeFiles) == 0 {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")).Render("当前没有排除模式\n\n"))
	} else {
		for i, pattern := range m.Cfg.Config.ExcludeFiles {
			var desc string
			if i < len(m.ExcludeDescriptions) {
				desc = m.ExcludeDescriptions[i]
			} else {
				desc = pattern
			}

			cursor := "  "
			style := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF41"))

			if m.ExcludeListChoice == i {
				cursor = "▸ "
				style = style.Bold(true).Foreground(lipgloss.Color("#FF00FF"))
			}

			line := fmt.Sprintf("%s%s", cursor, desc)
			b.WriteString(style.Render(line) + "\n")
		}
		b.WriteString("\n")
	}

	numPatterns := len(m.Cfg.Config.ExcludeFiles)

	addCursor := "  "
	addStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF"))
	if m.ExcludeListChoice == numPatterns {
		addCursor = "▸ "
		addStyle = addStyle.Bold(true).Foreground(lipgloss.Color("#FF00FF"))
	}
	b.WriteString(addStyle.Render(addCursor+"[添加新模式]") + "\n")

	resetCursor := "  "
	resetStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFAA00"))
	if m.ExcludeListChoice == numPatterns+1 {
		resetCursor = "▸ "
		resetStyle = resetStyle.Bold(true).Foreground(lipgloss.Color("#FF00FF"))
	}
	b.WriteString(resetStyle.Render(resetCursor+"[重置为默认]") + "\n\n")

	if m.ExcludeErrorMsg != "" {
		errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0040"))
		if strings.Contains(m.ExcludeErrorMsg, "成功") || strings.Contains(m.ExcludeErrorMsg, "已重置") {
			errStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF41"))
		}
		b.WriteString(errStyle.Render("⚠ "+m.ExcludeErrorMsg) + "\n\n")
	}

	hints := []string{
		"↑/↓ 选择",
		"Enter 编辑/执行",
		"d/x 删除",
		"q/Esc 返回",
	}
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	b.WriteString(hintStyle.Render(strings.Join(hints, " │ ")) + "\n")

	return b.String()
}

// renderExcludeEdit 渲染编辑排除模式界面
func (m Model) renderExcludeEdit() string {
	var b strings.Builder

	b.WriteString(logoStyle.Render(asciiLogo) + "\n")
	b.WriteString(scanLineStyle.Render(scanLine) + "\n\n")

	title := RenderGradientTitle("✏️  编辑排除模式 ✏️")
	b.WriteString(title + "\n\n")

	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")).Bold(true)
	b.WriteString(labelStyle.Render("原模式: ") + m.Cfg.Config.ExcludeFiles[m.ExcludeEditIndex] + "\n\n")

	b.WriteString(labelStyle.Render("新模式: "))
	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF00FF")).
		Background(lipgloss.Color("#1A1A2E")).
		Padding(0, 1)
	b.WriteString(inputStyle.Render(m.ExcludeEditInput+"█") + "\n\n")

	exampleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	examples := []string{
		"示例:",
		"  *.userdb        (通配符)",
		"  ^sync/.*$       (正则)",
		"  user.yaml       (精确)",
	}
	b.WriteString(exampleStyle.Render(strings.Join(examples, "\n")) + "\n\n")

	if m.ExcludeErrorMsg != "" {
		errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0040"))
		b.WriteString(errStyle.Render("⚠ "+m.ExcludeErrorMsg) + "\n\n")
	}

	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	b.WriteString(hintStyle.Render("Enter 保存 │ Esc 取消") + "\n")

	return b.String()
}

// renderExcludeAdd 渲染添加排除模式界面
func (m Model) renderExcludeAdd() string {
	var b strings.Builder

	b.WriteString(logoStyle.Render(asciiLogo) + "\n")
	b.WriteString(scanLineStyle.Render(scanLine) + "\n\n")

	title := RenderGradientTitle("➕ 添加排除模式 ➕")
	b.WriteString(title + "\n\n")

	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")).Bold(true)
	b.WriteString(labelStyle.Render("新模式: "))
	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF00FF")).
		Background(lipgloss.Color("#1A1A2E")).
		Padding(0, 1)
	b.WriteString(inputStyle.Render(m.ExcludeEditInput+"█") + "\n\n")

	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF41"))
	b.WriteString(infoStyle.Render("支持三种模式类型:") + "\n\n")

	exampleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	examples := []string{
		"1. 通配符模式 (最简单):",
		"   *.userdb           - 所有 userdb 文件",
		"   dicts/*.txt        - dicts 目录下所有 txt 文件",
		"   sync/**/*.yaml     - sync 目录下所有 yaml 文件",
		"",
		"2. 正则表达式 (高级):",
		"   ^sync/.*$          - sync 目录下所有文件",
		"   .*\\.custom\\.yaml$ - 以 .custom.yaml 结尾",
		"",
		"3. 精确匹配:",
		"   installation.yaml  - 只匹配这个文件",
		"   user.yaml          - 只匹配这个文件",
	}
	b.WriteString(exampleStyle.Render(strings.Join(examples, "\n")) + "\n\n")

	if m.ExcludeErrorMsg != "" {
		errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0040"))
		b.WriteString(errStyle.Render("⚠ "+m.ExcludeErrorMsg) + "\n\n")
	}

	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	b.WriteString(hintStyle.Render("Enter 添加 │ Esc 取消") + "\n")

	return b.String()
}

// InitExcludeView 初始化排除文件视图
func (m *Model) InitExcludeView() {
	m.ExcludeListChoice = 0
	m.ExcludeErrorMsg = ""
	m.ExcludeEditInput = ""

	descriptions, err := m.Cfg.GetExcludePatternDescriptions()
	if err != nil {
		m.ExcludeErrorMsg = fmt.Sprintf("加载失败: %v", err)
		m.ExcludeDescriptions = []string{}
	} else {
		m.ExcludeDescriptions = descriptions
	}
}
