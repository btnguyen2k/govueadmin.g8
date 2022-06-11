<template>
  <div>
    <CRow>
      <CCol sm="12">
        <CCard>
          <CCardHeader
            ><CCardTitle>{{
              $t('message.delete_blog_post')
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
                  <CCol>
                    <CFormCheck
                      inline
                      disabled
                      :label="$t('message.blog_public')"
                      v-model="post.is_public"
                    />
                    <small>({{ $t('message.blog_public_msg') }})</small>
                  </CCol>
                </CRow>
                <CRow class="mb-2">
                  <CCol>
                    <CFormInput
                      type="text"
                      v-model="post.title"
                      :placeholder="$t('message.blog_title_msg')"
                      disabled
                    />
                  </CCol>
                </CRow>
                <CRow class="mb-2">
                  <div class="col-sm-12">
                    <CFormTextarea
                      rows="10"
                      type="text"
                      v-model="post.content"
                      :placeholder="$t('message.blog_content_msg')"
                      disabled
                    />
                  </div>
                </CRow>
              </div>
            </CCardBody>
            <CCardFooter v-if="found">
              <CButton type="submit" color="danger" style="width: 96px">
                <CIcon name="cil-trash" />
                {{ $t('message.action_delete') }}
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
import utils from '@/utils/app_utils'
import marked from 'marked'
import DOMPurify from 'dompurify'

export default {
  name: 'DeletePost',
  computed: {
    previewContent() {
      const html = marked(this.post.content)
      return DOMPurify.sanitize(html)
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
    }
  },
  methods: {
    doCancel() {
      router.push({ name: 'MyBlog' })
    },
    doSubmit(e) {
      e.preventDefault()
      clientUtils.apiDoDelete(
        clientUtils.apiPost + '/' + this.$route.params.id,
        (apiRes) => {
          if (apiRes.status != 200) {
            this.errorMsg = apiRes.status + ': ' + apiRes.message
          } else {
            utils.localStorageSet(utils.lskeyLoginSessionLastCheck, null)
            this.$router.push({
              name: 'MyBlog',
              params: {
                flashMsg: this.$i18n.t('message.blog_deleted_msg', {
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
