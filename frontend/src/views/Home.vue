<template>
  <main class="container">
    <ul class="nav nav-tabs">
      <li class="nav-item">
        <a class="nav-link" :class="{ active: filter==='uploaded'}" @click="filter='uploaded'" href="#">Uploaded</a>
      </li>
      <li class="nav-item">
        <a class="nav-link" :class="{ active: filter==='machined'}" @click="filter='machined'" href="#">Machined</a>
      </li>
      <li class="nav-item">
        <a class="nav-link" :class="{ active: filter==='locked'}" @click="filter='locked'" href="#">Locked</a>
      </li>
    </ul>
    <div class="table-responsive">
      <table class="table table-sm">
        <thead>
        <tr>
          <th class="col-md-2"></th>
          <th class="col-md-6">
            Filename
          </th>
          <th class="col-md-1">
            Filesize
          </th>
          <th class="col-md-1">
            Sending time
          </th>
          <th class="col-md-2">
          </th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="file in filteredFiles" :key="file.ID">
          <td class="text-center"><img class="preview" :src="'/svg/' + file.ID + '.svg'"></td>
          <td class="align-middle">{{ file.name }}</td>
          <td class="align-middle">{{ file.size | prettyBytes }}</td>
          <td class="align-middle">{{ Math.round(file.size / 115200 * 100)/100 }}s</td>
          <td class="text-right align-middle">
            <a href="#" @click="run(file.ID)" class="btn btn-success btn-lg" :class="{disabled: sending}">
              <svg width="1em" height="1em" viewBox="0 0 16 16" class="bi bi-play-fill" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
                <path d="M11.596 8.697l-6.363 3.692c-.54.313-1.233-.066-1.233-.697V4.308c0-.63.692-1.01 1.233-.696l6.363 3.692a.802.802 0 0 1 0 1.393z"/>
              </svg></a>&nbsp;
            <a href="#" @click="lock(file.ID)" class="btn btn-primary btn-lg" v-if="file.type !== 'locked'" :class="{disabled: gcode.progress}">
              <svg width="1em" height="1em" viewBox="0 0 16 16" class="bi bi-shield-lock" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
                <path fill-rule="evenodd" d="M5.443 1.991a60.17 60.17 0 0 0-2.725.802.454.454 0 0 0-.315.366C1.87 7.056 3.1 9.9 4.567 11.773c.736.94 1.533 1.636 2.197 2.093.333.228.626.394.857.5.116.053.21.089.282.11A.73.73 0 0 0 8 14.5c.007-.001.038-.005.097-.023.072-.022.166-.058.282-.111.23-.106.525-.272.857-.5a10.197 10.197 0 0 0 2.197-2.093C12.9 9.9 14.13 7.056 13.597 3.159a.454.454 0 0 0-.315-.366c-.626-.2-1.682-.526-2.725-.802C9.491 1.71 8.51 1.5 8 1.5c-.51 0-1.49.21-2.557.491zm-.256-.966C6.23.749 7.337.5 8 .5c.662 0 1.77.249 2.813.525a61.09 61.09 0 0 1 2.772.815c.528.168.926.623 1.003 1.184.573 4.197-.756 7.307-2.367 9.365a11.191 11.191 0 0 1-2.418 2.3 6.942 6.942 0 0 1-1.007.586c-.27.124-.558.225-.796.225s-.526-.101-.796-.225a6.908 6.908 0 0 1-1.007-.586 11.192 11.192 0 0 1-2.417-2.3C2.167 10.331.839 7.221 1.412 3.024A1.454 1.454 0 0 1 2.415 1.84a61.11 61.11 0 0 1 2.772-.815z"/>
                <path d="M9.5 6.5a1.5 1.5 0 0 1-1 1.415l.385 1.99a.5.5 0 0 1-.491.595h-.788a.5.5 0 0 1-.49-.595l.384-1.99a1.5 1.5 0 1 1 2-1.415z"/>
              </svg></a>&nbsp;
            <a href="#" @click="remove(file.ID)" class="btn btn-danger btn-lg" v-if="file.type !== 'locked'" :class="{disabled: gcode.progress}">
              <svg width="1em" height="1em" viewBox="0 0 16 16" class="bi bi-trash-fill" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
                <path fill-rule="evenodd" d="M2.5 1a1 1 0 0 0-1 1v1a1 1 0 0 0 1 1H3v9a2 2 0 0 0 2 2h6a2 2 0 0 0 2-2V4h.5a1 1 0 0 0 1-1V2a1 1 0 0 0-1-1H10a1 1 0 0 0-1-1H7a1 1 0 0 0-1 1H2.5zm3 4a.5.5 0 0 1 .5.5v7a.5.5 0 0 1-1 0v-7a.5.5 0 0 1 .5-.5zM8 5a.5.5 0 0 1 .5.5v7a.5.5 0 0 1-1 0v-7A.5.5 0 0 1 8 5zm3 .5a.5.5 0 0 0-1 0v7a.5.5 0 0 0 1 0v-7z"/>
              </svg></a>
          </td>
        </tr>
        </tbody>
      </table>
    </div>
  </main>
</template>

<script>
import axios from 'axios';
import store from '../store/';
import {mapState} from "vuex";

export default {
  name: 'Home',
  data() {
    return {
      filter: "uploaded",
    }
  },
  computed: {
    filteredFiles: function() {
      let currentFilter = this.filter;
      return this.files.filter(file => file.type === currentFilter)
    },
    ...mapState('files', ['files']),
    ...mapState('files', ['gcode'])
  },
  methods: {
    run(id) {
      if(this.gcode.progress != 0) {
        return
      }
      axios.get(`/api/${id}/run`).then(response => {
        if(response.data != null) {
          console.log(response.data);
        }
        for(let i = 0; this.files.length < i; i++) {
          if(this.files[i].ID === id) {
            this.files[i].Type = "machined";
            break;
          }
        }
      }).catch(err => {
        console.log(err)
      })
    },
    lock(id) {
      if(this.gcode.progress != 0) {
        return
      }
      axios.get(`/api/${id}/lock`).then(response => {
        if(response.data != null) {
          console.log(response.data);
        }
        for(let i = 0; this.files.length < i; i++) {
          if(this.files[i].ID === id) {
            this.files[i].type = "locked";
            break;
          }
        }
      }).catch(err => {
        console.log(err)
      })
    },
    remove(id) {
      if(this.gcode.progress != 0) {
        return
      }
      store.dispatch("files/remove", { id })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
