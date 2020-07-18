<template>
  <div class="game">
    <DrawPile
      :n="state.draw_pile_left"
      @draw="draw"
    />
    <DiscardPile :cards="state.discard_cards" />
    <div
      class="btn"
      v-if="state.player_id === state.host_id"
      @click="leave"
    >Leave</div>
    <div
      class="btn"
      v-if="state.player_id === state.host_id"
      @click="reset(1)"
    >Reset Pile</div>
    <div
      class="btn"
      v-if="state.player_id === state.host_id"
      @click="reset(0)"
    >Reset Game</div>
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
  },
  methods: {
    draw () {
      this.$emit('draw')
    },
    leave () {
      this.$emit('leave')
    },
    reset (mode) {
      this.$emit('reset', mode)
    }
  }
}
</script>

<style scoped>
.game {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.btn {
  width: 50%;
  margin-top: 10px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  border: 2px solid #555555;
  cursor: pointer;
}

.btn:hover {
  color: #ffffff;
  background-color: #555555;
}
</style>
