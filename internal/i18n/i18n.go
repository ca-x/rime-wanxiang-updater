package i18n

import (
	"fmt"
	"strings"
)

type Locale string

const (
	LocaleZhCN Locale = "zh-CN"
	LocaleEn   Locale = "en"
)

const DefaultLocale = LocaleZhCN

var catalogs = map[Locale]map[string]string{
	LocaleZhCN: {
		"wizard.title":                             "初始化向导",
		"menu.auto_update.title":                   "自动更新",
		"menu.auto_update.desc":                    "依次检查方案、词库和模型，适合日常维护。",
		"menu.dict_update.title":                   "词库更新",
		"menu.dict_update.desc":                    "只刷新词库文件，适合实时词条更新。",
		"menu.scheme_update.title":                 "方案更新",
		"menu.scheme_update.desc":                  "更新完整方案包，适合升级到新的版本发布。",
		"menu.model_update.title":                  "模型更新",
		"menu.model_update.desc":                   "更新语法模型文件，不影响其他资源。",
		"menu.config.title":                        "查看配置",
		"menu.config.desc":                         "检查下载源、自动更新、代理和 Hook 等设置。",
		"menu.theme.title":                         "切换主题 (%s)",
		"menu.theme.desc":                          "切换配色方案，快速预览当前终端下的主题效果。",
		"menu.custom.title":                        "自定义",
		"menu.custom.desc":                         "调整程序 TUI 界面，并在支持的平台写入主题 patch。",
		"menu.wizard.title":                        "设置向导",
		"menu.wizard.desc":                         "重新选择方案、辅助码和下载源。",
		"menu.quit.title":                          "退出程序",
		"menu.quit.desc":                           "结束当前会话并返回终端。",
		"menu.summary.scheme":                      "当前方案:",
		"menu.summary.version":                     "版本:",
		"menu.summary.source":                      "下载源:",
		"menu.summary.engine":                      "引擎:",
		"menu.summary.theme":                       "主题:",
		"menu.summary.auto_update":                 "自动更新:",
		"menu.auto_update.disabled":                "关闭",
		"menu.auto_update.enabled":                 "已启用",
		"menu.auto_update.countdown":               "自动更新将在 %d 秒后开始... (按 Esc 取消)",
		"menu.auto_update.in":                      "%d 秒后自动开始",
		"menu.auto_update.cancelled":               "已取消自动更新",
		"menu.hint":                                "[1-8] 快捷执行 | J/K 或方向键移动 | Enter 确认 | Q 退出",
		"updating.stage.preparing":                 "准备中",
		"updating.stage":                           "当前阶段: %s",
		"updating.state":                           "状态:",
		"updating.source":                          "来源:",
		"updating.file":                            "文件:",
		"updating.url":                             "下载地址:",
		"updating.progress":                        "进度:",
		"updating.speed":                           "速度:",
		"updating.notice":                          "更新过程中暂不支持取消。完成后会自动返回结果页。",
		"updating.hint":                            "[Ctrl+C] 退出程序 | 其余按键在更新中不会中断任务",
		"result.failure":                           "更新失败",
		"result.success":                           "更新完成",
		"result.skipped":                           "已经是最新状态",
		"result.updated_count":                     "已更新:",
		"result.skipped_count":                     "已跳过:",
		"result.updated_count.value":               "%d 项",
		"result.skipped_count.value":               "%d 项",
		"result.updated_components":                "已更新组件",
		"result.unchanged_components":              "未变更组件",
		"result.hint":                              "按任意键返回主菜单。",
		"ui.badge.failure":                         "失败",
		"ui.badge.success":                         "完成",
		"ui.badge.skipped":                         "已跳过",
		"ui.hint.nav":                              "↑↓ / J K",
		"ui.hint.select":                           "Enter 选择",
		"ui.hint.shortcuts":                        "1-8 快捷操作",
		"ui.hint.edit":                             "Enter 编辑",
		"ui.hint.back":                             "Esc 返回",
		"ui.hint.apply_theme":                      "Enter 应用主题",
		"ui.hint.delete":                           "D / X 删除",
		"ui.hint.save":                             "Enter 保存",
		"ui.hint.add":                              "Enter 添加",
		"ui.hint.exit":                             "Ctrl+C 退出",
		"ui.hint.live_progress":                    "下载详情实时刷新",
		"ui.hint.switch_option":                    "方向键切换选项",
		"ui.hint.menu_return":                      "Enter 返回菜单",
		"ui.hint.about":                            "A 关于",
		"ui.hint.quit":                             "Q 退出",
		"boot.version":                             "Rime Wanxiang Updater · %s",
		"boot.step.init":                           "初始化系统",
		"boot.step.model":                          "加载更新模块",
		"boot.step.connect":                        "连接发布源",
		"boot.step.hardware":                       "扫描运行环境: %s",
		"boot.step.files":                          "挂载工作目录",
		"boot.step.channel":                        "建立安全通道",
		"boot.step.ready":                          "系统就绪",
		"boot.launch":                              "正在进入主界面",
		"boot.exit.line1":                          "本次会话已结束",
		"boot.exit.line2":                          "下次更新再见",
		"wizard.scheme_type":                       "选择方案版本:",
		"wizard.scheme_base":                       "万象基础版",
		"wizard.scheme_pro":                        "万象增强版（支持辅助码）",
		"wizard.variant":                           "选择辅助码方案:",
		"wizard.download_source":                   "选择下载源:",
		"wizard.source.cnb":                        "CNB 镜像（推荐，国内访问更快）",
		"wizard.source.github":                     "GitHub 官方源",
		"wizard.hint.1_2":                          "[1-2] 选择 | [Q] 退出",
		"wizard.hint.1_7":                          "[1-7] 选择 | [Q] 退出",
		"menu.title":                               "主控制面板",
		"config.title":                             "系统配置",
		"config.field.engine":                      "引擎",
		"config.field.scheme_type_name":            "方案类型",
		"config.field.scheme_file":                 "方案文件",
		"config.field.dict_file":                   "词库文件",
		"config.path":                              "配置路径: %s",
		"config.help":                              "使用方向键选择，Enter 编辑",
		"config.hint":                              "J/K 或方向键移动 | Enter 编辑 | Q/Esc 返回",
		"config.field.manage_engines":              "管理更新引擎",
		"config.field.language":                    "界面语言",
		"config.field.use_mirror":                  "使用镜像",
		"config.field.auto_update":                 "自动更新",
		"config.field.auto_update_secs":            "自动更新倒计时(秒)",
		"config.field.proxy_enabled":               "代理启用",
		"config.field.proxy_type":                  "代理类型",
		"config.field.proxy_address":               "代理地址",
		"config.field.pre_hook":                    "更新前 Hook",
		"config.field.post_hook":                   "更新后 Hook",
		"config.field.exclude":                     "管理排除文件",
		"config.field.theme_adaptive":              "自适应主题",
		"config.field.theme_light":                 "浅色主题",
		"config.field.theme_dark":                  "深色主题",
		"config.field.theme_fixed":                 "固定主题",
		"config.field.fcitx_compat":                "Fcitx 兼容(同步到 ~/.config/fcitx/rime)",
		"config.field.fcitx_use_link":              "同步方式",
		"config.edit.title":                        "编辑配置",
		"config.edit.item":                         "配置项:",
		"config.edit.current":                      "当前值:",
		"config.edit.hint.save":                    "[Enter] 保存 | [Esc] 取消 | [Backspace] 删除",
		"config.edit.option.on":                    "启用",
		"config.edit.option.off":                   "禁用",
		"config.edit.hint.bool":                    "[1] %s  [2] %s | 方向键切换",
		"config.edit.hint.language":                "[1] 简体中文  [2] English | 方向键切换",
		"config.edit.hint.countdown":               "输入倒计时秒数 (1-60 秒)",
		"config.edit.hint.fcitx_compat":            "启用后同步到 ~/.config/fcitx/rime/，用于兼容外部插件 | [1] 启用  [2] 禁用",
		"config.edit.hint.fcitx_link":              "[1] 软链接(推荐，自动同步，节省空间)  [2] 复制文件(独立，更安全)",
		"config.edit.hint.proxy_type":              "输入代理类型: http/https/socks5",
		"config.edit.hint.proxy_addr":              "输入代理地址，例如 127.0.0.1:7890",
		"config.edit.hint.pre_hook":                "脚本路径，例如 ~/backup.sh；更新前执行，失败会取消更新",
		"config.edit.hint.post_hook":               "脚本路径，例如 ~/notify.sh；更新后执行，失败不影响更新结果",
		"config.edit.hint.theme":                   "启用后根据终端明暗自动切换主题 | [1] 启用  [2] 禁用",
		"config.option.enable":                     "启用",
		"config.option.disable":                    "禁用",
		"config.language.zh":                       "简体中文",
		"config.language.en":                       "English",
		"config.value.unset":                       "(未设置)",
		"config.value.all_engines":                 "全部引擎",
		"config.value.enabled":                     "启用",
		"config.value.disabled":                    "禁用",
		"config.value.copy":                        "复制文件",
		"config.value.link":                        "软链接",
		"theme.select.title":                       "选择主题",
		"theme.select.dark":                        "选择深色主题",
		"theme.select.light":                       "选择浅色主题",
		"theme.current":                            "当前: %s",
		"theme.adaptive.current":                   " | 自适应模式已启用 (检测: %s背景)",
		"theme.bg.dark":                            "暗色",
		"theme.bg.light":                           "亮色",
		"theme.current_marker":                     " (当前使用)",
		"theme.quick_hint":                         "快速切换会关闭自适应模式 | [Enter] 选择 | [Q]/[Esc] 取消",
		"theme.hint":                               "J/K 或方向键移动 | [Enter] 选择 | [Q]/[Esc] 取消",
		"custom.menu.title":                        "自定义",
		"custom.menu.subtitle":                     "这里包含程序 TUI 界面与 Rime 主题 patch 的快捷入口。",
		"custom.program_tui.title":                 "程序 TUI 界面",
		"custom.program_tui.desc":                  "复用当前程序的界面主题切换器，仅影响本更新器界面。",
		"custom.theme_patch.title":                 "主题 Patch",
		"custom.theme_patch.desc":                  "先多选写入主题预设，再从已选主题中设置默认主题。",
		"custom.theme_patch.target":                "写入目标: %s",
		"custom.theme_patch.hint":                  "Space 勾选主题 | Enter 写入预设并进入下一步",
		"custom.theme_patch.search_label":          "检索: ",
		"custom.theme_patch.search_placeholder":    "直接输入关键字过滤主题",
		"custom.theme_patch.empty":                 "没有匹配的主题，请继续输入或删除关键字。",
		"custom.theme_patch.page":                  "第 %d/%d 页 · 结果 %d/%d",
		"custom.theme_patch.page_empty":            "结果 0/%d",
		"custom.theme_patch.hint.nav":              "↑↓ 选择",
		"custom.theme_patch.hint.search":           "直接输入检索",
		"custom.theme_patch.hint.clear":            "Backspace 删除",
		"custom.theme_patch.hint.toggle":           "Space 勾选",
		"custom.theme_patch.hint.next":             "Enter 下一步",
		"custom.theme_patch.selected_count":        "已选主题: %d",
		"custom.theme_patch.default_title":         "选择默认主题",
		"custom.theme_patch.default_empty":         "请先在上一步选择至少一个主题。",
		"custom.theme_patch.default_hint":          "从已选主题中选择一个作为默认主题。",
		"custom.theme_patch.deploy_title":          "重新部署",
		"custom.theme_patch.deploy_body":           "默认主题已设置为 %s。\n按 Enter 立即重新部署，按其他键返回。",
		"custom.theme_patch.hint.deploy":           "Enter 自动部署",
		"custom.theme_patch.hint.return":           "其他键返回",
		"custom.result.patch_path_error":           "无法确定主题 patch 目标: %v",
		"custom.result.patch_write_error":          "写入主题 patch 失败: %v",
		"custom.result.patch_deploy_success":       "已将默认主题设置为 %s，并完成重新部署。\nPatch 文件: %s",
		"custom.result.patch_deploy_error":         "主题 patch 已保存，但重新部署失败: %v",
		"custom.fcitx_theme.title":                 "Fcitx5 主题",
		"custom.fcitx_theme.desc":                  "同步内置的 Fcitx5 主题集合，再分别设置浅色和深色主题。",
		"custom.fcitx_theme.hint":                  "Space 勾选要保留的主题，程序会同步复制或删除内置主题目录，并保留当前浅色/深色主题供预选。",
		"custom.fcitx_theme.search_label":          "检索: ",
		"custom.fcitx_theme.search_placeholder":    "直接输入关键字过滤主题",
		"custom.fcitx_theme.empty":                 "没有匹配的主题，请继续输入或删除关键字。",
		"custom.fcitx_theme.page":                  "第 %d/%d 页 · 结果 %d/%d",
		"custom.fcitx_theme.page_empty":            "结果 0/%d",
		"custom.fcitx_theme.selected_count":        "已选主题: %d",
		"custom.fcitx_theme.current_light":         "当前浅色",
		"custom.fcitx_theme.current_dark":          "当前深色",
		"custom.fcitx_theme.selected_light":        "已选浅色",
		"custom.fcitx_theme.selected_dark":         "已选深色",
		"custom.fcitx_theme.hint.nav":              "↑↓ 选择",
		"custom.fcitx_theme.hint.search":           "直接输入检索",
		"custom.fcitx_theme.hint.clear":            "Backspace 删除",
		"custom.fcitx_theme.hint.toggle":           "Space 勾选",
		"custom.fcitx_theme.hint.next":             "Enter 下一步",
		"custom.fcitx_theme.default_title_light":   "选择浅色主题",
		"custom.fcitx_theme.default_title_dark":    "选择深色主题",
		"custom.fcitx_theme.default_empty":         "请先在上一步选择至少一个主题。",
		"custom.fcitx_theme.default_hint_light":    "从已勾选主题中选择一个作为浅色模式下的 Fcitx5 主题。",
		"custom.fcitx_theme.default_hint_dark":     "从已勾选主题中选择一个作为深色模式下的 Fcitx5 主题。",
		"custom.fcitx_theme.deploy_title":          "重载 Fcitx5",
		"custom.fcitx_theme.deploy_body":           "浅色主题已设置为 %s，深色主题已设置为 %s。\n程序会启用跟随系统深色模式。\n按 Enter 立即重载 Fcitx5，按其他键返回。",
		"custom.fcitx_theme.hint.deploy":           "Enter 重载",
		"custom.fcitx_theme.hint.return":           "其他键返回",
		"custom.result.fcitx_theme_error":          "Fcitx5 主题操作失败: %v",
		"custom.result.fcitx_theme_deploy_success": "Fcitx5 主题已设置为浅色 %s / 深色 %s，并完成重载。",
		"custom.result.fcitx_theme_deploy_error":   "Fcitx5 主题已设置，但重载失败: %v",
		"about.title":                              "关于界面",
		"about.hero":                               "Rime Wanxiang Updater / Control Surface",
		"about.subtitle":                           "作者、主页与一点程序性浪漫，都在这一页。",
		"about.label.en":                           "Author EN",
		"about.label.zh":                           "作者中文",
		"about.label.home":                         "Homepage",
		"about.name.en":                            "czyt",
		"about.name.zh":                            "虫子樱桃",
		"about.homepage":                           "https://github.com/czyt",
		"about.body":                               "这个界面服务于万象方案更新，也保留一点终端审美。\n冷启动要稳，交互要快，细节要有锋芒。",
		"about.footer":                             "Esc / Q 返回主菜单",
		"exclude.title":                            "排除文件管理",
		"exclude.help":                             "支持三种模式: 通配符(*.yaml) | 正则(^sync/.*$) | 精确(user.yaml)",
		"exclude.empty":                            "当前没有排除模式",
		"exclude.add":                              "[添加新模式]",
		"exclude.reset":                            "[重置为默认]",
		"exclude.hint":                             "↑/↓ 选择 │ Enter 编辑/执行 │ D/X 删除 │ Q/Esc 返回",
		"exclude.edit.title":                       "编辑排除模式",
		"exclude.original":                         "原模式: ",
		"exclude.new":                              "新模式: ",
		"exclude.examples":                         "示例:\n  *.userdb        (通配符)\n  ^sync/.*$       (正则)\n  user.yaml       (精确)",
		"exclude.edit.hint":                        "Enter 保存 │ Esc 取消",
		"exclude.add.title":                        "添加排除模式",
		"exclude.add.help":                         "支持三种模式类型:",
		"exclude.add.examples":                     "1. 通配符模式 (最简单):\n   *.userdb           - 所有 userdb 文件\n   dicts/*.txt        - dicts 目录下所有 txt 文件\n   sync/**/*.yaml     - sync 目录下所有 yaml 文件\n\n2. 正则表达式 (高级):\n   ^sync/.*$          - sync 目录下所有文件\n   .*\\.custom\\.yaml$ - 以 .custom.yaml 结尾\n\n3. 精确匹配:\n   installation.yaml  - 只匹配这个文件\n   user.yaml          - 只匹配这个文件",
		"exclude.add.hint":                         "Enter 添加 │ Esc 取消",
		"exclude.error.delete":                     "删除失败: %v",
		"exclude.error.reset":                      "重置失败: %v",
		"exclude.error.empty":                      "模式不能为空",
		"exclude.error.save":                       "保存失败: %v",
		"exclude.error.add":                        "添加失败: %v",
		"exclude.error.load":                       "加载失败: %v",
		"exclude.reset.done":                       "已重置为默认排除模式",
		"engine.title":                             "选择要更新的引擎",
		"engine.help":                              "使用空格或回车切换选择，按 S 保存",
		"engine.hint":                              "[Space/Enter] 切换 | [S] 保存 | [Q/Esc] 取消",
		"engine.prompt.title":                      "多引擎检测",
		"engine.prompt.message":                    "检测到您安装了多个输入法引擎：%s",
		"engine.prompt.question":                   "您希望如何处理更新？",
		"engine.prompt.manage":                     "进入设置选择要更新的引擎",
		"engine.prompt.all":                        "更新所有已安装的引擎",
		"engine.prompt.hint":                       "[1-2] 选择 | [Q/Esc] 取消",
		"fcitx.title":                              "Fcitx 目录冲突",
		"fcitx.detected":                           "检测到目录已存在: %s",
		"fcitx.question":                           "请选择如何处理:",
		"fcitx.delete":                             "直接删除",
		"fcitx.backup":                             "备份后删除",
		"fcitx.no_prompt":                          "不再提示，记住我的选择",
		"fcitx.hint":                               "[1-2] 或方向键选择 | [Space/Enter] 切换/确认 | [Esc] 取消",
		"updating.title":                           "正在更新",
		"result.title":                             "更新结果",
	},
	LocaleEn: {
		"wizard.title":                             "Setup Wizard",
		"menu.auto_update.title":                   "Auto Update",
		"menu.auto_update.desc":                    "Check scheme, dictionary, and model in one pass.",
		"menu.dict_update.title":                   "Dictionary Update",
		"menu.dict_update.desc":                    "Refresh dictionary files only for rolling lexicon updates.",
		"menu.scheme_update.title":                 "Scheme Update",
		"menu.scheme_update.desc":                  "Update the full scheme package for new releases.",
		"menu.model_update.title":                  "Model Update",
		"menu.model_update.desc":                   "Update the grammar model without touching other assets.",
		"menu.config.title":                        "Settings",
		"menu.config.desc":                         "Review source, auto update, proxy, and hook settings.",
		"menu.theme.title":                         "Theme (%s)",
		"menu.theme.desc":                          "Switch theme and preview the current terminal palette.",
		"menu.custom.title":                        "Customize",
		"menu.custom.desc":                         "Adjust the program TUI and write supported theme patch files.",
		"menu.wizard.title":                        "Setup Wizard",
		"menu.wizard.desc":                         "Re-select scheme, helper code, and download source.",
		"menu.quit.title":                          "Quit",
		"menu.quit.desc":                           "Leave the updater and return to the terminal.",
		"menu.summary.scheme":                      "Scheme:",
		"menu.summary.version":                     "Version:",
		"menu.summary.source":                      "Source:",
		"menu.summary.engine":                      "Engine:",
		"menu.summary.theme":                       "Theme:",
		"menu.summary.auto_update":                 "Auto update:",
		"menu.auto_update.disabled":                "Off",
		"menu.auto_update.enabled":                 "Enabled",
		"menu.auto_update.countdown":               "Auto update starts in %ds. Press Esc to cancel.",
		"menu.auto_update.in":                      "starts in %ds",
		"menu.auto_update.cancelled":               "Auto update cancelled",
		"menu.hint":                                "[1-8] Quick action | J/K or arrows to move | Enter to confirm | Q to quit",
		"updating.stage.preparing":                 "Preparing",
		"updating.stage":                           "Stage: %s",
		"updating.state":                           "Status:",
		"updating.source":                          "Source:",
		"updating.file":                            "File:",
		"updating.url":                             "Download URL:",
		"updating.progress":                        "Progress:",
		"updating.speed":                           "Speed:",
		"updating.notice":                          "Cancellation is not supported during updates. Results appear when the task finishes.",
		"updating.hint":                            "[Ctrl+C] Exit program | Other keys will not interrupt the update",
		"result.failure":                           "Update failed",
		"result.success":                           "Update complete",
		"result.skipped":                           "Already up to date",
		"result.updated_count":                     "Updated:",
		"result.skipped_count":                     "Skipped:",
		"result.updated_count.value":               "%d item(s)",
		"result.skipped_count.value":               "%d item(s)",
		"result.updated_components":                "Updated components",
		"result.unchanged_components":              "Unchanged components",
		"result.hint":                              "Press any key to return to the main menu.",
		"ui.badge.failure":                         "Failed",
		"ui.badge.success":                         "Done",
		"ui.badge.skipped":                         "Skipped",
		"ui.hint.nav":                              "↑↓ / J K",
		"ui.hint.select":                           "Enter Select",
		"ui.hint.shortcuts":                        "1-8 Quick actions",
		"ui.hint.edit":                             "Enter Edit",
		"ui.hint.back":                             "Esc Back",
		"ui.hint.apply_theme":                      "Enter Apply theme",
		"ui.hint.delete":                           "D / X Delete",
		"ui.hint.save":                             "Enter Save",
		"ui.hint.add":                              "Enter Add",
		"ui.hint.exit":                             "Ctrl+C Exit",
		"ui.hint.live_progress":                    "Live download details",
		"ui.hint.switch_option":                    "Arrows switch options",
		"ui.hint.menu_return":                      "Enter Main menu",
		"ui.hint.about":                            "A About",
		"ui.hint.quit":                             "Q Quit",
		"boot.version":                             "Rime Wanxiang Updater · %s",
		"boot.step.init":                           "Initializing system",
		"boot.step.model":                          "Loading update modules",
		"boot.step.connect":                        "Connecting to release source",
		"boot.step.hardware":                       "Scanning runtime: %s",
		"boot.step.files":                          "Mounting workspace",
		"boot.step.channel":                        "Establishing secure channel",
		"boot.step.ready":                          "System ready",
		"boot.launch":                              "Launching main interface",
		"boot.exit.line1":                          "Session complete",
		"boot.exit.line2":                          "See you next update",
		"wizard.scheme_type":                       "Choose a scheme edition:",
		"wizard.scheme_base":                       "Wanxiang Base",
		"wizard.scheme_pro":                        "Wanxiang Pro (with helper code)",
		"wizard.variant":                           "Choose a helper-code layout:",
		"wizard.download_source":                   "Choose a download source:",
		"wizard.source.cnb":                        "CNB Mirror (recommended for domestic access)",
		"wizard.source.github":                     "GitHub",
		"wizard.hint.1_2":                          "[1-2] Select | [Q] Quit",
		"wizard.hint.1_7":                          "[1-7] Select | [Q] Quit",
		"menu.title":                               "Control Panel",
		"config.title":                             "Settings",
		"config.field.engine":                      "Engine",
		"config.field.scheme_type_name":            "Scheme type",
		"config.field.scheme_file":                 "Scheme file",
		"config.field.dict_file":                   "Dictionary file",
		"config.path":                              "Config path: %s",
		"config.help":                              "Use arrow keys to select, then press Enter to edit.",
		"config.hint":                              "J/K or arrows to move | Enter to edit | Q/Esc to go back",
		"config.field.manage_engines":              "Manage update engines",
		"config.field.language":                    "Interface language",
		"config.field.use_mirror":                  "Use mirror",
		"config.field.auto_update":                 "Auto update",
		"config.field.auto_update_secs":            "Auto update countdown (s)",
		"config.field.proxy_enabled":               "Proxy enabled",
		"config.field.proxy_type":                  "Proxy type",
		"config.field.proxy_address":               "Proxy address",
		"config.field.pre_hook":                    "Pre-update hook",
		"config.field.post_hook":                   "Post-update hook",
		"config.field.exclude":                     "Manage excluded files",
		"config.field.theme_adaptive":              "Adaptive theme",
		"config.field.theme_light":                 "Light theme",
		"config.field.theme_dark":                  "Dark theme",
		"config.field.theme_fixed":                 "Fixed theme",
		"config.field.fcitx_compat":                "Fcitx compatibility (sync to ~/.config/fcitx/rime)",
		"config.field.fcitx_use_link":              "Sync mode",
		"config.edit.title":                        "Edit Setting",
		"config.edit.item":                         "Field:",
		"config.edit.current":                      "Current value:",
		"config.edit.hint.save":                    "[Enter] Save | [Esc] Cancel | [Backspace] Delete",
		"config.edit.option.on":                    "Enable",
		"config.edit.option.off":                   "Disable",
		"config.edit.hint.bool":                    "[1] %s  [2] %s | Arrow keys to toggle",
		"config.edit.hint.language":                "[1] Simplified Chinese  [2] English | Arrow keys to toggle",
		"config.edit.hint.countdown":               "Enter a countdown in seconds (1-60)",
		"config.edit.hint.fcitx_compat":            "Sync to ~/.config/fcitx/rime/ for external plugin compatibility | [1] Enable  [2] Disable",
		"config.edit.hint.fcitx_link":              "[1] Symlink (recommended, auto-sync, smaller footprint)  [2] Copy files (isolated, safer)",
		"config.edit.hint.proxy_type":              "Enter proxy type: http/https/socks5",
		"config.edit.hint.proxy_addr":              "Enter proxy address, for example 127.0.0.1:7890",
		"config.edit.hint.pre_hook":                "Script path, for example ~/backup.sh; runs before updates and cancels on failure",
		"config.edit.hint.post_hook":               "Script path, for example ~/notify.sh; runs after updates and does not change the final result",
		"config.edit.hint.theme":                   "Switch themes automatically based on terminal background | [1] Enable  [2] Disable",
		"config.option.enable":                     "Enable",
		"config.option.disable":                    "Disable",
		"config.language.zh":                       "Simplified Chinese",
		"config.language.en":                       "English",
		"config.value.unset":                       "(not set)",
		"config.value.all_engines":                 "All engines",
		"config.value.enabled":                     "Enabled",
		"config.value.disabled":                    "Disabled",
		"config.value.copy":                        "Copy files",
		"config.value.link":                        "Symlink",
		"theme.select.title":                       "Choose Theme",
		"theme.select.dark":                        "Choose Dark Theme",
		"theme.select.light":                       "Choose Light Theme",
		"theme.current":                            "Current: %s",
		"theme.adaptive.current":                   " | Adaptive theme is enabled (detected %s background)",
		"theme.bg.dark":                            "dark",
		"theme.bg.light":                           "light",
		"theme.current_marker":                     " (current)",
		"theme.quick_hint":                         "Quick switch disables adaptive mode | [Enter] Select | [Q]/[Esc] Cancel",
		"theme.hint":                               "J/K or arrows to move | [Enter] Select | [Q]/[Esc] Cancel",
		"custom.menu.title":                        "Customize",
		"custom.menu.subtitle":                     "Shortcuts for the program TUI and supported Rime theme patch flows.",
		"custom.program_tui.title":                 "Program TUI",
		"custom.program_tui.desc":                  "Reuse the current theme picker for this updater interface only.",
		"custom.theme_patch.title":                 "Theme Patch",
		"custom.theme_patch.desc":                  "First write multiple preset themes, then choose one selected theme as default.",
		"custom.theme_patch.target":                "Target file: %s",
		"custom.theme_patch.hint":                  "Space toggles themes | Enter writes presets and continues",
		"custom.theme_patch.search_label":          "Search: ",
		"custom.theme_patch.search_placeholder":    "Type to filter themes",
		"custom.theme_patch.empty":                 "No matching themes. Keep typing or delete keywords.",
		"custom.theme_patch.page":                  "Page %d/%d · Results %d/%d",
		"custom.theme_patch.page_empty":            "Results 0/%d",
		"custom.theme_patch.hint.nav":              "Up/Down Select",
		"custom.theme_patch.hint.search":           "Type to Search",
		"custom.theme_patch.hint.clear":            "Backspace Delete",
		"custom.theme_patch.hint.toggle":           "Space Toggle",
		"custom.theme_patch.hint.next":             "Enter Next",
		"custom.theme_patch.selected_count":        "Selected themes: %d",
		"custom.theme_patch.default_title":         "Choose Default Theme",
		"custom.theme_patch.default_empty":         "Select at least one theme in the previous step first.",
		"custom.theme_patch.default_hint":          "Choose one selected theme as the default theme.",
		"custom.theme_patch.deploy_title":          "Redeploy",
		"custom.theme_patch.deploy_body":           "The default theme is now %s.\nPress Enter to redeploy now, or press any other key to return.",
		"custom.theme_patch.hint.deploy":           "Enter Redeploy",
		"custom.theme_patch.hint.return":           "Any key Return",
		"custom.result.patch_path_error":           "Unable to resolve the theme patch target: %v",
		"custom.result.patch_write_error":          "Failed to write the theme patch: %v",
		"custom.result.patch_deploy_success":       "Default theme set to %s and redeploy completed.\nPatch file: %s",
		"custom.result.patch_deploy_error":         "The theme patch was saved, but redeploy failed: %v",
		"custom.fcitx_theme.title":                 "Fcitx5 Theme",
		"custom.fcitx_theme.desc":                  "Sync bundled Fcitx5 themes, then set separate light and dark themes.",
		"custom.fcitx_theme.hint":                  "Space toggles themes to keep. Bundled theme directories will be copied or removed, and the current light/dark themes are loaded for preselection.",
		"custom.fcitx_theme.search_label":          "Search: ",
		"custom.fcitx_theme.search_placeholder":    "Type to filter themes",
		"custom.fcitx_theme.empty":                 "No matching themes. Keep typing or delete keywords.",
		"custom.fcitx_theme.page":                  "Page %d/%d · Results %d/%d",
		"custom.fcitx_theme.page_empty":            "Results 0/%d",
		"custom.fcitx_theme.selected_count":        "Selected themes: %d",
		"custom.fcitx_theme.current_light":         "Current light",
		"custom.fcitx_theme.current_dark":          "Current dark",
		"custom.fcitx_theme.selected_light":        "Selected light",
		"custom.fcitx_theme.selected_dark":         "Selected dark",
		"custom.fcitx_theme.hint.nav":              "Up/Down Select",
		"custom.fcitx_theme.hint.search":           "Type to Search",
		"custom.fcitx_theme.hint.clear":            "Backspace Delete",
		"custom.fcitx_theme.hint.toggle":           "Space Toggle",
		"custom.fcitx_theme.hint.next":             "Enter Next",
		"custom.fcitx_theme.default_title_light":   "Choose Light Theme",
		"custom.fcitx_theme.default_title_dark":    "Choose Dark Theme",
		"custom.fcitx_theme.default_empty":         "Select at least one theme in the previous step first.",
		"custom.fcitx_theme.default_hint_light":    "Choose one checked theme for Fcitx5 light mode.",
		"custom.fcitx_theme.default_hint_dark":     "Choose one checked theme for Fcitx5 dark mode.",
		"custom.fcitx_theme.deploy_title":          "Reload Fcitx5",
		"custom.fcitx_theme.deploy_body":           "Light theme is set to %s and dark theme is set to %s.\nFollow-system dark mode will be enabled.\nPress Enter to reload Fcitx5 now, or press any other key to return.",
		"custom.fcitx_theme.hint.deploy":           "Enter Reload",
		"custom.fcitx_theme.hint.return":           "Any key Return",
		"custom.result.fcitx_theme_error":          "Fcitx5 theme operation failed: %v",
		"custom.result.fcitx_theme_deploy_success": "Fcitx5 themes set to light %s / dark %s and reload completed.",
		"custom.result.fcitx_theme_deploy_error":   "The Fcitx5 theme was set, but reload failed: %v",
		"about.title":                              "About",
		"about.hero":                               "Rime Wanxiang Updater / Control Surface",
		"about.subtitle":                           "Author, homepage, and a little terminal drama on one screen.",
		"about.label.en":                           "Author EN",
		"about.label.zh":                           "Author ZH",
		"about.label.home":                         "Homepage",
		"about.name.en":                            "czyt",
		"about.name.zh":                            "虫子樱桃",
		"about.homepage":                           "https://github.com/czyt",
		"about.body":                               "Built for Wanxiang maintenance, with enough attitude to avoid a flat utility screen.\nFast paths matter. Clear feedback matters. The interface should, too.",
		"about.footer":                             "Esc / Q returns to the main menu",
		"exclude.title":                            "Excluded Files",
		"exclude.help":                             "Supported patterns: wildcard (*.yaml) | regex (^sync/.*$) | exact match (user.yaml)",
		"exclude.empty":                            "No exclusion patterns configured",
		"exclude.add":                              "[Add new pattern]",
		"exclude.reset":                            "[Reset to default]",
		"exclude.hint":                             "↑/↓ Move │ Enter Edit/Run │ D/X Delete │ Q/Esc Back",
		"exclude.edit.title":                       "Edit Exclusion Pattern",
		"exclude.original":                         "Original: ",
		"exclude.new":                              "New pattern: ",
		"exclude.examples":                         "Examples:\n  *.userdb        (wildcard)\n  ^sync/.*$       (regex)\n  user.yaml       (exact)",
		"exclude.edit.hint":                        "Enter Save │ Esc Cancel",
		"exclude.add.title":                        "Add Exclusion Pattern",
		"exclude.add.help":                         "Supported pattern types:",
		"exclude.add.examples":                     "1. Wildcard (simple):\n   *.userdb           - all userdb files\n   dicts/*.txt        - all txt files under dicts\n   sync/**/*.yaml     - all yaml files under sync\n\n2. Regex (advanced):\n   ^sync/.*$          - every file under sync\n   .*\\.custom\\.yaml$ - files ending with .custom.yaml\n\n3. Exact match:\n   installation.yaml  - only this file\n   user.yaml          - only this file",
		"exclude.add.hint":                         "Enter Add │ Esc Cancel",
		"exclude.error.delete":                     "Delete failed: %v",
		"exclude.error.reset":                      "Reset failed: %v",
		"exclude.error.empty":                      "Pattern cannot be empty",
		"exclude.error.save":                       "Save failed: %v",
		"exclude.error.add":                        "Add failed: %v",
		"exclude.error.load":                       "Load failed: %v",
		"exclude.reset.done":                       "Restored default exclusion patterns",
		"engine.title":                             "Choose Engines to Update",
		"engine.help":                              "Use Space or Enter to toggle, then press S to save.",
		"engine.hint":                              "[Space/Enter] Toggle | [S] Save | [Q/Esc] Cancel",
		"engine.prompt.title":                      "Multiple Engines Detected",
		"engine.prompt.message":                    "Multiple input engines were detected: %s",
		"engine.prompt.question":                   "How should updates be handled?",
		"engine.prompt.manage":                     "Open settings and choose the update engines",
		"engine.prompt.all":                        "Update every installed engine",
		"engine.prompt.hint":                       "[1-2] Select | [Q/Esc] Cancel",
		"fcitx.title":                              "Fcitx Directory Conflict",
		"fcitx.detected":                           "Existing directory detected: %s",
		"fcitx.question":                           "Choose how to continue:",
		"fcitx.delete":                             "Delete directly",
		"fcitx.backup":                             "Backup then delete",
		"fcitx.no_prompt":                          "Remember this choice and stop asking",
		"fcitx.hint":                               "[1-2] or arrows to choose | [Space/Enter] Toggle/Confirm | [Esc] Cancel",
		"updating.title":                           "Updating",
		"result.title":                             "Update Result",
	},
}

