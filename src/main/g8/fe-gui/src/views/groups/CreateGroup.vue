<template>
    <div>
        <CRow>
            <CCol sm="12">
                <CCard>
                    <CCardHeader>Create New Group</CCardHeader>
                    <CForm @submit.prevent="doSubmit" method="post">
                        <CCardBody>
                            <CInput
                                    id="id" name="id"
                                    type="text"
                                    description="Please enter group id, format [0-9a-z_]+, must be unique."
                                    label="Id"
                                    placeholder="Enter group id"
                                    horizontal
                                    required
                                    :is-valid="validatorGroupId"
                                    was-validated
                            />
                            <CInput
                                    id="name" name="name"
                                    type="text"
                                    description="Please enter group name"
                                    label="Name"
                                    placeholder="Enter group name..."
                                    horizontal
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
    let patternGroupId = /^[0-9a-z_]$/

    export default {
        name: 'Forms',
        data() {
            return {
                // selected: [], // Must be an array reference!
                // show: true,
                // horizontal: {label: 'col-3', input: 'col-9'},
                // options: ['Option 1', 'Option 2', 'Option 3'],
                // selectOptions: [
                //     'Option 1', 'Option 2', 'Option 3',
                //     {
                //         value: ['some value', 'another value'],
                //         label: 'Selected option'
                //     }
                // ],
                // selectedOption: ['some value', 'another value'],
                //
                // formCollapsed: true,
                // checkboxNames: ['Checkboxes', 'Inline Checkboxes',
                //     'Checkboxes - custom', 'Inline Checkboxes - custom'],
                // radioNames: ['Radios', 'Inline Radios',
                //     'Radios - custom', 'Inline Radios - custom']
            }
        },
        methods: {
            doCancel() {
                router.push({name: "Groups"})
            },
            doSubmit(e) {
                e.preventDefault()
            },
            validatorGroupId(val) {
                return val ? patternGroupId.test(val) : false
            },
        }
    }
</script>
