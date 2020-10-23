//#GovueAdmin-Customized
import Vue from 'vue'
import Router from 'vue-router'

// Groups - #GovueAdmin-Customized
const Groups = () => import('@/views/gva/groups/Groups')
const CreateGroup = () => import('@/views/gva/groups/CreateGroup')
const EditGroup = () => import('@/views/gva/groups/EditGroup')
const DeleteGroup = () => import('@/views/gva/groups/DeleteGroup')

// Users - #GovueAdmin-Customized
const Users = () => import('@/views/gva/users/Users')
const CreateUser = () => import('@/views/gva/users/CreateUser')
const EditUser = () => import('@/views/gva/users/EditUser')
const DeleteUser = () => import('@/views/gva/users/DeleteUser')

// Containers
const TheContainer = () => import('@/containers/TheContainer')

// Views
const Dashboard = () => import('@/views/Dashboard')
const MyBlog = () => import('@/views/gva/blog/MyBlog')
const CreatePost = () => import('@/views/gva/blog/CreatePost')
const EditPost = () => import('@/views/gva/blog/EditPost')

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

function configRoutes() {
    return [
        {
            path: '/',
            redirect: {name: "Dashboard"},
            name: 'Home',
            component: TheContainer,
            children: [
                {
                    path: 'dashboard', name: 'Dashboard', component: Dashboard
                },
                {
                    path: 'posts', meta: {label: 'Blog'},
                    component: {
                        render(c) {
                            return c('router-view')
                        }
                    },
                    children: [
                        {
                            path: '', meta: {label: 'My Blog'}, name: 'MyBlog', component: MyBlog, props: true, //[props=true] to pass flashMsg
                        },
                        {
                            path: '_create', meta: {label: 'Create Blog Post'}, name: 'CreatePost', component: CreatePost,
                        },
                        {
                            path: '_edit/:id', meta: {label: 'Edit Blog Post'}, name: 'EditPost', component: EditPost,
                        },
                        // {
                        //     path: '_delete/:id',
                        //     meta: {label: 'Delete Group'},
                        //     name: 'DeleteGroup',
                        //     component: DeleteGroup,
                        // },
                    ]
                },
                {
                    path: 'groups', meta: {label: 'Groups'},
                    component: {
                        render(c) {
                            return c('router-view')
                        }
                    },
                    children: [
                        {
                            path: '', meta: {label: 'Group List'}, name: 'Groups', component: Groups, props: true,
                        },
                        {
                            path: '_create',
                            meta: {label: 'Create New Group'},
                            name: 'CreateGroup',
                            component: CreateGroup,
                        },
                        {
                            path: '_edit/:id', meta: {label: 'Edit Group'}, name: 'EditGroup', component: EditGroup,
                        },
                        {
                            path: '_delete/:id',
                            meta: {label: 'Delete Group'},
                            name: 'DeleteGroup',
                            component: DeleteGroup,
                        },
                    ]
                },
                {
                    path: 'users', meta: {label: 'Users'},
                    component: {
                        render(c) {
                            return c('router-view')
                        }
                    },
                    children: [
                        {
                            path: '', meta: {label: 'User List'}, name: 'Users', component: Users, props: true,
                        },
                        {
                            path: '_create',
                            meta: {label: 'Create New User'},
                            name: 'CreateUser',
                            component: CreateUser,
                        },
                        {
                            path: '_edit/:username', meta: {label: 'Edit User'}, name: 'EditUser', component: EditUser,
                        },
                        {
                            path: '_delete/:username',
                            meta: {label: 'Delete User'},
                            name: 'DeleteUser',
                            component: DeleteUser,
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
