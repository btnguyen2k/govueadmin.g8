//#GovueAdmin-Customized
import Vue from 'vue'
import VueI18n from 'vue-i18n'

const messages = {
    en: {
        message: {
            action_create: 'Create',
            action_back: 'Back',
            login: 'Login',
            login_info: 'Please sign in to continue',
            login_social: 'Login with social account',
            username: 'Username',
            password: 'Password',
            demo_msg: "This is instance is for demo purpose only. Login with default account <strong>admin@local/s3cr3t</strong>.<br/>Or you can login with your <u>social account</u> via \"Login with social account\" link (your social account credential <u>will not</u> be stored on the server).",
            wait: 'Please wait...',
            error_parse_login_token: 'Error parsing login-token',
            dashboard: 'Dashboard',
            blog: 'Blog',
            my_blog: 'My blog',
            create_blog_post: 'Create blog post',
            home: 'Home',
            blog_public: 'Public',
            blog_public_msg: 'Other people can see, comment and vote your public posts',
            blog_title: 'Title',
            blog_title_msg: "My blog post's awesome title",
            blog_editor: 'Editor',
            blog_preview: 'Preview',
            blog_content: 'Content',
            blog_content_msg: "My blog post's awesome content (Markdown supported)",
            blog_created_msg: 'Blog post "{title}" has been created successfully.',
        }
    },
    vi: {
        message: {
            action_create: 'Tạo',
            action_back: 'Quay lại',
            login: 'Đăng nhập',
            login_info: 'Đăng nhập để tiếp tục',
            login_social: 'Đăng nhập với tài khoản mxh',
            username: 'Tên đăng nhập',
            password: 'Mật mã',
            demo_msg: 'Bản triển khai này dành do mục đích thử nghiệm. Đăng nhập với tài khoản <strong>admin@local/s3cr3t</strong>.<br/>Hoặc đăng nhập với <i>tài khoản mxh</i> (nhấn vào đường dẫn "Đăng nhập với tài khoản mxh").',
            wait: 'Vui lòng giờ giây lát...',
            error_parse_login_token: 'Có lỗi khi xử lý login-token',
            dashboard: 'Trang nhà',
            blog: 'Bài viết',
            my_blog: 'Bài viết của tôi',
            create_blog_post: 'Viết bài mới',
            home: 'Trang gốc',
            blog_public: 'Công cộng',
            blog_public_msg: 'Các thành viên khác có thể thấy, nhận xét hoặc đánh giá bài viết của bạn',
            blog_title: 'Tựa đề',
            blog_title_msg: 'Tựa đề của bài viết',
            blog_editor: 'Soạn thảo',
            blog_preview: 'Xem trước',
            blog_content: 'Nội dung',
            blog_content_msg: 'Nội dung bài viết (hỗ trợ Markdown)',
            blog_created_msg: 'Bài viết "{title}" đã được tạo.',
        }
    }
}

Vue.use(VueI18n)

const i18n = new VueI18n({
    locale: 'vi',
    messages: messages
})

export default i18n
