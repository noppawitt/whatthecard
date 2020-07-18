<template>
  <div class="home">
    <h1 class="title">What the Card!</h1>
    <Button
      title="Create Room"
      @click="createRoom"
    />
    <Button
      title="Join Room"
      @click="joinRoom"
    />
  </div>
</template>

<script>
import Button from '../components/Button.vue'

export default {
  name: "Home",
  components: {
    Button
  },
  methods: {
    createRoom () {
      fetch('/room', {
        method: 'POST',
      })
        .then(res => res.json())
        .then(({ room_id }) => {
          this.goToRoom(room_id)
        })
        .catch(console.error)
    },
    goToRoom (id) {
      this.$router.push(`/room/${id}`)
    },
    joinRoom () {
      const roomId = window.prompt('Enter room id')
      if (!roomId) {
        this.$router.push('/')
        return
      }
      this.goToRoom(roomId)
    }
  }
}
</script>

<style scoped>
.home {
  height: 100%;
  display: flex;
  flex-direction: column;
  text-align: center;
  justify-content: center;
  align-items: center;
}

.title {
  margin: 0 0 20px 0;
  font-size: 2.5em;
}
</style>
