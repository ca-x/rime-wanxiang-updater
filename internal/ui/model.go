package ui

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"rime-wanxiang-updater/internal/config"
	"rime-wanxiang-updater/internal/types"
	"rime-wanxiang-updater/internal/updater"
	"rime-wanxiang-updater/internal/version"
)

// ViewState 视图状态
type ViewState int

const (
	ViewWizard ViewState = iota
	ViewMenu
	ViewUpdating
	ViewConfig
	ViewConfigEdit // 新增：配置编辑
	ViewResult     // 新增：显示更新结果
)

// WizardStep 向导步骤
type WizardStep int

const (
	WizardSchemeType WizardStep = iota
	WizardSchemeVariant
	WizardDownloadSource
	WizardComplete
)

// Model Bubble Tea 模型
type Model struct {
	cfg              *config.Manager
	state            ViewState
	wizardStep       WizardStep
	menuChoice       int
	configChoice     int    // 配置菜单选择
	editingKey       string // 正在编辑的配置键
	editingValue     string // 编辑中的值
	schemeChoice     string
	variantChoice    string
	mirrorChoice     bool // 是否使用镜像
	updating         bool
	progress         progress.Model
	progressMsg      string
	downloadSource   string                 // 下载源
	downloadFileName string                 // 下载文件名
	downloaded       int64                  // 已下载字节
	totalSize        int64                  // 总大小字节
	downloadSpeed    float64                // 下载速度
	isDownloading    bool                   // 是否在下载中
	progressChan     chan UpdateMsg         // 进度通道
	completionChan   chan UpdateCompleteMsg // 完成通道
	err              error
	resultMsg        string // 结果消息
	resultSuccess    bool   // 是否成功
	resultSkipped    bool   // 是否跳过更新（已是最新版本）
	width            int
	height           int
}

// NewModel 创建新模型
func NewModel(cfg *config.Manager) Model {
	p := progress.New(progress.WithDefaultGradient())
	p.Width = 60 // 设置默认宽度

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
	message      string
	percent      float64
	source       string  // 下载源
	fileName     string  // 文件名
	downloaded   int64   // 已下载字节
	total        int64   // 总大小字节
	speed        float64 // 下载速度 MB/s
	downloadMode bool    // 是否在下载模式
}

type UpdateCompleteMsg struct {
	err           error
	updateType    string // 更新类型：词库、方案、模型、自动
	skipped       bool   // 是否跳过更新（已是最新版本）
	statusMessage string // 状态消息（包含版本信息）
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
		case ViewConfigEdit:
			return m.handleConfigEditInput(msg)
		case ViewResult:
			return m.handleResultInput(msg)
		case ViewUpdating:
			// 允许用户强制退出更新过程
			switch msg.String() {
			case "q", "esc":
				// 强制返回菜单
				m.state = ViewMenu
				m.updating = false
				// 清理 channel（如果存在）
				m.progressChan = nil
				m.completionChan = nil
				return m, nil
			case "ctrl+c":
				return m, tea.Quit
			}
			return m, nil
		}

	case UpdateMsg:
		m.progressMsg = msg.message

		// 更新下载信息
		if msg.downloadMode {
			m.isDownloading = true
			m.downloadSource = msg.source
			m.downloadFileName = msg.fileName
			m.downloaded = msg.downloaded
			m.totalSize = msg.total
			m.downloadSpeed = msg.speed
		} else {
			m.isDownloading = false
		}

		// 更新进度条 - 移除 >= 0 的检查，允许 0 值
		cmd := m.progress.SetPercent(msg.percent)
		// 继续监听下一个进度消息
		if m.progressChan != nil && m.completionChan != nil {
			return m, tea.Batch(cmd, listenForProgress(m.progressChan, m.completionChan))
		}
		return m, cmd

	case UpdateCompleteMsg:
		m.updating = false
		m.state = ViewResult // 切换到结果视图

		// 清理 channel
		m.progressChan = nil
		m.completionChan = nil

		// 保存 skipped 状态
		m.resultSkipped = msg.skipped

		if msg.err != nil {
			m.resultSuccess = false
			m.resultMsg = fmt.Sprintf("%s更新失败: %v", msg.updateType, msg.err)
		} else if msg.skipped {
			m.resultSuccess = true
			// 如果有状态消息，使用它（包含版本号）
			if msg.statusMessage != "" {
				m.resultMsg = fmt.Sprintf("%s%s", msg.updateType, msg.statusMessage)
			} else {
				m.resultMsg = fmt.Sprintf("%s已是最新版本，无需更新", msg.updateType)
			}
		} else {
			m.resultSuccess = true
			m.resultMsg = fmt.Sprintf("%s更新完成！", msg.updateType)
		}
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
			m.wizardStep = WizardDownloadSource
			return m, nil
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
			m.wizardStep = WizardDownloadSource
			return m, nil
		}

	case WizardDownloadSource:
		switch msg.String() {
		case "1":
			m.mirrorChoice = true
			return m.completeWizard()
		case "2":
			m.mirrorChoice = false
			return m.completeWizard()
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

// completeWizard 完成向导
func (m Model) completeWizard() (tea.Model, tea.Cmd) {
	// 保存镜像选择
	m.cfg.Config.UseMirror = m.mirrorChoice

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
		m.progressMsg = "检查所有更新..."
		return m, m.runAutoUpdate()
	case "2":
		m.state = ViewUpdating
		m.progressMsg = "检查词库更新..."
		return m, m.runDictUpdate()
	case "3":
		m.state = ViewUpdating
		m.progressMsg = "检查方案更新..."
		return m, m.runSchemeUpdate()
	case "4":
		m.state = ViewUpdating
		m.progressMsg = "检查模型更新..."
		return m, m.runModelUpdate()
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
			m.progressMsg = "检查所有更新..."
			return m, m.runAutoUpdate()
		case 1:
			m.state = ViewUpdating
			m.progressMsg = "检查词库更新..."
			return m, m.runDictUpdate()
		case 2:
			m.state = ViewUpdating
			m.progressMsg = "检查方案更新..."
			return m, m.runSchemeUpdate()
		case 3:
			m.state = ViewUpdating
			m.progressMsg = "检查模型更新..."
			return m, m.runModelUpdate()
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
		m.configChoice = 0
		return m, nil
	case "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.configChoice > 0 {
			m.configChoice--
		}
	case "down", "j":
		// 可编辑的配置项数量（不包括 Engine 和排除文件）
		maxChoice := 3 // UseMirror, ProxyEnabled, AutoUpdate

		// Linux 平台添加 fcitx 兼容性配置
		if runtime.GOOS == "linux" {
			maxChoice++ // FcitxCompat
			// 只有启用了 fcitx 兼容，才能选择软链接选项
			if m.cfg.Config.FcitxCompat {
				maxChoice++ // FcitxUseLink
			}
		}

		if m.cfg.Config.ProxyEnabled {
			maxChoice += 2 // ProxyType, ProxyAddress
		}

		// Hook 脚本配置
		maxChoice += 2 // PreUpdateHook, PostUpdateHook

		if m.configChoice < maxChoice {
			m.configChoice++
		}
	case "enter":
		// 根据选择进入编辑模式
		return m.startConfigEdit()
	}
	return m, nil
}

