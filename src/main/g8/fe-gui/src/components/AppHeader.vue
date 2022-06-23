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
      <a href="/" class="c-header-brand mx-auto d-lg-none" style="font-weight: bolder; font-size: x-large">{{
        appName
      }}</a>
      <CHeaderNav class="d-none d-md-flex me-auto">
        <CNavItem class="px-2">
          <CNavLink :href="this.$router.resolve({ name: 'Dashboard' }).href">{{ $t('message.dashboard') }}</CNavLink>
        </CNavItem>
        <CNavItem class="px-2">
          <CNavLink :href="this.$router.resolve({ name: 'MyBlog' }).href">{{ $t('message.my_blog') }}</CNavLink>
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
        <CDropdown variant="nav-item">
          <CDropdownToggle placement="bottom-end" class="py-0" :caret="false">
            <CNavLink href="#">
              <CIcon name="cil-flag-alt" />
            </CNavLink>
            <CDropdownMenu class="pt-0">
              <CDropdownItem v-for="locale in $i18n.availableLocales" :key="locale" @click="doSwitchLanguage(locale)">
                <CIcon :name="$i18n.messages[locale]._flag" />
                <span class="px-2">{{ $i18n.messages[locale]._name }}</span>
              </CDropdownItem>
            </CDropdownMenu>
          </CDropdownToggle>
        </CDropdown>
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
import { swichLanguage } from '@/i18n'

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
  methods: {
    doSwitchLanguage(locale) {
      const answer = confirm(this.$i18n.t('message.switch_language_msg'))
      if (answer) {
        swichLanguage(locale, true)
      }
    },
  },
}
</script>
