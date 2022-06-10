<template>
  <CRow>
    <CCol sm="12">
      <CCard accent-color="info">
        <CCardHeader>
          <CCardTitle class="float-start">{{
            $tc('message.blog_posts', blogPostList.length, {
              count: blogPostList.length,
            })
          }}</CCardTitle>
          <div class="float-end">
            <CButton class="btn-sm btn-primary" @click="clickCreateBlogPost">
              <CIcon name="cil-image-plus" />
              {{ $t('message.create_blog_post') }}
            </CButton>
          </div>
        </CCardHeader>
        <CCardBody>
          <CAlert color="success" v-if="flashMsg">{{ flashMsg }}</CAlert>
          <CTable>
            <CTableHead>
              <CTableRow>
                <CTableHeaderCell scope="col">&nbsp;</CTableHeaderCell>
                <CTableHeaderCell scope="col">{{
                  $t('message.blog_tcreated')
                }}</CTableHeaderCell>
                <CTableHeaderCell scope="col">{{
                  $t('message.blog_title')
                }}</CTableHeaderCell>
                <CTableHeaderCell scope="col">{{
                  $t('message.blog_comments')
                }}</CTableHeaderCell>
                <CTableHeaderCell scope="col">{{
                  $t('message.blog_votes') + '↑'
                }}</CTableHeaderCell>
                <CTableHeaderCell scope="col">{{
                  $t('message.blog_votes') + '↓'
                }}</CTableHeaderCell>
                <CTableHeaderCell scope="col" style="text-align: center">{{
                  $t('message.actions')
                }}</CTableHeaderCell>
              </CTableRow>
            </CTableHead>
            <CTableBody>
              <template v-for="blog in blogPostList" :key="blog.id">
                <CTableRow>
                  <CTableDataCell
                    style="font-size: smaller; white-space: nowrap"
                    ><CIcon
                      :name="`${
                        blog.is_public ? 'cil-check' : 'cil-check-alt'
                      }`"
                      :style="`color: ${blog.is_public ? 'green' : 'grey'}`"
                  /></CTableDataCell>
                  <CTableDataCell
                    style="font-size: smaller; white-space: nowrap"
                  >
                    {{
                      '{timestamp} (GMT{tz})'
                        .replace('{timestamp}', blog.t_created.substring(0, 19))
                        .replace('{tz}', blog.t_created.substring(26))
                    }}
                  </CTableDataCell>
                  <CTableDataCell style="font-size: smaller">{{
                    blog.title
                  }}</CTableDataCell>
                  <CTableDataCell
                    style="font-size: smaller; text-align: center"
                  >
                    {{ blog.num_comments }}
                  </CTableDataCell>
                  <CTableDataCell
                    style="font-size: smaller; text-align: center"
                  >
                    {{ blog.num_votes_up }}
                  </CTableDataCell>
                  <CTableDataCell
                    style="font-size: smaller; text-align: center"
                  >
                    {{ blog.num_votes_down }}
                  </CTableDataCell>
                  <CTableDataCell
                    style="
                      font-size: smaller;
                      white-space: nowrap;
                      text-align: center;
                    "
                  >
                    <CLink
                      @click="clickEditBlogPost(blog.id)"
                      :label="$t('message.action_edit')"
                      class="btn btn-sm btn-primary"
                    >
                      <CIcon name="cil-pencil" />
                    </CLink>
                    &nbsp;
                    <CLink
                      @click="clickDeleteBlogPost(blog.id)"
                      :label="$t('message.action_delete')"
                      class="btn btn-sm btn-danger"
                    >
                      <CIcon name="cil-trash" />
                    </CLink>
                  </CTableDataCell>
                </CTableRow>
              </template>
            </CTableBody>
          </CTable>
        </CCardBody>
        <CCardFooter>
          <CButton class="btn-sm btn-primary" @click="clickCreateBlogPost">
            <CIcon name="cil-image-plus" />
            {{ $t('message.create_blog_post') }}
          </CButton>
        </CCardFooter>
      </CCard>
    </CCol>
  </CRow>
</template>

<script>
import { reactive } from 'vue'
import clientUtils from '@/utils/api_client'

export default {
  name: 'MyBlog',
  data: () => {
    let blogPostList = reactive([])
    clientUtils.apiDoGet(
      clientUtils.apiMyBlog,
      (apiRes) => {
        if (apiRes.status == 200) {
          blogPostList.push(...apiRes.data)
        } else {
          console.error('Getting blog post list was unsuccessful: ' + apiRes)
        }
      },
      (err) => {
        console.error('Error getting blog post list: ' + err)
      },
    )
    return {
      blogPostList: blogPostList,
    }
  },
  props: ['flashMsg'],
  methods: {
    clickCreateBlogPost() {
      this.$router.push({ name: 'CreatePost' })
    },
    clickEditBlogPost(id) {
      this.$router.push({ name: 'EditPost', params: { id: id.toString() } })
    },
    clickDeleteBlogPost(id) {
      this.$router.push({ name: 'DeletePost', params: { id: id.toString() } })
    },
  },
}
</script>
