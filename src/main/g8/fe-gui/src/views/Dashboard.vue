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
                            <CLink class="card-header-action btn-minimize" @click="isCollapsed = !isCollapsed">
                                <CIcon :name="`cil-chevron-${isCollapsed ? 'bottom' : 'top'}`"/>
                            </CLink>
                        </div>
                    </CCardHeader>
                    <CCollapse :show="isCollapsed" :duration="400">
                        <CCardBody>
                            <CTableWrapper :items="getShuffledUsersData()">
                                <template #header>
                                    <CIcon name="cil-grid"/>
                                    Simple Table
                                    <div class="card-header-actions">
                                        <a
                                                href="https://coreui.io/vue/docs/components/nav"
                                                class="card-header-action"
                                                rel="noreferrer noopener"
                                                target="_blank"
                                        >
                                            <small class="text-muted">docs</small>
                                        </a>
                                    </div>
                                </template>
                            </CTableWrapper>
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

        <!--        <CCard>-->
        <!--            <CCardBody>-->
        <!--                <CRow>-->
        <!--                    <CCol sm="5">-->
        <!--                        <h4 id="traffic" class="card-title mb-0">Traffic</h4>-->
        <!--                        <div class="small text-muted">November 2017</div>-->
        <!--                    </CCol>-->
        <!--                    <CCol sm="7" class="d-none d-md-block">-->
        <!--                        <CButton color="primary" class="float-right">-->
        <!--                            <CIcon name="cil-cloud-download"/>-->
        <!--                        </CButton>-->
        <!--                        <CButtonGroup class="float-right mr-3">-->
        <!--                            <CButton-->
        <!--                                    color="outline-secondary"-->
        <!--                                    v-for="(value, key) in ['Day', 'Month', 'Year']"-->
        <!--                                    :key="key"-->
        <!--                                    class="mx-0"-->
        <!--                                    :pressed="value === selected ? true : false"-->
        <!--                                    @click="selected = value"-->
        <!--                            >-->
        <!--                                {{value}}-->
        <!--                            </CButton>-->
        <!--                        </CButtonGroup>-->
        <!--                    </CCol>-->
        <!--                </CRow>-->
        <!--                <MainChartExample style="height:300px;margin-top:40px;"/>-->
        <!--            </CCardBody>-->
        <!--            <CCardFooter>-->
        <!--                <CRow class="text-center">-->
        <!--                    <CCol md sm="12" class="mb-sm-2 mb-0">-->
        <!--                        <div class="text-muted">Visits</div>-->
        <!--                        <strong>29.703 Users (40%)</strong>-->
        <!--                        <CProgress-->
        <!--                                class="progress-xs mt-2"-->
        <!--                                :precision="1"-->
        <!--                                color="success"-->
        <!--                                :value="40"-->
        <!--                        />-->
        <!--                    </CCol>-->
        <!--                    <CCol md sm="12" class="mb-sm-2 mb-0 d-md-down-none">-->
        <!--                        <div class="text-muted">Unique</div>-->
        <!--                        <strong>24.093 Users (20%)</strong>-->
        <!--                        <CProgress-->
        <!--                                class="progress-xs mt-2"-->
        <!--                                :precision="1"-->
        <!--                                color="info"-->
        <!--                                :value="20"-->
        <!--                        />-->
        <!--                    </CCol>-->
        <!--                    <CCol md sm="12" class="mb-sm-2 mb-0">-->
        <!--                        <div class="text-muted">Pageviews</div>-->
        <!--                        <strong>78.706 Views (60%)</strong>-->
        <!--                        <CProgress-->
        <!--                                class="progress-xs mt-2"-->
        <!--                                :precision="1"-->
        <!--                                color="warning"-->
        <!--                                :value="60"-->
        <!--                        />-->
        <!--                    </CCol>-->
        <!--                    <CCol md sm="12" class="mb-sm-2 mb-0">-->
        <!--                        <div class="text-muted">New Users</div>-->
        <!--                        <strong>22.123 Users (80%)</strong>-->
        <!--                        <CProgress-->
        <!--                                class="progress-xs mt-2"-->
        <!--                                :precision="1"-->
        <!--                                color="danger"-->
        <!--                                :value="80"-->
        <!--                        />-->
        <!--                    </CCol>-->
        <!--                    <CCol md sm="12" class="mb-sm-2 mb-0 d-md-down-none">-->
        <!--                        <div class="text-muted">Bounce Rate</div>-->
        <!--                        <strong>Average Rate (40.15%)</strong>-->
        <!--                        <CProgress-->
        <!--                                class="progress-xs mt-2"-->
        <!--                                :precision="1"-->
        <!--                                :value="40"-->
        <!--                        />-->
        <!--                    </CCol>-->
        <!--                </CRow>-->
        <!--            </CCardFooter>-->
        <!--        </CCard>-->
        <!--        <WidgetsBrand/>-->
        <!--        <CRow>-->
        <!--            <CCol md="12">-->
        <!--                <CCard>-->
        <!--                    <CCardHeader>-->
        <!--                        Traffic &amp; Sales-->
        <!--                    </CCardHeader>-->
        <!--                    <CCardBody>-->
        <!--                        <CRow>-->
        <!--                            <CCol sm="12" lg="6">-->
        <!--                                <CRow>-->
        <!--                                    <CCol sm="6">-->
        <!--                                        <CCallout color="info">-->
        <!--                                            <small class="text-muted">New Clients</small><br>-->
        <!--                                            <strong class="h4">9,123</strong>-->
        <!--                                        </CCallout>-->
        <!--                                    </CCol>-->
        <!--                                    <CCol sm="6">-->
        <!--                                        <CCallout color="danger">-->
        <!--                                            <small class="text-muted">Recurring Clients</small><br>-->
        <!--                                            <strong class="h4">22,643</strong>-->
        <!--                                        </CCallout>-->
        <!--                                    </CCol>-->
        <!--                                </CRow>-->
        <!--                                <hr class="mt-0">-->
        <!--                                <div class="progress-group mb-4">-->
        <!--                                    <div class="progress-group-prepend">-->
        <!--                    <span class="progress-group-text">-->
        <!--                      Monday-->
        <!--                    </span>-->
        <!--                                    </div>-->
        <!--                                    <div class="progress-group-bars">-->
        <!--                                        <CProgress-->
        <!--                                                class="progress-xs"-->
        <!--                                                color="info"-->
        <!--                                                :value="34"-->
        <!--                                        />-->
        <!--                                        <CProgress-->
        <!--                                                class="progress-xs"-->
        <!--                                                color="danger"-->
        <!--                                                :value="78"-->
        <!--                                        />-->
        <!--                                    </div>-->
        <!--                                </div>-->
        <!--                                <div class="progress-group mb-4">-->
        <!--                                    <div class="progress-group-prepend">-->
        <!--                    <span class="progress-group-text">-->
        <!--                      Tuesday-->
        <!--                    </span>-->
        <!--                                    </div>-->
        <!--                                    <div class="progress-group-bars">-->
        <!--                                        <CProgress-->
        <!--                                                class="progress-xs"-->
        <!--                                                :value="56"-->
        <!--                                                color="info"-->
        <!--                                        />-->
        <!--                                        <CProgress-->
        <!--                                                class="progress-xs"-->
        <!--                                                :value="94"-->
        <!--                                                color="danger"-->
        <!--                                        />-->
        <!--                                    </div>-->
        <!--                                </div>-->
        <!--                                <div class="progress-group mb-4">-->
        <!--                                    <div class="progress-group-prepend">-->
        <!--                    <span class="progress-group-text">-->
        <!--                      Wednesday-->
        <!--                    </span>-->
        <!--                                    </div>-->
        <!--                                    <div class="progress-group-bars">-->
        <!--                                        <CProgress-->
        <!--                                                class="progress-xs"-->
        <!--                                                :value="12"-->
        <!--                                                color="info"-->
        <!--                                        />-->
        <!--                                        <CProgress-->
        <!--                                                class="progress-xs"-->
        <!--                                                :value="67"-->
        <!--                                                color="danger"-->
        <!--                                        />-->
        <!--                                    </div>-->
        <!--                                </div>-->
        <!--                                <div class="progress-group mb-4">-->
        <!--                                    <div class="progress-group-prepend">-->
        <!--                    <span class="progress-group-text">-->
        <!--                      Thursday-->
        <!--                    </span>-->
        <!--                                    </div>-->
        <!--                                    <div class="progress-group-bars">-->
        <!--                                        <CProgress-->
        <!--                                                class="progress-xs"-->
        <!--                                                :value="43"-->
        <!--                                                color="info"-->
        <!--                                        />-->
        <!--                                        <CProgress-->
        <!--                                                class="progress-xs"-->
        <!--                                                :value="91"-->
        <!--                                                color="danger"-->
        <!--                                        />-->
        <!--                                    </div>-->
        <!--                                </div>-->
        <!--                                <div class="progress-group mb-4">-->
        <!--                                    <div class="progress-group-prepend">-->
        <!--                    <span class="progress-group-text">-->
        <!--                      Friday-->
        <!--                    </span>-->
        <!--                                    </div>-->
        <!--                                    <div class="progress-group-bars">-->
        <!--                                        <CProgress-->
        <!--                                                class="progress-xs"-->
        <!--                                                :value="22"-->
        <!--                                                color="info"-->
        <!--                                        />-->
        <!--                                        <CProgress-->
        <!--                                                class="progress-xs"-->
        <!--                                                :value="73"-->
        <!--                                                color="danger"-->
        <!--                                        />-->
        <!--                                    </div>-->
        <!--                                </div>-->
        <!--                                <div class="progress-group mb-4">-->
        <!--                                    <div class="progress-group-prepend">-->
        <!--                    <span class="progress-group-text">-->
        <!--                      Saturday-->
        <!--                    </span>-->
        <!--                                    </div>-->
        <!--                                    <div class="progress-group-bars">-->
        <!--                                        <CProgress-->
        <!--                                                class="progress-xs"-->
        <!--                                                :value="53"-->
        <!--                                                color="info"-->
        <!--                                        />-->
        <!--                                        <CProgress-->
        <!--                                                class="progress-xs"-->
        <!--                                                :value="82"-->
        <!--                                                color="danger"-->
        <!--                                        />-->
        <!--                                    </div>-->
        <!--                                </div>-->
        <!--                                <div class="progress-group mb-4">-->
        <!--                                    <div class="progress-group-prepend">-->
        <!--                    <span class="progress-group-text">-->
        <!--                      Sunday-->
        <!--                    </span>-->
        <!--                                    </div>-->
        <!--                                    <div class="progress-group-bars">-->
        <!--                                        <CProgress-->
        <!--                                                class="progress-xs"-->
        <!--                                                :value="9"-->
        <!--                                                color="info"-->
        <!--                                        />-->
        <!--                                        <CProgress-->
        <!--                                                class="progress-xs"-->
        <!--                                                :value="69"-->
        <!--                                                color="danger"-->
        <!--                                        />-->
        <!--                                    </div>-->
        <!--                                </div>-->
        <!--                                <div class="legend text-center">-->
        <!--                                    <small>-->
        <!--                                        <sup>-->
        <!--                                            <CBadge shape="pill" color="info">&nbsp;</CBadge>-->
        <!--                                        </sup>-->
        <!--                                        New clients-->
        <!--                                        &nbsp;&nbsp;-->
        <!--                                        <sup>-->
        <!--                                            <CBadge shape="pill" color="danger">&nbsp;</CBadge>-->
        <!--                                        </sup>-->
        <!--                                        Recurring clients-->
        <!--                                    </small>-->
        <!--                                </div>-->
        <!--                            </CCol>-->
        <!--                            <CCol sm="12" lg="6">-->
        <!--                                <CRow>-->
        <!--                                    <CCol sm="6">-->
        <!--                                        <CCallout color="warning">-->
        <!--                                            <small class="text-muted">Pageviews</small><br>-->
        <!--                                            <strong class="h4">78,623</strong>-->
        <!--                                        </CCallout>-->
        <!--                                    </CCol>-->
        <!--                                    <CCol sm="6">-->
        <!--                                        <CCallout color="success">-->
        <!--                                            <small class="text-muted">Organic</small><br>-->
        <!--                                            <strong class="h4">49,123</strong>-->
        <!--                                        </CCallout>-->
        <!--                                    </CCol>-->
        <!--                                </CRow>-->
        <!--                                <hr class="mt-0">-->
        <!--                                <ul class="horizontal-bars type-2">-->
        <!--                                    <div class="progress-group">-->
        <!--                                        <div class="progress-group-header">-->
        <!--                                            <CIcon name="cil-user" class="progress-group-icon"/>-->
        <!--                                            <span class="title">Male</span>-->
        <!--                                            <span class="ml-auto font-weight-bold">43%</span>-->
        <!--                                        </div>-->
        <!--                                        <div class="progress-group-bars">-->
        <!--                                            <CProgress-->
        <!--                                                    class="progress-xs"-->
        <!--                                                    :value="43"-->
        <!--                                                    color="warning"-->
        <!--                                            />-->
        <!--                                        </div>-->
        <!--                                    </div>-->
        <!--                                    <div class="progress-group mb-5">-->
        <!--                                        <div class="progress-group-header">-->
        <!--                                            <CIcon name="cil-user-female" class="progress-group-icon"/>-->
        <!--                                            <span class="title">Female</span>-->
        <!--                                            <span class="ml-auto font-weight-bold">37%</span>-->
        <!--                                        </div>-->
        <!--                                        <div class="progress-group-bars">-->
        <!--                                            <CProgress-->
        <!--                                                    class="progress-xs"-->
        <!--                                                    :value="37"-->
        <!--                                                    color="warning"-->
        <!--                                            />-->
        <!--                                        </div>-->
        <!--                                    </div>-->
        <!--                                    <div class="progress-group">-->
        <!--                                        <div class="progress-group-header">-->
        <!--                                            <CIcon name="cil-globe-alt" class="progress-group-icon"/>-->
        <!--                                            <span class="title">Organic Search</span>-->
        <!--                                            <span class="ml-auto font-weight-bold">-->
        <!--                        191,235 <span class="text-muted small">(56%)</span>-->
        <!--                      </span>-->
        <!--                                        </div>-->
        <!--                                        <div class="progress-group-bars">-->
        <!--                                            <CProgress-->
        <!--                                                    class="progress-xs"-->
        <!--                                                    :value="56"-->
        <!--                                                    color="success"-->
        <!--                                            />-->
        <!--                                        </div>-->
        <!--                                    </div>-->
        <!--                                    <div class="progress-group">-->
        <!--                                        <div class="progress-group-header">-->
        <!--                                            <CIcon-->
        <!--                                                    name="cib-facebook"-->
        <!--                                                    height="17"-->
        <!--                                                    class="progress-group-icon"-->
        <!--                                            />-->
        <!--                                            <span class="title">Facebook</span>-->
        <!--                                            <span class="ml-auto font-weight-bold">-->
        <!--                        51,223 <span class="text-muted small">(15%)</span>-->
        <!--                      </span>-->
        <!--                                        </div>-->
        <!--                                        <div class="progress-group-bars">-->
        <!--                                            <CProgress-->
        <!--                                                    class="progress-xs"-->
        <!--                                                    :value="15"-->
        <!--                                                    color="success"-->
        <!--                                            />-->
        <!--                                        </div>-->
        <!--                                    </div>-->
        <!--                                    <div class="progress-group">-->
        <!--                                        <div class="progress-group-header">-->
        <!--                                            <CIcon-->
        <!--                                                    name="cib-twitter"-->
        <!--                                                    height="17"-->
        <!--                                                    class="progress-group-icon"-->
        <!--                                            />-->
        <!--                                            <span class="title">Twitter</span>-->
        <!--                                            <span class="ml-auto font-weight-bold">-->
        <!--                        37,564 <span class="text-muted small">(11%)</span>-->
        <!--                      </span>-->
        <!--                                        </div>-->
        <!--                                        <div class="progress-group-bars">-->
        <!--                                            <CProgress-->
        <!--                                                    class="progress-xs"-->
        <!--                                                    :value="11"-->
        <!--                                                    color="success"-->
        <!--                                            />-->
        <!--                                        </div>-->
        <!--                                    </div>-->
        <!--                                    <div class="progress-group">-->
        <!--                                        <div class="progress-group-header">-->
        <!--                                            <CIcon-->
        <!--                                                    name="cib-linkedin"-->
        <!--                                                    height="17"-->
        <!--                                                    class="progress-group-icon"-->
        <!--                                            />-->
        <!--                                            <span class="title">LinkedIn</span>-->
        <!--                                            <span class="ml-auto font-weight-bold">-->
        <!--                        27,319 <span class="text-muted small">&nbsp;(8%)</span>-->
        <!--                      </span>-->
        <!--                                        </div>-->
        <!--                                        <div class="progress-group-bars">-->
        <!--                                            <CProgress-->
        <!--                                                    class="progress-xs"-->
        <!--                                                    :value="8"-->
        <!--                                                    color="success"-->
        <!--                                            />-->
        <!--                                        </div>-->
        <!--                                    </div>-->
        <!--                                    <div class="divider text-center">-->
        <!--                                        <CButton color="link" size="sm" class="text-muted">-->
        <!--                                            <CIcon name="cil-options"/>-->
        <!--                                        </CButton>-->
        <!--                                    </div>-->
        <!--                                </ul>-->
        <!--                            </CCol>-->
        <!--                        </CRow>-->
        <!--                        <br/>-->
        <!--                        <CDataTable-->
        <!--                                class="mb-0 table-outline"-->
        <!--                                hover-->
        <!--                                :items="tableItems"-->
        <!--                                :fields="tableFields"-->
        <!--                                head-color="light"-->
        <!--                                no-sorting-->
        <!--                        >-->
        <!--                            <td slot="avatar" class="text-center" slot-scope="{item}">-->
        <!--                                <div class="c-avatar">-->
        <!--                                    <img :src="item.avatar.url" class="c-avatar-img" alt="">-->
        <!--                                    <span-->
        <!--                                            class="c-avatar-status"-->
        <!--                                            :class="`bg-${item.avatar.status || 'secondary'}`"-->
        <!--                                    ></span>-->
        <!--                                </div>-->
        <!--                            </td>-->
        <!--                            <td slot="user" slot-scope="{item}">-->
        <!--                                <div>{{item.user.name}}</div>-->
        <!--                                <div class="small text-muted">-->
        <!--                  <span>-->
        <!--                    <template v-if="item.user.new">New</template>-->
        <!--                    <template v-else>Recurring</template>-->
        <!--                  </span> | Registered: {{item.user.registered}}-->
        <!--                                </div>-->
        <!--                            </td>-->
        <!--                            <td-->
        <!--                                    slot="country"-->
        <!--                                    slot-scope="{item}"-->
        <!--                                    class="text-center"-->
        <!--                            >-->
        <!--                                <CIcon-->
        <!--                                        :name="item.country.flag"-->
        <!--                                        height="25"-->
        <!--                                />-->
        <!--                            </td>-->
        <!--                            <td slot="usage" slot-scope="{item}">-->
        <!--                                <div class="clearfix">-->
        <!--                                    <div class="float-left">-->
        <!--                                        <strong>{{item.usage.value}}%</strong>-->
        <!--                                    </div>-->
        <!--                                    <div class="float-right">-->
        <!--                                        <small class="text-muted">{{item.usage.period}}</small>-->
        <!--                                    </div>-->
        <!--                                </div>-->
        <!--                                <CProgress-->
        <!--                                        class="progress-xs"-->
        <!--                                        v-model="item.usage.value"-->
        <!--                                        :color="color(item.usage.value)"-->
        <!--                                />-->
        <!--                            </td>-->
        <!--                            <td-->
        <!--                                    slot="payment"-->
        <!--                                    slot-scope="{item}"-->
        <!--                                    class="text-center"-->
        <!--                            >-->
        <!--                                <CIcon-->
        <!--                                        :name="item.payment.icon"-->
        <!--                                        height="25"-->
        <!--                                />-->
        <!--                            </td>-->
        <!--                            <td slot="activity" slot-scope="{item}">-->
        <!--                                <div class="small text-muted">Last login</div>-->
        <!--                                <strong>{{item.activity}}</strong>-->
        <!--                            </td>-->
        <!--                        </CDataTable>-->
        <!--                    </CCardBody>-->
        <!--                </CCard>-->
        <!--            </CCol>-->
        <!--        </CRow>-->
    </div>
