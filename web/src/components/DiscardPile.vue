<template>
  <div class="discard-pile">
    <CardBack
      v-for="i in nToRender"
      :key="i"
      :style="{marginLeft:`${i*2}px`, marginTop:`${i}px`}"
    />
  </div>
</template>

<script>
import CardBack from './CardBack.vue'

export default {
  name: 'Discardpile',
  components: {
    CardBack
  },
  props: {
    n: Number
  },
  computed: {
    nToRender () {
      return this.n > 20 ? 20 : this.n
    }
  }
}
</script>

<style>
.discard-pile {
  position: relative;
  width: 150px;
  height: 225px;
}
</style>