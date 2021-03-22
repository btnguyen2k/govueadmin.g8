<template>
  <CRow>
    <CCol sm="12">
      <CCard accent-color="info">
        <CCardHeader>
          <strong>{{ $tc('message.blog_posts', blogPostList.data.length, {count: blogPostList.data.length}) }}</strong>
          <div class="card-header-actions">
            <CButton class="btn-sm btn-primary" @click="clickCreateBlogPost">
              <CIcon name="cil-image-plus"/>
              {{ $t('message.create_blog_post') }}
            </CButton>
          </div>
        </CCardHeader>
        <CCardBody>
          <p v-if="flashMsg" class="alert alert-success">{{ flashMsg }}</p>
          <CDataTable :items="blogPostList.data" :fields="[
              {key:'public',label:''},
              {key:'created',label:$t('message.blog_tcreated')},
              {key:'title',label:$t('message.blog_title')},
              {key:'num_comments',label:$t('message.blog_comments')},
              {key:'num_votes_up',label:$t('message.blog_votes')+' ↑',_style:'white-space: nowrap'},
              {key:'num_votes_down',label:$t('message.blog_votes')+' ↓',_style:'white-space: nowrap'},
              {key:'actions',label:$t('message.actions'),_style:'text-align: center'}
            ]">
            <template #public="{item}">
              <td>
                <CIcon :name="`${item.is_public?'cil-check':'cil-check-alt'}`"
                       :style="`color: ${item.is_public?'green':'grey'}`"/>
              </td>
            </template>
            <template #created="{item}">
              <td style="font-size: smaller; white-space: nowrap">{{item.t_created.substring(0,19)}} (GMT{{item.t_created.substring(26)}})</td>
            </template>
            <template #title="{item}">
              <td style="font-size: smaller">{{item.title}}</td>
            </template>
            <template #num_comments="{item}">
              <td style="font-size: smaller; text-align: center">{{item.num_comments}}</td>
            </template>
            <template #num_votes_up="{item}">
              <td style="font-size: smaller; text-align: center">{{item.num_votes_up}}</td>
            </template>
            <template #num_votes_down="{item}">
              <td style="font-size: smaller; text-align: center">{{item.num_votes_down}}</td>
            </template>
            <template #actions="{item}">
              <td style="font-size: smaller; white-space: nowrap; text-align: center">
                <CLink @click="clickEditBlogPost(item.id)" :label="$t('message.action_edit')" class="btn btn-sm btn-primary">
                  <CIcon name="cil-pencil"/>
                </CLink>
                &nbsp;
                <CLink @click="clickDeleteBlogPost(item.id)" :label="$t('message.action_delete')" class="btn btn-sm btn-danger">
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