func Normalize(raw string) Locale {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	switch normalized {
	case "", "zh", "zh-cn", "zh-hans":
		return LocaleZhCN
	case "en", "en-us", "en-gb":
		return LocaleEn
	default:
		return DefaultLocale
	}
}

func Text(locale Locale, key string, args ...any) string {
	msg := catalogs[DefaultLocale][key]
	if localeCatalog, ok := catalogs[Normalize(string(locale))]; ok {
		if translated, ok := localeCatalog[key]; ok {
			msg = translated
		}
	}
	if msg == "" {
		msg = key
	}
	if len(args) == 0 {
		return msg
	}
	return fmt.Sprintf(msg, args...)
}

func LocaleName(locale Locale, target Locale) string {
	switch Normalize(string(locale)) {
	case LocaleEn:
		return Text(target, "config.language.en")
	default:
		return Text(target, "config.language.zh")
	}
}

func Component(locale Locale, component string) string {
	if Normalize(string(locale)) != LocaleEn {
		return component
	}

	switch component {
	case "方案":
		return "Scheme"
	case "词库":
		return "Dictionary"
	case "模型":
		return "Model"
	case "检查":
		return "Check"
	case "完成":
		return "Done"
	case "准备":
		return "Prepare"
	case "部署":
		return "Deploy"
	case "恢复":
		return "Recover"
	default:
		return component
	}
}

