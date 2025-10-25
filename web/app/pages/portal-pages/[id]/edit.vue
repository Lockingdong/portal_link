<template>
  <div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <div class="mb-8">
      <NuxtLink to="/dashboard" class="text-primary-600 hover:text-primary-700 flex items-center gap-2 mb-4">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
        返回儀表板
      </NuxtLink>
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">編輯 Portal Page</h1>
    </div>

    <div v-if="pageLoading" class="text-center py-12">
      <p class="text-gray-500 dark:text-gray-400">載入中...</p>
    </div>

    <form v-else @submit.prevent="handleSubmit" class="space-y-6">
      <div v-if="error" class="rounded-md bg-red-50 dark:bg-red-900/20 p-4">
        <p class="text-sm text-red-800 dark:text-red-200">{{ error }}</p>
      </div>

      <!-- 基本資訊 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6 space-y-6">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white">基本資訊</h2>

        <div>
          <label for="slug" class="block text-sm font-medium mb-1 text-gray-700 dark:text-gray-300">
            自訂網址 (Slug)
          </label>
          <div class="mt-1 flex rounded-md shadow-sm">
            <span class="inline-flex items-center px-3 rounded-l-md border border-r-0 border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-700 text-gray-500 dark:text-gray-400 text-sm">
              portallink.com/
            </span>
            <input
              id="slug"
              v-model="form.slug"
              type="text"
              pattern="^[a-z0-9]+(-[a-z0-9]+)*$"
              minlength="3"
              maxlength="50"
              class="flex-1 min-w-0 block w-full px-3 py-2 rounded-none rounded-r-md border border-gray-300 dark:border-gray-600 focus:ring-primary-500 focus:border-primary-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
            />
          </div>
        </div>

        <div>
          <label for="title" class="block text-sm font-medium mb-1 text-gray-700 dark:text-gray-300">
            頁面標題
          </label>
          <input
            id="title"
            v-model="form.title"
            type="text"
            maxlength="100"
            class="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:ring-primary-500 focus:border-primary-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
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

      <!-- 連結管理 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6 space-y-6">
        <div class="flex justify-between items-center">
          <h2 class="text-xl font-semibold text-gray-900 dark:text-white">連結管理</h2>
          <button
            type="button"
            @click="addLink"
            class="px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 text-sm flex items-center gap-2"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            新增連結
          </button>
        </div>

        <div v-if="form.links.length === 0" class="text-center py-8 text-gray-500 dark:text-gray-400">
          尚無連結,點擊「新增連結」開始新增
        </div>

        <div v-else class="space-y-4">
          <div
            v-for="(link, index) in form.links"
            :key="index"
            class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 space-y-4"
          >
            <div class="flex justify-between items-center">
              <span class="text-sm font-medium text-gray-500 dark:text-gray-400">連結 #{{ index + 1 }}</span>
              <button
                type="button"
                @click="removeLink(index)"
                class="text-red-600 hover:text-red-700 text-sm"
              >
                刪除
              </button>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <label :for="`link-title-${index}`" class="block text-sm font-medium mb-1 text-gray-700 dark:text-gray-300">
                  標題 <span class="text-red-500">*</span>
                </label>
                <input
                  :id="`link-title-${index}`"
                  v-model="link.title"
                  type="text"
                  required
                  maxlength="100"
                  class="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:ring-primary-500 focus:border-primary-500 bg-white dark:bg-gray-800 text-sm text-gray-900 dark:text-white"
                />
              </div>

              <div>
                <label :for="`link-url-${index}`" class="block text-sm font-medium mb-1 text-gray-700 dark:text-gray-300">
                  URL <span class="text-red-500">*</span>
                </label>
                <input
                  :id="`link-url-${index}`"
                  v-model="link.url"
                  type="url"
                  required
                  class="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:ring-primary-500 focus:border-primary-500 bg-white dark:bg-gray-800 text-sm text-gray-900 dark:text-white"
                />
              </div>

              <div class="col-span-2">
                <label :for="`link-description-${index}`" class="block text-sm font-medium mb-1 text-gray-700 dark:text-gray-300">
                  描述
                </label>
                <input
                  :id="`link-description-${index}`"
                  v-model="link.description"
                  type="text"
                  maxlength="500"
                  class="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:ring-primary-500 focus:border-primary-500 bg-white dark:bg-gray-800 text-sm text-gray-900 dark:text-white"
                />
              </div>

              <div>
                <label :for="`link-icon-${index}`" class="block text-sm font-medium mb-1 text-gray-700 dark:text-gray-300">
                  圖示 URL
                </label>
                <input
                  :id="`link-icon-${index}`"
                  v-model="link.icon_url"
                  type="url"
                  class="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:ring-primary-500 focus:border-primary-500 bg-white dark:bg-gray-800 text-sm text-gray-900 dark:text-white"
                />
              </div>

              <div>
                <label :for="`link-order-${index}`" class="block text-sm font-medium mb-1 text-gray-700 dark:text-gray-300">
                  顯示順序 <span class="text-red-500">*</span>
                </label>
                <input
                  :id="`link-order-${index}`"
                  v-model.number="link.display_order"
                  type="number"
                  min="1"
                  required
                  class="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:ring-primary-500 focus:border-primary-500 bg-white dark:bg-gray-800 text-sm text-gray-900 dark:text-white"
                />
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 操作按鈕 -->
      <div class="flex gap-4">
        <button
          type="submit"
          :disabled="loading"
          class="flex-1 px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <span v-if="loading">儲存中...</span>
          <span v-else>儲存變更</span>
        </button>
        <NuxtLink
          :to="`/${currentPage?.slug}`"
          target="_blank"
          class="px-4 py-2 text-white border border-gray-300 dark:border-gray-600 rounded-md hover:bg-gray-50 dark:hover:bg-gray-800"
        >
          預覽
        </NuxtLink>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import type { UpdatePortalPageRequest, LinkRequest, FindPortalPageByIDResponse } from '~/types/api'

