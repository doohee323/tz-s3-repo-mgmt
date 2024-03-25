import HomePage from "@/components/Home-page.vue";
import NotfoundPage from "@/components/Notfound-page.vue";
import LoginPage from "@/components/Login-page.vue";
import S3Repos from "@/components/s3Repo/S3Repos-page.vue";
import Upload from "@/components/s3Repo/Upload-page.vue";
import Users from "@/components/user/Users-page.vue";
import UserAdd from "@/components/user/Add-page.vue";
import UserEdit from "@/components/user/Edit-page.vue";

export default [
    {
        path: '/',
        component: HomePage,
        name: 'Home',
    },
    {
        path: '/users',
        component: Users,
        name: 'Users',
    },
    {
        path: '/user/add',
        component: UserAdd,
        name: 'Add User',
    },
    {
        path: '/user/edit',
        component: UserEdit,
        name: 'Edit User',
    },
    {
        path: '/s3Repos',
        component: S3Repos,
        name: 'S3Repos',
    },
    {
        path: '/upload',
        component: Upload,
        name: 'Upload Image',
    },
    {
        path: '/login',
        component: LoginPage,
        name: 'Login',
    },
    {
        path: '*', beforeEnter: (to, from, next) => {
            next('/')
        }
    },
    {
        path: '/404',
        name: 'Not Found',
        component: NotfoundPage
    },
]