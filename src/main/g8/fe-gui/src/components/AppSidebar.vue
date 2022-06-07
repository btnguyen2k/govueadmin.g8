<!-- #GovueAdmin-Customized -->
<template>
  <CSidebar
    position="fixed"
    :unfoldable="sidebarUnfoldable"
    :visible="sidebarVisible"
    @visible-change="
      (value) =>
        this.$store.commit({
          type: 'updateSidebarVisible',
          value: value,
        })
    "
  >
    <CSidebarBrand>
      <span
        class="sidebar-brand-full"
        style="color: #fff; font-weight: bolder; font-size: x-large"
        >{{ appName }}</span
      >
      <span
        class="sidebar-brand-narrow"
        style="color: #fff; font-weight: bolder; font-size: large"
        >{{ appInitial }}</span
      >
    </CSidebarBrand>
    <AppSidebarNav />
    <CSidebarToggler
      class="d-none d-lg-flex"
      @click="this.$store.commit('toggleUnfoldable')"
    />
  </CSidebar>
</template>

<script>
import { computed } from 'vue'
import { useStore } from 'vuex'
import { AppSidebarNav } from './AppSidebarNav'
import cfg from '@/utils/app_config'

export default {
  name: 'AppSidebar',
  components: {
    AppSidebarNav,
  },
  setup() {
    const store = useStore()
    return {
      sidebarUnfoldable: computed(() => store.state.sidebarUnfoldable),
      sidebarVisible: computed(() => store.state.sidebarVisible),
      appName: cfg.APP_CONFIG.app.name,
      appInitial: cfg.APP_CONFIG.app.initial,
    }
  },
}
</script>
