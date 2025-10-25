# Portal Link Web Frontend

Portal Link 的前端應用程式，使用 Nuxt 3 開發。

## 功能特色

- 🔐 使用者註冊與登入
- 📄 建立和管理個人化 Portal Pages
- 🔗 管理多個連結
- 🎨 支援淺色/深色主題
- 📱 響應式設計
- ⚡ 快速的頁面載入

## 技術棧

- **框架**: Nuxt 3
- **語言**: TypeScript
- **樣式**: Tailwind CSS
- **狀態管理**: Pinia
- **工具**: VueUse

## 快速開始

### 環境要求

- Node.js 18+
- npm 或 pnpm

### 安裝

```bash
npm install
```

### 配置環境變數

複製 `.env.example` 為 `.env` 並設定 API 基礎 URL：

```bash
cp .env.example .env
```

編輯 `.env` 檔案：

```env
NUXT_PUBLIC_API_BASE=http://localhost:8080/api/v1
```

### 開發模式

```bash
npm run dev
```

應用程式將在 http://localhost:3000 啟動。

### 建置生產版本

```bash
npm run build
```

### 預覽生產版本

```bash
npm run preview
```

## 專案結構

```
app/
├── assets/          # 靜態資源（CSS、圖片等）
├── components/      # Vue 元件
├── composables/     # Composable 函式
├── layouts/         # 佈局檔案
├── middleware/      # 路由中介層
├── pages/           # 頁面檔案（自動路由）
├── stores/          # Pinia stores
├── types/           # TypeScript 型別定義
└── utils/           # 工具函式
```

## 主要頁面

- `/` - 首頁
- `/signup` - 註冊頁面
- `/signin` - 登入頁面
- `/dashboard` - 使用者儀表板
- `/portal-pages/new` - 建立新的 Portal Page
- `/portal-pages/[id]/edit` - 編輯 Portal Page
- `/[slug]` - 公開的 Portal Page 檢視

## API 整合

前端透過 `useApi()` composable 與後端 API 通訊。所有 API 請求都包含：

- 自動添加 Bearer Token（若已登入）
- 統一的錯誤處理
- TypeScript 型別安全

## 身份驗證

使用 Cookie 儲存 access token，有效期為 1 天。

受保護的路由使用 `auth` middleware，未登入使用者會被重導向至登入頁面。

## 開發注意事項

1. 所有 API 型別定義位於 `app/types/api.ts`
2. 使用 Tailwind CSS 進行樣式設計
3. 支援深色模式（透過 `@nuxtjs/color-mode`）
4. 使用 Composition API 和 `<script setup>` 語法

## 瀏覽器支援

- Chrome (最新版本)
- Firefox (最新版本)
- Safari (最新版本)
- Edge (最新版本)

## 授權

MIT
