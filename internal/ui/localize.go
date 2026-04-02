package ui

import (
	"strings"

	"rime-wanxiang-updater/internal/i18n"
)

func (m Model) locale() i18n.Locale {
	if m.Cfg == nil || m.Cfg.Config == nil {
		return i18n.DefaultLocale
	}

	return i18n.Normalize(m.Cfg.Config.Language)
}

func (m Model) t(key string, args ...any) string {
	return i18n.Text(m.locale(), key, args...)
}

func (m Model) runtimeText(text string) string {
	return i18n.RuntimeText(m.locale(), text)
}

func (m Model) componentLabel(component string) string {
	return i18n.Component(m.locale(), component)
}

func (m Model) sourceLabel(source string) string {
	return i18n.Source(m.locale(), source)
}

func (m Model) schemeLabel(scheme string) string {
	return i18n.Scheme(m.locale(), scheme)
}

func (m Model) languageLabel(value string) string {
	return i18n.LocaleName(i18n.Normalize(value), m.locale())
}

func (m Model) boolLabel(value bool) string {
	if value {
		return m.t("config.value.enabled")
	}

	return m.t("config.value.disabled")
}

func (m Model) configFieldLabel(key string) string {
	labels := map[string]string{
		"manage_update_engines": "config.field.manage_engines",
		"use_mirror":            "config.field.use_mirror",
		"auto_update":           "config.field.auto_update",
		"auto_update_countdown": "config.field.auto_update_secs",
		"proxy_enabled":         "config.field.proxy_enabled",
		"proxy_type":            "config.field.proxy_type",
		"proxy_address":         "config.field.proxy_address",
		"pre_update_hook":       "config.field.pre_hook",
		"post_update_hook":      "config.field.post_hook",
		"exclude_file_manager":  "config.field.exclude",
		"theme_adaptive":        "config.field.theme_adaptive",
		"theme_light":           "config.field.theme_light",
		"theme_dark":            "config.field.theme_dark",
		"theme_fixed":           "config.field.theme_fixed",
		"fcitx_compat":          "config.field.fcitx_compat",
		"fcitx_use_link":        "config.field.fcitx_use_link",
		"language":              "config.field.language",
	}

	if labelKey, ok := labels[key]; ok {
		return m.t(labelKey)
	}

	return key
}

func (m Model) localizedValue(value string) string {
	switch strings.TrimSpace(value) {
	case "", "(未设置)":
		return m.t("config.value.unset")
	case "全部引擎":
		return m.t("config.value.all_engines")
	case "true":
		return m.t("config.value.enabled")
	case "false":
		return m.t("config.value.disabled")
	case "启用":
		return m.t("config.value.enabled")
	case "禁用":
		return m.t("config.value.disabled")
	case "复制文件":
		return m.t("config.value.copy")
	case "软链接":
		return m.t("config.value.link")
	default:
		return m.runtimeText(value)
	}
}
