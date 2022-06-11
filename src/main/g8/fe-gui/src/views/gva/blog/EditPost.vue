<template>
  <div>
    <CRow>
      <CCol sm="12">
        <CCard>
          <CCardHeader
            ><CCardTitle>{{
              $t('message.edit_blog_post')
            }}</CCardTitle></CCardHeader
          >
          <CForm @submit.prevent="doSubmit" method="post">
            <CCardBody>
              <CAlert color="danger" v-if="!found">{{
                $t('message.error_blog_post_not_found', {
                  id: this.$route.params.id,
                })
              }}</CAlert>
              <CAlert color="danger" v-if="errorMsg != ''">{{
                errorMsg
              }}</CAlert>
              <div v-if="found">
                <CRow class="mb-2">
                  <CCol :sm="{ offset: 3, size: 9 }">
                    <CFormCheck
                      inline
                      :label="$t('message.blog_public')"
                      v-model="post.is_public"
                    />
                    <small>({{ $t('message.blog_public_msg') }})</small>
                  </CCol>
                </CRow>
                <CRow class="mb-2">
                  <CFormLabel class="col-sm-3 col-form-label">{{
                    $t('message.blog_title')
                  }}</CFormLabel>
                  <div class="col-sm-9">
                    <CFormInput
                      type="text"
                      v-model="post.title"
                      :placeholder="$t('message.blog_title_msg')"
                      required
                    />
                  </div>
                </CRow>
                <CNav variant="tabs">
                  <CNavItem active>
                    <CNavLink
                      href="javascript:void(0);"
                      :active="tabPaneActiveKey === 1"
                      @click="switchTab(1)"
                    >
                      <CIcon name="cib-markdown" />&nbsp;{{
                        $t('message.blog_editor')
                      }}
                    </CNavLink>
                  </CNavItem>
                  <CNavItem active>
                    <CNavLink
                      href="javascript:void(0);"
                      :active="tabPaneActiveKey === 2"
                      @click="switchTab(2)"
                    >
                      <CIcon name="cil-calculator" />&nbsp;{{
                        $t('message.blog_preview')
                      }}
                    </CNavLink>
                  </CNavItem>
                </CNav>
                <CTabContent>
                  <CTabPane
                    role="tabpanel"
                    aria-labelledby="editor-tab"
                    :visible="tabPaneActiveKey === 1"
                  >
                    <CRow class="mt-2">
                      <div class="col-sm-12">
                        <CFormTextarea
                          rows="10"
                          type="text"
                          v-model="post.content"
                          :placeholder="$t('message.blog_content_msg')"
                          required
                        />
                      </div>
                    </CRow>
                  </CTabPane>
                  <CTabPane
                    role="tabpanel"
                    aria-labelledby="preview-tab"
                    :visible="tabPaneActiveKey === 2"
                  >
                    <CRow class="mt-2">
                      <div class="col-sm-12" v-html="previewContent" />
                    </CRow>
                  </CTabPane>
                </CTabContent>
              </div>
            </CCardBody>
            <CCardFooter v-if="found">
              <CButton type="submit" color="primary" style="width: 96px">
                <CIcon name="cil-save" />
                {{ $t('message.action_save') }}
              </CButton>
              <CButton
                type="button"
                color="info"
                class="ms-2"
                style="width: 96px"
                @click="doCancel"
              >
                <CIcon name="cil-arrow-circle-left" />
                {{ $t('message.action_back') }}
              </CButton>
            </CCardFooter>
          </CForm>
        </CCard>
      </CCol>
    </CRow>
  </div>
</template>

<script>
import router from '@/router'
import clientUtils from '@/utils/api_client'
import marked from 'marked'
import DOMPurify from 'dompurify'

export default {
  name: 'EditPost',
  computed: {
    previewContent() {
      return this.found
        ? DOMPurify.sanitize(marked(this.post.content), {
            ADD_ATTR: ['target'],
          })
        : ''
    },
  },
  data() {
    clientUtils.apiDoGet(
      clientUtils.apiPost + '/' + this.$route.params.id,
      (apiRes) => {
        this.found = apiRes.status == 200
        if (this.found) {
          this.post = apiRes.data
        }
      },
      (err) => {
        this.errorMsg = err
      },
    )
    return {
      post: {},
      errorMsg: '',
      found: false,
      tabPaneActiveKey: 1,
    }
  },
  methods: {
    switchTab(tabId) {
      this.tabPaneActiveKey = tabId
    },
    doCancel() {
      router.push({ name: 'MyBlog' })
    },
    doSubmit(e) {
      e.preventDefault()
      let data = {
        is_public: this.post.is_public,
        title: this.post.title,
        content: this.post.content,
      }
      console.log(this.post)
      clientUtils.apiDoPut(
        clientUtils.apiPost + '/' + this.$route.params.id,
        data,
        (apiRes) => {
          if (apiRes.status != 200) {
            this.errorMsg = apiRes.status + ': ' + apiRes.message
          } else {
            this.$router.push({
              name: 'MyBlog',
              params: {
                flashMsg: this.$i18n.t('message.blog_updated_msg', {
                  title: this.post.title,
                }),
              },
            })
          }
        },
        (err) => {
          this.errorMsg = err
        },
      )
    },
  },
}
</script>
