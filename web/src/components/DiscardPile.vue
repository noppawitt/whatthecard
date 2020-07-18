<template>
  <div class="discard-pile">
    <CardFace
      v-for="i in nToRender"
      :key="i"
      :style="{marginLeft:`${i*2}px`, marginTop:`${i}px`}"
      :text="i === nToRender ? topCard.text : ''"
      :author="i === nToRender ? topCard.author : ''"
    />
  </div>
</template>

<script>
import CardFace from './CardFace.vue'

export default {
  name: 'Discardpile',
  components: {
    CardFace
  },
  props: {
    cards: Array,
  },
  computed: {
    nToRender () {
      return this.cards.length > 20 ? 20 : this.cards.length
    },
    topCard () {
      return this.cards[this.cards.length - 1]
    }
  }
}
</script>

<style scoped>
.discard-pile {
  position: relative;
  width: 190px;
  height: 245px;
  margin-bottom: 10px;
}
</style>