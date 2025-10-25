<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">我的 Portal Pages</h1>
      <NuxtLink
        to="/portal-pages/new"
        class="px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 flex items-center gap-2"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        建立新頁面
      </NuxtLink>
    </div>

    <div v-if="loading" class="text-center py-12">
      <p class="text-gray-500 dark:text-gray-400">載入中...</p>
    </div>

    <div v-else-if="error" class="rounded-md bg-red-50 dark:bg-red-900/20 p-4">
      <p class="text-sm text-red-800 dark:text-red-200">{{ error }}</p>
    </div>

    <div v-else-if="portalPages.length === 0" class="text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400 dark:text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900 dark:text-white">尚無 Portal Page</h3>
      <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">開始建立您的第一個 Portal Page</p>
      <div class="mt-6">
        <NuxtLink
          to="/portal-pages/new"
          class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700"
        >
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          建立新頁面
        </NuxtLink>
      </div>
    </div>

    <div v-else class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
      <div
        v-for="page in portalPages"
        :key="page.id"
        class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6 hover:shadow-lg transition"
      >
        <h3 class="text-xl font-semibold mb-2 text-gray-900 dark:text-white">{{ page.title }}</h3>
        <p class="text-sm text-gray-500 dark:text-gray-400 mb-4">{{ page.slug }}</p>
        <div class="flex gap-2">
          <NuxtLink
            :to="`/portal-pages/${page.id}/edit`"
            class="flex-1 px-3 py-2 text-sm bg-primary-600 text-white rounded-md hover:bg-primary-700 text-center"
          >
            編輯
          </NuxtLink>
          <NuxtLink
            :to="`/${page.slug}`"
            target="_blank"
            class="flex-1 px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-md hover:bg-gray-50 dark:hover:bg-gray-700 text-center text-gray-700 dark:text-gray-200"
          >
            檢視
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { PortalPageSummary } from '~/types/api'

definePageMeta({
  layout: 'default',
  middleware: 'auth'
})

const api = useApi()

const loading = ref(true)
const error = ref<string | null>(null)
const portalPages = ref<PortalPageSummary[]>([])

onMounted(async () => {
  try {
    const response = await api.listPortalPages()
    portalPages.value = response.portal_pages
  } catch (err: any) {
    error.value = err.message || '載入失敗'
  } finally {
    loading.value = false
  }
})
</script>
