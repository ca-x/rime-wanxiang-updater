package i18n

import "testing"

func TestNormalize(t *testing.T) {
	tests := []struct {
		input string
		want  Locale
	}{
		{input: "", want: DefaultLocale},
		{input: "zh-CN", want: LocaleZhCN},
		{input: "zh", want: LocaleZhCN},
		{input: "en", want: LocaleEn},
		{input: "en-US", want: LocaleEn},
		{input: "fr", want: DefaultLocale},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := Normalize(tt.input); got != tt.want {
				t.Fatalf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestTextFallsBackToDefaultLocale(t *testing.T) {
	if got := Text(Locale("fr"), "menu.auto_update.title"); got != "自动更新" {
		t.Fatalf("Text fallback = %q, want %q", got, "自动更新")
	}
}

func TestTextFormatsTranslatedTemplate(t *testing.T) {
	if got := Text(LocaleEn, "menu.auto_update.countdown", 5); got != "Auto update starts in 5s. Press Esc to cancel." {
		t.Fatalf("Text(LocaleEn, countdown) = %q", got)
	}
}
