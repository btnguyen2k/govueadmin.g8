<template>
    <div>
        <CRow>
            <CCol sm="6" lg="3">
                <CWidgetDropdown color="primary"
                                 v-bind:text="systemInfo.cpu.cores+' core(s) / '+systemInfo.cpu.load+' load'"
                                 header="CPU">
                    <template #default>
                        <CDropdown color="transparent p-0" placement="bottom-end">
                            <template #toggler-content>
                                <CIcon name="cil-settings"/>
                            </template>
                            <CDropdownItem>Action</CDropdownItem>
                            <CDropdownItem>Another action</CDropdownItem>
                            <CDropdownItem>Something else here...</CDropdownItem>
                            <CDropdownItem disabled>Disabled action</CDropdownItem>
                        </CDropdown>
                    </template>
                    <template #footer>
                        <CChartLineSimple pointed class="mt-3 mx-3" style="height:70px"
                                          :data-points="systemInfo.cpu.history_load"
                                          point-hover-background-color="primary"
                                          label="Load"
                        />
                    </template>
                </CWidgetDropdown>
            </CCol>
            <CCol sm="6" lg="3">
                <CWidgetDropdown color="info" v-bind:text="systemInfo.go_routines.num+''" header="Go Routines">
                    <template #default>
                        <CDropdown color="transparent p-0" placement="bottom-end">
                            <template #toggler-content>
                                <CIcon name="cil-settings"/>
                            </template>
                            <CDropdownItem>Action</CDropdownItem>
                            <CDropdownItem>Another action</CDropdownItem>
                            <CDropdownItem>Something else here...</CDropdownItem>
                            <CDropdownItem disabled>Disabled action</CDropdownItem>
                        </CDropdown>
                    </template>
                    <template #footer>
                        <CChartLineSimple pointed class="mt-3 mx-3" style="height:70px"
                                          :data-points="systemInfo.go_routines.history"
                                          point-hover-background-color="info"
                                          :options="{ elements: { line: { tension: 0.00001 }}}"
                                          label="AppMemory"
                        />
                    </template>
                </CWidgetDropdown>
            </CCol>
            <CCol sm="6" lg="3">
                <CWidgetDropdown color="warning"
                                 v-bind:text="systemInfo.app_memory.usedMb+' Mb'"
                                 header="Used AppMemory">
                    <template #default>
                        <CDropdown color="transparent p-0" placement="bottom-end" :caret="false">
                            <template #toggler-content>
                                <CIcon name="cil-location-pin"/>
                            </template>
                            <CDropdownItem>Action</CDropdownItem>
                            <CDropdownItem>Another action</CDropdownItem>
                            <CDropdownItem>Something else here...</CDropdownItem>
                            <CDropdownItem disabled>Disabled action</CDropdownItem>
                        </CDropdown>
                    </template>
                    <template #footer>
                        <CChartLineSimple class="mt-3" style="height:70px" background-color="rgba(255,255,255,.2)"
                                          :data-points="systemInfo.app_memory.history_usedMb"
                                          :options="{ elements: { line: { borderWidth: 2.5 }}}"
                                          point-hover-background-color="warning" label="GoRoutines"
                        />
                    </template>
                </CWidgetDropdown>
            </CCol>
            <CCol sm="6" lg="3">
                <CWidgetDropdown color="danger" v-bind:text="systemInfo.memory.freeGb+' Gb'" header="Free SystemMemory">
                    <template #default>
                        <CDropdown color="transparent p-0" placement="bottom-end">
                            <template #toggler-content>
                                <CIcon name="cil-settings"/>
                            </template>
                            <CDropdownItem>Action</CDropdownItem>
                            <CDropdownItem>Another action</CDropdownItem>
                            <CDropdownItem>Something else here...</CDropdownItem>
                            <CDropdownItem disabled>Disabled action</CDropdownItem>
                        </CDropdown>
                    </template>
                    <template #footer>
                        <CChartLineSimple class="mt-3" style="height:70px" background-color="rgb(250, 152, 152)"
                                          :data-points="systemInfo.memory.history_freeGb"
                                          :options="{ elements: { line: { borderWidth: 2.5 }}}"
                                          point-hover-background-color="danger" label="Free SystemMemory"
                        />
                    </template>
                </CWidgetDropdown>
            </CCol>
        </CRow>
        <CRow>
            <CCol sm="12" md="6">
                <CCard accent-color="info">
                    <CCardHeader>
                        <strong>Group ({{groupList.data.length}})</strong>
                        <div class="card-header-actions">
                            <CLink class="card-header-action btn-minimize" @click="clickAddGroup">
                                <CIcon name="cil-library-add"/>
                            </CLink>
                            <CLink class="card-header-action btn-minimize"
                                   @click="isCollapsedGroups = !isCollapsedGroups">
                                <CIcon :name="`cil-chevron-${isCollapsedGroups ? 'bottom' : 'top'}`"/>
                            </CLink>
                        </div>
                    </CCardHeader>
                    <CCollapse :show="isCollapsedGroups" :duration="400">
                        <CCardBody>
                            <CDataTable :items="groupList.data" :fields="['id','name','actions']">
                                <template #actions="{item}">
                                    <td>
                                        <CLink @click="clickEditGroup(item.id)" label="Edit" class="btn-sm btn-primary">
                                            <CIcon name="cil-pencil"/>
                                        </CLink>
                                        &nbsp;
                                        <CLink @click="clickDeleteGroup(item.id)" label="Delete" class="btn-sm btn-danger">
                                            <CIcon name="cil-trash"/>
                                        </CLink>
                                    </td>
                                </template>
                            </CDataTable>
                        </CCardBody>
                    </CCollapse>
                </CCard>
            </CCol>
            <CCol sm="12" md="6">
                <CCard accent-color="success">
                    <CCardHeader>
                        <strong>User ({{userList.data.length}})</strong>
                        <div class="card-header-actions">
                            <CLink class="card-header-action btn-minimize" @click="clickAddUser">
                                <CIcon name="cil-playlist-add"/>
                            </CLink>
                            <CLink class="card-header-action btn-minimize"
                                   @click="isCollapsedUsers = !isCollapsedUsers">
                                <CIcon :name="`cil-chevron-${isCollapsedUsers ? 'bottom' : 'top'}`"/>
                            </CLink>
                        </div>
                    </CCardHeader>
                    <CCollapse :show="isCollapsedUsers" :duration="400">
                        <CCardBody>
                            <CDataTable :items="userList.data" :fields="['username','name',{key:'gid',label:'Group'},'actions']">
                                <template #actions="{item}">
                                    <td>
                                        <CLink @click="clickEditUser(item.username)" label="Edit" class="btn-sm btn-primary">
                                            <CIcon name="cil-pencil"/>
                                        </CLink>
                                        &nbsp;
                                        <CLink @click="clickDeleteUser(item.username)" label="Delete" class="btn-sm btn-danger">
                                            <CIcon name="cil-trash"/>
                                        </CLink>
                                    </td>
                                </template>
                            </CDataTable>
                        </CCardBody>
                    </CCollapse>
                </CCard>
            </CCol>
        </CRow>
    </div>
