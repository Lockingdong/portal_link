<template>
  <div>
    <div class="relative overflow-hidden bg-gradient-to-br from-primary-500 to-primary-700 dark:from-primary-800 dark:to-primary-900">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24 text-center">
        <h1 class="text-5xl font-bold text-white mb-6">
          歡迎使用 Portal Link
        </h1>
        <p class="text-xl text-primary-100 mb-8 max-w-2xl mx-auto">
          個人化的連結整合平台，把所有重要的社群、網站、商店連結集中到一個專屬頁面
        </p>
        <div class="flex gap-4 justify-center">
          <NuxtLink
            to="/signup"
            class="px-8 py-3 bg-white text-primary-600 rounded-lg font-semibold hover:bg-gray-100 transition"
          >
            開始使用
          </NuxtLink>
          <NuxtLink
            to="/signin"
            class="px-8 py-3 bg-primary-600 text-white rounded-lg font-semibold hover:bg-primary-500 transition border-2 border-white"
          >
            登入帳號
          </NuxtLink>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
      <div class="grid md:grid-cols-3 gap-8">
        <div class="text-center p-6">
          <div class="w-16 h-16 bg-primary-100 dark:bg-primary-900 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg class="w-8 h-8 text-primary-600 dark:text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <h3 class="text-xl font-semibold mb-2 text-gray-900 dark:text-white">快速建立</h3>
          <p class="text-gray-600 dark:text-gray-400">幾分鐘內就能建立專屬的連結頁面</p>
        </div>

        <div class="text-center p-6">
          <div class="w-16 h-16 bg-primary-100 dark:bg-primary-900 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg class="w-8 h-8 text-primary-600 dark:text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4" />
            </svg>
          </div>
          <h3 class="text-xl font-semibold mb-2 text-gray-900 dark:text-white">自訂設計</h3>
          <p class="text-gray-600 dark:text-gray-400">自由設定頁面主題和排版</p>
        </div>

        <div class="text-center p-6">
          <div class="w-16 h-16 bg-primary-100 dark:bg-primary-900 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg class="w-8 h-8 text-primary-600 dark:text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <h3 class="text-xl font-semibold mb-2 text-gray-900 dark:text-white">簡單分享</h3>
          <p class="text-gray-600 dark:text-gray-400">只需分享一個連結，展示所有內容</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'default'
})
</script>
