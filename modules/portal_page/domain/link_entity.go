package domain

// Link 實體代表使用者在 Portal Page 中展示的個別連結項目
// Link 是 Portal Page 聚合內的實體，必須透過 Portal Page（聚合根）來管理
type Link struct {
}

// LinkParams 用於建立或更新 Link 的參數
type LinkParams Link

// NewLink 建立新的 Link 實體（私有方法，只能透過 PortalPage 聚合根調用）
func NewLink(params LinkParams) *Link {
	return &Link{}
}
