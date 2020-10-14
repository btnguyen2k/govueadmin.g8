<template>
    <div>
        <CRow>
            <CCol sm="12">
                <CCard>
                    <CCardHeader>Edit Group</CCardHeader>
                    <CForm @submit.prevent="doSubmit" method="post">
                        <CCardBody>
                            <p v-if="!found" class="alert alert-danger">Group [{{this.$route.params.id}}] not found</p>
                            <p v-if="errorMsg!=''" class="alert alert-danger">{{errorMsg}}</p>
                            <CInput v-if="found"
                                    type="text"
                                    v-model="group.id"
                                    label="Id"
                                    placeholder="Enter group id"
                                    horizontal
                                    readonly="readonly"
                            />
                            <CInput v-if="found"
                                    type="text"
                                    v-model="group.name"
                                    label="Name"
                                    description="Please enter group name"
                                    placeholder="Enter group name..."
                                    horizontal
                                    required
                                    was-validated
                            />
                        </CCardBody>
                        <CCardFooter>
                            <CButton v-if="found" type="submit" color="primary" style="width: 96px">
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
    import clientUtils from "@/utils/api_client";

    export default {
        name: 'EditGroup',
        data() {
            clientUtils.apiDoGet(clientUtils.apiGroup + "/" + this.$route.params.id,
                (apiRes) => {
                    this.found = apiRes.status == 200
                    if (apiRes.status == 200) {
                        this.group = apiRes.data
                    }
                },
                (err) => {
                    this.errorMsg = err
                })
            return {
                group: {id: "", name: ""},
                errorMsg: "",
                found: true,
            }
        },
        methods: {
            doCancel() {
                router.push("/groups")
            },
            doSubmit(e) {
                e.preventDefault()
                let data = {name: this.group.name}
                clientUtils.apiDoPut(
                    clientUtils.apiGroup + "/" + this.$route.params.id, data,
                    (apiRes) => {
                        if (apiRes.status != 200) {
                            this.errorMsg = apiRes.status + ": " + apiRes.message
                        } else {
                            this.$router.push({
                                name: "Groups",
                                params: {flashMsg: "Group [" + this.group.id + "] has been updated successfully."},
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
