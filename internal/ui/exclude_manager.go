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
		m.state = ViewConfig
		m.excludeErrorMsg = ""
		return m, nil

	case "up", "k":
		if m.excludeListChoice > 0 {
			m.excludeListChoice--
		}

	case "down", "j":
		maxChoice := len(m.cfg.Config.ExcludeFiles) + 2 // +2 for "添加新模式" and "重置为默认"
		if m.excludeListChoice < maxChoice {
			m.excludeListChoice++
		}

	case "enter", " ":
		return m.handleExcludeListSelect()

	case "d", "x":
		// 删除当前选中的模式
		if m.excludeListChoice < len(m.cfg.Config.ExcludeFiles) {
			if err := m.cfg.RemoveExcludePattern(m.excludeListChoice); err != nil {
				m.excludeErrorMsg = fmt.Sprintf("删除失败: %v", err)
			} else {
				m.excludeErrorMsg = ""
				// 更新描述列表
				m.excludeDescriptions, _ = m.cfg.GetExcludePatternDescriptions()
				// 调整光标位置
				if m.excludeListChoice >= len(m.cfg.Config.ExcludeFiles) && m.excludeListChoice > 0 {
					m.excludeListChoice--
				}
			}
		}
	}

	return m, nil
}

// handleExcludeListSelect 处理排除列表选择
func (m Model) handleExcludeListSelect() (Model, tea.Cmd) {
	numPatterns := len(m.cfg.Config.ExcludeFiles)

	if m.excludeListChoice < numPatterns {
		// 编辑现有模式
		m.excludeEditIndex = m.excludeListChoice
		m.excludeEditInput = m.cfg.Config.ExcludeFiles[m.excludeListChoice]
		m.state = ViewExcludeEdit
		m.excludeErrorMsg = ""
	} else if m.excludeListChoice == numPatterns {
		// 添加新模式
		m.excludeEditInput = ""
		m.state = ViewExcludeAdd
		m.excludeErrorMsg = ""
	} else if m.excludeListChoice == numPatterns+1 {
		// 重置为默认
		if err := m.cfg.ResetExcludePatterns(); err != nil {
			m.excludeErrorMsg = fmt.Sprintf("重置失败: %v", err)
		} else {
			m.excludeErrorMsg = "已重置为默认排除模式"
			m.excludeDescriptions, _ = m.cfg.GetExcludePatternDescriptions()
			m.excludeListChoice = 0
		}
	}

	return m, nil
}

// handleExcludeEditInput 处理排除模式编辑输入
func (m Model) handleExcludeEditInput(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = ViewExcludeList
		m.excludeErrorMsg = ""
		return m, nil

	case "enter":
		// 保存编辑
		if strings.TrimSpace(m.excludeEditInput) == "" {
			m.excludeErrorMsg = "模式不能为空"
			return m, nil
		}

		// 先删除旧的，再添加新的
		if err := m.cfg.RemoveExcludePattern(m.excludeEditIndex); err != nil {
			m.excludeErrorMsg = fmt.Sprintf("删除失败: %v", err)
			return m, nil
		}

		if err := m.cfg.AddExcludePattern(m.excludeEditInput); err != nil {
			// 如果添加失败，尝试恢复原来的
			m.excludeErrorMsg = fmt.Sprintf("保存失败: %v", err)
			return m, nil
		}

		m.excludeDescriptions, _ = m.cfg.GetExcludePatternDescriptions()
		m.state = ViewExcludeList
		m.excludeErrorMsg = ""
		return m, nil

	case "backspace":
		if len(m.excludeEditInput) > 0 {
			m.excludeEditInput = m.excludeEditInput[:len(m.excludeEditInput)-1]
		}

	default:
		if len(msg.String()) == 1 {
			m.excludeEditInput += msg.String()
		}
	}

	return m, nil
}

// handleExcludeAddInput 处理添加排除模式输入
func (m Model) handleExcludeAddInput(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = ViewExcludeList
		m.excludeErrorMsg = ""
		return m, nil

	case "enter":
		if strings.TrimSpace(m.excludeEditInput) == "" {
			m.excludeErrorMsg = "模式不能为空"
			return m, nil
		}

		if err := m.cfg.AddExcludePattern(m.excludeEditInput); err != nil {
			m.excludeErrorMsg = fmt.Sprintf("添加失败: %v", err)
			return m, nil
		}

		m.excludeDescriptions, _ = m.cfg.GetExcludePatternDescriptions()
		m.state = ViewExcludeList
		m.excludeErrorMsg = ""
		m.excludeEditInput = ""
		return m, nil

	case "backspace":
		if len(m.excludeEditInput) > 0 {
			m.excludeEditInput = m.excludeEditInput[:len(m.excludeEditInput)-1]
		}

	default:
		if len(msg.String()) == 1 {
			m.excludeEditInput += msg.String()
		}
	}

	return m, nil
}

