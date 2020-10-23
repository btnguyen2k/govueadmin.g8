<template>
  <div>
    <CRow v-for="post in blogPostList.data">
      <CCol sm="12">
        <CCard accent-color="primary">
          <CCardHeader><strong>Compete. Learn. Win. Sign up for the Azure SQL Championship Today</strong></CCardHeader>
          <CCardBody>
            The last six months have seen an unprecedented acceleration in digital transformation. There has never been a more important time for database administrators (DBAs) and developers to prepare for an increasingly cloud-centric future. However, we’d like to make this interesting and fun. That’s why we’ve partnered with PASS, a worldwide community of over 300,000 data professionals, to create an interactive learning experience that allows you to test your talent and build your skillset for the future—all while earning to chance to win prizes. Sign up now for the Azure SQL Championship, starting October 12, 2020.
          </CCardBody>
        </CCard>
      </CCol>
    </CRow>
  </div>
</template>

<script>
import clientUtils from "@/utils/api_client"

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
  }
}
</script>
