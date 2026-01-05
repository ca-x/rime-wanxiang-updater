package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"rime-wanxiang-updater/internal/config"
	"rime-wanxiang-updater/internal/types"
	"rime-wanxiang-updater/internal/updater"
)

// ViewState 视图状态
type ViewState int

const (
	ViewWizard ViewState = iota
	ViewMenu
	ViewUpdating
	ViewConfig
)

// WizardStep 向导步骤
type WizardStep int

const (
	WizardSchemeType WizardStep = iota
	WizardSchemeVariant
	WizardComplete
)

// Model Bubble Tea 模型
type Model struct {
	cfg           *config.Manager
	state         ViewState
	wizardStep    WizardStep
	menuChoice    int
	schemeChoice  string
	variantChoice string
	updating      bool
	progress      progress.Model
	progressMsg   string
	err           error
	width         int
	height        int
}

// NewModel 创建新模型
func NewModel(cfg *config.Manager) Model {
	p := progress.New(progress.WithDefaultGradient())

	// 检查是否需要首次配置
	state := ViewMenu
	wizardStep := WizardSchemeType
	if cfg.Config.SchemeType == "" || cfg.Config.SchemeFile == "" || cfg.Config.DictFile == "" {
		state = ViewWizard
	}

	return Model{
		cfg:        cfg,
		state:      state,
		wizardStep: wizardStep,
		progress:   p,
	}
}

// Init 初始化
func (m Model) Init() tea.Cmd {
	return nil
}

// UpdateMsg 更新消息类型
type UpdateMsg struct {
	message string
	percent float64
}

type UpdateCompleteMsg struct {
	err error
}

// Update 更新模型
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.progress.Width = msg.Width - 4
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case ViewWizard:
			return m.handleWizardInput(msg)
		case ViewMenu:
			return m.handleMenuInput(msg)
		case ViewConfig:
			return m.handleConfigInput(msg)
		case ViewUpdating:
			// 更新中不接受输入
			return m, nil
		}

	case UpdateMsg:
		m.progressMsg = msg.message
		if msg.percent >= 0 {
			cmd := m.progress.SetPercent(msg.percent)
			return m, cmd
		}
		return m, nil

	case UpdateCompleteMsg:
		m.updating = false
		m.err = msg.err
		m.state = ViewMenu
		return m, nil
	}

	return m, nil
}

// handleWizardInput 处理向导输入
func (m Model) handleWizardInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.wizardStep {
	case WizardSchemeType:
		switch msg.String() {
		case "1":
			m.cfg.Config.SchemeType = "base"
			m.schemeChoice = "base"
			return m.completeWizard()
		case "2":
			m.cfg.Config.SchemeType = "pro"
			m.wizardStep = WizardSchemeVariant
			return m, nil
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case WizardSchemeVariant:
		key := msg.String()
		if key == "q" || key == "ctrl+c" {
			return m, tea.Quit
		}
		if variant, ok := types.SchemeMap[key]; ok {
			m.schemeChoice = variant
			return m.completeWizard()
		}
	}
	return m, nil
}

// completeWizard 完成向导
func (m Model) completeWizard() (tea.Model, tea.Cmd) {
	// 获取实际文件名
	schemeFile, dictFile, err := m.cfg.GetActualFilenames(m.schemeChoice)
	if err != nil {
		m.err = err
		return m, nil
	}

	m.cfg.Config.SchemeFile = schemeFile
	m.cfg.Config.DictFile = dictFile

	// 保存配置
	if err := m.cfg.SaveConfig(); err != nil {
		m.err = err
		return m, nil
	}

	m.wizardStep = WizardComplete
	m.state = ViewMenu
	return m, nil
}

// handleMenuInput 处理菜单输入
func (m Model) handleMenuInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "1":
		m.state = ViewUpdating
		m.progressMsg = "检查词库更新..."
		return m, m.runDictUpdate()
	case "2":
		m.state = ViewUpdating
		m.progressMsg = "检查方案更新..."
		return m, m.runSchemeUpdate()
	case "3":
		m.state = ViewUpdating
		m.progressMsg = "检查模型更新..."
		return m, m.runModelUpdate()
	case "4":
		m.state = ViewUpdating
		m.progressMsg = "检查所有更新..."
		return m, m.runAutoUpdate()
	case "5":
		m.state = ViewConfig
		return m, nil
	case "6", "q", "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.menuChoice > 0 {
			m.menuChoice--
		}
	case "down", "j":
		if m.menuChoice < 5 {
			m.menuChoice++
		}
	case "enter":
		switch m.menuChoice {
		case 0:
			m.state = ViewUpdating
			m.progressMsg = "检查词库更新..."
			return m, m.runDictUpdate()
		case 1:
			m.state = ViewUpdating
			m.progressMsg = "检查方案更新..."
			return m, m.runSchemeUpdate()
		case 2:
			m.state = ViewUpdating
			m.progressMsg = "检查模型更新..."
			return m, m.runModelUpdate()
		case 3:
			m.state = ViewUpdating
			m.progressMsg = "检查所有更新..."
			return m, m.runAutoUpdate()
		case 4:
			m.state = ViewConfig
			return m, nil
		case 5:
			return m, tea.Quit
		}
	}
	return m, nil
}

