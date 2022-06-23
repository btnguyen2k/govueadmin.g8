// #GoVueAdmin-Customized
import { createI18n } from 'vue-i18n'
import utils from '@/utils/app_utils'

const messages = {
  en: {
    _name: 'English',
    _flag: 'cif-gb',
    _demo_msg:
      'This is instance is for demo purpose only. ' +
      "Login with default account <strong>admin{'@'}local/s3cr3t</strong>." +
      '<br/>Or you can login with your <u>social account</u> via "Login with social account" link ' +
      '(your social account credential <u>will not</u> be stored on the server).',

    message: {
      language: 'Language',
      switch_language_msg:
        'Changing the display language requires the current page to be reloaded. Are you sure you wish to change the display language?',

      actions: 'Actions',
      action_create: 'Create',
      action_save: 'Save',
      action_back: 'Back',
      action_edit: 'Edit',
      action_delete: 'Delete',

      wait: 'Please wait...',

      login: 'Login',
      login_info: 'Please sign in to continue',
      login_social: 'Login with social account',
      username: 'Username',
      password: 'Password',
      error_parse_login_token: 'Error parsing login-token',

      home: 'Home',
      dashboard: 'Dashboard',
      blog: 'Blog',
      my_blog: 'My blog',

      create_blog_post: 'Create blog post',
      edit_blog_post: 'Edit blog post',
      delete_blog_post: 'Delete blog post',
      blog_posts: 'Blog post ({count}) | Blog post ({count}) | Blog posts ({count})',

      error_blog_post_not_found: 'Blog post "{id}" not found!',

      blog_public: 'Public',
      blog_public_msg: 'Other people can see, comment and vote your public posts',
      blog_tcreated: 'Created',
      blog_title: 'Title',
      blog_title_msg: "My awesome blog post's title",
      blog_comments: 'Comments',
      blog_votes: 'Votes',
      blog_editor: 'Editor',
      blog_preview: 'Preview',
      blog_content: 'Content',
      blog_content_msg: "My awesome blog post's content (Markdown supported)",
      blog_created_msg: 'Blog post "{title}" has been created successfully.',
      blog_updated_msg: 'Blog post "{title}" has been updated successfully.',
      blog_deleted_msg: 'Blog post "{title}" has been deleted successfully.',
    },
  },
  vi: {
    _name: 'Tiếng Việt',
    _flag: 'cif-vn',
    _demo_msg:
      'Bản triển khai này dành do mục đích thử nghiệm. ' +
      "Đăng nhập với tài khoản <strong>admin{'@'}local/s3cr3t</strong>." +
      'Hoặc đăng nhập với <i>tài khoản mxh</i> (nhấn vào đường dẫn "Đăng nhập với tài khoản mxh").',

    message: {
      language: 'Ngôn ngữ',
      switch_language_msg:
        'Thay đổi ngôn ngữ hiển thị sẽ tải lại trang hiện thời. Bạn có chắc bạn muốn thay đổi ngôn ngữ hiển thị?',

      actions: 'Hành động',
      action_create: 'Tạo',
      action_save: 'Lưu',
      action_back: 'Quay lại',
      action_edit: 'Sửa',
      action_delete: 'Xoá',

      wait: 'Vui lòng giờ giây lát...',

      login: 'Đăng nhập',
      login_info: 'Đăng nhập để tiếp tục',
      login_social: 'Đăng nhập với tài khoản mxh',
      username: 'Tên đăng nhập',
      password: 'Mật mã',
      error_parse_login_token: 'Có lỗi khi xử lý login-token!',

      home: 'Trang gốc',
      dashboard: 'Trang nhà',
      blog: 'Bài viết',
      my_blog: 'Bài viết của tôi',
      create_blog_post: 'Viết bài mới',
      edit_blog_post: 'Chỉnh sửa bài viết',
      delete_blog_post: 'Xoá bài viết',
      blog_posts: 'Bài viết của tôi ({count})',

      error_blog_post_not_found: 'Bài viết "{id}" không tồn tại!',

      blog_public: 'Công cộng',
      blog_public_msg: 'Các thành viên khác có thể thấy, nhận xét hoặc bầu chọn bài viết của bạn',
      blog_tcreated: 'Thời gian tạo',
      blog_title: 'Tựa đề',
      blog_title_msg: 'Tựa đề của bài viết',
      blog_comments: 'Nhận xét',
      blog_votes: 'Phiếu bầu',
      blog_editor: 'Soạn thảo',
      blog_preview: 'Xem trước',
      blog_content: 'Nội dung',
      blog_content_msg: 'Nội dung bài viết (hỗ trợ Markdown)',
      blog_created_msg: 'Bài viết "{title}" đã được tạo.',
      blog_updated_msg: 'Bài viết "{title}" đã được cập nhật.',
      blog_deleted_msg: 'Bài viết "{title}" đã được xoá.',
    },
  },
}

let savedLocale = utils.localStorageGet('_l')
savedLocale = savedLocale ? (messages[savedLocale] ? savedLocale : 'en') : 'en'
const _i18n = createI18n({
  locale: savedLocale,
  fallbackLocale: 'en',
  messages: messages,
})

/* i18n.global is readonly, we need to clone a new instance and make it reactive */
const i18n = { ..._i18n }
import { reactive, watchEffect } from 'vue'
i18n.global = reactive(i18n.global)
let oldLocale = i18n.global.locale
watchEffect(() => {
  if (i18n.global.locale !== oldLocale) {
    utils.localStorageSet('_l', i18n.global.locale)
    oldLocale = i18n.global.locale
  }
})

export default i18n

export function swichLanguage(locale, refreshPage) {
  if (locale !== oldLocale) {
    i18n.global.locale = locale
    if (refreshPage) {
      window.location.reload(false)
    }
  }
}