// startConfigEdit 开始编辑配置
func (m Model) startConfigEdit() (tea.Model, tea.Cmd) {
	configItems := []string{"use_mirror", "auto_update", "proxy_enabled"}

	// Linux 平台添加 fcitx 兼容性配置
	if runtime.GOOS == "linux" {
		configItems = append(configItems, "fcitx_compat")
		// 只有启用了 fcitx 兼容，才显示软链接选项
		if m.cfg.Config.FcitxCompat {
			configItems = append(configItems, "fcitx_use_link")
		}
	}

	if m.cfg.Config.ProxyEnabled {
		configItems = append(configItems, "proxy_type", "proxy_address")
	}

	// Hook 脚本配置
	configItems = append(configItems, "pre_update_hook", "post_update_hook")

	if m.configChoice < len(configItems) {
		m.editingKey = configItems[m.configChoice]

		// 设置初始编辑值
		switch m.editingKey {
		case "use_mirror":
			if m.cfg.Config.UseMirror {
				m.editingValue = "true"
			} else {
				m.editingValue = "false"
			}
		case "auto_update":
			if m.cfg.Config.AutoUpdate {
				m.editingValue = "true"
			} else {
				m.editingValue = "false"
			}
		case "proxy_enabled":
			if m.cfg.Config.ProxyEnabled {
				m.editingValue = "true"
			} else {
				m.editingValue = "false"
			}
		case "fcitx_compat":
			if m.cfg.Config.FcitxCompat {
				m.editingValue = "true"
			} else {
				m.editingValue = "false"
			}
		case "fcitx_use_link":
			if m.cfg.Config.FcitxUseLink {
				m.editingValue = "true"
			} else {
				m.editingValue = "false"
			}
		case "proxy_type":
			m.editingValue = m.cfg.Config.ProxyType
		case "proxy_address":
			m.editingValue = m.cfg.Config.ProxyAddress
		case "pre_update_hook":
			m.editingValue = m.cfg.Config.PreUpdateHook
		case "post_update_hook":
			m.editingValue = m.cfg.Config.PostUpdateHook
		}

		m.state = ViewConfigEdit
	}
	return m, nil
}