</template>

<script>
    import {CChartLineSimple} from './charts/index.js'
    import clientUtils from "@/utils/api_client"

    export default {
        name: 'Dashboard',
        components: {
            CChartLineSimple,
        },
        data() {
            let systemInfo = {
                cpu: {cores: -1, load: -1.0, history_load: []},
                memory: {free: 0, freeGb: 0.0, history_freeGb: []},
                app_memory: {usedMb: 0.0, history_usedMb: []},
                go_routines: {num: 0, history: []},
            }
            clientUtils.apiDoGet(
                clientUtils.apiSystemInfo,
                (apiRes) => {
                    if (apiRes.status == 200) {
                        systemInfo.cpu = apiRes.data.cpu
                        systemInfo.memory = apiRes.data.memory
                        systemInfo.app_memory = apiRes.data.app_memory
                        systemInfo.go_routines = apiRes.data.go_routines
                    } else {
                        console.error("Getting system info was unsuccessful: " + JSON.stringify(apiRes))
                    }
                },
                (err) => {
                    console.error("Error getting system info: " + err)
                }
            )

            let groupList = {data: []}
            clientUtils.apiDoGet(clientUtils.apiGroupList,
                (apiRes) => {
                    if (apiRes.status == 200) {
                        groupList.data = apiRes.data
                    } else {
                        console.error("Getting group list was unsuccessful: " + JSON.stringify(apiRes))
                    }
                },
                (err) => {
                    console.error("Error getting group list: " + err)
                })

            let userList = {data: []}
            clientUtils.apiDoGet(clientUtils.apiUserList,
                (apiRes) => {
                    if (apiRes.status == 200) {
                        userList.data = apiRes.data
                    } else {
                        console.error("Getting user list was unsuccessful: " + JSON.stringify(apiRes))
                    }
                },
                (err) => {
                    console.error("Error getting user list: " + err)
                })

            return {
                isCollapsedGroups: true,
                isCollapsedUsers: true,
                systemInfo: systemInfo,
                groupList: groupList,
                userList: userList,
            }
        },
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
