<template>
  <CRow>
    <CCol sm="12">
      <CCard accent-color="info">
        <CCardHeader>
          <strong>Blog Post ({{ blogPostList.data.length }})</strong>
          <div class="card-header-actions">
            <CButton class="btn-sm btn-primary" @click="clickCreateBlogPost">
              <CIcon name="cil-image-plus"/>
              Create Blog Post
            </CButton>
          </div>
        </CCardHeader>
        <CCardBody>
          <p v-if="flashMsg" class="alert alert-success">{{ flashMsg }}</p>
          <CDataTable :items="blogPostList.data" :fields="[
              {key:'public',label:''},
              'Created',
              'title',
              {key:'num_comments',label:'Comments'},
              {key:'num_votes_up',label:'Votes ↑'},
              {key:'num_votes_down',label:'Votes ↓'},
              'actions']">
            <template #public="{item}">
              <td>
                <CIcon :name="`${item.is_public?'cil-check':'cil-check-alt'}`"
                       :style="`color: ${item.is_public?'green':'grey'}`"/>
              </td>
            </template>
            <template #Created="{item}">
              <td>{{item.t_created.substring(0,19)}} (GMT{{item.t_created.substring(26)}})</td>
            </template>
            <template #actions="{item}">
              <td>
                <CLink @click="clickEditBlogPost(item.id)" label="Edit" class="btn btn-primary">
                  <CIcon name="cil-pencil"/>
                </CLink>
                &nbsp;
                <CLink @click="clickDeleteBlogPost(item.id)" label="Delete" class="btn btn-danger">
                  <CIcon name="cil-trash"/>
                </CLink>
              </td>
            </template>
          </CDataTable>
        </CCardBody>
      </CCard>
    </CCol>
  </CRow>
</template>

<script>
import clientUtils from "@/utils/api_client";

export default {
  name: 'MyBlog',
  data: () => {
    let blogPostList = {data: []}
    clientUtils.apiDoGet(clientUtils.apiMyBlog,
        (apiRes) => {
          if (apiRes.status == 200) {
            blogPostList.data = apiRes.data
            blogPostList.data.forEach((item) => {
              console.log(item)
            })
          } else {
            console.error("Getting blog post list was unsuccessful: " + apiRes)
          }
        },
        (err) => {
          console.error("Error getting blog post list: " + err)
        })
    return {
      blogPostList: blogPostList,
    }
  },
  props: ["flashMsg"],
  methods: {
    clickCreateBlogPost() {
      this.$router.push({name: "CreatePost"})
    },
    clickEditBlogPost(id) {
      this.$router.push({name: "EditPost", params: {id: id.toString()}})
    },
    clickDeleteBlogPost(id) {
      this.$router.push({name: "DeletePost", params: {id: id.toString()}})
    },
  }
}
</script>
