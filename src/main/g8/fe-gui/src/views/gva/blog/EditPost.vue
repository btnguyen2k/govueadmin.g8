<template>
  <div>
    <CRow>
      <CCol sm="12">
        <CCard>
          <CCardHeader><h5>{{ $t('message.edit_blog_post') }}</h5></CCardHeader>
          <CForm @submit.prevent="doSubmit" method="post">
            <CCardBody>
              <p v-if="!found" class="alert alert-danger">Blog Post [{{ this.$route.params.id }}] not found</p>
              <p v-if="errorMsg!=''" class="alert alert-danger">{{ errorMsg }}</p>
              <div class="form-group form-row" v-if="found">
                <CCol :sm="{offset:3,size:9}" class="form-inline">
                  <CInputCheckbox inline :label="$t('message.blog_public')" :checked.sync="post.is_public"/>
                  <small>({{ $t('message.blog_public_msg') }})</small>
                </CCol>
              </div>
              <CInput v-if="found"
                      type="text"
                      v-model="post.title"
                      :label="$t('message.blog_title')"
                      :placeholder="$t('message.blog_title_msg')"
                      horizontal
                      required
                      was-validated
              />
              <CTabs v-if="found">
                <CTab active>
                  <template slot="title">
                    <CIcon name="cib-markdown"/>
                    {{ $t('message.blog_editor') }}
                  </template>
                  <CTextarea
                      rows="10"
                      type="text"
                      v-model="post.content"
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
              <CButton v-if="found" type="submit" color="primary" style="width: 96px">
                <CIcon name="cil-save"/>
                {{ $t('message.action_save') }}
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
  name: 'EditPost',
  computed: {
    previewContent() {
      return this.found ? DOMPurify.sanitize(marked(this.post.content), {ADD_ATTR: ['target']}) : ''
    }
  },
  data() {
    clientUtils.apiDoGet(clientUtils.apiPost + "/" + this.$route.params.id,
        (apiRes) => {
          this.found = apiRes.status == 200
          if (this.found) {
            this.post = apiRes.data
          }
        },
        (err) => {
          this.errorMsg = err
        })
    return {
      post: {},
      errorMsg: "",
      found: false,
    }
  },
  methods: {
    doCancel() {
      router.push({name: "MyBlog"})
    },
    doSubmit(e) {
      e.preventDefault()
      let data = {is_public: this.post.is_public, title: this.post.title, content: this.post.content}
      clientUtils.apiDoPut(
          clientUtils.apiPost + "/" + this.$route.params.id, data,
          (apiRes) => {
            if (apiRes.status != 200) {
              this.errorMsg = apiRes.status + ": " + apiRes.message
            } else {
              this.$router.push({
                name: "MyBlog",
                params: {flashMsg: this.$i18n.t('message.blog_updated_msg', {title: this.post.title})},
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
