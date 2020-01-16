<template>
    <div>
        <CRow>
            <CCol sm="12">
                <CCard>
                    <CCardHeader>Create New Group</CCardHeader>
                    <CForm @submit.prevent="doSubmit" method="post">
                        <CCardBody>
                            <p v-if="erroMsg!=''" class="alert alert-danger">{{erroMsg}}</p>
                            <CInput
                                    id="id" name="id"
                                    type="text"
                                    v-model="form.id"
                                    label="Id"
                                    placeholder="Enter group id"
                                    horizontal
                                    :is-valid="validatorGroupId"
                                    invalid-feedback="Please enter group id, format [0-9a-z_]+, must be unique."
                                    valid-feedback="Please enter group id, format [0-9a-z_]+, must be unique."
                            />
                            <CInput
                                    id="name" name="name"
                                    type="text"
                                    v-model="form.name"
                                    label="Name"
                                    description="Please enter group name"
                                    placeholder="Enter group name..."
                                    horizontal
                                    required
                                    was-validated
                            />
                        </CCardBody>
                        <CCardFooter>
                            <CButton type="submit" color="primary" style="width: 96px">
                                <CIcon name="cil-save"/>
                                Submit
                            </CButton>
                            &nbsp;
                            <CButton type="reset" color="warning" style="width: 96px">
                                <CIcon name="cil-ban"/>
                                Reset
                            </CButton>
                            &nbsp;
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
    import utils from "@/utils/app_utils";

    let patternGroupId = /^[0-9a-z_]+$/

    export default {
        name: 'CreateGroup',
        data() {
            return {
                form: {id: "", name: ""},
                erroMsg: "",
            }
        },
        methods: {
            doCancel() {
                router.push("/groups")
            },
            doSubmit(e) {
                e.preventDefault()
                let data = {id: this.form.id, name: this.form.name}
                clientUtils.apiDoPost(
                    clientUtils.apiGroups, data,
                    (apiRes) => {
                        console.log(apiRes)
                        if (apiRes.status != 200) {
                            this.erroMsg = apiRes.status + ": " + apiRes.message
                        } else {
                            this.$router.push({
                                name: "Groups",
                                params: {flashMsg: "Group [" + this.form.id + "] has been created successfully."},
                            })
                        }
                    },
                    (err) => {
                        this.erroMsg = err
                    }
                )
            },
            validatorGroupId(val) {
                return val ? patternGroupId.test(val.toString()) : false
            },
        }
    }
</script>
