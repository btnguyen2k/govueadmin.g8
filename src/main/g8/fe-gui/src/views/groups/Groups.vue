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
                        <!--
                        <CLink class="card-header-action btn-minimize" @click="clickAddGroup">
                            <CIcon name="cil-library-add"/>
                        </CLink>
                        -->
                    </div>
                </CCardHeader>
                <CCardBody>
                    <p v-if="flashMsg" class="alert alert-success">{{flashMsg}}</p>
                    <CDataTable :items="groupList.data" :fields="['id','name','actions']">
                        <template #actions="{item}">
                            <td>{{item}}</td>
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
            clientUtils.apiDoGet(clientUtils.apiGroups,
                (apiRes) => {
                    if (apiRes.status == 200) {
                        groupList.data = apiRes.data
                    } else {
                    }
                },
                (err) => {
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
        }
    }
</script>
