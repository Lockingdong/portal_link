-- +goose Up
-- +goose StatementBegin
CREATE TABLE portal_pages (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    bio TEXT,
    profile_image_url VARCHAR(500),
    theme VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for portal_pages
CREATE INDEX idx_portal_pages_user_id ON portal_pages (user_id);
CREATE UNIQUE INDEX idx_portal_pages_slug ON portal_pages (slug);
CREATE INDEX idx_portal_pages_created_at ON portal_pages (created_at);

-- Add comments for portal_pages
COMMENT ON TABLE portal_pages IS 'Portal Page 資料表';
COMMENT ON COLUMN portal_pages.id IS 'Portal Page 的唯一標識符';
COMMENT ON COLUMN portal_pages.user_id IS '擁有此頁面的使用者 ID';
COMMENT ON COLUMN portal_pages.slug IS '頁面的 URL 識別名稱，必須是唯一的';
COMMENT ON COLUMN portal_pages.title IS '頁面標題或顯示名稱';
COMMENT ON COLUMN portal_pages.bio IS '使用者的個人簡介或描述';
COMMENT ON COLUMN portal_pages.profile_image_url IS '個人頭像圖片的 URL';
COMMENT ON COLUMN portal_pages.theme IS '頁面主題設定';
COMMENT ON COLUMN portal_pages.created_at IS '建立時間';
COMMENT ON COLUMN portal_pages.updated_at IS '更新時間';

CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    portal_page_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    url VARCHAR(500) NOT NULL,
    description TEXT,
    icon_url VARCHAR(500),
    display_order INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for links
CREATE INDEX idx_links_portal_page_id ON links (portal_page_id);
CREATE INDEX idx_links_created_at ON links (created_at);

-- Add comments for links
COMMENT ON TABLE links IS 'Link 資料表';
COMMENT ON COLUMN links.id IS 'Link 的唯一標識符';
COMMENT ON COLUMN links.portal_page_id IS '所屬的 Portal Page ID';
COMMENT ON COLUMN links.title IS '連結的顯示標題';
COMMENT ON COLUMN links.url IS '連結的目標 URL';
COMMENT ON COLUMN links.description IS '連結的描述或說明（選填）';
COMMENT ON COLUMN links.icon_url IS '連結的圖示 URL（選填）';
COMMENT ON COLUMN links.display_order IS '連結在頁面上的顯示順序';
COMMENT ON COLUMN links.created_at IS '建立時間';
COMMENT ON COLUMN links.updated_at IS '更新時間';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS links;
DROP TABLE IF EXISTS portal_pages;
-- +goose StatementEnd

