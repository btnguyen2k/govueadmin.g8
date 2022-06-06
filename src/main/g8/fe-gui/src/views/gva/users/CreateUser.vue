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
                                    type="text"
                                    v-model="form.username"
                                    label="Username"
                                    placeholder="Enter user's username..."
                                    horizontal
                                    :is-valid="validatorUsername"
                                    invalid-feedback="Enter user's username, format [0-9a-z_]+, must be unique."
                                    valid-feedback="Enter user's username, format [0-9a-z_]+, must be unique."
                            />
                            <CInput
                                    type="password"
                                    v-model="form.password"
                                    label="Password"
                                    description="Enter user's password"
                                    placeholder="Enter user's password..."
                                    horizontal
                                    required
                                    :is-valid="validatorPassword"
                                    invalid-feedback="Password must match the confirmed one"
                            />
                            <CInput
                                    type="password"
                                    v-model="form.password2"
                                    label="Confirmed Password"
                                    description="Confirm user's password"
                                    placeholder="Enter user's password again..."
                                    horizontal
                                    required
                                    :is-valid="validatorPassword"
                                    invalid-feedback="Password must match the confirmed one"
                            />
                            <CInput
                                    type="text"
                                    v-model="form.name"
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
                                    :value.sync="form.groupId"
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
            let form = {
                username: "", password: "", password2: "",
                name: "",
                groupId: groupList.data.length > 0 ? groupList.data[0].value : ""
            }
            clientUtils.apiDoGet(clientUtils.apiGroupList,
                (apiRes) => {
                    if (apiRes.status == 200) {
                        apiRes.data.every(function (e) {
                            if (form.groupId == "") {
                                form.groupId = e.id
                            }
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
                form: form,
                errorMsg: "",
            }
        },
        methods: {
            doCancel() {
                router.push("/users")
            },
            doSubmit(e) {
                e.preventDefault()
                let data = {
                    username: this.form.username,
                    password: this.form.password, password2: this.form.password2,
                    name: this.form.name, group_id: this.form.groupId
                }
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
            validatorPassword(val) {
                return this.form.password == this.form.password2
            },
        }
    }
</script>
