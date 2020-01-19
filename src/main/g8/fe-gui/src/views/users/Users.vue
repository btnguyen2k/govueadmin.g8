<template>
    <CRow>
        <CCol sm="12">
            <CCard accent-color="info">
                <CCardHeader>
                    <strong>User ({{userList.data.length}})</strong>
                    <div class="card-header-actions">
                        <CButton class="btn-sm btn-primary" @click="clickAddUser">
                            <CIcon name="cil-playlist-add"/>
                            Create New User
                        </CButton>
                    </div>
                </CCardHeader>
                <CCardBody>
                    <p v-if="flashMsg" class="alert alert-success">{{flashMsg}}</p>
                    <CDataTable :items="userList.data" :fields="['username','name',{key:'gid',label:'Group'},'actions']">
                        <template #actions="{item}">
                            <td>
                                <CLink @click="clickEditUser(item.username)" label="Edit" class="btn btn-primary">
                                    <CIcon name="cil-pencil"/>
                                </CLink>
                                &nbsp;
                                <CLink @click="clickDeleteUser(item.username)" label="Delete" class="btn btn-danger">
                                    <CIcon name="cil-trash"/>
                                </CLink>
                            </td>
                        </template>
                    </CDataTable>
                </CCardBody>
            </CCard>
        </CCol>
    </CRow>
</template>

<script>
    import clientUtils from "@/utils/api_client";

    export default {
        name: 'Users',
        data: () => {
            let userList = {data: []}
            clientUtils.apiDoGet(clientUtils.apiUserList,
                (apiRes) => {
                    if (apiRes.status == 200) {
                        userList.data = apiRes.data
                    } else {
                        console.error("Getting user list was unsuccessful: " + apiRes)
                    }
                },
                (err) => {
                    console.error("Error getting user list: " + err)
                })

            return {
                userList: userList,
            }
        },
        props: ["flashMsg"],
        methods: {
            clickAddUser(e) {
                this.$router.push({name: "CreateUser"})
            },
            clickEditUser(username) {
                this.$router.push({name: "EditUser", params: {username: username.toString()}})
            },
            clickDeleteUser(username) {
                this.$router.push({name: "DeleteUser", params: {username: username.toString()}})
            },
        }
    }
</script>
