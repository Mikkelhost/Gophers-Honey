import Vue from 'vue';
import Router from 'vue-router';
import VueCookies from 'vue-cookies'

import HomePage from './Pages/Home'
import LoginPage from './Pages/Login'
import Signup from "./Pages/Signup";


Vue.use(Router);
Vue.use(VueCookies)

export const router = new Router({
    mode: 'history',
    routes: [
        { name: 'home', path: '/', component: HomePage },
        { name: 'login', path: '/login', component: LoginPage },
        { name: 'signup', path: '/signup', component: Signup },
        // otherwise redirect to home
        { path: '*', redirect: '/' }
    ]
});

router.beforeEach((to, from, next) => {
    // redirect to login page if not logged in and trying to access a restricted page
    const publicPages = ['/login', '/signup'];
    const authRequired = !publicPages.includes(to.path);
    const loggedIn = Vue.$cookies.get("session")

    if (authRequired && !loggedIn) {
        return next({
            path: '/login',
            query: { returnUrl: to.path }
        });
    }

    next();
})