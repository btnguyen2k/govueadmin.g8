<!-- #GovueAdmin-Customized -->
<template>
  <CFooter :fixed="false">
    <div>© {{ appName }} version {{ appVersion }}.</div>
    <div class="mfs-auto">
      <span class="mr-1">Powered by</span
      ><a href="https://github.com/btnguyen2k/govueadmin.g8" target="_blank"
        >GoVueAdmin</a
      >
      |
      <span class="mr-1">Template by</span
      ><a href="https://coreui.io/vue" target="_blank">CoreUI for Vue</a>
    </div>
  </CFooter>
</template>

<script>
import cfg from '@/utils/app_config'

export default {
  name: 'TheFooter',
  data() {
    return {
      appName: cfg.APP_CONFIG.app.name,
      appVersion: cfg.APP_CONFIG.app.version,
    }
  },
}
</script>
