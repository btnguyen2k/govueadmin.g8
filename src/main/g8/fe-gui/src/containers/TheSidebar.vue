<!-- #GovueAdmin-Customized -->
<template>
  <CSidebar fixed :minimize="minimize" :show="show"
            @update:show="(value) => $store.commit('set', ['sidebarShow', value])">
    <div class="c-sidebar-brand">
      <a href="/">
        <span class="c-sidebar-brand-full"
              style="color: #fff; font-weight: bolder; font-size: x-large">{{ appName }}</span>
        <span class="c-sidebar-brand-minimized"
              style="color: #fff; font-weight: bolder; font-size: large">{{ appInitial }}</span>
      </a>
    </div>
    <CRenderFunction flat :content-to-render="$options.nav"/>
    <CSidebarMinimizer class="d-md-down-none" @click.native="$store.commit('set', ['sidebarMinimize', !minimize])"/>
  </CSidebar>
</template>

<script>
import nav from './_nav.gva'
import cfg from '@/utils/app_config'

export default {
  name: 'TheSidebar',
  nav,
  computed: {
    show() {
      return this.$store.state.sidebarShow
    },
    minimize() {
      return this.$store.state.sidebarMinimize
    },
    appName: cfg.APP_CONFIG.app.name,
    appInitial: cfg.APP_CONFIG.app.initial,
  },
}
</script>
