<!-- #GovueAdmin-Customized -->
<template>
  <CHeader fixed with-subheader light>
    <CToggler
        in-header
        class="ml-3 d-lg-none"
        @click="$store.commit('toggleSidebarMobile')"
    />
    <CToggler
        in-header
        class="ml-3 d-md-down-none"
        @click="$store.commit('toggleSidebarDesktop')"
    />
    <a href="/" class="c-header-brand mx-auto d-lg-none"
       style="font-weight: bolder; font-size: x-large">{{ appName }}</a>
    <!--    <CHeaderBrand class="mx-auto d-lg-none" to="/">-->
    <!--      <CIcon name="logo" height="48" alt="Logo"/>-->
    <!--    </CHeaderBrand>-->
    <CHeaderNav class="d-md-down-none mr-auto">
      <CHeaderNavItem class="px-3">
        <CHeaderNavLink :to="{name:'Dashboard'}">
          {{ $t('message.dashboard') }}
        </CHeaderNavLink>
      </CHeaderNavItem>
      <CHeaderNavItem class="px-3">
        <CHeaderNavLink :to="{name:'MyBlog'}">
          {{ $t('message.my_blog') }}
        </CHeaderNavLink>
      </CHeaderNavItem>
      <CHeaderNavItem class="px-3">
        <CHeaderNavLink :to="{name:'CreatePost'}">
          {{ $t('message.create_blog_post') }}
        </CHeaderNavLink>
      </CHeaderNavItem>
    </CHeaderNav>
    <CHeaderNav class="mr-4">
      <CHeaderNavItem class="d-md-down-none mx-2">
        <CHeaderNavLink>
          <CIcon name="cil-bell"/>
        </CHeaderNavLink>
      </CHeaderNavItem>
      <CHeaderNavItem class="d-md-down-none mx-2">
        <CHeaderNavLink>
          <CIcon name="cil-list"/>
        </CHeaderNavLink>
      </CHeaderNavItem>
      <CHeaderNavItem class="d-md-down-none mx-2">
        <CHeaderNavLink>
          <CIcon name="cil-envelope-open"/>
        </CHeaderNavLink>
      </CHeaderNavItem>
      <TheHeaderDropdownAccnt/>
    </CHeaderNav>
    <CSubheader class="px-3">
      <CBreadcrumbRouter class="border-0 mb-0"/>
    </CSubheader>
  </CHeader>
</template>

<script>
import TheHeaderDropdownAccnt from './TheHeaderDropdownAccnt'
import cfg from '@/utils/app_config'

export default {
  name: 'TheHeader',
  computed: {
    appName() {
      return cfg.APP_CONFIG.app.name
    },
  },
  components: {
    TheHeaderDropdownAccnt
  }
}
</script>
