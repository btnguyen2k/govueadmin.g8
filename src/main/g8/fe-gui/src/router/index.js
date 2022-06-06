// #GovueAdmin-Customized
import { h, resolveComponent } from 'vue'
import { createRouter, createWebHashHistory } from 'vue-router'

import DefaultLayout from '@/layouts/DefaultLayout'

// let router = new Router({
//   mode: 'hash', // https://router.vuejs.org/api/#mode
//   linkActiveClass: 'active',
//   scrollBehavior: () => ({ y: 0 }),
//   base: '/app/',
//   routes: configRoutes(),
// })
const router = createRouter({
  history: createWebHashHistory(process.env.BASE_URL),
  linkActiveClass: 'active',
  scrollBehavior: () => {
    // always scroll to top
    return { top: 0 }
  },
  routes: configRoutes(),
})

import appConfig from '@/utils/app_config'
import utils from '@/utils/app_utils'
import clientUtils from '@/utils/api_client'

router.beforeEach((to, from, next) => {
  if (!to.matched.some((record) => record.meta.allowGuest)) {
    let session = utils.loadLoginSession()
    if (session == null) {
      //redirect to login page if not logged in
      return next({
        name: 'Login',
        query: { returnUrl: router.resolve(to, from).href },
      })
    }
    let lastUserTokenCheck = utils.localStorageGetAsInt(
      utils.lskeyLoginSessionLastCheck,
    )
    if (lastUserTokenCheck + 60 < utils.getUnixTimestamp()) {
      lastUserTokenCheck = utils.getUnixTimestamp()
      let token = session.token
      clientUtils.apiDoPost(
        clientUtils.apiVerifyLoginToken,
        { app: appConfig.APP_ID, token: token },
        (apiRes) => {
          if (apiRes.status != 200) {
            //redirect to login page if session verification failed
            console.error(
              'Session verification failed: ' + JSON.stringify(apiRes),
            )
            return next({
              name: 'Login',
              query: { returnUrl: router.resolve(to, from).href },
            })
            // return next({name: "Login", query: {returnUrl: to.fullPath}})
          } else {
            utils.localStorageSet(
              utils.lskeyLoginSessionLastCheck,
              lastUserTokenCheck,
            )
            next()
          }
        },
        (err) => {
          console.error('Session verification error: ' + err)
          //redirect to login page if cannot verify session
          return next({
            name: 'Login',
            query: { returnUrl: router.resolve(to, from).href },
          })
          // return next({name: "Login", query: {returnUrl: to.fullPath}})
        },
      )
    } else {
      next()
    }
  } else {
    next()
  }
})

export default router

import i18n from '@/i18n'

function configRoutes() {
  return [
    {
      path: '/',
      name: 'Home',
      meta: { label: i18n.global.t('message.home') },
      component: DefaultLayout,
      redirect: { name: 'Dashboard' },
      children: [
        {
          path: 'dashboard',
          name: 'Dashboard',
          meta: { label: i18n.global.t('message.dashboard') },
          component: () => import('@/views/gva/Dashboard'),
        },
        {
          path: 'posts',
          meta: { label: i18n.global.t('message.blog') },
          component: {
            render() {
              return h(resolveComponent('router-view'))
            },
          },
          // component: {
          //   render(c) {
          //     return c('router-view')
          //   },
          // },
          children: [
            {
              path: '',
              meta: { label: i18n.global.t('message.my_blog') },
              name: 'MyBlog',
              component: () => import('@/views/gva/blog/MyBlog'),
              props: true, //[props=true] to pass flashMsg
            },
            {
              path: '_create',
              meta: { label: i18n.global.t('message.create_blog_post') },
              name: 'CreatePost',
              component: () => import('@/views/gva/blog/CreatePost'),
            },
            {
              path: '_edit/:id',
              meta: { label: i18n.global.t('message.edit_blog_post') },
              name: 'EditPost',
              component: () => import('@/views/gva/blog/EditPost'),
            },
            {
              path: '_delete/:id',
              meta: { label: i18n.global.t('message.delete_blog_post') },
              name: 'DeletePost',
              component: () => import('@/views/gva/blog/DeletePost'),
            },
          ],
        },
      ],
    },
    {
      path: '/pages',
      redirect: { name: 'Page404' },
      name: 'Pages',
      component: {
        render() {
          return h(resolveComponent('router-view'))
        },
      },
      meta: {
        allowGuest: true, //do not required login to view
      },
      children: [
        {
          path: '404',
          name: 'Page404',
          component: () => import('@/views/pages/Page404'),
        },
        {
          path: '500',
          name: 'Page500',
          component: () => import('@/views/pages/Page500'),
        },
        {
          path: 'login',
          name: 'Login',
          component: () => import('@/views/gva/pages/Login'),
          props: (route) => ({ returnUrl: route.query.returnUrl }),
          params: (route) => ({ returnUrl: route.query.returnUrl }),
        },
        {
          path: 'register',
          name: 'Register',
          component: () => import('@/views/pages/Register'),
        },
      ],
    },
  ]
}
