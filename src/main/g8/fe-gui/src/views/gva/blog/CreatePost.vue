<template>
  <div>
    <CRow>
      <CCol sm="12">
        <CCard>
          <CCardHeader><h5>{{ $t('message.create_blog_post') }}</h5></CCardHeader>
          <CForm @submit.prevent="doSubmit" method="post">
            <CCardBody>
              <p v-if="errorMsg!=''" class="alert alert-danger">{{ errorMsg }}</p>
              <div class="form-group form-row">
                <CCol :sm="{offset:3,size:9}" class="form-inline">
                  <CInputCheckbox inline :label="$t('message.blog_public')" :checked.sync="form.isPublic"/>
                  <small>({{ $t('message.blog_public_msg') }})</small>
                </CCol>
              </div>
              <CInput
                  type="text"
                  v-model="form.title"
                  :label="$t('message.blog_title')"
                  :placeholder="$t('message.blog_title_msg')"
                  horizontal
                  required
                  was-validated
              />
              <CTabs>
                <CTab active>
                  <template slot="title">
                    <CIcon name="cib-markdown"/>
                    {{ $t('message.blog_editor') }}
                  </template>
                  <CTextarea
                      rows="10"
                      type="text"
                      v-model="form.content"
                      :label="$t('message.blog_content')"
                      :placeholder="$t('message.blog_content_msg')"
                      horizontal
                      required
                      was-validated
                  />
                </CTab>
                <CTab>
                  <template slot="title">
                    <CIcon name="cil-calculator"/>
                    {{ $t('message.blog_preview') }}
                  </template>
                  <div v-html="previewContent"></div>
                </CTab>
              </CTabs>
            </CCardBody>
            <CCardFooter>
              <CButton type="submit" color="primary" style="width: 96px">
                <CIcon name="cil-save"/>
                {{ $t('message.action_create') }}
              </CButton>
              <CButton type="button" color="info" class="ml-2" style="width: 96px" @click="doCancel">
                <CIcon name="cil-arrow-circle-left"/>
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
import router from "@/router"
import clientUtils from "@/utils/api_client"
import marked from "marked"
import DOMPurify from "dompurify"

export default {
  name: 'CreatePost',
  computed: {
    previewContent() {
      const html = marked(this.form.content)
      return DOMPurify.sanitize(html, {ADD_ATTR: ['target']})
    }
  },
  data() {
    return {
      form: {title: "", content: "", isPublic: false},
      errorMsg: "",
    }
  },
  methods: {
    doCancel() {
      router.push({name: "MyBlog"})
    },
    doSubmit(e) {
      e.preventDefault()
      let data = {is_public: this.form.isPublic, title: this.form.title, content: this.form.content}
      clientUtils.apiDoPost(
          clientUtils.apiMyBlog, data,
          (apiRes) => {
            if (apiRes.status != 200) {
              this.errorMsg = apiRes.status + ": " + apiRes.message
            } else {
              this.$router.push({
                name: "MyBlog",
                params: {flashMsg: this.$i18n.t('message.blog_created_msg', {title: this.form.title})},
              })
            }
          },
          (err) => {
            console.error(err)
            this.errorMsg = err
          }
      )
    },
  }
}
</script>
