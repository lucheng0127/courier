<script setup lang="ts">
import { ConfigProvider } from 'ant-design-vue'

const customTheme = {
  token: {
    colorPrimary: '#10A37F',
    borderRadius: 6,
    fontSize: 14
  }
}
</script>

<template>
  <ConfigProvider :theme="customTheme">
    <router-view />
  </ConfigProvider>
</template>

<style>
/* 全局样式 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  font-size: 14px;
  line-height: 1.5;
  color: #111827;
  background-color: #F9FAFB;
}

#app {
  min-height: 100vh;
}
</style>
