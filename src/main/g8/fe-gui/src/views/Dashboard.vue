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
                        <strong>Group (1)</strong>
                        <div class="card-header-actions">
                            <CLink class="card-header-action btn-minimize" @click="clickAddGroup">
                                <CIcon name="cil-library-add"/>
                            </CLink>
                            <CLink class="card-header-action btn-minimize" @click="isCollapsed = !isCollapsed">
                                <CIcon :name="`cil-chevron-${isCollapsed ? 'bottom' : 'top'}`"/>
                            </CLink>
                        </div>
                    </CCardHeader>
                    <CCollapse :show="isCollapsed" :duration="400">
                        <CCardBody>
                            <CDataTable
                                    :items="groupList.data"
                                    :fields="['id','name','actions']"
                            >
                                <template #actions="{item}">
                                    <td>{{item}}</td>
                                </template>
                            </CDataTable>
                            <!--
                            <CDataTable
                                    :hover="hover"
                                    :striped="striped"
                                    :bordered="bordered"
                                    :small="small"
                                    :fixed="fixed"
                                    :items="items"
                                    :fields="fields"
                                    :items-per-page="small ? 10 : 5"
                                    :dark="dark"
                                    pagination
                            >
                                <template #status="{item}">
                                    <td>
                                        <CBadge :color="getBadge(item.status)">{{item.status}}</CBadge>
                                    </td>
                                </template>
                            </CDataTable>
                            -->
                        </CCardBody>
                    </CCollapse>
                </CCard>
            </CCol>
            <CCol sm="12" md="6">
                <CCard accent-color="success">
                    <CCardHeader>
                        <strong>User (2)</strong>
                        <div class="card-header-actions">
                            <CLink class="card-header-action btn-minimize" @click="isCollapsed = !isCollapsed">
                                <CIcon :name="`cil-chevron-${isCollapsed ? 'bottom' : 'top'}`"/>
                            </CLink>
                        </div>
                    </CCardHeader>
                    <CCollapse :show="isCollapsed" :duration="400">
                        <CCardBody>
                            {{loremIpsum}}
                        </CCardBody>
                    </CCollapse>
                </CCard>
            </CCol>
        </CRow>
    </div>
</template>

<script>
    import MainChartExample from './charts/MainChartExample'
    import WidgetsDropdown from './widgets/WidgetsDropdown'
    import WidgetsBrand from './widgets/WidgetsBrand'
    import {CChartLineSimple, CChartBarSimple} from './charts/index.js'
    import clientUtils from "@/utils/api_client"
    import CTableWrapper from '@/views/base/Table.vue'

    export default {
        name: 'Dashboard',
        components: {
            MainChartExample,
            WidgetsDropdown,
            WidgetsBrand,
            CChartLineSimple, CChartBarSimple, CTableWrapper,
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
                    }
                },
                (err) => {
                }
            )

            let groupList = {data: []}
            clientUtils.apiDoGet(clientUtils.apiGroupList,
                (apiRes) => {
                    if (apiRes.status == 200) {
                        groupList.data = apiRes.data
                    } else {
                    }
                },
                (err) => {
                })
            return {
                isCollapsed: true,
                loremIpsum: 'Lorem ipsum dolor sit amet, consectetuer adipiscing elit, sed diam nonummy nibh euismod tincidunt ut laoreet dolore magna aliquam erat volutpat. Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat.',

                systemInfo: systemInfo,
                groupList: groupList,

                selected: 'Month',
                tableItems: [
                    {
                        avatar: {url: 'img/avatars/1.jpg', status: 'success'},
                        user: {name: 'Yiorgos Avraamu', new: true, registered: 'Jan 1, 2015'},
                        country: {name: 'USA', flag: 'cif-us'},
                        usage: {value: 50, period: 'Jun 11, 2015 - Jul 10, 2015'},
                        payment: {name: 'Mastercard', icon: 'cib-cc-mastercard'},
                        activity: '10 sec ago'
                    },
                    {
                        avatar: {url: 'img/avatars/2.jpg', status: 'danger'},
                        user: {name: 'Avram Tarasios', new: false, registered: 'Jan 1, 2015'},
                        country: {name: 'Brazil', flag: 'cif-br'},
                        usage: {value: 22, period: 'Jun 11, 2015 - Jul 10, 2015'},
                        payment: {name: 'Visa', icon: 'cib-cc-visa'},
                        activity: '5 minutes ago'
                    },
                    {
                        avatar: {url: 'img/avatars/3.jpg', status: 'warning'},
                        user: {name: 'Quintin Ed', new: true, registered: 'Jan 1, 2015'},
                        country: {name: 'India', flag: 'cif-in'},
                        usage: {value: 74, period: 'Jun 11, 2015 - Jul 10, 2015'},
                        payment: {name: 'Stripe', icon: 'cib-stripe'},
                        activity: '1 hour ago'
                    },
                    {
                        avatar: {url: 'img/avatars/4.jpg', status: ''},
                        user: {name: 'Enéas Kwadwo', new: true, registered: 'Jan 1, 2015'},
                        country: {name: 'France', flag: 'cif-fr'},
                        usage: {value: 98, period: 'Jun 11, 2015 - Jul 10, 2015'},
                        payment: {name: 'PayPal', icon: 'cib-paypal'},
                        activity: 'Last month'
                    },
                    {
                        avatar: {url: 'img/avatars/5.jpg', status: 'success'},
                        user: {name: 'Agapetus Tadeáš', new: true, registered: 'Jan 1, 2015'},
                        country: {name: 'Spain', flag: 'cif-es'},
                        usage: {value: 22, period: 'Jun 11, 2015 - Jul 10, 2015'},
                        payment: {name: 'Google Wallet', icon: 'cib-google-pay'},
                        activity: 'Last week'
                    },
                    {
                        avatar: {url: 'img/avatars/6.jpg', status: 'danger'},
                        user: {name: 'Friderik Dávid', new: true, registered: 'Jan 1, 2015'},
                        country: {name: 'Poland', flag: 'cif-pl'},
                        usage: {value: 43, period: 'Jun 11, 2015 - Jul 10, 2015'},
                        payment: {name: 'Amex', icon: 'cib-cc-amex'},
                        activity: 'Last week'
                    }
                ],
                tableFields: [
                    {key: 'avatar', label: '', _classes: 'text-center'},
                    {key: 'user'},
                    {key: 'country', _classes: 'text-center'},
                    {key: 'usage'},
                    {key: 'payment', label: 'Payment method', _classes: 'text-center'},
                    {key: 'activity'},
                ]
            }
        },
        methods: {
            clickAddGroup(e) {
                this.$router.push({name:"CreateGroup"})
            },
        }
    }
</script>
