<template>
  <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <div class="mb-8">
      <NuxtLink to="/dashboard" class="text-primary-600 hover:text-primary-700 flex items-center gap-2 mb-4">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
        返回儀表板
      </NuxtLink>
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">建立 Portal Page</h1>
    </div>

    <form @submit.prevent="handleSubmit" class="space-y-6">
      <div v-if="error" class="rounded-md bg-red-50 dark:bg-red-900/20 p-4">
        <p class="text-sm text-red-800 dark:text-red-200">{{ error }}</p>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6 space-y-6">
        <div>
          <label for="slug" class="block text-sm font-medium mb-1 text-gray-700 dark:text-gray-300">
            自訂網址 (Slug) <span class="text-red-500">*</span>
          </label>
          <div class="mt-1 flex rounded-md shadow-sm">
            <span class="inline-flex items-center px-3 rounded-l-md border border-r-0 border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-700 text-gray-500 dark:text-gray-400 text-sm">
              portallink.com/
            </span>
            <input
              id="slug"
              v-model="form.slug"
              type="text"
              required
              pattern="^[a-z0-9]+(-[a-z0-9]+)*$"
              minlength="3"
              maxlength="50"
              class="flex-1 min-w-0 block w-full px-3 py-2 rounded-none rounded-r-md border border-gray-300 dark:border-gray-600 focus:ring-primary-500 focus:border-primary-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
              placeholder="your-name"
            />
          </div>
          <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
            只能使用小寫字母、數字和連字號，長度 3-50 字元
          </p>
        </div>

        <div>
          <label for="title" class="block text-sm font-medium mb-1 text-gray-700 dark:text-gray-300">
            頁面標題 <span class="text-red-500">*</span>
          </label>
          <input
            id="title"
            v-model="form.title"
            type="text"
            required
            maxlength="100"
            class="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:ring-primary-500 focus:border-primary-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
            placeholder="我的 Portal Page"
          />
        </div>

        <div>
          <label for="bio" class="block text-sm font-medium mb-1 text-gray-700 dark:text-gray-300">
            個人簡介
          </label>
          <textarea
            id="bio"
            v-model="form.bio"
            rows="3"
            maxlength="500"
            class="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:ring-primary-500 focus:border-primary-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
            placeholder="歡迎來到我的頁面..."
          />
          <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
            {{ form.bio?.length || 0 }} / 500
          </p>
        </div>

        <div>
          <label for="profile_image_url" class="block text-sm font-medium mb-1 text-gray-700 dark:text-gray-300">
            頭像圖片 URL
          </label>
          <input
            id="profile_image_url"
            v-model="form.profile_image_url"
            type="url"
            class="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:ring-primary-500 focus:border-primary-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
            placeholder="https://example.com/avatar.jpg"
          />
        </div>

        <div>
          <label for="theme" class="block text-sm font-medium mb-1 text-gray-700 dark:text-gray-300">
            主題
          </label>
          <select
            id="theme"
            v-model="form.theme"
            class="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:ring-primary-500 focus:border-primary-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
          >
            <option value="light">淺色</option>
            <option value="dark">深色</option>
          </select>
        </div>
      </div>

      <div class="flex gap-4">
        <button
          type="submit"
          :disabled="loading"
          class="flex-1 px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <span v-if="loading">建立中...</span>
          <span v-else>建立 Portal Page</span>
        </button>
        <NuxtLink
          to="/dashboard"
          class="px-4 py-2 text-white border border-gray-300 dark:border-gray-600 rounded-md hover:bg-gray-50 dark:hover:bg-gray-800"
        >
          取消
        </NuxtLink>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import type { CreatePortalPageRequest } from '~/types/api'

definePageMeta({
  layout: 'default',
  middleware: 'auth'
})

const api = useApi()
const router = useRouter()

const loading = ref(false)
const error = ref<string | null>(null)

const form = reactive<CreatePortalPageRequest>({
  slug: '',
  title: '',
  bio: '',
  profile_image_url: '',
  theme: 'light'
})

async function handleSubmit() {
  loading.value = true
  error.value = null

  try {
    const response = await api.createPortalPage(form)
    await router.push(`/portal-pages/${response.id}/edit`)
  } catch (err: any) {
    error.value = err.message || '建立失敗'
  } finally {
    loading.value = false
  }
}
</script>
