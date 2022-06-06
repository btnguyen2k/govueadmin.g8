<template>
  <div>
    <CRow v-for="postId in blogPostIdList.data">
      <CCol sm="12">
        <CCard accent-color="primary">
          <CCardHeader>
            <div class="c-avatar" style="vertical-align: middle">
              <img :src="avatar(blogPostMap[postId])" class="c-avatar-img" :title="displayName(blogPostMap[postId])"/>
            </div>
            <span style="padding-left: 8px; font-size: large; font-weight: bold">{{ blogPostMap[postId].title }}</span>
          </CCardHeader>
          <CCardBody>
            <p style="font-style: italic">by <strong>{{ displayName(blogPostMap[postId]) }}</strong> on
              {{ creationTime(blogPostMap[postId]) }}</p>
            <div v-html="renderMarkdown(blogPostMap[postId])"></div>
          </CCardBody>
          <CCardFooter>
            <div class="float-right">
              <CButton @click="doVote(postId,1)">
                <CIcon :height="18" name="cil-arrow-thick-top"
                       :style="`color: ${voteValue(blogPostMap[postId])>0?'blue':'grey'}`"/>
                {{ blogPostMap[postId].num_votes_up }}
              </CButton>
              <span style="padding-left: 16px"></span>
              <CButton @click="doVote(postId,-1)">
                <CIcon :height="18" name="cil-arrow-thick-bottom"
                       :style="`color: ${voteValue(blogPostMap[postId])<0?'red':'grey'}`"/>
                {{ blogPostMap[postId].num_votes_down }}
              </CButton>
            </div>
          </CCardFooter>
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
  mounted() {
    const vue = this
    let blogPostIdList = []
    clientUtils.apiDoGet(clientUtils.apiMyFeed,
        (apiRes) => {
          if (apiRes.status == 200) {
            apiRes.data.forEach(post => {
              vue.$set(vue.blogPostMap, post.id, post)
              blogPostIdList.push(post.id)
            })
            vue.blogPostIdList.data = blogPostIdList

            apiRes.data.forEach((post) => {
              clientUtils.apiDoGet(clientUtils.apiUserVoteForPost + "/" + post.id,
                  (apiRes) => {
                    if (apiRes.status == 200) {
                      vue.$set(vue.blogPostVotes, post.id, apiRes.data)
                    }
                  },
                  (err) => {
                    console.error("Error getting user vote for post: " + err)
                  })
            })
          } else {
            console.error("Getting user vote for post was unsuccessful: " + apiRes)
          }
        },
        (err) => {
          console.error("Error getting user feed: " + err)
        })
  },
  data() {
    return {
      blogPostVotes: {},
      blogPostIdList: {data: []},
      blogPostMap: {},
    }
  },
  methods: {
    voteValue(post) {
      return this.blogPostVotes[post.id]
    },
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
      return DOMPurify.sanitize(marked(post.content), {ADD_ATTR: ['target']})
    },
    doVote(postId, v) {
      const data = {vote: v}
      const vue = this
      clientUtils.apiDoPost(clientUtils.apiUserVoteForPost + "/" + postId, data,
          (apiRes) => {
            // console.log(apiRes)
            if (apiRes.status == 200 && apiRes.data.vote) {
              vue.blogPostVotes[postId] = apiRes.data.value
              vue.blogPostMap[postId].num_votes_up = apiRes.data.num_votes_up
              vue.blogPostMap[postId].num_votes_down = apiRes.data.num_votes_down
            }
          },
          (err) => {
            console.error("Error voting for post: " + err)
          })
    }
  }
}
</script>
