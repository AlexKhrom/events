import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

const state = {
    // titles: [
    //     '😴 уснул',
    //     '🤤 проснулся',
    //     '🍔 поел',
    // ],
    titles: JSON.parse(localStorage.getItem('titles') || '[]').filter(e => !e.deleted),
    events: JSON.parse(localStorage.getItem('events') || '[]').filter(e => !e.deleted),
    tasks: JSON.parse(localStorage.getItem('tasks') || '[]').filter(e => !e.deleted),
    user: {
        login: '',
        email: '',
        password: '',

    },
    async makeReq(url,method,body) {
        let response
        if(body===undefined){
            response = await fetch(url, {
                method: method,
            });
        }else{
            response = await fetch(url, {
                method: method,
                body:body
            });
        }

        if (response.ok) { // если HTTP-статус в диапазоне 200-299
            // получаем тело ответа (см. про этот метод ниже)
            // let respBody = await response.json()
            // console.log(respBody)
            console.log("okkkk")
            return response
        } else if (response.status > 299 && response.status < 500) {
            console.log("err")
            return 'err'
        } else if (response.status > 500) {
            console.log("some wrong on backend")
            return 'backend err'
        }
    },
}

const mutations = {
    set(state, [variable, value]) {
        state[variable] = value
    },

}

export default new Vuex.Store({
    state,
    mutations,
})
