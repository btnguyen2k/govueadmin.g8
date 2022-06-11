// #GoVueAdmin-Customized
import i18n from '@/i18n'

export default [
  {
    component: 'CNavItem',
    name: i18n.global.t('message.dashboard'),
    to: { name: 'Dashboard' },
    icon: 'cil-speedometer',
    badge: {
      color: 'primary',
      text: 'NEW',
    },
  },
  {
    component: 'CNavItem',
    name: i18n.global.t('message.my_blog'),
    to: { name: 'MyBlog' },
    icon: 'cil-address-book',
  },
  {
    component: 'CNavItem',
    name: i18n.global.t('message.create_blog_post'),
    to: { name: 'CreatePost' },
    icon: 'cil-image-plus',
  },
  // {
  //   component: 'CNavTitle',
  //   name: 'Theme',
  // },
  // {
  //   component: 'CNavItem',
  //   name: 'Colors',
  //   to: '/theme/colors',
  //   icon: 'cil-drop',
  // },
  // {
  //   component: 'CNavGroup',
  //   name: 'Base',
  //   to: '/base',
  //   icon: 'cil-puzzle',
  //   items: [
  //     {
  //       component: 'CNavItem',
  //       name: 'Accordion',
  //       to: '/base/accordion',
  //     },
  //     {
  //       component: 'CNavItem',
  //       name: 'Breadcrumbs',
  //       to: '/base/breadcrumbs',
  //     },
  //     {
  //       component: 'CNavItem',
  //       name: 'Cards',
  //       to: '/base/cards',
  //     },
  //   ],
  // },
]
