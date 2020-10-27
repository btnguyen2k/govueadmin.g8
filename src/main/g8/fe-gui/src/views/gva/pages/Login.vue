<!-- #GovueAdmin-Customized -->
<template>
  <div class="c-app flex-row align-items-center">
    <CContainer>
      <CRow class="justify-content-center">
        <CCol md="8">
          <CCardGroup>
            <CCard class="p-4">
              <CCardBody>
                <CForm @submit.prevent="doSubmit" method="post">
                  <h1>Login</h1>
                  <p v-if="errorMsg!=''" class="alert alert-danger">{{ errorMsg }}</p>
                  <p v-if="infoMsg!=''" class="text-muted">{{ infoMsg }}</p>
                  <CInput placeholder="Username" autocomplete="username email" name="username" id="username"
                          v-model="form.username">
                    <template #prepend-content>
                      <CIcon name="cil-user"/>
                    </template>
                  </CInput>
                  <CInput placeholder="Password" type="password" autocomplete="current-password" name="password"
                          id="password" v-model="form.password">
                    <template #prepend-content>
                      <CIcon name="cil-lock-locked"/>
                    </template>
                  </CInput>
                  <CRow>
                    <CCol col="4" class="text-left">
                      <CButton color="primary" class="px-4" type="submit">Login</CButton>
                    </CCol>
                    <CCol col="8" class="text-right">
                      <CButton color="link" class="px-0" @click="doClickLoginSocial">Login with social account</CButton>
                    </CCol>
                  </CRow>
                </CForm>
              </CCardBody>
            </CCard>
            <CCard color="primary" text-color="white" class="py-5 d-md-down-none" body-wrapper>
              <CCardBody>
                <h2>Demo</h2>
                <p>
                  This is instance is for demo purpose only. Login with default account <strong>admin@local/s3cr3t</strong>.
                  <br/>
                  Or you can login with your <u>social account</u> via "Login with social account" link
                  (your social account credential <u>will not</u> be stored on server).
                </p>
              </CCardBody>
            </CCard>
          </CCardGroup>
        </CCol>
      </CRow>
    </CContainer>
  </div>
</template>

<script>
import apiClient from "@/utils/api_client"
import utils from "@/utils/app_utils"
// import appConfig from "@/utils/app_config"
// import router from "@/router"

const defaultInfoMsg = "Please sign in to continue"
const waitInfoMsg = "Please wait..."
// const waitLoginInfoMsg = "Logging in, please wait..."
// const invalidReturnUrlErrMsg = "Error: invalid return url"

export default {
  name: 'Login',
  mounted() {
    if (this.$route.query.exterToken != undefined && this.$route.query.exterToken != "") {
      let data = {token: this.$route.query.exterToken, mode: "exter"}
      this._doLogin(data)
    }
    this.infoMsg = waitInfoMsg
    apiClient.apiDoGet(apiClient.apiInfo,
        (apiRes) => {
          if (apiRes.status != 200) {
            this.errorMsg = apiRes.message
          } else {
            this.exterAppId = apiRes.data.exter.app_id
            this.exterBaseUrl = apiRes.data.exter.base_url
            this.infoMsg = defaultInfoMsg
          }
        },
        (err) => {
          this.errorMsg = err
        })
  },
  computed: {
    returnUrl() {
      return this.$route.query.returnUrl ? this.$route.query.returnUrl : ''
    },
  },
  data() {
    return {
      exterAppId: String,
      exterBaseUrl: String,
      errorMsg: "",
      infoMsg: "",
      form: {username: "", password: ""},
    }
  },
  methods: {
    funcNotImplemented() {
      alert("Not implemented")
    },
    doClickLoginSocial() {
      let prUrl = this.$route.query.returnUrl ? this.$route.query.returnUrl : ''
      let rUrl = window.location.origin + this.$router.resolve({name: 'Login'}).href
          + '?returnUrl=' + prUrl.replaceAll('#', encodeURIComponent('#')).replaceAll('=', encodeURIComponent('#'))
          + '&exterToken=${token}'
      let cUrl = window.location.origin + this.$router.resolve({name: 'Login'}).href
      let url = this.exterBaseUrl + '/app/xlogin?app=' + this.exterAppId
          + '&returnUrl=' + encodeURIComponent(rUrl)
          + '&cancelUrl=' + encodeURIComponent(cUrl)
      window.location.href = url
    },
    _doLogin(data) {
      apiClient.apiDoPost(
          apiClient.apiLogin, data,
          (apiResp) => {
            if (apiResp.status != 200) {
              this.errorMsg = apiResp.status + ": " + apiResp.message
            } else {
              const jwt = utils.parseJwt(apiResp.data)
              if (!jwt) {
                this.errorMsg = 'Error parsing login-token.'
              } else {
                utils.saveLoginSession({uid: jwt.payloadObj.uid, name: jwt.payloadObj.name, token: apiResp.data})
                let rUrl = this.returnUrl
                if (rUrl == null || rUrl == "") {
                  this.$router.push(this.$router.resolve({name: 'Dashboard'}).location)
                } else {
                  window.location.href = rUrl
                }
              }
            }
          },
          (err) => {
            console.error('Error: ', err)
            this.errorMsg = err
          }
      )
    },
    doSubmit(e) {
      e.preventDefault()
      let data = {username: this.form.username, password: this.form.password, mode: "form"}
      this._doLogin(data)
    },
  }
}
</script>
