import axios from "axios";

export const namespaced = true
export const state = {
    interfaces: []
}
export const mutations = {
    SET_INFO(state, info) {
        state.interfaces = info.interfaces;
    },
}
export const actions = {
    getInfo({ commit, state }) {
        if (state.interfaces.length > 0) {
            return state
        }
        return axios.get("/api/info").then(response => {
            if(response) {
                commit('SET_INFO', response.data)
                return response.data
            } else {
                return {}
            }
        }).catch(() => {
            return {}
        })
    },
}
