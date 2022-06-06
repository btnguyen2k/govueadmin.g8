<template>
    <div>
        <CRow v-if="!found">
            <CCol sm="12">
                <CCard>
                    <CCardHeader>Edit User</CCardHeader>
                    <CCardBody>
                        <p class="alert alert-danger">User [{{this.$route.params.username}}] not found</p>
                    </CCardBody>
                    <CCardFooter>
                        <CButton type="button" color="info" class="ml-2" style="width: 96px" @click="doCancel">
                            <CIcon name="cil-arrow-circle-left"/>
                            Back
                        </CButton>
                    </CCardFooter>
                </CCard>
            </CCol>
        </CRow>
        <CRow v-if="found">
            <CCol sm="12">
                <CCard>
                    <CCardHeader>Edit User</CCardHeader>
                    <CForm @submit.prevent="doSubmit" method="post">
                        <CCardBody>
                            <p v-if="errorMsg!=''" class="alert alert-danger">{{errorMsg}}</p>
                            <CInput
                                    type="text"
                                    v-model="user.username"
                                    label="Username"
                                    placeholder="Enter user's username..."
                                    horizontal
                                    readonly="readonly"
                            />
                            <CInput
                                    type="text"
                                    v-model="user.name"
                                    label="Name"
                                    description="Enter user's name"
                                    placeholder="Enter user's name..."
                                    horizontal
                                    required
                                    was-validated
                            />
                            <CSelect
                                    label="Group"
                                    description="Assign user to group"
                                    :options="groupList"
                                    :value.sync="user.groupId"
                                    horizontal
                            />
                        </CCardBody>
                        <CCardFooter>
                            <CButton type="submit" color="primary" style="width: 96px">
                                <CIcon name="cil-save"/>
                                Save
                            </CButton>
                            <CButton type="button" color="info" class="ml-2" style="width: 96px" @click="doCancel">
                                <CIcon name="cil-arrow-circle-left"/>
                                Back
                            </CButton>
                        </CCardFooter>
                        <CCardHeader>Change Password</CCardHeader>
                        <CCardBody>
                            <CInput
                                    type="password"
                                    v-model="user.password"
                                    label="Current Password"
                                    description="Enter current password to verify password change"
                                    placeholder="Enter user's current password..."
                                    horizontal
                            />
                            <CInput
                                    type="password"
                                    v-model="user.newPassword"
                                    label="New Password"
                                    description="Enter new password to change"
                                    placeholder="Enter new password..."
                                    horizontal
                                    :is-valid="validatorPassword"
                                    invalid-feedback="New password must match the confirmed one"
                            />
                            <CInput
                                    type="password"
                                    v-model="user.newPassword2"
                                    label="Confirm New Password"
                                    description="Confirm new password"
                                    placeholder="Enter new password again..."
                                    horizontal
                                    :is-valid="validatorPassword"
                                    invalid-feedback="New password must match the confirmed one"
                            />
                        </CCardBody>
                        <CCardFooter>
                            <CButton type="submit" color="primary" style="width: 96px">
                                <CIcon name="cil-save"/>
                                Save
                            </CButton>
                            <CButton type="button" color="info" class="ml-2" style="width: 96px" @click="doCancel">
                                <CIcon name="cil-arrow-circle-left"/>
                                Back
                            </CButton>
                        </CCardFooter>
                    </CForm>
                </CCard>
            </CCol>
        </CRow>
    </div>
</template>

<script>
    import router from "@/router"
    import clientUtils from "@/utils/api_client"
    import utils from "@/utils/app_utils"

    export default {
        name: 'EditUser',
        data() {
            let groupList = {data: []}
            clientUtils.apiDoGet(clientUtils.apiGroupList,
                (apiRes) => {
                    if (apiRes.status == 200) {
                        apiRes.data.every(function (e) {
                            groupList.data.push({value: e.id, label: e.name})
                            return true
                        })
                        clientUtils.apiDoGet(clientUtils.apiUser + "/" + this.$route.params.username,
                            (apiRes) => {
                                this.found = apiRes.status == 200
                                if (apiRes.status == 200) {
                                    this.user.username = apiRes.data.username
                                    this.user.name = apiRes.data.name
                                    this.user.groupId = apiRes.data.gid
                                }
                            },
                            (err) => {
                                this.errorMsg = err
                            })
                    } else {
                        console.error("Getting group list was unsuccessful: " + apiRes)
                    }
                },
                (err) => {
                    console.error("Error getting group list: " + err)
                })
            return {
                user: {username: "", name: "", groupId: "", password: "", newPassword: "", newPassword2: ""},
                groupList: groupList.data,
                errorMsg: "",
                found: true,
            }
        },
        methods: {
            doCancel() {
                router.push("/users")
            },
            doSubmit(e) {
                e.preventDefault()
                let data = {
                    name: this.user.name,
                    group_id: this.user.groupId,
                    password: this.user.password,
                    new_password: this.user.newPassword,
                    new_password2: this.user.newPassword2,
                }
                clientUtils.apiDoPut(
                    clientUtils.apiUser + "/" + this.$route.params.username, data,
                    (apiRes) => {
                        if (apiRes.status != 200) {
                            this.errorMsg = apiRes.status + ": " + apiRes.message
                        } else {
                            utils.localStorageSet(utils.lskeyLoginSessionLastCheck, null)
                            this.$router.push({
                                name: "Users",
                                params: {flashMsg: "User [" + this.user.username + "] has been updated successfully."},
                            })
                        }
                    },
                    (err) => {
                        this.errorMsg = err
                    }
                )
            },
            validatorPassword(val) {
                return this.user.password == "" || this.user.newPassword == this.user.newPassword2
            },
        }
    }
</script>
