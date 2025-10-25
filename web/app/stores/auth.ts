import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', () => {
  const accessToken = useCookie('access_token', {
    maxAge: 60 * 60 * 24 // 1 天
  })

  const isAuthenticated = computed(() => !!accessToken.value)

  function setToken(token: string) {
    accessToken.value = token
  }

  function clearToken() {
    accessToken.value = null
  }

  function logout() {
    clearToken()
    navigateTo('/signin')
  }

  return {
    accessToken: readonly(accessToken),
    isAuthenticated,
    setToken,
    clearToken,
    logout
  }
})
