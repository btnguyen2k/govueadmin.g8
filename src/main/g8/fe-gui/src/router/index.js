//#GovueAdmin-Customized
import Vue from 'vue'
import Router from 'vue-router'

// Containers
const TheContainer = () => import('@/containers/TheContainer')

// Views
const Dashboard = () => import('@/views/gva/Dashboard')
const MyBlog = () => import('@/views/gva/blog/MyBlog')
const CreatePost = () => import('@/views/gva/blog/CreatePost')
const EditPost = () => import('@/views/gva/blog/EditPost')
const DeletePost = () => import('@/views/gva/blog/DeletePost')

// Views - Pages
const Login = () => import('@/views/gva/pages/Login')
const Page404 = () => import('@/views/pages/Page404')
const Page500 = () => import('@/views/pages/Page500')
const Register = () => import('@/views/pages/Register')

Vue.use(Router)

let router = new Router({
    mode: 'hash', // https://router.vuejs.org/api/#mode
    linkActiveClass: 'active',
    scrollBehavior: () => ({y: 0}),
    base: "/app/",
    routes: configRoutes()
})

import appConfig from "@/utils/app_config"
import utils from "@/utils/app_utils"
import clientUtils from "@/utils/api_client"

router.beforeEach((to, from, next) => {
    if (!to.matched.some(record => record.meta.allowGuest)) {
        let session = utils.loadLoginSession()
        if (session == null) {
            //redirect to login page if not logged in
            return next({name: "Login", query: {returnUrl: router.resolve(to, from).href}})
        }
        let lastUserTokenCheck = utils.localStorageGetAsInt(utils.lskeyLoginSessionLastCheck)
        if (lastUserTokenCheck + 60 < utils.getUnixTimestamp()) {
            lastUserTokenCheck = utils.getUnixTimestamp()
            let token = session.token
            clientUtils.apiDoPost(clientUtils.apiVerifyLoginToken, {app: appConfig.APP_ID, token: token},
                (apiRes) => {
                    if (apiRes.status != 200) {
                        //redirect to login page if session verification failed
                        console.error("Session verification failed: " + JSON.stringify(apiRes))
                        return next({name: "Login", query: {returnUrl: router.resolve(to, from).href}})
                        // return next({name: "Login", query: {returnUrl: to.fullPath}})
                    } else {
                        utils.localStorageSet(utils.lskeyLoginSessionLastCheck, lastUserTokenCheck)
                        next()
                    }
                },
                (err) => {
                    console.error("Session verification error: " + err)
                    //redirect to login page if cannot verify session
                    return next({name: "Login", query: {returnUrl: router.resolve(to, from).href}})
                    // return next({name: "Login", query: {returnUrl: to.fullPath}})
                })
        } else {
            next()
        }
    } else {
        next()
    }
})

export default router

import i18n from '../i18n'

function configRoutes() {
    return [
        {
            path: '/',
            redirect: {name: "Dashboard"},
            name: 'Home', meta: {label: i18n.t('message.home')},
            component: TheContainer,
            children: [
                {
                    path: 'dashboard',
                    name: 'Dashboard',
                    meta: {label: i18n.t('message.dashboard')},
                    component: Dashboard
                },
                {
                    path: 'posts', meta: {label: i18n.t('message.blog')},
                    component: {
                        render(c) {
                            return c('router-view')
                        }
                    },
                    children: [
                        {
                            path: '',
                            meta: {label: i18n.t('message.my_blog')},
                            name: 'MyBlog',
                            component: MyBlog,
                            props: true, //[props=true] to pass flashMsg
                        },
                        {
                            path: '_create',
                            meta: {label: i18n.t('message.create_blog_post')},
                            name: 'CreatePost',
                            component: CreatePost,
                        },
                        {
                            path: '_edit/:id',
                            meta: {label: i18n.t('message.edit_blog_post')},
                            name: 'EditPost',
                            component: EditPost,
                        },
                        {
                            path: '_delete/:id',
                            meta: {label: i18n.t('message.delete_blog_post')},
                            name: 'DeletePost',
                            component: DeletePost,
                        },
                    ]
                },
            ]
        },
        {
            path: '/pages', redirect: {name: "Page404"}, name: 'Pages',
            component: {
                render(c) {
                    return c('router-view')
                }
            },
            meta: {
                allowGuest: true //do not required login to view
            },
            children: [
                {
                    path: '404', name: 'Page404', component: Page404
                },
                {
                    path: '500', name: 'Page500', component: Page500
                },
                {
                    path: 'login', name: 'Login', component: Login,
                    props: (route) => ({returnUrl: route.query.returnUrl}),
                    params: (route) => ({returnUrl: route.query.returnUrl}),
                },
                {
                    path: 'register', name: 'Register', component: Register
                }
            ]
        }
    ]
}
