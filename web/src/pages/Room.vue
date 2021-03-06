<template>
  <div class="room">
    <h2 class="header">Room: {{ roomId }}</h2>
    <WaitingRoom
      v-if="state.phase === 'WAITING_PHASE'"
      :state="state"
      @setCardsPerPlayer="setCardsPerPlayer"
      @start="start"
    />
    <SubmitCard
      v-else-if="state.phase === 'SUBMIT_PHASE'"
      :state="state"
      @submit="submitCard"
      @leave="leave"
    />
    <div v-else-if="state.phase === 'PLAY_PHASE'">
      <Game
        :state="state"
        @draw="draw"
        @leave="leave"
        @reset="reset"
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
import { WEBSOCKET_SCHEME } from '../config'

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
    setCardsPerPlayer (n) {
      this.sendJSON({ name: 'set_cards_per_player', payload: { cards_per_player: n } })
    },
    start () {
      this.sendJSON({ name: 'start' })
    },
    submitCard (text) {
      this.sendJSON({ name: 'add_card', payload: { text } })
    },
    draw () {
      this.sendJSON({ name: 'draw_card' })
    },
    leave () {
      this.$router.push('/')
    },
    reset (mode) {
      const confirm = window.confirm('Are you sure?')
      if (confirm) {
        this.sendJSON({ name: 'reset', payload: { mode } })
      }
    }
  },
  mounted () {
    const name = window.prompt('What is your name?')
    if (!name) {
      this.$router.push('/')
      return
    }
    this.roomId = this.$route.params.id.toLowerCase()
    if (!this.roomId) {
      this.$router.push('/')
      return
    }
    this.ws = new WebSocket(`${WEBSOCKET_SCHEME}://${window.location.host}/ws/room/${this.roomId}?player_name=${name}`)
    this.ws.addEventListener('message', (event) => {
      this.state = JSON.parse(event.data)
    })
  },
  destroyed () {
    this.ws.close()
  }
}
</script>

<style scoped>
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