// handleConfigEditInput 处理配置编辑输入
func (m Model) handleConfigEditInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	isBooleanField := m.editingKey == "use_mirror" || m.editingKey == "auto_update" || m.editingKey == "proxy_enabled" ||
		m.editingKey == "fcitx_compat" || m.editingKey == "fcitx_use_link"

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		// 取消编辑
		m.state = ViewConfig
		m.editingKey = ""
		m.editingValue = ""
		return m, nil
	case "enter":
		// 保存编辑
		return m.saveConfigEdit()
	case "backspace":
		if !isBooleanField && len(m.editingValue) > 0 {
			m.editingValue = m.editingValue[:len(m.editingValue)-1]
		}
	default:
		// 对于布尔值，使用数字选择
		if isBooleanField {
			key := msg.String()
			switch key {
			case "1":
				m.editingValue = "true"
			case "2":
				m.editingValue = "false"
			case "left", "right", "up", "down":
				// 使用方向键切换
				if m.editingValue == "true" {
					m.editingValue = "false"
				} else {
					m.editingValue = "true"
				}
			}
		} else {
			// 其他配置项允许输入
			if len(msg.String()) == 1 {
				m.editingValue += msg.String()
			}
		}
	}
	return m, nil
}

// saveConfigEdit 保存配置编辑
func (m Model) saveConfigEdit() (tea.Model, tea.Cmd) {
	// 更新配置
	switch m.editingKey {
	case "use_mirror":
		m.cfg.Config.UseMirror = m.editingValue == "true"
	case "auto_update":
		m.cfg.Config.AutoUpdate = m.editingValue == "true"
	case "proxy_enabled":
		m.cfg.Config.ProxyEnabled = m.editingValue == "true"
	case "fcitx_compat":
		oldValue := m.cfg.Config.FcitxCompat
		m.cfg.Config.FcitxCompat = m.editingValue == "true"

		// 如果从启用变为禁用，重置选择索引避免越界
		if oldValue && !m.cfg.Config.FcitxCompat {
			m.configChoice = 0
		}

		// 如果启用 fcitx 兼容，立即同步
		if m.cfg.Config.FcitxCompat {
			if err := m.cfg.SyncToFcitxDir(); err != nil {
				m.err = err
			}
		}
	case "fcitx_use_link":
		m.cfg.Config.FcitxUseLink = m.editingValue == "true"
		// 如果已启用 fcitx 兼容，重新同步以应用链接方式的改变
		if m.cfg.Config.FcitxCompat {
			if err := m.cfg.SyncToFcitxDir(); err != nil {
				m.err = err
			}
		}
	case "proxy_type":
		m.cfg.Config.ProxyType = m.editingValue
	case "proxy_address":
		m.cfg.Config.ProxyAddress = m.editingValue
	case "pre_update_hook":
		m.cfg.Config.PreUpdateHook = m.editingValue
	case "post_update_hook":
		m.cfg.Config.PostUpdateHook = m.editingValue
	}

	// 保存到文件
	if err := m.cfg.SaveConfig(); err != nil {
		m.err = err
		m.state = ViewConfig
		return m, nil
	}

	// 返回配置视图
	m.state = ViewConfig
	m.editingKey = ""
	m.editingValue = ""
	return m, nil
}

// handleResultInput 处理结果页面输入
func (m Model) handleResultInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// 按任意键返回主菜单
	if msg.String() == "ctrl+c" {
		return m, tea.Quit
	}
	m.state = ViewMenu
	return m, nil
}

// runDictUpdate 运行词库更新
func (m *Model) runDictUpdate() tea.Cmd {
	// 创建通道
	m.progressChan = make(chan UpdateMsg, 100)
	m.completionChan = make(chan UpdateCompleteMsg, 1)

	// 启动更新 goroutine
	go func() {
		dictUpdater := updater.NewDictUpdater(m.cfg)

		// 进度回调
		progressFunc := func(message string, percent float64, source string, fileName string, downloaded int64, total int64, speed float64, downloadMode bool) {
			select {
			case m.progressChan <- UpdateMsg{
				message:      message,
				percent:      percent,
				source:       source,
				fileName:     fileName,
				downloaded:   downloaded,
				total:        total,
				speed:        speed,
				downloadMode: downloadMode,
			}:
			default:
				// Channel 满了，跳过
			}
		}

		// 检查是否需要更新
		status, err := dictUpdater.GetStatus()
		if err != nil {
			m.completionChan <- UpdateCompleteMsg{err: err, updateType: "词库", skipped: false, statusMessage: ""}
			close(m.progressChan)
			return
		}

		// 如果不需要更新，直接返回
		if !status.NeedsUpdate {
			progressFunc("词库已是最新版本，跳过更新", 1.0, "", "", 0, 0, 0, false)
			m.completionChan <- UpdateCompleteMsg{err: nil, updateType: "词库", skipped: true, statusMessage: status.Message}
			close(m.progressChan)
			return
		}

		// 执行更新
		if err = dictUpdater.Run(progressFunc); err == nil {
			err = dictUpdater.Deploy()
		}

		// 发送完成消息
		m.completionChan <- UpdateCompleteMsg{err: err, updateType: "词库", skipped: false, statusMessage: ""}
		close(m.progressChan)
	}()

	// 返回监听命令
	return listenForProgress(m.progressChan, m.completionChan)
}

