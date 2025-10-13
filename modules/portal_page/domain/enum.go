package domain

type Theme string

const (
	ThemeLight Theme = "light"
	ThemeDark  Theme = "dark"
)

// IsValid 檢查 Theme 是否為有效值
func (t Theme) IsValid() bool {
	switch t {
	case ThemeLight, ThemeDark:
		return true
	default:
		return false
	}
}

// GetDefaultTheme 返回預設的 Theme
func GetDefaultTheme() Theme {
	return ThemeLight
}
