<template>
    <CRow>
        <CCol sm="12">
            <CCard accent-color="info">
                <CCardHeader>
                    <strong>Group ({{groupList.data.length}})</strong>
                    <div class="card-header-actions">
                        <CButton class="btn-sm btn-primary" @click="clickAddGroup">
                            <CIcon name="cil-library-add"/>
                            Create New Group
                        </CButton>
                    </div>
                </CCardHeader>
                <CCardBody>
                    <p v-if="flashMsg" class="alert alert-success">{{flashMsg}}</p>
                    <CDataTable :items="groupList.data" :fields="['id','name','actions']">
                        <template #actions="{item}">
                            <td>
                                <CLink @click="clickEditGroup(item.id)" label="Edit" class="btn btn-primary">
                                    <CIcon name="cil-pencil"/>
                                </CLink>
                                &nbsp;
                                <CLink @click="clickDeleteGroup(item.id)" label="Delete" class="btn btn-danger">
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
        name: 'Groups',
        data: () => {
            let groupList = {data: []}
            clientUtils.apiDoGet(clientUtils.apiGroupList,
                (apiRes) => {
                    if (apiRes.status == 200) {
                        groupList.data = apiRes.data
                    } else {
                        console.error("Getting group list was unsuccessful: " + apiRes)
                    }
                },
                (err) => {
                    console.error("Error getting group list: " + err)
                })

            return {
                groupList: groupList,
            }
        },
        props: ["flashMsg"],
        methods: {
            clickAddGroup(e) {
                this.$router.push({name: "CreateGroup"})
            },
            clickEditGroup(id) {
                this.$router.push({name: "EditGroup", params: {id: id.toString()}})
            },
            clickDeleteGroup(id) {
                this.$router.push({name: "DeleteGroup", params: {id: id.toString()}})
            },
        }
    }
</script>
