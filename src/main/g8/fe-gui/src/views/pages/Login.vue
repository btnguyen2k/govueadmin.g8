<template>
    <CContainer class="d-flex align-items-center min-vh-100">
        <CRow class="justify-content-center">
            <CCol md="8">
                <CCardGroup>
                    <CCard class="p-4">
                        <CCardBody>
                            <CForm @submit.prevent="doSubmit" method="post">
                                <h1>Login</h1>
                                <p v-if="erroMsg!=''" class="alert alert-danger">{{erroMsg}}</p>
                                <p class="text-muted">Please sign in to continue</p>
                                <CInput placeholder="Username" autocomplete="username email" name="username"
                                        id="username" v-model="form.username">
                                    <template #prepend-content>
                                        <CIcon name="cil-user"/>
                                    </template>
                                </CInput>
                                <CInput placeholder="Password" type="password" autocomplete="current-password"
                                        name="password" id="password" v-model="form.password">
                                    <template #prepend-content>
                                        <CIcon name="cil-lock-locked"/>
                                    </template>
                                </CInput>
                                <CRow>
                                    <CCol col="6">
                                        <CButton color="primary" class="px-4" type="submit">
                                            Login
                                        </CButton>
                                    </CCol>
                                    <CCol col="6" class="text-right">
                                        <CButton color="link" class="px-0" @click="jsNotImplemented">Forgot password?
                                        </CButton>
                                    </CCol>
                                </CRow>
                            </CForm>
                        </CCardBody>
                    </CCard>
                    <CCard
                            color="primary"
                            text-color="white"
                            class="text-center py-5 d-md-down-none"
                            style="width:44%"
                            body-wrapper
                    >
                        <h2>Sign up</h2>
                        <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut
                            labore et dolore magna aliqua.</p>
                        <CButton color="primary" class="active mt-3" :to="pageRegister">
                            Register Now!
                        </CButton>
                    </CCard>
                </CCardGroup>
            </CCol>
        </CRow>
    </CContainer>
</template>

<script>
    import Register from "@/views/pages/Register"
    import clientUtils from "@/utils/api_client"
    import utils from "@/utils/app_utils"

    export default {
        name: 'Login',
        data() {
            return {
                returnUrl: "/",
                pageRegister: Register,
                form: {username: "", password: ""},
                erroMsg: "",
            }
        },
        created() {
            this.returnUrl = this.$route.query.returnUrl
        },
        methods: {
            jsNotImplemented(e) {
                e.preventDefault()
                alert("Not implemented")
            },
            doSubmit(e) {
                e.preventDefault()
                let data = {username: this.form.username, password: this.form.password}
                clientUtils.apiDoPost(
                    clientUtils.apiLogin, data,
                    (apiRes) => {
                        if (apiRes.status != 200) {
                            this.erroMsg = apiRes.status + ": " + apiRes.message
                        } else {
                            utils.saveLoginSession({
                                uid: apiRes.data.uid,
                                token: apiRes.data.token,
                                expiry: apiRes.data.expiry,
                            })
                            this.$router.push(this.returnUrl)
                        }
                    },
                    (err) => {
                        this.erroMsg = err
                    }
                )
            },
        }
    }
</script>
