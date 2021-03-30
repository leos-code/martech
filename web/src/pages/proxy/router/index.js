import Vue from 'vue';
import VueRouter from 'vue-router';

Vue.use(VueRouter);

const routes = [
    {
        path: '/',
        name: 'Index',
        component: () => import('@proxy/views/index.vue'),
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import('@proxy/views/Login/index.vue'),
    },
    {
        path: '/login/callback',
        name: 'LoginCallback',
        component: () => import('@proxy/views/Login/callback.vue'),
    },
    {
        path: '/advertiser',
        name: 'Advertiser',
        component: () => import('@proxy/views/Advertiser/index.vue'),
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
