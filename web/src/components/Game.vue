<template>
  <div class="game">
    <DrawPile
      :n="state.draw_pile_left"
      @draw="$emit('draw')"
    />
    <DiscardPile :cards="state.discard_cards" />
  </div>
</template>

<script>
import DrawPile from './DrawPile.vue'
import DiscardPile from './DiscardPile.vue'

export default {
  name: 'WaitingRoom',
  components: {
    DrawPile,
    DiscardPile
  },
  props: {
    state: Object
  }
}
</script>

<style>
.game {
  display: flex;
  flex-direction: column;
  align-items: center;
}
</style>
