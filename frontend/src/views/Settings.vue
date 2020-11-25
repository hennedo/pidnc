<template>
  <main class="container">
    <h1>Settings
      <a class="nav-link" href="https://github.com/hennedo/pidnc" target="_blank">
        <svg width="1em" height="1em" viewBox="0 0 16 16" class="bi bi-info-circle" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
          <path fill-rule="evenodd" d="M8 15A7 7 0 1 0 8 1a7 7 0 0 0 0 14zm0 1A8 8 0 1 0 8 0a8 8 0 0 0 0 16z"/>
          <path d="M8.93 6.588l-2.29.287-.082.38.45.083c.294.07.352.176.288.469l-.738 3.468c-.194.897.105 1.319.808 1.319.545 0 1.178-.252 1.465-.598l.088-.416c-.2.176-.492.246-.686.246-.275 0-.375-.193-.304-.533L8.93 6.588z"/>
          <circle cx="8" cy="4.5" r="1"/>
        </svg>
      </a></h1>
    <form @submit.prevent="save()">
      <div class="mb-3">
        <label class="form-check-label">Serial Port</label>
        <select class="form-select" v-model="settings.serialPort">
          <option value="">---</option>
          <option v-for="port in ports" :key="port" :value="port" >{{ port }}</option>
        </select>
      </div>
      <button type="submit" class="btn btn-primary">Submit</button>
    </form>
  </main>
</template>

<script>
import store from '../store/';
import {mapState} from "vuex";

export default {
  name: 'Settings',
  data() {
    return {
      filter: "uploaded",
      ports: []
    }
  },
  computed: {
    ...mapState('info', ['settings'])
  },
  mounted() {
    store.dispatch('info/getPorts').then(ports => {
      this.ports = ports;
    })
  },
  methods: {
    save() {
      store.dispatch('info/updateSettings', { settings: this.settings })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