// listenForProgress 持续监听进度更新
func listenForProgress(progressChan chan UpdateMsg, completeChan chan UpdateCompleteMsg) tea.Cmd {
	return func() tea.Msg {
		select {
		case msg, ok := <-progressChan:
			if ok {
				return msg
			}
			// Channel 已关闭，尝试非阻塞地读取完成消息
			select {
			case msg := <-completeChan:
				return msg
			default:
				// 如果 completeChan 也没有消息，继续等待
				return <-completeChan
			}
		case msg := <-completeChan:
			return msg
		}
	}
}

// runSchemeUpdate 运行方案更新
func (m *Model) runSchemeUpdate() tea.Cmd {
	// 创建通道
	m.progressChan = make(chan UpdateMsg, 100)
	m.completionChan = make(chan UpdateCompleteMsg, 1)

	// 启动更新 goroutine
	go func() {
		schemeUpdater := updater.NewSchemeUpdater(m.cfg)

		// 进度回调
		progressFunc := func(message string, percent float64, source string, fileName string, downloaded int64, total int64, speed float64, downloadMode bool) {
			select {
			case m.progressChan <- UpdateMsg{
				message:      message,
				percent:      percent,
				source:       source,
				fileName:     fileName,
				downloaded:   downloaded,
				total:        total,
				speed:        speed,
				downloadMode: downloadMode,
			}:
			default:
				// Channel 满了，跳过
			}
		}

		// 检查是否需要更新
		status, err := schemeUpdater.GetStatus()
		if err != nil {
			m.completionChan <- UpdateCompleteMsg{err: err, updateType: "方案", skipped: false, statusMessage: ""}
			close(m.progressChan)
			return
		}

		// 如果不需要更新，直接返回
		if !status.NeedsUpdate {
			progressFunc("方案已是最新版本，跳过更新", 1.0, "", "", 0, 0, 0, false)
			m.completionChan <- UpdateCompleteMsg{err: nil, updateType: "方案", skipped: true, statusMessage: status.Message}
			close(m.progressChan)
			return
		}

		// 执行更新
		if err = schemeUpdater.Run(progressFunc); err == nil {
			err = schemeUpdater.Deploy()
		}

		// 发送完成消息
		m.completionChan <- UpdateCompleteMsg{err: err, updateType: "方案", skipped: false, statusMessage: ""}
		close(m.progressChan)
	}()

	// 返回监听命令
	return listenForProgress(m.progressChan, m.completionChan)
}

// runModelUpdate 运行模型更新
func (m *Model) runModelUpdate() tea.Cmd {
	// 创建通道
	m.progressChan = make(chan UpdateMsg, 100)
	m.completionChan = make(chan UpdateCompleteMsg, 1)

	// 启动更新 goroutine
	go func() {
		modelUpdater := updater.NewModelUpdater(m.cfg)

		// 进度回调
		progressFunc := func(message string, percent float64, source string, fileName string, downloaded int64, total int64, speed float64, downloadMode bool) {
			select {
			case m.progressChan <- UpdateMsg{
				message:      message,
				percent:      percent,
				source:       source,
				fileName:     fileName,
				downloaded:   downloaded,
				total:        total,
				speed:        speed,
				downloadMode: downloadMode,
			}:
			default:
				// Channel 满了，跳过
			}
		}

		// 执行更新
		var err error
		if err = modelUpdater.Run(progressFunc); err == nil {
			err = modelUpdater.Deploy()
		}

		// 发送完成消息
		m.completionChan <- UpdateCompleteMsg{err: err, updateType: "模型", skipped: false, statusMessage: ""}
		close(m.progressChan)
	}()

	// 返回监听命令
	return listenForProgress(m.progressChan, m.completionChan)
}

// runAutoUpdate 运行自动更新
func (m *Model) runAutoUpdate() tea.Cmd {
	// 创建通道
	m.progressChan = make(chan UpdateMsg, 100)
	m.completionChan = make(chan UpdateCompleteMsg, 1)

	// 启动更新 goroutine
	go func() {
		combined := updater.NewCombinedUpdater(m.cfg)

		// 进度回调 - 完整版本，包含下载详情
		progressFunc := func(component, message string, percent float64, source string, fileName string, downloaded int64, total int64, speed float64, downloadMode bool) {
			select {
			case m.progressChan <- UpdateMsg{
				message:      fmt.Sprintf("[%s] %s", component, message),
				percent:      percent,
				source:       source,
				fileName:     fileName,
				downloaded:   downloaded,
				total:        total,
				speed:        speed,
				downloadMode: downloadMode,
			}:
			default:
				// Channel 满了，跳过
			}
		}

		// 检查所有更新
		progressFunc("检查", "正在检查所有更新...", 0.0, "", "", 0, 0, 0, false)
		if err := combined.FetchAllUpdates(); err != nil {
			m.completionChan <- UpdateCompleteMsg{err: err, updateType: "自动", skipped: false, statusMessage: ""}
			close(m.progressChan)
			return
		}

		// 检查是否有任何更新
		if !combined.HasAnyUpdate() {
			progressFunc("完成", "所有组件已是最新版本", 1.0, "", "", 0, 0, 0, false)
			m.completionChan <- UpdateCompleteMsg{err: nil, updateType: "所有组件", skipped: true, statusMessage: "已是最新版本"}
			close(m.progressChan)
			return
		}

		// 执行所有更新
		err := combined.RunAllWithProgress(progressFunc)

		// 发送完成消息
		m.completionChan <- UpdateCompleteMsg{err: err, updateType: "自动", skipped: false, statusMessage: ""}
		close(m.progressChan)
	}()

	// 返回监听命令
	return listenForProgress(m.progressChan, m.completionChan)
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
	case ViewConfigEdit:
		return m.renderConfigEdit()
	case ViewResult:
		return m.renderResult()
	}
	return ""
}

