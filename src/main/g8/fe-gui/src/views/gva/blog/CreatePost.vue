<template>
  <div>
    <CRow>
      <CCol sm="12">
        <CCard>
          <CCardHeader><h5>Create Blog Post</h5></CCardHeader>
          <CForm @submit.prevent="doSubmit" method="post">
            <CCardBody>
              <p v-if="errorMsg!=''" class="alert alert-danger">{{ errorMsg }}</p>
              <div class="form-group form-row">
                <CCol :sm="{offset:3,size:9}" class="form-inline">
                  <CInputCheckbox inline label="Public" :checked.sync="form.isPublic"/>
                </CCol>
              </div>
              <CInput
                  type="text"
                  v-model="form.title"
                  label="Title"
                  placeholder="My blog post's awesome title"
                  horizontal
                  required
                  was-validated
              />
              <CTabs>
                <CTab active>
                  <template slot="title">
                    <CIcon name="cib-markdown"/> Editor
                  </template>
                  <CTextarea
                      rows="10"
                      type="text"
                      v-model="form.content"
                      label="Content (Markdown supported)"
                      placeholder="My blog post's awesome content"
                      horizontal
                      required
                      was-validated
                  />
                </CTab>
                <CTab>
                  <template slot="title">
                    <CIcon name="cil-calculator"/> Preview
                  </template>
                  <div v-html="previewContent"></div>
                </CTab>
              </CTabs>
            </CCardBody>
            <CCardFooter>
              <CButton type="submit" color="primary" style="width: 96px">
                <CIcon name="cil-save"/>
                Create
              </CButton>
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
import clientUtils from "@/utils/api_client"
import marked from "marked"
import DOMPurify from 'dompurify'

export default {
  name: 'CreatePost',
  computed: {
    previewContent() {
      const html = marked(this.form.content)
      return DOMPurify.sanitize(html)
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
                params: {flashMsg: "Blog post [" + this.form.title + "] has been created successfully."},
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
