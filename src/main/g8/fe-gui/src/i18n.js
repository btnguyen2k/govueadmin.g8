//#GovueAdmin-Customized
import Vue from 'vue'
import VueI18n from 'vue-i18n'

const messages = {
    en: {
        message: {
            actions: 'Actions',
            action_create: 'Create',
            action_save: 'Save',
            action_back: 'Back',
            action_edit: 'Edit',
            action_delete: 'Delete',

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
            edit_blog_post: 'Edit blog post',
            delete_blog_post: 'Delete blog post',
            home: 'Home',
            blog_public: 'Public',
            blog_public_msg: 'Other people can see, comment and vote your public posts',
            blog_tcreated: 'Created',
            blog_title: 'Title',
            blog_title_msg: "My blog post's awesome title",
            blog_comments: 'Comments',
            blog_votes: 'Votes',
            blog_editor: 'Editor',
            blog_preview: 'Preview',
            blog_content: 'Content',
            blog_content_msg: "My blog post's awesome content (Markdown supported)",
            blog_created_msg: 'Blog post "{title}" has been created successfully.',
            blog_updated_msg: 'Blog post "{title}" has been updated successfully.',
            blog_posts: 'Blog post ({count}) | Blog post ({count}) | Blog posts ({count})',
        }
    },
    vi: {
        message: {
            actions: 'Hành động',
            action_create: 'Tạo',
            action_save: 'Lưu',
            action_back: 'Quay lại',
            action_edit: 'Sửa',
            action_delete: 'Xoá',

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
            edit_blog_post: 'Chỉnh sửa bài viết',
            delete_blog_post: 'Xoá bài viét',
            home: 'Trang gốc',
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
            blog_posts: 'Bài viết của tôi ({count})',
        }
    }
}

Vue.use(VueI18n)

const i18n = new VueI18n({
    locale: 'vi',
    messages: messages
})

export default i18n
