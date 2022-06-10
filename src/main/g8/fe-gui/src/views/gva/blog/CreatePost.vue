<template>
  <div>
    <CRow>
      <CCol sm="12">
        <CCard>
          <CCardHeader
            ><CCardTitle>{{
              $t('message.create_blog_post')
            }}</CCardTitle></CCardHeader
          >
          <CForm @submit.prevent="doSubmit" method="post">
            <CCardBody>
              <p v-if="errorMsg != ''" class="alert alert-danger">
                {{ errorMsg }}
              </p>
              <CRow class="mb-2">
                <CCol :sm="{ offset: 3, size: 9 }">
                  <CFormCheck
                    inline
                    :label="$t('message.blog_public')"
                    v-model:checked="form.isPublic"
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
                    v-model="form.title"
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
                        v-model="form.content"
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
            </CCardBody>
            <CCardFooter>
              <CButton type="submit" color="primary" style="width: 96px">
                <CIcon name="cil-save" />
                {{ $t('message.action_create') }}
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
  name: 'CreatePost',
  computed: {
    previewContent() {
      const html = marked(this.form.content)
      return DOMPurify.sanitize(html, { ADD_ATTR: ['target'] })
    },
  },
  data() {
    return {
      form: { title: '', content: '', isPublic: false },
      errorMsg: '',
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
        is_public: this.form.isPublic,
        title: this.form.title,
        content: this.form.content,
      }
      clientUtils.apiDoPost(
        clientUtils.apiMyBlog,
        data,
        (apiRes) => {
          if (apiRes.status != 200) {
            this.errorMsg = apiRes.status + ': ' + apiRes.message
          } else {
            this.$router.push({
              name: 'MyBlog',
              params: {
                flashMsg: this.$i18n.t('message.blog_created_msg', {
                  title: this.form.title,
                }),
              },
            })
          }
        },
        (err) => {
          console.error(err)
          this.errorMsg = err
        },
      )
    },
  },
}
</script>
