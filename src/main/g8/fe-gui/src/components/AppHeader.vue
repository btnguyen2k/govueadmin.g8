<!-- #GoVueAdmin-Customized -->
<template>
  <CHeader position="sticky" class="mb-4">
    <CContainer fluid>
      <CHeaderToggler class="ps-1" @click="this.$store.commit('toggleSidebar')">
        <CIcon icon="cil-menu" size="lg" />
      </CHeaderToggler>
      <!--      <CHeaderBrand class="mx-auto d-lg-none" to="/">-->
      <!--        <CIcon :icon="logo" height="48" alt="Logo" />-->
      <!--      </CHeaderBrand>-->
      <a
        href="/"
        class="c-header-brand mx-auto d-lg-none"
        style="font-weight: bolder; font-size: x-large"
        >{{ appName }}</a
      >
      <CHeaderNav class="d-none d-md-flex me-auto">
        <CNavItem class="px-2">
          <CNavLink :href="this.$router.resolve({ name: 'Dashboard' }).href">{{
            $t('message.dashboard')
          }}</CNavLink>
        </CNavItem>
        <CNavItem class="px-2">
          <CNavLink :href="this.$router.resolve({ name: 'MyBlog' }).href">{{
            $t('message.my_blog')
          }}</CNavLink>
        </CNavItem>
        <CNavItem class="px-2">
          <CNavLink :href="this.$router.resolve({ name: 'CreatePost' }).href">{{
            $t('message.create_blog_post')
          }}</CNavLink>
        </CNavItem>
      </CHeaderNav>
      <CHeaderNav>
        <CNavItem>
          <CNavLink href="#">
            <CIcon class="mx-2" icon="cil-bell" size="lg" />
          </CNavLink>
        </CNavItem>
        <CNavItem>
          <CNavLink href="#">
            <CIcon class="mx-2" icon="cil-list" size="lg" />
          </CNavLink>
        </CNavItem>
        <CNavItem>
          <CNavLink href="#">
            <CIcon class="mx-2" icon="cil-envelope-open" size="lg" />
          </CNavLink>
        </CNavItem>
        <AppHeaderDropdownAccnt />
      </CHeaderNav>
    </CContainer>
    <CHeaderDivider />
    <CContainer fluid>
      <AppBreadcrumb />
    </CContainer>
  </CHeader>
</template>

<script>
import AppBreadcrumb from './AppBreadcrumb'
import AppHeaderDropdownAccnt from './AppHeaderDropdownAccnt'
import { logo } from '@/assets/brand/logo'
import cfg from '@/utils/app_config'

export default {
  name: 'AppHeader',
  components: {
    AppBreadcrumb,
    AppHeaderDropdownAccnt,
  },
  computed: {
    appName() {
      return cfg.APP_CONFIG.app.name
    },
  },
  setup() {
    return {
      logo,
    }
  },
}
</script>
