<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <nav class="bg-white dark:bg-gray-800 shadow">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <NuxtLink to="/" class="text-2xl font-bold text-primary-600 dark:text-primary-400">
              Portal Link
            </NuxtLink>
          </div>
          <div class="flex items-center gap-4">
            <template v-if="isAuthenticated">
              <NuxtLink
                to="/dashboard"
                class="text-gray-700 dark:text-gray-300 hover:text-primary-600 dark:hover:text-primary-400"
              >
                儀表板
              </NuxtLink>
              <button
                @click="handleLogout"
                class="px-4 py-2 text-sm font-medium text-white bg-primary-600 rounded-md hover:bg-primary-700"
              >
                登出
              </button>
            </template>
            <template v-else>
              <NuxtLink
                to="/signin"
                class="text-gray-700 dark:text-gray-300 hover:text-primary-600 dark:hover:text-primary-400"
              >
                登入
              </NuxtLink>
              <NuxtLink
                to="/signup"
                class="px-4 py-2 text-sm font-medium text-white bg-primary-600 rounded-md hover:bg-primary-700"
              >
                註冊
              </NuxtLink>
            </template>
          </div>
        </div>
      </div>
    </nav>

    <main>
      <slot />
    </main>
  </div>
</template>

<script setup lang="ts">
const { isAuthenticated, logout } = useAuth()

function handleLogout() {
  logout()
}
</script>
