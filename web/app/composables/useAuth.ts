import type { SignUpRequest, SignInRequest } from '~/types/api'

export function useAuth() {
  const authStore = useAuthStore()
  const api = useApi()
  const router = useRouter()

  const loading = ref(false)
  const error = ref<string | null>(null)

  async function signUp(data: SignUpRequest) {
    loading.value = true
    error.value = null

    try {
      const response = await api.signUp(data)
      authStore.setToken(response.access_token)
      await router.push('/dashboard')
      return true
    } catch (err: any) {
      error.value = err.message || '註冊失敗'
      return false
    } finally {
      loading.value = false
    }
  }

  async function signIn(data: SignInRequest) {
    loading.value = true
    error.value = null

    try {
      const response = await api.signIn(data)
      authStore.setToken(response.access_token)
      await router.push('/dashboard')
      return true
    } catch (err: any) {
      error.value = err.message || '登入失敗'
      return false
    } finally {
      loading.value = false
    }
  }

  function logout() {
    authStore.logout()
  }

  return {
    loading: readonly(loading),
    error: readonly(error),
    isAuthenticated: authStore.isAuthenticated,
    signUp,
    signIn,
    logout
  }
}
