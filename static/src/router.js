import vueRouter from 'vue-router'
import Login from './components/login/Login'
import SignUp from './components/login/SignUp'
import Main from "@/components/Main";
import ListTasks from "@/components/tasks/ListTasks";
import AddEvent from "@/components/AddEvent";
import ListEvents from "@/components/ListEvents";

// import Main from './components/Main'

export default new vueRouter({
    mode:'history',
    routes: [
        {
            path: '/',
            component: Main,
        },
        {
            path: '/login',
            component: Login,
        },
        {
            path: '/signUp',
            component: SignUp,
        },
        {
            path: '/add',
            component: AddEvent,
        },
        {
            path: '/events',
            component: ListEvents,
        },
        {
            path: '/tasks',
            component: ListTasks,
        }
    ]
})