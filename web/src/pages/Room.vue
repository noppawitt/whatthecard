<template>
  <div class="room">
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
      <Game
        :state="state"
        @draw="draw"
      />
    </div>
    <div v-else>
      Loading
    </div>
  </div>
</template>

<script>
import WaitingRoom from '../components/WaitingRoom.vue'
import SubmitCard from '../components/SubmitCard.vue'
import Game from '../components/Game.vue'

export default {
  name: 'Room',
  components: {
    WaitingRoom,
    SubmitCard,
    Game
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
    },
    draw () {
      this.sendJSON({ name: 'draw_card' })
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
.room {
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
