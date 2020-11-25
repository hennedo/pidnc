<template>
  <footer>
    <div class="progress" style="height: 30px">
      <div class="progress-bar progress-bar-striped progress-bar-animated" role="progressbar" :style="{width: progress + '%'}" :aria-valuenow="progress" aria-valuemin="0" aria-valuemax="100"><slot></slot> {{progress}}%</div>
    </div>
  </footer>
</template>

<script>

export default {
  name: 'Progressbar',
  props: {
    progress: Number,
    text: String,
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  footer {
    width: 100%;
    position: absolute;
    bottom: 0;
  }
</style>
