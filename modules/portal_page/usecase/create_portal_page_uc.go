package usecase

import (
	"context"
)

// CreatePortalPageParams 建立 Portal Page 用例的輸入參數
type CreatePortalPageParams struct {
}

// CreatePortalPageResult 建立 Portal Page 用例的輸出結果
type CreatePortalPageResult struct {
}

// CreatePortalPageUC 建立 Portal Page 用例
type CreatePortalPageUC struct {
}

func NewCreatePortalPageUC() *CreatePortalPageUC {
	return &CreatePortalPageUC{}
}

func (c *CreatePortalPageUC) Execute(ctx context.Context, params *CreatePortalPageParams) (*CreatePortalPageResult, error) {
	return nil, nil
}
