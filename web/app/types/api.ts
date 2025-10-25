// API 型別定義，基於 OpenAPI 規範

// ===== 使用者相關 =====

export interface SignUpRequest {
  name: string
  email: string
  password: string
}

export interface SignUpResponse {
  access_token: string
}

export interface SignInRequest {
  email: string
  password: string
}

export interface SignInResponse {
  access_token: string
}

// ===== Portal Page 相關 =====

export interface CreatePortalPageRequest {
  slug: string
  title: string
  bio?: string
  profile_image_url?: string
  theme?: 'light' | 'dark'
}

export interface CreatePortalPageResponse {
  id: number
}

export interface UpdatePortalPageRequest {
  slug?: string
  title?: string
  bio?: string
  profile_image_url?: string
  theme?: 'light' | 'dark'
  links: LinkRequest[]
}

export interface UpdatePortalPageResponse {
  id: number
}

export interface LinkRequest {
  id?: number
  title: string
  url: string
  description?: string
  icon_url?: string
  display_order: number
}

export interface LinkDetail {
  id: number
  title: string
  url: string
  description?: string
  icon_url?: string
  display_order: number
}

export interface FindPortalPageByIDResponse {
  id: number
  slug: string
  title: string
  bio?: string
  profile_image_url?: string
  theme: 'light' | 'dark'
  links: LinkDetail[]
}

export interface FindPortalPageBySlugResponse {
  id: number
  slug: string
  title: string
  bio?: string
  profile_image_url?: string
  theme: 'light' | 'dark'
  links: LinkDetail[]
}

export interface PortalPageSummary {
  id: number
  slug: string
  title: string
}

export interface ListPortalPagesResponse {
  portal_pages: PortalPageSummary[]
}

// ===== 錯誤回應 =====

export interface ErrorResponse {
  error: string
  message: string
}

// ===== API 錯誤代碼 =====

export type ErrorCode =
  | 'ErrInvalidParams'
  | 'ErrEmailExists'
  | 'ErrInvalidCredentials'
  | 'ErrUnauthorized'
  | 'ErrForbidden'
  | 'ErrNotFound'
  | 'ErrPortalPageNotFound'
  | 'ErrSlugExists'
  | 'ErrInternal'
