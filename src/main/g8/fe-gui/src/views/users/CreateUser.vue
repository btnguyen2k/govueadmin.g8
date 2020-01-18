<template>
    <div>
        <CRow>
            <CCol sm="12">
                <CCard>
                    <CCardHeader>Create New User</CCardHeader>
                    <CForm @submit.prevent="doSubmit" method="post">
                        <CCardBody>
                            <p v-if="errorMsg!=''" class="alert alert-danger">{{errorMsg}}</p>
                            <CInput
                                    id="username" name="username"
                                    type="text"
                                    v-model="form.username"
                                    label="Username"
                                    placeholder="Enter user's username..."
                                    horizontal
                                    :is-valid="validatorUsername"
                                    invalid-feedback="Please enter user's username, format [0-9a-z_]+, must be unique."
                                    valid-feedback="Please enter user's username, format [0-9a-z_]+, must be unique."
                            />
                            <CInput
                                    id="name" name="name"
                                    type="text"
                                    v-model="form.name"
                                    label="Name"
                                    description="Please enter user's name"
                                    placeholder="Enter user's name..."
                                    horizontal
                                    required
                                    was-validated
                            />
                            <CSelect
                                    id="groupId" name="groupId"
                                    v-model="form.groupId"
                                    label="Group"
                                    description="Please assign user to group"
                                    :options="groupList"
                                    horizontal
                            />
                        </CCardBody>
                        <CCardFooter>
                            <CButton type="submit" color="primary" style="width: 96px">
                                <CIcon name="cil-save"/>
                                Create
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
    import clientUtils from "@/utils/api_client";

    let patternUsername = /^[0-9a-z_]+$/

    export default {
        name: 'CreateUser',
        data() {
            let groupList = {data: []}
            clientUtils.apiDoGet(clientUtils.apiGroupList,
                (apiRes) => {
                    if (apiRes.status == 200) {
                        apiRes.data.every(function (e) {
                            groupList.data.push({value: e.id, label: e.name})
                            return true
                        })
                    } else {
                        console.error("Getting group list was unsuccessful: " + apiRes)
                    }
                },
                (err) => {
                    console.error("Error getting group list: " + err)
                })
            return {
                groupList: groupList.data,
                form: {username: "", name: "", groupId: ""},
                errorMsg: "",
            }
        },
        methods: {
            doCancel() {
                router.push("/users")
            },
            doSubmit(e) {
                e.preventDefault()
                let data = {username: this.form.username, name: this.form.name, group_id: this.form.groupId}
                clientUtils.apiDoPost(
                    clientUtils.apiUserList, data,
                    (apiRes) => {
                        if (apiRes.status != 200) {
                            this.errorMsg = apiRes.status + ": " + apiRes.message
                        } else {
                            this.$router.push({
                                name: "Users",
                                params: {flashMsg: "User [" + this.form.username + "] has been created successfully."},
                            })
                        }
                    },
                    (err) => {
                        console.error(err)
                        this.errorMsg = err
                    }
                )
            },
            validatorUsername(val) {
                return val ? patternUsername.test(val.toString()) : false
            },
        }
    }
</script>
