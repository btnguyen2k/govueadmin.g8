<template>
  <div>
    <CRow>
      <CCol sm="12">
        <CCard>
          <CCardHeader><h5>Delete Blog Post</h5></CCardHeader>
          <CForm @submit.prevent="doSubmit" method="post">
            <CCardBody>
              <p v-if="!found" class="alert alert-danger">Blog Post [{{ this.$route.params.id }}] not found</p>
              <p v-if="errorMsg!=''" class="alert alert-danger">{{ errorMsg }}</p>
              <div class="form-group form-row" v-if="found">
                <CCol :sm="{offset:3,size:9}" class="form-inline">
                  <CInputCheckbox inline label="Public" :checked.sync="post.is_public" disabled="disabled"/>
                  <small>(Other people can see, comment and vote your public posts)</small>
                </CCol>
              </div>
              <CInput v-if="found"
                      type="text"
                      v-model="post.title"
                      label="Title"
                      placeholder="My blog post's awesome title"
                      horizontal
                      readonly="readonly"
              />
              <CTextarea v-if="found"
                         rows="10"
                         type="text"
                         v-model="post.content"
                         label="Content (Markdown supported)"
                         placeholder="My blog post's awesome content"
                         horizontal
                         readonly="readonly"
              />
            </CCardBody>
            <CCardFooter>
              <CButton v-if="found" type="submit" color="danger" style="width: 96px">
                <CIcon name="cil-trash"/>
                Delete
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
import utils from "@/utils/app_utils"
import marked from "marked"
import DOMPurify from "dompurify"

export default {
  name: 'DeletePost',
  computed: {
    previewContent() {
      const html = marked(this.post.content)
      return DOMPurify.sanitize(html)
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
      clientUtils.apiDoDelete(
          clientUtils.apiPost + "/" + this.$route.params.id,
          (apiRes) => {
            if (apiRes.status != 200) {
              this.errorMsg = apiRes.status + ": " + apiRes.message
            } else {
              utils.localStorageSet(utils.lskeyLoginSessionLastCheck, null)
              this.$router.push({
                name: "MyBlog",
                params: {flashMsg: "Blog post [" + this.post.title + "] has been deleted successfully."},
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