func Source(locale Locale, source string) string {
	if Normalize(string(locale)) != LocaleEn {
		return source
	}

	switch source {
	case "CNB 镜像", "CNB 镜像（推荐，国内访问更快）":
		return "CNB Mirror"
	case "GitHub 官方源":
		return "GitHub"
	default:
		return source
	}
}

func Scheme(locale Locale, scheme string) string {
	switch Normalize(string(locale)) {
	case LocaleEn:
		switch scheme {
		case "base":
			return "Base"
		case "moqi":
			return "Moqi"
		case "flypy":
			return "Flypy"
		case "zrm":
			return "Ziranma"
		case "tiger":
			return "Tiger"
		case "wubi":
			return "Wubi"
		case "hanxin":
			return "Hanxin"
		case "shouyou":
			return "Shouyou"
		case "基础版":
			return "Base"
		case "增强版-墨奇码":
			return "Pro - Moqi"
		case "增强版-小鹤双拼":
			return "Pro - Flypy"
		case "增强版-自然码":
			return "Pro - Ziranma"
		case "增强版-虎码":
			return "Pro - Tiger"
		case "增强版-五笔":
			return "Pro - Wubi"
		case "增强版-汉信":
			return "Pro - Hanxin"
		case "增强版-手语":
			return "Pro - Shouyou"
		case "墨奇码":
			return "Moqi"
		case "小鹤双拼":
			return "Flypy"
		case "自然码":
			return "Ziranma"
		case "虎码":
			return "Tiger"
		case "五笔":
			return "Wubi"
		case "汉信":
			return "Hanxin"
		case "手语":
			return "Shouyou"
		default:
			return scheme
		}
	default:
		switch scheme {
		case "base":
			return "基础版"
		case "moqi":
			return "墨奇码"
		case "flypy":
			return "小鹤双拼"
		case "zrm":
			return "自然码"
		case "tiger":
			return "虎码"
		case "wubi":
			return "五笔"
		case "hanxin":
			return "汉信"
		case "shouyou":
			return "手语"
		default:
			return scheme
		}
	}
}

