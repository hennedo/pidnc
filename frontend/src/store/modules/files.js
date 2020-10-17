import axios from "axios";

export const namespaced = true
export const state = {
    files: [],
    settings: {}
}
export const mutations = {
    SET_FILES(state, files) {
        state.files = files;
    },
    ADD_FILE(state, file) {
        state.files.push(file);
    },
    REMOVE_FILE(state, index) {
        state.files.splice(index, 1)
    }
}
export const actions = {
    getAll({ commit, state }) {
        if (state.files.length > 0) {
            return state.files
        }
        return axios.get("/api/files").then(response => {
            if(response) {
                commit('SET_FILES', response.data)
                return response.data
            } else {
                return []
            }
        }).catch(() => {
            return []
        })
    },
    remove({ commit, state }, { id }) {
        axios.get(`/api/${id}/delete`).then(() => {
            for(let i = 0; i < state.files.length; i++) {
                if(state.files[i].ID === id) {
                    commit('REMOVE_FILE', i)
                    break;
                }
            }
        }).catch(err => {
            console.log(err)
        })
    },
    uploadFile({ commit }, { file }) {
        let formdata = new FormData();
        formdata.append("gcode", file)
        axios.post( '/api/upload',
            formdata,
            {
                headers: {
                    'Content-Type': 'multipart/form-data'
                }
            }
        ).then( response => {
            console.log(response.data)
            commit('ADD_FILE', response.data)
        })
        .catch(function(){
            console.log('FAILURE!!');
        });
    }
}
