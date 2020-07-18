<template>
  <div class="draw-pile">
    <CardBack
      v-for="i in nToRender"
      :key="i"
      :style="{marginLeft:`${i*2}px`, marginTop:`${i}px`}"
      @click="onCardClick(i)"
    />
  </div>
</template>

<script>
import CardBack from './CardBack.vue'

export default {
  name: 'DrawPile',
  components: {
    CardBack
  },
  props: {
    n: Number
  },
  methods: {
    onCardClick (i) {
      if (i !== this.n) {
        return
      }
      this.$emit('draw')
    }
  },
  computed: {
    nToRender () {
      return this.n > 20 ? 20 : this.n
    }
  }
}
</script>

<style>
.draw-pile {
  position: relative;
  width: 140px;
  height: 170px;
  margin-bottom: 10px;
}
</style>