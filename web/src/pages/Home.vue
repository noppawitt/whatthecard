<template>
  <div class="home">
    <h1 class="title">What the Card!</h1>
    <Button
      title="Create Room"
      @click="createRoom"
    />
    <Button title="Join Room" />
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
      fetch('room', {
        method: 'POST',
      })
        .then(res => res.json())
        .then(({ room_id }) => {
          this.$router.push(`/room/${room_id}`)
        })
        .catch(console.log)
    }
  }
}
</script>

<style>
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
