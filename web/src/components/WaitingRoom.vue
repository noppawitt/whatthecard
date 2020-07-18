<template>
  <div class="waiting-room">
    <PlayerSlot
      v-for="p in state.players"
      :key=p.id
      :name="p.name"
    />
    <div class="row">
      <label for="cards-per-player">cards per player</label>
      <select
        id="cards-per-player"
        v-model="cardsPerPlayer"
        @change="setCardsPerPlayer"
      >
        <option
          v-for="i in 20"
          :key="i"
          :value="i"
        >{{ i }}</option>
      </select>
    </div>
    <div
      class="start-btn"
      v-if="state.player_id === state.host_id"
      @click="start"
    >Start</div>
  </div>
</template>

<script>
import PlayerSlot from './PlayerSlot.vue'

export default {
  name: 'WaitingRoom',
  components: {
    PlayerSlot
  },
  props: {
    state: Object
  },
  data () {
    return {
      cardsPerPlayer: 5
    }
  },
  methods: {
    start () {
      this.$emit('start')
    },
    setCardsPerPlayer () {
      this.$emit('setCardsPerPlayer', this.cardsPerPlayer)
    }
  }
}
</script>

<style scoped>
.waiting-room {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.start-btn {
  width: 80%;
  margin-top: 20px;
  padding: 5px 0;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  border: 2px solid #555555;
  cursor: pointer;
}

.start-btn:hover {
  color: #ffffff;
  background-color: #555555;
}

.row {
  width: 80%;
  display: flex;
}

label {
  margin-right: 10px;
}

#cards-per-player {
  padding: 4px;
}
</style>
