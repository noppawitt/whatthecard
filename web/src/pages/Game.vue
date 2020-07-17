<template>
  <div class="game">
    <h2 class="header">Room: {{ roomId }}</h2>
    <WaitingRoom
      v-if="state.phase === 'WAITING_PHASE'"
      :state="state"
      @start="start"
    />
    <SubmitCard
      v-else-if="state.phase === 'SUBMIT_PHASE'"
      :state="state"
      @submit="submitCard"
    />
    <div v-else-if="state.phase === 'PLAY_PHASE'">
      <DrawPile :n="20" />
      <DiscardPile :n="20" />
    </div>
    <div v-else>
      Loading
    </div>
  </div>
</template>

<script>
import DrawPile from '../components/DrawPile.vue'
import DiscardPile from '../components/DiscardPile.vue'
import WaitingRoom from '../components/WaitingRoom.vue'
import SubmitCard from '../components/SubmitCard.vue'

export default {
  name: 'Game',
  components: {
    DrawPile,
    DiscardPile,
    WaitingRoom,
    SubmitCard
  },
  data () {
    return {
      roomId: '',
      state: {}
    }
  },
  methods: {
    sendJSON (o) {
      this.ws.send(JSON.stringify(o))
    },
    start () {
      this.sendJSON({ name: 'start' })
    },
    submitCard (text) {
      this.sendJSON({ name: 'add_card', payload: { text } })
    }
  },
  mounted () {
    const name = window.prompt('What is your name?')
    this.roomId = this.$route.params.id
    this.ws = new WebSocket(`ws://localhost:4000/ws/room/${this.roomId}?player_name=${name}`)
    this.ws.addEventListener('message', (event) => {
      this.state = JSON.parse(event.data)
    })
  },
  destroyed () {
    this.ws.close()
  }
}
</script>

<style>
.game {
  position: relative;
  height: 100%;
  width: 100%;
  max-width: 600px;
}

.header {
  height: 20%;
  margin: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  text-align: center;
}

.main {
  height: 80%;
}
</style>
