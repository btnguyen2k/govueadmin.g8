//#GovueAdmin-Customized
import i18n from '../i18n'

export default [
    {
        _name: 'CSidebarNav',
        _children: [
            {
                _name: 'CSidebarNavItem',
                //name: 'Dashboard',
                name: i18n.t('message.dashboard'),
                to: {name: 'Dashboard'},
                icon: 'cil-wallpaper',
                // badge: {
                //     color: 'primary',
                //     text: 'NEW'
                // }
            },
            {
                _name: 'CSidebarNavItem',
                // name: 'My Blog',
                name: i18n.t('message.my_blog'),
                to: {name: 'MyBlog'},
                icon: 'cil-address-book',
            },
            {
                _name: 'CSidebarNavItem',
                // name: 'Create Blog Post',
                name: i18n.t('message.create_blog_post'),
                to: {name: 'CreatePost'},
                icon: 'cil-image-plus',
            },
        ]
    }
]