func RuntimeText(locale Locale, text string) string {
	if Normalize(string(locale)) != LocaleEn {
		return text
	}

	exact := map[string]string{
		"检查所有更新...":      "Checking all updates...",
		"检查词库更新...":      "Checking dictionary updates...",
		"检查方案更新...":      "Checking scheme updates...",
		"检查模型更新...":      "Checking model updates...",
		"正在检查所有更新...":    "Checking all updates...",
		"所有组件已是最新版本":     "All components are already up to date.",
		"已是最新版本":         "Already up to date.",
		"本地文件已是最新版本":     "The local file is already up to date.",
		"正在校验本地文件...":    "Validating local file...",
		"正在计算文件校验和...":   "Calculating file checksum...",
		"正在清理旧文件...":     "Cleaning old files...",
		"正在应用更新...":      "Applying update...",
		"正在终止相关进程...":    "Stopping related processes...",
		"正在解压方案文件...":    "Extracting scheme package...",
		"正在解压词库文件...":    "Extracting dictionary package...",
		"正在同步到其他引擎...":   "Syncing to other engines...",
		"正在同步词库到其他引擎...": "Syncing dictionary to other engines...",
		"正在保存文件...":      "Saving files...",
		"执行更新前 hook...":  "Running pre-update hook...",
		"执行更新后 hook...":  "Running post-update hook...",
		"更新完成！":          "Update complete!",
		"正在保存模型文件...":    "Saving model file...",
		"正在更新方案...":      "Updating scheme...",
		"正在更新词库...":      "Updating dictionary...",
		"正在更新模型...":      "Updating model...",
		"尝试重启服务...":      "Trying to restart services...",
		"所有更新已完成":        "All updates completed.",
		"严重错误":           "Fatal error",
	}
	if translated, ok := exact[text]; ok {
		return translated
	}

	prefixes := []struct {
		old string
		new string
	}{
		{"检查更新失败: ", "Failed to check updates: "},
		{"更新失败: ", "Update failed: "},
		{"获取状态失败: ", "Failed to read status: "},
		{"发现新版本: ", "New version available: "},
		{"检测到可用版本: ", "Available version detected: "},
		{"检测到可用模型: ", "Available model detected: "},
		{"正在检查方案更新 [", "Checking scheme updates ["},
		{"正在检查词库更新 [", "Checking dictionary updates ["},
		{"正在检查模型更新 [", "Checking model updates ["},
		{"准备从 ", "Preparing download from "},
		{"下载失败: ", "Download failed: "},
		{"终止进程失败: ", "Failed to stop processes: "},
		{"解压失败: ", "Extraction failed: "},
		{"处理嵌套目录失败: ", "Failed to process nested directory: "},
		{"同步到其他引擎失败: ", "Failed to sync to other engines: "},
		{"同步词库到其他引擎失败: ", "Failed to sync dictionary to other engines: "},
		{"重命名失败: ", "Rename failed: "},
		{"post-update hook 失败: ", "Post-update hook failed: "},
		{"pre-update hook 失败，已取消更新: ", "Pre-update hook failed, update cancelled: "},
		{"fcitx 同步失败: ", "Fcitx sync failed: "},
		{"下载中: ", "Downloading: "},
		{"下载完成: ", "Download complete: "},
		{"正在部署到 ", "Deploying to "},
		{"部署失败: ", "Deploy failed: "},
		{"更新过程中出现错误: ", "Update finished with errors: "},
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(text, prefix.old) {
			return prefix.new + strings.TrimPrefix(text, prefix.old)
		}
	}

	replacements := strings.NewReplacer(
		"下载方案...", "download scheme...",
		"下载词库...", "download dictionary...",
		"下载模型...", "download model...",
		" (关键文件缺失)", " (required file missing)",
		" (方案已切换，需要更新)", " (scheme changed, update required)",
		" (无版本记录，将重新安装)", " (missing local record, reinstall required)",
		" (当前版本: ", " (current version: ",
		"已切换方案 (从 ", "Scheme changed (from ",
		"未安装", "Not installed",
		"未知版本", "Unknown version",
		"安装和更新完成！", "Install and update complete!",
		"安装完成！", "Installation complete!",
		"更新完成！", "Update complete!",
		"所有组件已是最新版本", "All components are already up to date",
	)
	return replacements.Replace(text)
}