</template>

<script>
    import MainChartExample from './charts/MainChartExample'
    import WidgetsDropdown from './widgets/WidgetsDropdown'
    import WidgetsBrand from './widgets/WidgetsBrand'
    import {CChartLineSimple, CChartBarSimple} from './charts/index.js'
    import clientUtils from "@/utils/api_client"
    import usersData from '@/views/users/UsersData'
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
                    console.log(apiRes)
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
            return {
                isCollapsed: true,
                loremIpsum: 'Lorem ipsum dolor sit amet, consectetuer adipiscing elit, sed diam nonummy nibh euismod tincidunt ut laoreet dolore magna aliquam erat volutpat. Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat.',
                systemInfo: systemInfo,
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
            color(value) {
                let $color
                if (value <= 25) {
                    $color = 'info'
                } else if (value > 25 && value <= 50) {
                    $color = 'success'
                } else if (value > 50 && value <= 75) {
                    $color = 'warning'
                } else if (value > 75 && value <= 100) {
                    $color = 'danger'
                }
                return $color
            },

            shuffleArray(array) {
                for (let i = array.length - 1; i > 0; i--) {
                    let j = Math.floor(Math.random() * (i + 1))
                    let temp = array[i]
                    array[i] = array[j]
                    array[j] = temp
                }
                return array
            },

            getShuffledUsersData() {
                return this.shuffleArray(usersData.slice(0))
            }
        }
    }
</script>