// handleConfigInput 处理配置输入
func (m Model) handleConfigInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		m.state = ViewMenu
		return m, nil
	case "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}

// runDictUpdate 运行词库更新
func (m Model) runDictUpdate() tea.Cmd {
	return func() tea.Msg {
		dictUpdater := updater.NewDictUpdater(m.cfg)
		if err := dictUpdater.Run(); err != nil {
			return UpdateCompleteMsg{err: err}
		}

		if err := dictUpdater.Deploy(); err != nil {
			return UpdateCompleteMsg{err: err}
		}

		return UpdateCompleteMsg{err: nil}
	}
}

// runSchemeUpdate 运行方案更新
func (m Model) runSchemeUpdate() tea.Cmd {
	return func() tea.Msg {
		schemeUpdater := updater.NewSchemeUpdater(m.cfg)
		if err := schemeUpdater.Run(); err != nil {
			return UpdateCompleteMsg{err: err}
		}

		if err := schemeUpdater.Deploy(); err != nil {
			return UpdateCompleteMsg{err: err}
		}

		return UpdateCompleteMsg{err: nil}
	}
}

// runModelUpdate 运行模型更新
func (m Model) runModelUpdate() tea.Cmd {
	return func() tea.Msg {
		modelUpdater := updater.NewModelUpdater(m.cfg)
		if err := modelUpdater.Run(); err != nil {
			return UpdateCompleteMsg{err: err}
		}

		if err := modelUpdater.Deploy(); err != nil {
			return UpdateCompleteMsg{err: err}
		}

		return UpdateCompleteMsg{err: nil}
	}
}

// runAutoUpdate 运行自动更新
func (m Model) runAutoUpdate() tea.Cmd {
	return func() tea.Msg {
		combined := updater.NewCombinedUpdater(m.cfg)
		if err := combined.FetchAllUpdates(); err != nil {
			return UpdateCompleteMsg{err: err}
		}

		if !combined.HasAnyUpdate() {
			return UpdateCompleteMsg{err: fmt.Errorf("所有组件均为最新版本")}
		}

		if err := combined.RunAll(); err != nil {
			return UpdateCompleteMsg{err: err}
		}

		return UpdateCompleteMsg{err: nil}
	}
}

// View 渲染视图
func (m Model) View() string {
	switch m.state {
	case ViewWizard:
		return m.renderWizard()
	case ViewMenu:
		return m.renderMenu()
	case ViewUpdating:
		return m.renderUpdating()
	case ViewConfig:
		return m.renderConfig()
	}
	return ""
}

// renderWizard 渲染向导
func (m Model) renderWizard() string {
	var b strings.Builder

	// 标题
	title := headerStyle.Render("Rime 万象输入法更新工具 " + types.VERSION)
	b.WriteString("\n" + title + "\n\n")

	// 错误信息
	if m.err != nil {
		errorMsg := errorStyle.Render("❌ 错误: " + m.err.Error())
		b.WriteString(errorMsg + "\n\n")
	}

	switch m.wizardStep {
	case WizardSchemeType:
		wizardTitle := titleStyle.Render("🔧 首次运行配置向导")
		b.WriteString(wizardTitle + "\n\n")

		question := infoBoxStyle.Render("请选择方案版本:")
		b.WriteString(question + "\n\n")

		b.WriteString(menuItemStyle.Render("[1] 万象基础版") + "\n")
		b.WriteString(menuItemStyle.Render("[2] 万象增强版（支持各种辅助码）") + "\n\n")

		hint := hintStyle.Render("请输入选择 (1-2, q 退出)")
		b.WriteString(hint)

	case WizardSchemeVariant:
		question := infoBoxStyle.Render("请选择辅助码方案:")
		b.WriteString(question + "\n\n")

		for k, v := range types.SchemeMap {
			b.WriteString(menuItemStyle.Render(fmt.Sprintf("[%s] %s", k, v)) + "\n")
		}

		hint := hintStyle.Render("\n请输入选择 (1-7, q 退出)")
		b.WriteString(hint)
	}

	return containerStyle.Render(b.String())
}

