<template>
  <div>
    <CRow v-for="post in blogPostList.data">
      <CCol sm="12">
        <CCard accent-color="primary">
          <CCardHeader>
            <div class="c-avatar" style="vertical-align: middle">
              <img :src="avatar(post)" class="c-avatar-img" :title="displayName(post)"/>
            </div>
            &nbsp;&nbsp;&nbsp;&nbsp;<span style="font-size: large; font-weight: bold">{{post.title}}</span>
          </CCardHeader>
          <CCardBody>
            <p style="font-style: italic">by <strong>{{displayName(post)}}</strong> on {{creationTime(post)}})</p>
            <div v-html="renderMarkdown(post)"></div>
          </CCardBody>
        </CCard>
      </CCol>
    </CRow>
  </div>
</template>

<script>
import clientUtils from "@/utils/api_client"
import marked from "marked"
import DOMPurify from "dompurify"

export default {
  name: 'Dashboard',
  data() {
    let blogPostList = {data: []}
    clientUtils.apiDoGet(clientUtils.apiMyFeed,
        (apiRes) => {
          if (apiRes.status == 200) {
            blogPostList.data = apiRes.data
          } else {
            console.error("Getting user feed was unsuccessful: " + apiRes)
          }
        },
        (err) => {
          console.error("Error getting user feed: " + err)
        })
    return {
      blogPostList: blogPostList,
    }
  },
  methods: {
    avatar(post) {
      return "https://www.gravatar.com/avatar/" + post.owner.id.trim().toLowerCase().md5() + "?s=40"
    },
    displayName(post) {
      return post.owner.display_name
    },
    creationTime(post) {
      return post.t_created.substring(0, 19) + ' (GMT' + post.t_created.substring(26) + ')'
    },
    renderMarkdown(post) {
      // return marked(post.content)
      return DOMPurify.sanitize(marked(post.content), {ADD_ATTR:['target']})
    }
  }
}
</script>
