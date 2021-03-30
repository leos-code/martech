import Vue from 'vue';
import VueRouter from 'vue-router';

Vue.use(VueRouter);

const routes = [
    {
        path: '/',
        redirect: '/login',
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import('@portal/views/Login/index.vue'),
    },
    {
        path: '/register',
        name: 'Register',
        component: () => import('@portal/views/Register/index.vue'),
    },
    // 404
    {
        path: '*',
        redirect: '/',
    },
];

const router = new VueRouter({
    routes,
});

export default router;