// renderWizard 渲染向导
func (m Model) renderWizard() string {
	var b strings.Builder

	// ASCII Logo
	logo := logoStyle.Render(asciiLogo)
	b.WriteString(logo + "\n")

	// 启动序列状态
	bootSeq := RenderBootSequence(version.GetVersion())
	b.WriteString(bootSeq + "\n")

	// 扫描线效果
	b.WriteString(scanLineStyle.Render(scanLine) + "\n\n")

	// 错误信息
	if m.err != nil {
		errorMsg := errorStyle.Render("⚠ 严重错误 ⚠ " + m.err.Error())
		b.WriteString(errorMsg + "\n\n")
	}

	switch m.wizardStep {
	case WizardSchemeType:
		wizardTitle := RenderGradientTitle("⚡ 初始化向导 ⚡")
		b.WriteString(wizardTitle + "\n\n")

		question := infoBoxStyle.Render("▸ 选择方案版本:")
		b.WriteString(question + "\n\n")

		b.WriteString(menuItemStyle.Render("  [1] ► 万象基础版") + "\n")
		b.WriteString(menuItemStyle.Render("  [2] ► 万象增强版（支持辅助码）") + "\n\n")

		b.WriteString(gridStyle.Render(gridLine) + "\n")
		hint := hintStyle.Render("[>] Input: 1-2 | [Q] Quit")
		b.WriteString(hint)

	case WizardSchemeVariant:
		wizardTitle := RenderGradientTitle("⚡ 初始化向导 ⚡")
		b.WriteString(wizardTitle + "\n\n")

		question := infoBoxStyle.Render("▸ 选择辅助码方案:")
		b.WriteString(question + "\n\n")

		for k, v := range types.SchemeMap {
			b.WriteString(menuItemStyle.Render(fmt.Sprintf("  [%s] ► %s", k, v)) + "\n")
		}

		b.WriteString("\n" + gridStyle.Render(gridLine) + "\n")
		hint := hintStyle.Render("[>] Input: 1-7 | [Q] Quit")
		b.WriteString(hint)

	case WizardDownloadSource:
		wizardTitle := RenderGradientTitle("⚡ 初始化向导 ⚡")
		b.WriteString(wizardTitle + "\n\n")

		question := infoBoxStyle.Render("▸ 选择下载源:")
		b.WriteString(question + "\n\n")

		b.WriteString(menuItemStyle.Render("  [1] ► CNB 镜像（推荐，国内访问更快）") + "\n")
		b.WriteString(menuItemStyle.Render("  [2] ► GitHub 官方源") + "\n\n")

		b.WriteString(gridStyle.Render(gridLine) + "\n")
		hint := hintStyle.Render("[>] Input: 1-2 | [Q] Quit")
		b.WriteString(hint)
	}

	return containerStyle.Render(b.String())
}

// renderMenu 渲染菜单
func (m Model) renderMenu() string {
	var b strings.Builder

	// ASCII Logo
	logo := logoStyle.Render(asciiLogo)
	b.WriteString(logo + "\n")

	// 启动序列状态
	bootSeq := RenderBootSequence(version.GetVersion())
	b.WriteString(bootSeq + "\n")

	// 扫描线效果
	b.WriteString(scanLineStyle.Render(scanLine) + "\n\n")

	// 主菜单标题
	menuTitle := RenderGradientTitle("⚡ 主控制面板 ⚡")
	b.WriteString(menuTitle + "\n\n")

	// 菜单项
	menuItems := []struct {
		icon string
		text string
	}{
		{"▣", "自动更新"},
		{"▣", "词库更新"},
		{"▣", "方案更新"},
		{"▣", "模型更新"},
		{"▣", "查看配置"},
		{"▣", "退出程序"},
	}

	for i, item := range menuItems {
		itemText := fmt.Sprintf(" %s  [%d] %s", item.icon, i+1, item.text)
		if i == m.menuChoice {
			b.WriteString(selectedMenuItemStyle.Render("►"+itemText) + "\n")
		} else {
			b.WriteString(menuItemStyle.Render(" "+itemText) + "\n")
		}
	}

	// 网格线
	b.WriteString("\n" + gridStyle.Render(gridLine) + "\n")

	// 提示
	hint := hintStyle.Render("[>] Input: 1-6 | Navigate: J/K or Arrow Keys | [Q] Quit")
	b.WriteString(hint + "\n\n")

	// 状态栏
	statusBar := RenderStatusBar(
		fmt.Sprintf("v%s", version.GetVersion()),
		m.cfg.Config.Engine,
		func() string {
			if m.cfg.Config.UseMirror {
				return "CNB镜像"
			}
			return "GitHub"
		}(),
	)
	b.WriteString(statusBar)

	return containerStyle.Render(b.String())
}

