//#GovueAdmin-Customized
export default [
    {
        _name: 'CSidebarNav',
        _children: [
            {
                _name: 'CSidebarNavItem',
                name: 'Dashboard',
                to: '/dashboard',
                icon: 'cil-speedometer',
                badge: {
                    color: 'primary',
                    text: 'NEW'
                }
            },
            {
                _name: 'CSidebarNavItem',
                name: 'Groups',
                to: '/groups',
                icon: 'cil-address-book',
            },
            {
                _name: 'CSidebarNavItem',
                name: 'Users',
                to: '/users',
                icon: 'cil-user',
            },
        ]
    }
]