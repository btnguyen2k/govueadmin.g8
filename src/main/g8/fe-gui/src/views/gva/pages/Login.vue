<!-- #GoVueAdmin-Customized -->
<template>
  <div class="bg-light min-vh-100 d-flex flex-row align-items-center">
    <CContainer>
      <CRow class="justify-content-center">
        <CCol :md="8">
          <CCardGroup>
            <CCard class="p-4">
              <CCardBody>
                <CForm @submit.prevent="doSubmit" method="post">
                  <h1>{{ $t('message.login') }}</h1>
                  <p v-if="errorMsg != ''" class="alert alert-danger">
                    {{ errorMsg }}
                  </p>
                  <p v-if="infoMsg != ''" class="text-muted">{{ infoMsg }}</p>
                  <CInputGroup class="mb-3">
                    <CInputGroupText>
                      <CIcon icon="cil-user" />
                    </CInputGroupText>
                    <CFormInput
                      :placeholder="$t('message.username')"
                      autocomplete="username email"
                      name="username"
                      id="username"
                      v-model="form.username"
                    />
                  </CInputGroup>
                  <CInputGroup class="mb-4">
                    <CInputGroupText>
                      <CIcon icon="cil-lock-locked" />
                    </CInputGroupText>
                    <CFormInput
                      :placeholder="$t('message.password')"
                      type="password"
                      autocomplete="current-password"
                      name="password"
                      id="password"
                      v-model="form.password"
                    />
                  </CInputGroup>
                  <CRow>
                    <CCol :xs="5">
                      <CButton size="sm" color="primary" class="px-4" type="submit">{{ $t('message.login') }}</CButton>
                    </CCol>
                    <CCol :xs="7" class="text-right">
                      <CButton size="sm" color="link" class="px-2" @click="doClickLoginSocial">{{
                        $t('message.login_social')
                      }}</CButton>
                    </CCol>
                  </CRow>
                  <CRow class="py-2">
                    <CCol sm="auto">
                      <CFormLabel class="col-form-label col-form-label-sm">{{ $t('message.language') }}</CFormLabel>
                    </CCol>
                    <CCol sm="auto">
                      <CFormSelect size="sm" v-model="$root.$i18n.locale" :options="languageOptions" />
                    </CCol>
                  </CRow>
                </CForm>
              </CCardBody>
            </CCard>
            <CCard v-if="demoMode" color="primary" text-color="white" class="py-5 d-md-down-none" body-wrapper>
              <CCardBody>
                <h2>{{ $t('_demo') }}</h2>
                <p v-html="$t('_demo_msg')"></p>
              </CCardBody>
            </CCard>
          </CCardGroup>
        </CCol>
      </CRow>
    </CContainer>
  </div>
</template>

<script>
import apiClient from '@/utils/api_client'
import utils from '@/utils/app_utils'

export default {
  name: 'Login',
  mounted() {
    if (this.$route.query.exterToken != undefined && this.$route.query.exterToken != '') {
      let data = { token: this.$route.query.exterToken, mode: 'exter' }
      this._doLogin(data)
    }
    this.infoMsgSwitch = 1
    apiClient.apiDoGet(
      apiClient.apiInfo,
      (apiRes) => {
        if (apiRes.status != 200) {
          this.errorMsg = apiRes.message
        } else {
          this.exterAppId = apiRes.data.exter.app_id
          this.exterBaseUrl = apiRes.data.exter.base_url
          this.infoMsgSwitch = 2
          this.demoMode = apiRes.data.demo_mode
          if (apiRes.data.demo_mode) {
            this.form.username = apiRes.data.demo.user_id
            this.form.password = apiRes.data.demo.user_pwd
          }
        }
      },
      (err) => {
        this.errorMsg = err
      },
    )
  },
  computed: {
    infoMsg() {
      if (this.infoMsgSwitch == 0) {
        return ''
      }
      if (this.infoMsgSwitch == 1) {
        return this.$i18n.t('message.wait')
      }
      return this.$i18n.t('message.login_info')
    },
    parseLoginTokenErrMsg() {
      return this.$i18n.t('message.error_parse_login_token')
    },
    returnUrl() {
      return this.$route.query.returnUrl ? this.$route.query.returnUrl : ''
    },
    languageOptions() {
      let result = []
      this.$i18n.availableLocales.forEach((locale) => {
        result.push({ value: locale, label: this.$i18n.messages[locale]._name })
      })
      return result
    },
  },
  data() {
    return {
      exterAppId: String,
      exterBaseUrl: String,
      errorMsg: '',
      infoMsgSwitch: 0,
      form: { username: '', password: '' },
      demoMode: false,
    }
  },
  methods: {
    funcNotImplemented() {
      alert('Not implemented')
    },
    doClickLoginSocial() {
      let prUrl = this.$route.query.returnUrl ? this.$route.query.returnUrl : ''
      let rUrl =
        window.location.origin +
        this.$router.resolve({ name: 'Login' }).href +
        '?returnUrl=' +
        prUrl.replaceAll('#', encodeURIComponent('#')).replaceAll('=', encodeURIComponent('#')) +
        '&exterToken=${token}'
      let cUrl = window.location.origin + this.$router.resolve({ name: 'Login' }).href
      let url =
        this.exterBaseUrl +
        '/app/xlogin?app=' +
        this.exterAppId +
        '&returnUrl=' +
        encodeURIComponent(rUrl) +
        '&cancelUrl=' +
        encodeURIComponent(cUrl)
      window.location.href = url
    },
    _doLogin(data) {
      apiClient.apiDoPost(
        apiClient.apiLogin,
        data,
        (apiResp) => {
          if (apiResp.status != 200) {
            this.errorMsg = apiResp.status + ': ' + apiResp.message
          } else {
            const jwt = utils.parseJwt(apiResp.data)
            if (!jwt) {
              this.errorMsg = this.parseLoginTokenErrMsg
            } else {
              utils.saveLoginSession({
                uid: jwt.payloadObj.uid,
                name: jwt.payloadObj.name,
                token: apiResp.data,
              })
              let rUrl = this.returnUrl
              if (rUrl == null || rUrl == '') {
                this.$router.push({ name: 'Dashboard' })
                // this.$router.push(
                //   this.$router.resolve({ name: 'Dashboard' }).location,
                // )
              } else {
                window.location.href = rUrl
              }
            }
          }
        },
        (err) => {
          console.error('Error: ', err)
          this.errorMsg = err
        },
      )
    },
    doSubmit(e) {
      e.preventDefault()
      let data = {
        username: this.form.username,
        password: this.form.password,
        mode: 'form',
      }
      this._doLogin(data)
    },
  },
}
</script>
