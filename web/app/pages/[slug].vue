<template>
  <div
    :class="[
      'min-h-screen flex items-center justify-center p-4',
      page?.theme === 'dark'
        ? 'bg-gray-900 text-white'
        : 'bg-gradient-to-br from-primary-50 to-primary-100'
    ]"
  >
    <div v-if="loading" class="text-center">
      <p class="text-gray-500 dark:text-gray-400">載入中...</p>
    </div>

    <div v-else-if="error" class="text-center max-w-md">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <h2 class="mt-4 text-xl font-semibold">找不到此頁面</h2>
      <p class="mt-2 text-gray-600 dark:text-gray-400">{{ error }}</p>
      <NuxtLink
        to="/"
        class="mt-6 inline-block px-6 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700"
      >
        返回首頁
      </NuxtLink>
    </div>

    <div v-else-if="page" class="w-full max-w-2xl">
      <div
        :class="[
          'rounded-2xl shadow-xl p-8',
          page.theme === 'dark'
            ? 'bg-gray-800'
            : 'bg-white'
        ]"
      >
        <!-- 頭像與基本資訊 -->
        <div class="text-center mb-8">
          <div v-if="page.profile_image_url" class="mb-4">
            <img
              :src="page.profile_image_url"
              :alt="page.title"
              class="w-24 h-24 rounded-full mx-auto object-cover border-4 border-primary-500"
            />
          </div>
          <div v-else class="w-24 h-24 rounded-full mx-auto bg-primary-500 flex items-center justify-center text-white text-3xl font-bold mb-4">
            {{ page.title.charAt(0).toUpperCase() }}
          </div>

          <h1 class="text-3xl font-bold mb-2">{{ page.title }}</h1>
          <p v-if="page.bio" class="text-gray-600 dark:text-gray-300">
            {{ page.bio }}
          </p>
        </div>

        <!-- 連結列表 -->
        <div class="space-y-4">
          <a
            v-for="link in page.links"
            :key="link.id"
            :href="link.url"
            target="_blank"
            rel="noopener noreferrer"
            :class="[
              'block p-4 rounded-lg transition-all duration-200 hover:scale-105 hover:shadow-lg',
              page.theme === 'dark'
                ? 'bg-gray-700 hover:bg-gray-600'
                : 'bg-gray-50 hover:bg-gray-100'
            ]"
          >
            <div class="flex items-center gap-4">
              <div v-if="link.icon_url" class="flex-shrink-0">
                <img
                  :src="link.icon_url"
                  :alt="link.title"
                  class="w-12 h-12 rounded-lg object-cover"
                />
              </div>
              <div v-else class="flex-shrink-0 w-12 h-12 rounded-lg bg-primary-500 flex items-center justify-center text-white font-bold">
                {{ link.title.charAt(0).toUpperCase() }}
              </div>

              <div class="flex-1 min-w-0">
                <h3 class="font-semibold text-lg truncate">{{ link.title }}</h3>
                <p v-if="link.description" class="text-sm text-gray-600 dark:text-gray-400 truncate">
                  {{ link.description }}
                </p>
              </div>

              <svg
                class="w-5 h-5 text-gray-400 flex-shrink-0"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
              </svg>
            </div>
          </a>

          <div v-if="page.links.length === 0" class="text-center py-8 text-gray-500 dark:text-gray-400">
            目前尚無連結
          </div>
        </div>

        <!-- 頁尾 -->
        <div class="mt-8 pt-6 border-t border-gray-200 dark:border-gray-700 text-center">
          <p class="text-sm text-gray-500 dark:text-gray-400">
            使用
            <NuxtLink to="/" class="text-primary-600 hover:text-primary-700 font-medium">
              Portal Link
            </NuxtLink>
            建立
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { FindPortalPageBySlugResponse } from '~/types/api'

definePageMeta({
  layout: 'portal'
})

const route = useRoute()
const api = useApi()

const slug = computed(() => route.params.slug as string)

const loading = ref(true)
const error = ref<string | null>(null)
const page = ref<FindPortalPageBySlugResponse | null>(null)

onMounted(async () => {
  try {
    const response = await api.getPortalPageBySlug(slug.value)
    page.value = response

    // 設定頁面標題
    useHead({
      title: `${response.title} - Portal Link`,
      meta: [
        { name: 'description', content: response.bio || `查看 ${response.title} 的所有連結` }
      ]
    })
  } catch (err: any) {
    error.value = err.message || '載入失敗'
  } finally {
    loading.value = false
  }
})
</script>
