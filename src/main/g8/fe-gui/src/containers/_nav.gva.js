//#GovueAdmin-Customized
export default [
    {
        _name: 'CSidebarNav',
        _children: [
            {
                _name: 'CSidebarNavItem',
                name: 'Dashboard',
                to: {name: 'Dashboard'},
                icon: 'cil-wallpaper',
                // badge: {
                //     color: 'primary',
                //     text: 'NEW'
                // }
            },
            {
                _name: 'CSidebarNavItem',
                name: 'My Blog',
                to: {name: 'MyPosts'},
                icon: 'cil-address-book',
            },
            {
                _name: 'CSidebarNavItem',
                name: 'Create Blog Post',
                to: {name: 'CreatePost'},
                icon: 'cil-image-plus',
            },
        ]
    }
]