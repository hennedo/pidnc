import axios from "axios";

export const namespaced = true
export const state = {
    interfaces: [],
    settings: {}
}
export const mutations = {
    SET_INFO(state, info) {
        state.interfaces = info.interfaces;
        state.settings = info.settings;
    },
}
export const actions = {
    getInfo({ commit, state }) {
        if (state.interfaces.length > 0) {
            return state
        }
        return axios.get("/api/settings").then(response => {
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
    updateSettings({commit}, { settings}) {
        return axios.post('/api/settings', settings).then(response => {
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
    getPorts() {
        return axios.get('/api/serialPorts').then(response => {
            if (response) {
                return response.data
            } else {
                return []
            }
        }).catch(() => {
            return []
        })
    }
}