// renderMenu 渲染菜单
func (m Model) renderMenu() string {
	var b strings.Builder

	// 标题
	title := headerStyle.Render("Rime 万象输入法更新工具 " + types.VERSION)
	b.WriteString("\n" + title + "\n\n")

	// 消息显示
	if m.err != nil {
		if m.err.Error() == "所有组件均为最新版本" {
			msg := successStyle.Render("✓ " + m.err.Error())
			b.WriteString(msg + "\n\n")
		} else {
			msg := errorStyle.Render("❌ 错误: " + m.err.Error())
			b.WriteString(msg + "\n\n")
		}
		m.err = nil
	}

	// 主菜单标题
	menuTitle := titleStyle.Render("📋 主菜单")
	b.WriteString(menuTitle + "\n\n")

	// 菜单项
	menuItems := []string{
		"📚 词库下载",
		"⚙️  方案下载",
		"🤖 模型下载",
		"🔄 自动更新",
		"🔧 修改配置",
		"❌ 退出程序",
	}

	for i, item := range menuItems {
		if i == m.menuChoice {
			b.WriteString(selectedMenuItemStyle.Render(fmt.Sprintf("▶ [%d] %s", i+1, item)) + "\n")
		} else {
			b.WriteString(menuItemStyle.Render(fmt.Sprintf("  [%d] %s", i+1, item)) + "\n")
		}
	}

	// 提示
	hint := hintStyle.Render("\n请输入选择 (1-6, ↑↓/jk 导航, q 退出)")
	b.WriteString(hint)

	return containerStyle.Render(b.String())
}

// renderUpdating 渲染更新中
func (m Model) renderUpdating() string {
	var b strings.Builder

	// 标题
	title := headerStyle.Render("正在更新...")
	b.WriteString("\n" + title + "\n\n")

	// 进度消息
	msg := progressMsgStyle.Render(m.progressMsg)
	b.WriteString(msg + "\n\n")

	// 进度条
	progressBar := infoBoxStyle.Render(m.progress.View())
	b.WriteString(progressBar + "\n\n")

	// 提示
	hint := hintStyle.Render("请稍候...")
	b.WriteString(hint)

	return containerStyle.Render(b.String())
}

// renderConfig 渲染配置
func (m Model) renderConfig() string {
	var b strings.Builder

	// 标题
	title := headerStyle.Render("当前配置")
	b.WriteString("\n" + title + "\n\n")

	// 配置项
	configs := []struct {
		key   string
		value string
	}{
		{"引擎", m.cfg.Config.Engine},
		{"方案类型", m.cfg.Config.SchemeType},
		{"方案文件", m.cfg.Config.SchemeFile},
		{"词库文件", m.cfg.Config.DictFile},
		{"使用镜像", fmt.Sprintf("%v", m.cfg.Config.UseMirror)},
		{"GitHub Token", m.cfg.Config.GithubToken},
		{"排除文件", fmt.Sprintf("%v", m.cfg.Config.ExcludeFiles)},
		{"自动更新", fmt.Sprintf("%v", m.cfg.Config.AutoUpdate)},
		{"代理启用", fmt.Sprintf("%v", m.cfg.Config.ProxyEnabled)},
	}

	if m.cfg.Config.ProxyEnabled {
		configs = append(configs,
			struct {
				key   string
				value string
			}{"代理类型", m.cfg.Config.ProxyType},
			struct {
				key   string
				value string
			}{"代理地址", m.cfg.Config.ProxyAddress},
		)
	}

	var configContent strings.Builder
	for _, cfg := range configs {
		key := configKeyStyle.Render(cfg.key + ":")
		value := configValueStyle.Render(cfg.value)
		configContent.WriteString(key + " " + value + "\n")
	}

	configBox := infoBoxStyle.Render(configContent.String())
	b.WriteString(configBox + "\n")

	// 配置文件路径
	pathInfo := hintStyle.Render("配置文件路径: " + m.cfg.ConfigPath)
	b.WriteString(pathInfo + "\n\n")

	// 提示信息
	hint1 := warningStyle.Render("⚠ 提示: 可以手动编辑配置文件来修改设置")
	b.WriteString(hint1 + "\n\n")

	hint2 := hintStyle.Render("按 q 或 ESC 返回主菜单")
	b.WriteString(hint2)

	return containerStyle.Render(b.String())
}