definePageMeta({
  layout: 'default',
  middleware: 'auth'
})

const route = useRoute()
const router = useRouter()
const api = useApi()

const pageId = computed(() => Number(route.params.id))

const pageLoading = ref(true)
const loading = ref(false)
const error = ref<string | null>(null)
const currentPage = ref<FindPortalPageByIDResponse | null>(null)

const form = reactive<UpdatePortalPageRequest>({
  slug: '',
  title: '',
  bio: '',
  profile_image_url: '',
  theme: 'light',
  links: []
})

onMounted(async () => {
  try {
    const response = await api.getPortalPageById(pageId.value)
    currentPage.value = response

    form.slug = response.slug
    form.title = response.title
    form.bio = response.bio || ''
    form.profile_image_url = response.profile_image_url || ''
    form.theme = response.theme
    form.links = response.links.map(link => ({
      id: link.id,
      title: link.title,
      url: link.url,
      description: link.description || '',
      icon_url: link.icon_url || '',
      display_order: link.display_order
    }))
  } catch (err: any) {
    error.value = err.message || '載入失敗'
  } finally {
    pageLoading.value = false
  }
})

function addLink() {
  const maxOrder = form.links.length > 0
    ? Math.max(...form.links.map(l => l.display_order))
    : 0

  form.links.push({
    title: '',
    url: '',
    description: '',
    icon_url: '',
    display_order: maxOrder + 1
  })
}

function removeLink(index: number) {
  form.links.splice(index, 1)
}

async function handleSubmit() {
  loading.value = true
  error.value = null

  try {
    await api.updatePortalPage(pageId.value, form)
    await router.push('/dashboard')
  } catch (err: any) {
    error.value = err.message || '儲存失敗'
  } finally {
    loading.value = false
  }
}
</script>
