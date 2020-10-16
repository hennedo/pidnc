<template>
  <div class="container">
    <h2>IP Addresses</h2>
    <dl>
      <template v-for="iface of interfaces" >
        <dt :key="iface.name">{{ iface.name }}</dt>
        <dd v-for="ip of iface.ip" :key="ip">{{ ip }}</dd>
      </template>
    </dl>
  </div>
</template>

<script>
import store from '../store/';
import {mapState} from "vuex";

export default {
  name: 'IP',
  mounted() {
    store.dispatch('info/getInfo')
  },
  computed: {
    ...mapState('info', ['interfaces'])
  },
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  dl, dt, dd {
    display: inline-block;
  }
  dt {
    padding-right: 1rem;
  }
  dd::after {
    content: ","
  }
  dd:last-of-type::after {
    content: ""
  }
</style>
