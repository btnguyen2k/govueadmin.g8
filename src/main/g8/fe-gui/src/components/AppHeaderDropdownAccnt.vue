<!-- #GoVueAdmin-Customized -->
<template>
  <CDropdown variant="nav-item">
    <CDropdownToggle placement="bottom-end" class="py-0" :caret="false">
      <CAvatar :src="avatarUrl" size="md" />
    </CDropdownToggle>
    <CDropdownMenu class="pt-0">
      <CDropdownHeader component="h6" class="bg-light fw-semibold py-2"> Account </CDropdownHeader>
      <CDropdownItem @click="funcNotImplemented">
        <CIcon icon="cil-bell" /> Updates
        <CBadge color="info" class="ms-auto">{{ itemsCount }}</CBadge>
      </CDropdownItem>
      <CDropdownItem @click="funcNotImplemented">
        <CIcon icon="cil-envelope-open" /> Messages
        <CBadge color="success" class="ms-auto">{{ itemsCount }}</CBadge>
      </CDropdownItem>
      <CDropdownItem @click="funcNotImplemented">
        <CIcon icon="cil-task" /> Tasks
        <CBadge color="danger" class="ms-auto">{{ itemsCount }}</CBadge>
      </CDropdownItem>
      <CDropdownItem @click="funcNotImplemented">
        <CIcon icon="cil-comment-square" /> Comments
        <CBadge color="warning" class="ms-auto">{{ itemsCount }}</CBadge>
      </CDropdownItem>
      <CDropdownHeader component="h6" class="bg-light fw-semibold py-2"> Settings </CDropdownHeader>
      <CDropdownItem @click="funcNotImplemented"> <CIcon icon="cil-user" /> Profile </CDropdownItem>
      <CDropdownItem @click="funcNotImplemented"> <CIcon icon="cil-settings" /> Settings </CDropdownItem>
      <CDropdownItem @click="funcNotImplemented">
        <CIcon icon="cil-dollar" /> Payments
        <CBadge color="secondary" class="ms-auto">{{ itemsCount }}</CBadge>
      </CDropdownItem>
      <CDropdownItem @click="funcNotImplemented">
        <CIcon icon="cil-file" /> Projects
        <CBadge color="primary" class="ms-auto">{{ itemsCount }}</CBadge>
      </CDropdownItem>
      <CDropdownDivider />
      <CDropdownItem @click="funcNotImplemented"> <CIcon icon="cil-shield-alt" /> Lock Account </CDropdownItem>
      <CDropdownItem @click="doLogout"> <CIcon icon="cil-lock-locked" /> Logout </CDropdownItem>
    </CDropdownMenu>
  </CDropdown>
</template>

<script>
import MD5 from 'crypto-js/md5'
import utils from '@/utils/app_utils'

export default {
  name: 'AppHeaderDropdownAccnt',
  setup() {
    String.prototype.md5 = function () {
      return MD5(this)
    }
    let session = utils.loadLoginSession()
    let uid = session != null ? session.uid : ''
    return {
      itemsCount: 42,
      displayName: session != null ? session.name : uid,
      avatarUrl: 'https://www.gravatar.com/avatar/{aid}?s=40'.replace('{aid}', uid.trim().toLowerCase().md5()),
    }
  },
  methods: {
    funcNotImplemented() {
      alert('Not implemented')
    },
    doLogout() {
      utils.localStorageSet(utils.lskeyLoginSession, null)
      utils.localStorageSet(utils.lskeyLoginSessionLastCheck, null)
      this.$router.push({ name: 'Login' })
    },
  },
}
</script>

<style scoped>
a.dropdown-item {
  cursor: pointer;
}
</style>
