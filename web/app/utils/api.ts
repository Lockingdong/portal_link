import type {
  SignUpRequest,
  SignUpResponse,
  SignInRequest,
  SignInResponse,
  CreatePortalPageRequest,
  CreatePortalPageResponse,
  UpdatePortalPageRequest,
  UpdatePortalPageResponse,
  FindPortalPageByIDResponse,
  FindPortalPageBySlugResponse,
  ListPortalPagesResponse,
  ErrorResponse
} from '~/types/api'

export class ApiError extends Error {
  constructor(
    public statusCode: number,
    public errorCode: string,
    message: string
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

export function useApi() {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase as string

  async function request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const token = useCookie('access_token').value

    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers
    }

    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }

    try {
      const response = await fetch(`${apiBase}${endpoint}`, {
        ...options,
        headers
      })

      const data = await response.json()

      if (!response.ok) {
        const errorData = data as ErrorResponse
        throw new ApiError(
          response.status,
          errorData.error,
          errorData.message
        )
      }

      return data as T
    } catch (error) {
      if (error instanceof ApiError) {
        throw error
      }
      throw new ApiError(500, 'ErrInternal', '網路錯誤，請稍後再試')
    }
  }

  return {
    // 使用者 API
    signUp: (data: SignUpRequest) =>
      request<SignUpResponse>('/user/signup', {
        method: 'POST',
        body: JSON.stringify(data)
      }),

    signIn: (data: SignInRequest) =>
      request<SignInResponse>('/user/signin', {
        method: 'POST',
        body: JSON.stringify(data)
      }),

    // Portal Page API
    createPortalPage: (data: CreatePortalPageRequest) =>
      request<CreatePortalPageResponse>('/me/portal-pages', {
        method: 'POST',
        body: JSON.stringify(data)
      }),

    listPortalPages: () =>
      request<ListPortalPagesResponse>('/me/portal-pages', {
        method: 'GET'
      }),

    getPortalPageById: (id: number) =>
      request<FindPortalPageByIDResponse>(`/me/portal-pages/${id}`, {
        method: 'GET'
      }),

    updatePortalPage: (id: number, data: UpdatePortalPageRequest) =>
      request<UpdatePortalPageResponse>(`/me/portal-pages/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data)
      }),

    getPortalPageBySlug: (slug: string) =>
      request<FindPortalPageBySlugResponse>(`/portal-pages/${slug}`, {
        method: 'GET'
      })
  }
}
