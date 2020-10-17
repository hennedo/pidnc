<template>
  <main class="container">
    <h1>Settings</h1>
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