// renderExcludeList 渲染排除文件列表
func (m Model) renderExcludeList() string {
	var b strings.Builder

	// Logo
	b.WriteString(logoStyle.Render(asciiLogo) + "\n")
	b.WriteString(scanLineStyle.Render(scanLine) + "\n\n")

	// 标题
	title := RenderGradientTitle("📋 排除文件管理 📋")
	b.WriteString(title + "\n\n")

	// 说明
	helpText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FFFF")).
		Render("支持三种模式: 通配符(*.yaml) | 正则(^sync/.*$) | 精确(user.yaml)")
	b.WriteString(helpText + "\n\n")

	// 当前排除模式列表
	if len(m.cfg.Config.ExcludeFiles) == 0 {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")).Render("当前没有排除模式\n\n"))
	} else {
		for i, pattern := range m.cfg.Config.ExcludeFiles {
			var desc string
			if i < len(m.excludeDescriptions) {
				desc = m.excludeDescriptions[i]
			} else {
				desc = pattern
			}

			cursor := "  "
			style := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF41"))

			if m.excludeListChoice == i {
				cursor = "▸ "
				style = style.Bold(true).Foreground(lipgloss.Color("#FF00FF"))
			}

			line := fmt.Sprintf("%s%s", cursor, desc)
			b.WriteString(style.Render(line) + "\n")
		}
		b.WriteString("\n")
	}

	// 操作选项
	numPatterns := len(m.cfg.Config.ExcludeFiles)

	// 添加新模式
	addCursor := "  "
	addStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF"))
	if m.excludeListChoice == numPatterns {
		addCursor = "▸ "
		addStyle = addStyle.Bold(true).Foreground(lipgloss.Color("#FF00FF"))
	}
	b.WriteString(addStyle.Render(addCursor+"[添加新模式]") + "\n")

	// 重置为默认
	resetCursor := "  "
	resetStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFAA00"))
	if m.excludeListChoice == numPatterns+1 {
		resetCursor = "▸ "
		resetStyle = resetStyle.Bold(true).Foreground(lipgloss.Color("#FF00FF"))
	}
	b.WriteString(resetStyle.Render(resetCursor+"[重置为默认]") + "\n\n")

	// 错误消息
	if m.excludeErrorMsg != "" {
		errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0040"))
		if strings.Contains(m.excludeErrorMsg, "成功") || strings.Contains(m.excludeErrorMsg, "已重置") {
			errStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF41"))
		}
		b.WriteString(errStyle.Render("⚠ "+m.excludeErrorMsg) + "\n\n")
	}

	// 操作提示
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

	// 当前编辑的模式
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")).Bold(true)
	b.WriteString(labelStyle.Render("原模式: ") + m.cfg.Config.ExcludeFiles[m.excludeEditIndex] + "\n\n")

	// 输入框
	b.WriteString(labelStyle.Render("新模式: "))
	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF00FF")).
		Background(lipgloss.Color("#1A1A2E")).
		Padding(0, 1)
	b.WriteString(inputStyle.Render(m.excludeEditInput+"█") + "\n\n")

	// 示例
	exampleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	examples := []string{
		"示例:",
		"  *.userdb        (通配符)",
		"  ^sync/.*$       (正则)",
		"  user.yaml       (精确)",
	}
	b.WriteString(exampleStyle.Render(strings.Join(examples, "\n")) + "\n\n")

	// 错误消息
	if m.excludeErrorMsg != "" {
		errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0040"))
		b.WriteString(errStyle.Render("⚠ "+m.excludeErrorMsg) + "\n\n")
	}

	// 提示
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

	// 输入框
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")).Bold(true)
	b.WriteString(labelStyle.Render("新模式: "))
	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF00FF")).
		Background(lipgloss.Color("#1A1A2E")).
		Padding(0, 1)
	b.WriteString(inputStyle.Render(m.excludeEditInput+"█") + "\n\n")

	// 说明和示例
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

	// 错误消息
	if m.excludeErrorMsg != "" {
		errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0040"))
		b.WriteString(errStyle.Render("⚠ "+m.excludeErrorMsg) + "\n\n")
	}

	// 提示
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	b.WriteString(hintStyle.Render("Enter 添加 │ Esc 取消") + "\n")

	return b.String()
}

// InitExcludeView 初始化排除文件视图
func (m *Model) InitExcludeView() {
	m.excludeListChoice = 0
	m.excludeErrorMsg = ""
	m.excludeEditInput = ""

	// 加载描述
	descriptions, err := m.cfg.GetExcludePatternDescriptions()
	if err != nil {
		m.excludeErrorMsg = fmt.Sprintf("加载失败: %v", err)
		m.excludeDescriptions = []string{}
	} else {
		m.excludeDescriptions = descriptions
	}
}
