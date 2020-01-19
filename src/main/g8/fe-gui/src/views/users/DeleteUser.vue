<template>
    <div>
        <CRow>
            <CCol sm="12">
                <CCard>
                    <CCardHeader>Delete User</CCardHeader>
                    <CForm @submit.prevent="doSubmit" method="post">
                        <CCardBody>
                            <p v-if="!found" class="alert alert-danger">User [{{this.$route.params.username}}] not
                                found</p>
                            <p v-if="errorMsg!=''" class="alert alert-danger">{{errorMsg}}</p>
                            <CInput v-if="found"
                                    type="text"
                                    v-model="user.username"
                                    label="Username"
                                    placeholder="Enter user's username..."
                                    horizontal
                                    disabled="disabled"
                            />
                            <CInput v-if="found"
                                    type="text"
                                    v-model="user.name"
                                    label="Name"
                                    placeholder="Enter user's name..."
                                    horizontal
                                    disabled="disabled"
                            />
                            <CSelect v-if="found"
                                     label="Group"
                                     :options="groupList"
                                     :value.sync="user.groupId"
                                     horizontal
                                     disabled="disabled"
                            />
                        </CCardBody>
                        <CCardFooter>
                            <CButton v-if="found" type="submit" color="danger" style="width: 96px">
                                <CIcon name="cil-trash"/>
                                Delete
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
        name: 'DeleteUser',
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
                user: {username: "", name: "", groupId: ""},
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
                clientUtils.apiDoDelete(
                    clientUtils.apiUser + "/" + this.$route.params.username,
                    (apiRes) => {
                        if (apiRes.status != 200) {
                            this.errorMsg = apiRes.status + ": " + apiRes.message
                        } else {
                            utils.localStorageSet(utils.lskeyLoginSessionLastCheck, null)
                            this.$router.push({
                                name: "Users",
                                params: {flashMsg: "User [" + this.user.username + "] has been deleted successfully."},
                            })
                        }
                    },
                    (err) => {
                        this.errorMsg = err
                    }
                )
            },
        }
    }
</script>
