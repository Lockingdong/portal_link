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
      <div v-if="successMessage" class="rounded-md bg-green-50 dark:bg-green-900/30 border border-green-200 dark:border-green-700 p-4">
        <div class="flex items-center gap-2">
          <svg class="w-5 h-5 text-green-600 dark:text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <p class="text-sm font-medium text-green-800 dark:text-green-200">{{ successMessage }}</p>
        </div>
      </div>

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

        <div v-else ref="linksContainer" class="space-y-4">
          <div
            v-for="(link, index) in form.links"
            :key="`link-${index}`"
            :data-index="index"
            class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 space-y-4 bg-white dark:bg-gray-800 hover:border-primary-500 dark:hover:border-primary-400 transition-colors"
          >
            <div class="flex justify-between items-center">
              <div class="flex items-center gap-3">
                <div class="drag-handle cursor-grab active:cursor-grabbing p-1">
                  <svg class="w-5 h-5 text-gray-400 dark:text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8h16M4 16h16" />
                  </svg>
                </div>
                <span class="text-sm font-medium text-gray-500 dark:text-gray-400">連結 #{{ index + 1 }}</span>
              </div>
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

              <div class="col-span-2">
                <label :for="`link-icon-${index}`" class="block text-sm font-medium mb-1 text-gray-700 dark:text-gray-300">
                  圖示 URL
                </label>
                <input
                  :id="`link-icon-${index}`"
                  v-model="link.icon_url"
                  type="url"
                  class="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:ring-primary-500 focus:border-primary-500 bg-white dark:bg-gray-800 text-sm text-gray-900 dark:text-white"
                  placeholder="https://example.com/icon.png"
                />
                <div class="mt-2">
                  <p class="text-xs text-gray-500 dark:text-gray-400 mb-2">常用圖示：</p>
                  <div class="flex flex-wrap gap-2">
                    <button
                      v-for="icon in commonIcons"
                      :key="icon.name"
                      type="button"
                      @click="link.icon_url = icon.url"
                      class="flex items-center gap-1.5 px-2.5 py-1.5 text-xs border border-gray-300 dark:border-gray-600 rounded-md hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                      :title="icon.name"
                    >
                      <img :src="icon.url" :alt="icon.name" class="w-4 h-4" />
                      <span class="text-gray-700 dark:text-gray-300">{{ icon.name }}</span>
                    </button>
                  </div>
                </div>
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
          class="px-4 py-2 text-gray-700 dark:text-gray-200 border border-gray-300 dark:border-gray-600 rounded-md hover:bg-gray-50 dark:hover:bg-gray-700"
        >
          預覽
        </NuxtLink>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import type { UpdatePortalPageRequest, FindPortalPageByIDResponse } from '~/types/api'
import Sortable from 'sortablejs'

definePageMeta({
  layout: 'default',
  middleware: 'auth'
})

const route = useRoute()
const api = useApi()

const pageId = computed(() => Number(route.params.id))

const pageLoading = ref(true)
const loading = ref(false)
const error = ref<string | null>(null)
const successMessage = ref<string | null>(null)
const currentPage = ref<FindPortalPageByIDResponse | null>(null)

const form = reactive<UpdatePortalPageRequest>({
  slug: '',
  title: '',
  bio: '',
  profile_image_url: '',
  theme: 'light',
  links: []
})

// 拖動排序功能
const linksContainer = ref<HTMLElement | null>(null)
let sortableInstance: Sortable | null = null

// 常用圖示清單（使用 Simple Icons CDN）
const commonIcons = [
  { name: 'GitHub', url: 'https://cdn.simpleicons.org/github' },
  { name: 'Instagram', url: 'https://cdn.simpleicons.org/instagram' },
  { name: 'Threads', url: 'https://cdn.simpleicons.org/threads' },
  { name: 'Facebook', url: 'https://cdn.simpleicons.org/facebook' },
  { name: 'YouTube', url: 'https://cdn.simpleicons.org/youtube' },
  { name: 'Email', url: 'https://cdn.simpleicons.org/gmail' },
  { name: 'Website', url: 'https://cdn.simpleicons.org/googlechrome' },
]

onMounted(async () => {
  // 先載入資料
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

  // 資料載入後才初始化拖動功能
  await nextTick()
  if (linksContainer.value) {
    sortableInstance = new Sortable(linksContainer.value, {
      animation: 150,
      handle: '.drag-handle',
      ghostClass: 'opacity-50',
      onEnd: (evt: Sortable.SortableEvent) => {
        const oldIndex = evt.oldIndex
        const newIndex = evt.newIndex

        if (oldIndex !== undefined && newIndex !== undefined && oldIndex !== newIndex) {
          // 移動陣列元素
          const movedItem = form.links.splice(oldIndex, 1)[0]
          if (movedItem) {
            form.links.splice(newIndex, 0, movedItem)
            // 更新順序
            updateDisplayOrder()
          }
        }
      }
    })
  }
})

onUnmounted(() => {
  if (sortableInstance) {
    sortableInstance.destroy()
  }
})

function addLink() {
  form.links.push({
    title: '',
    url: '',
    description: '',
    icon_url: '',
    display_order: form.links.length + 1
  })
}

function removeLink(index: number) {
  form.links.splice(index, 1)
  updateDisplayOrder()
}

// 自動更新所有連結的 display_order
function updateDisplayOrder() {
  form.links.forEach((link, index) => {
    link.display_order = index + 1
  })
}

async function handleSubmit() {
  loading.value = true
  error.value = null
  successMessage.value = null

  try {
    // 提交前確保順序正確
    updateDisplayOrder()
    await api.updatePortalPage(pageId.value, form)

    // 重新載入資料以同步最新狀態
    const response = await api.getPortalPageById(pageId.value)
    currentPage.value = response

    // 顯示成功訊息
    successMessage.value = '儲存成功！'

    // 滾動到頂部以確保用戶看到成功訊息
    window.scrollTo({ top: 0, behavior: 'smooth' })

    // 5秒後自動隱藏成功訊息
    setTimeout(() => {
      successMessage.value = null
    }, 5000)
  } catch (err: any) {
    error.value = err.message || '儲存失敗'
    // 發生錯誤時也滾動到頂部
    window.scrollTo({ top: 0, behavior: 'smooth' })
  } finally {
    loading.value = false
  }
}
</script>
