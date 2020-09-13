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
                  <p class="text-muted">Please sign in to continue</p>
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
                    <CCol col="6" class="text-left">
                      <CButton color="primary" class="px-4" type="submit">Login</CButton>
                    </CCol>
                    <CCol col="6" class="text-right">
                      <CButton color="link" class="px-0" @click="funcNotImplemented">Forgot password?</CButton>
                      <!--                      <CButton color="link" class="d-lg-none" @click="funcNotImplemented">Register now!</CButton>-->
                    </CCol>
                  </CRow>
                </CForm>
              </CCardBody>
            </CCard>
            <CCard color="primary" text-color="white" class="text-center py-5 d-md-down-none" body-wrapper>
              <CCardBody>
                <h2>Demo</h2>
                <p>This is instance is for demo purpose only. Login with administrator account
                  <strong>admin/s3cr3t</strong>.
                  You can create/edit/delete other user group or user account. This special admin account,
                  however, can not be modified or deleted.</p>
                <!--                <CButton-->
                <!--                    color="light"-->
                <!--                    variant="outline"-->
                <!--                    size="lg"-->
                <!--                >-->
                <!--                  Register Now!-->
                <!--                </CButton>-->
              </CCardBody>
            </CCard>
          </CCardGroup>
        </CCol>
      </CRow>
    </CContainer>
  </div>
</template>

<script>
import clientUtils from "@/utils/api_client"
import utils from "@/utils/app_utils"
// import appConfig from "@/utils/app_config"
// import router from "@/router"

// const defaultInfoMsg = "Please sign in to continue"
// const waitInfoMsg = "Please wait..."
// const waitLoginInfoMsg = "Logging in, please wait..."
// const invalidReturnUrlErrMsg = "Error: invalid return url"

export default {
  name: 'Login',
  computed: {
    returnUrl() {
      return this.$route.query.returnUrl ? this.$route.query.returnUrl : "#"
    },
  },
  data() {
    return {
      errorMsg: "",
      form: {username: "", password: ""},
    }
  },
  methods: {
    funcNotImplemented() {
      alert("Not implemented")
    },
    doSubmit(e) {
      e.preventDefault()
      let data = {username: this.form.username, password: this.form.password}
      clientUtils.apiDoPost(
          clientUtils.apiLogin, data,
          (apiResp) => {
            if (apiResp.status != 200) {
              this.errorMsg = apiResp.status + ": " + apiResp.message
            } else {
              console.log(apiResp.data)
              // const jwt = utils.parseJwt(apiResp.data)
              // if (!jwt) {
              //   this.errorMsg = 'Error parsing login-token.'
              // } else {
              //   let rUrl = this.returnUrl
              //   if (rUrl == "" || rUrl == null || rUrl == '#') {
              //     rUrl = this.$router.resolve({name: 'Dashboard'}).href
              //   }
              //   utils.saveLoginSession({uid: jwt.payloadObj.uid, name: jwt.payloadObj.name, token: apiResp.data})
              //   // window.location.href = rUrl
              //   this.$router.push(rUrl)
              // }
            }
          },
          (err) => {
            console.error('Error: ', err)
            this.errorMsg = err
          }
      )
    },
  }
}
</script>