// renderUpdating 渲染更新中
func (m Model) renderUpdating() string {
	var b strings.Builder

	// ASCII Logo
	logo := logoStyle.Render(asciiLogo)
	b.WriteString(logo + "\n")

	// 启动序列状态
	bootSeq := RenderBootSequence(version.GetVersion())
	b.WriteString(bootSeq + "\n")

	// 处理状态指示器
	status := statusProcessingStyle.Render("⬢ 处理中 ⬢")
	b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(65).Render(status) + "\n\n")

	// 扫描线效果
	b.WriteString(scanLineStyle.Render(scanLine) + "\n\n")

	// 更新标题
	title := RenderGradientTitle("⚡ 正在更新 ⚡")
	b.WriteString(title + "\n\n")

	// 显示进度信息 - 只用一个框显示所有信息
	msgBox := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(neonGreen).
		Padding(1, 2).
		Width(60)

	var msgContent strings.Builder

	// 如果正在下载，显示详细信息
	if m.isDownloading {
		// 下载源和文件名
		if m.downloadSource != "" && m.downloadFileName != "" {
			msgContent.WriteString(configKeyStyle.Render("▸ ") +
				configValueStyle.Render(m.downloadSource) +
				configKeyStyle.Render(" > ") +
				configValueStyle.Render(m.downloadFileName) + "\n\n")
		}

		// 进度条（如果有总大小）
		if m.totalSize > 0 {
			downloadedMB := float64(m.downloaded) / 1024 / 1024
			totalMB := float64(m.totalSize) / 1024 / 1024

			// 进度数字和速度在一行
			progressLine := successStyle.Render(fmt.Sprintf("%.2f MB / %.2f MB", downloadedMB, totalMB))
			if m.downloadSpeed > 0 {
				progressLine += configKeyStyle.Render("  |  ") +
					neonGreenStyle.Render(fmt.Sprintf("%.2f MB/s", m.downloadSpeed))
			}
			msgContent.WriteString(progressLine)
		} else {
			// 没有总大小时显示基本消息
			msgContent.WriteString(progressMsgStyle.Render("▸ " + m.progressMsg))
		}
	} else {
		// 非下载状态显示普通消息
		msgContent.WriteString(progressMsgStyle.Render("▸ " + m.progressMsg))
	}

	b.WriteString(msgBox.Render(msgContent.String()) + "\n\n")

	// 进度条 - 只在下载时显示
	if m.isDownloading && m.totalSize > 0 {
		progressBox := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(neonCyan).
			Padding(0, 1)

		// 直接计算当前百分比
		percent := float64(m.downloaded) / float64(m.totalSize)

		// 使用 ViewAs 直接渲染指定百分比的进度条
		progressBar := progressBox.Render(m.progress.ViewAs(percent))
		b.WriteString(progressBar + "\n\n")
	}

	// 扫描线动画
	b.WriteString(scanLineStyle.Render(scanLine) + "\n\n")

	// 提示
	hint := hintStyle.Render("[...] Please wait... System is updating... | [Q]/[ESC] Cancel | [Ctrl+C] Quit")
	b.WriteString(hint)

	return containerStyle.Render(b.String())
}

