import Vue from 'vue';
import Router from 'vue-router';
import HomePage from './Pages/Home'
import LoginPage from './Pages/Login'
import SetupPage from './Pages/Setup'
import SettingsPage from './Pages/Settings'


Vue.use(Router);

export const router = new Router({
    mode: 'history',
    routes: [
        { name: 'home', path: '/', component: HomePage },
        { name: 'settings', path: '/settings', component: SettingsPage },
        { name: 'login', path: '/login', component: LoginPage },
        { name: 'setup', path: '/setup', component: SetupPage },
        // otherwise redirect to home
        { path: '*', redirect: '/' }
    ]
});


router.beforeEach(async(to, from, next) => {
    // redirect to login page if not logged in and trying to access a restricted page
    const publicPages = ['/login', '/signup', '/setup'];
    const authRequired = !publicPages.includes(to.path);
    const token = Vue.$cookies.get("token")
    var validJwt = null
    try {
       validJwt = Vue.jwtDec.decode(token)
    }
    catch(err){
        window.console.log(err)
    }

    if (authRequired && validJwt == null) {
        return next({
            path: '/login',
            query: { returnUrl: to.path }
        });
    }

    next();
})