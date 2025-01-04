import { createRouter, createWebHashHistory } from "vue-router";

export default createRouter({
    history: createWebHashHistory(),
    routes: [
        {
            path: '/',
            redirect: '/select'
        },
        {
            path: '/select',
            component: () => import('../views/select/index.vue')
        },
        {
            path: '/game',
            component: () => import('../views/game/index.vue')
        }
    ],
});