// renderConfig 渲染配置
func (m Model) renderConfig() string {
	var b strings.Builder

	// ASCII Logo
	logo := logoStyle.Render(asciiLogo)
	b.WriteString(logo + "\n")

	// 启动序列状态
	bootSeq := RenderBootSequence(version.GetVersion())
	b.WriteString(bootSeq + "\n")

	// 扫描线效果
	b.WriteString(scanLineStyle.Render(scanLine) + "\n\n")

	// 标题
	title := RenderGradientTitle("⚡ 系统配置 ⚡")
	b.WriteString(title + "\n\n")

	// 配置项 - 重新组织，标记可编辑项
	editableConfigs := []struct {
		key      string
		value    string
		editable bool
		index    int
	}{
		{"引擎", m.cfg.Config.Engine, false, -1},
		{"方案类型", m.cfg.Config.SchemeType, false, -1},
		{"方案文件", m.cfg.Config.SchemeFile, false, -1},
		{"词库文件", m.cfg.Config.DictFile, false, -1},
		{"使用镜像", fmt.Sprintf("%v", m.cfg.Config.UseMirror), true, 0},
		{"自动更新", fmt.Sprintf("%v", m.cfg.Config.AutoUpdate), true, 1},
		{"代理启用", fmt.Sprintf("%v", m.cfg.Config.ProxyEnabled), true, 2},
	}

	editIndex := 3

	// Linux 平台添加 fcitx 兼容性配置
	if runtime.GOOS == "linux" {
		editableConfigs = append(editableConfigs,
			struct {
				key      string
				value    string
				editable bool
				index    int
			}{"Fcitx兼容(同步到~/.config/fcitx/rime)", fmt.Sprintf("%v", m.cfg.Config.FcitxCompat), true, editIndex},
		)
		editIndex++

		// 只有启用了 fcitx 兼容，才显示软链接选项
		if m.cfg.Config.FcitxCompat {
			linkMethod := "复制文件"
			if m.cfg.Config.FcitxUseLink {
				linkMethod = "软链接"
			}
			editableConfigs = append(editableConfigs,
				struct {
					key      string
					value    string
					editable bool
					index    int
				}{"同步方式", linkMethod, true, editIndex},
			)
			editIndex++
		}
	}

	if m.cfg.Config.ProxyEnabled {
		editableConfigs = append(editableConfigs,
			struct {
				key      string
				value    string
				editable bool
				index    int
			}{"代理类型", m.cfg.Config.ProxyType, true, editIndex},
			struct {
				key      string
				value    string
				editable bool
				index    int
			}{"代理地址", m.cfg.Config.ProxyAddress, true, editIndex + 1},
		)
		editIndex += 2
	}

	// Hook 脚本配置
	preHookDisplay := m.cfg.Config.PreUpdateHook
	if preHookDisplay == "" {
		preHookDisplay = "(未设置)"
	}
	postHookDisplay := m.cfg.Config.PostUpdateHook
	if postHookDisplay == "" {
		postHookDisplay = "(未设置)"
	}

	editableConfigs = append(editableConfigs,
		struct {
			key      string
			value    string
			editable bool
			index    int
		}{"更新前Hook", preHookDisplay, true, editIndex},
		struct {
			key      string
			value    string
			editable bool
			index    int
		}{"更新后Hook", postHookDisplay, true, editIndex + 1},
	)

	var configContent strings.Builder
	for _, cfg := range editableConfigs {
		key := configKeyStyle.Render(cfg.key + ":")
		value := configValueStyle.Render(cfg.value)
		line := "  ▸ " + key + " " + value

		// 如果是可编辑且被选中，添加高亮
		if cfg.editable && cfg.index == m.configChoice {
			line = selectedMenuItemStyle.Render("►" + line)
		} else {
			line = menuItemStyle.Render(" " + line)
		}

		configContent.WriteString(line + "\n")
	}

	configBox := infoBoxStyle.Render(configContent.String())
	b.WriteString(configBox + "\n\n")

	// 配置文件路径
	pathBox := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(neonPurple).
		Padding(0, 1).
		Foreground(neonPurple)

	pathInfo := pathBox.Render("配置路径: " + m.cfg.ConfigPath)
	b.WriteString(pathInfo + "\n\n")

	// 提示信息
	hint1 := warningStyle.Render("[!] Use Arrow Keys to select, Enter to edit")
	b.WriteString(hint1 + "\n\n")

	b.WriteString(gridStyle.Render(gridLine) + "\n")

	hint2 := hintStyle.Render("[>] Navigate: J/K or Arrow Keys | [Enter] Edit | [Q]/[ESC] Back")
	b.WriteString(hint2)

	return containerStyle.Render(b.String())
}

