package domain

type PortalPageParams PortalPage

// PortalPage 實體代表使用者的個人化連結整合頁面
// PortalPage 是聚合根（Aggregate Root），負責管理其內部的所有 Link 實體
type PortalPage struct {
}

// NewPortalPage 建立新的 PortalPage 實體
func NewPortalPage(params PortalPageParams) *PortalPage {
	return &PortalPage{}
}