// renderConfigEdit 渲染配置编辑
func (m Model) renderConfigEdit() string {
	var b strings.Builder

	// ASCII Logo
	logo := logoStyle.Render(asciiLogo)
	b.WriteString(logo + "\n")

	// 启动序列状态
	bootSeq := RenderBootSequence(version.GetVersion())
	b.WriteString(bootSeq + "\n")

	// 扫描线效果
	b.WriteString(scanLineStyle.Render(scanLine) + "\n\n")

	// 标题
	title := RenderGradientTitle("⚡ 编辑配置 ⚡")
	b.WriteString(title + "\n\n")

	// 获取配置项名称
	var configName string
	var inputHint string
	isBooleanField := false
	switch m.editingKey {
	case "use_mirror":
		configName = "使用镜像"
		inputHint = "Select: [1] Enable  [2] Disable | Arrow keys to toggle"
		isBooleanField = true
	case "auto_update":
		configName = "自动更新"
		inputHint = "Select: [1] Enable  [2] Disable | Arrow keys to toggle"
		isBooleanField = true
	case "proxy_enabled":
		configName = "代理启用"
		inputHint = "Select: [1] Enable  [2] Disable | Arrow keys to toggle"
		isBooleanField = true
	case "fcitx_compat":
		configName = "Fcitx兼容"
		inputHint = "启用后将同步配置到 ~/.config/fcitx/rime/ 以兼容外部插件 | [1] Enable  [2] Disable"
		isBooleanField = true
	case "fcitx_use_link":
		configName = "同步方式"
		inputHint = "[1] 软链接(推荐,自动同步,节省空间)  [2] 复制文件(独立,更安全)"
		isBooleanField = true
	case "proxy_type":
		configName = "代理类型"
		inputHint = "Input proxy type: http/https/socks5"
	case "proxy_address":
		configName = "代理地址"
		inputHint = "Input proxy address (e.g. 127.0.0.1:7890)"
	case "pre_update_hook":
		configName = "更新前Hook"
		inputHint = "脚本路径(如~/backup.sh),更新前执行,失败将取消更新"
	case "post_update_hook":
		configName = "更新后Hook"
		inputHint = "脚本路径(如~/notify.sh),更新后执行,失败不影响更新结果"
	}

	// 编辑框
	editBox := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(neonMagenta).
		Padding(1, 2).
		Width(60)

	var editContent strings.Builder
	editContent.WriteString(configKeyStyle.Render("配置项: ") + configValueStyle.Render(configName) + "\n\n")

	// 对于布尔值，显示选项选择界面
	if isBooleanField {
		trueSelected := m.editingValue == "true"
		falseSelected := m.editingValue == "false"

		var trueOption, falseOption string
		if trueSelected {
			trueOption = selectedMenuItemStyle.Render("► [1] Enable (true)")
		} else {
			trueOption = menuItemStyle.Render("  [1] Enable (true)")
		}

		if falseSelected {
			falseOption = selectedMenuItemStyle.Render("► [2] Disable (false)")
		} else {
			falseOption = menuItemStyle.Render("  [2] Disable (false)")
		}

		editContent.WriteString(trueOption + "\n")
		editContent.WriteString(falseOption + "\n\n")
	} else {
		// 非布尔值显示输入框
		editContent.WriteString(configKeyStyle.Render("当前值: "))
		valueWithCursor := m.editingValue + blinkStyle.Render("_")
		editContent.WriteString(successStyle.Render(valueWithCursor) + "\n\n")
	}

	editContent.WriteString(hintStyle.Render(inputHint))

	editBoxRendered := editBox.Render(editContent.String())
	b.WriteString(editBoxRendered + "\n\n")

	// 网格线
	b.WriteString(gridStyle.Render(gridLine) + "\n\n")

	// 提示
	hint := hintStyle.Render("[>] [Enter] Save | [ESC] Cancel | [Backspace] Delete")
	b.WriteString(hint)

	return containerStyle.Render(b.String())
}

// renderResult 渲染更新结果
func (m Model) renderResult() string {
	var b strings.Builder

	// ASCII Logo
	logo := logoStyle.Render(asciiLogo)
	b.WriteString(logo + "\n")

	// 启动序列状态
	bootSeq := RenderBootSequence(version.GetVersion())
	b.WriteString(bootSeq + "\n")

	// 扫描线效果
	b.WriteString(scanLineStyle.Render(scanLine) + "\n\n")

	// 结果标题
	title := RenderGradientTitle("⚡ 更新结果 ⚡")
	b.WriteString(title + "\n\n")

	// 结果消息 - 根据成功/失败使用不同样式
	var resultBox lipgloss.Style
	var icon string

	if m.resultSuccess {
		resultBox = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(neonGreen).
			Padding(2, 3).
			Width(60)
		icon = "✓"
	} else {
		resultBox = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(glitchRed).
			Padding(2, 3).
			Width(60)
		icon = "✗"
	}

	// 消息内容
	var msgContent strings.Builder
	if m.resultSuccess {
		msgContent.WriteString(successStyle.Render(fmt.Sprintf("%s %s", icon, m.resultMsg)))
		// 只有在实际执行了更新时才显示"更新已成功应用到系统"
		if !m.resultSkipped {
			msgContent.WriteString("\n\n")
			msgContent.WriteString(configValueStyle.Render("更新已成功应用到系统"))
		}
	} else {
		msgContent.WriteString(errorStyle.Render(fmt.Sprintf("%s %s", icon, m.resultMsg)))
		msgContent.WriteString("\n\n")
		msgContent.WriteString(configValueStyle.Render("请检查错误信息并重试"))
	}

	resultMessage := resultBox.Render(msgContent.String())
	b.WriteString(resultMessage + "\n\n")

	// 网格线
	b.WriteString(gridStyle.Render(gridLine) + "\n\n")

	// 提示
	hint := blinkStyle.Render("[>] Press any key to return to main menu...")
	b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(65).Render(hint))

	return containerStyle.Render(b.String())
}
