<template>
  <div class="submit-card">
    <div v-if="!finishedSubmit">
      <p>{{ me.number_of_submitted_cards }}/{{ state.cards_per_player }}</p>
      <div
        class="card-form"
        contenteditable
        autofocus
        ref="cardForm"
      ></div>
      <div
        class="submit-btn"
        @click="submit"
      >Submit</div>
    </div>
    <div v-else>
      Waiting for other players
    </div>
  </div>
</template>

<script>
export default {
  name: 'SubmitCard',
  props: {
    state: Object
  },
  methods: {
    submit () {
      const cardText = this.$refs.cardForm.innerText
      this.$refs.cardForm.innerText = ''
      this.$emit('submit', cardText)
    }
  },
  computed: {
    me () {
      return this.state.players.find(p => p.id === this.state.player_id)
    },
    finishedSubmit () {
      return this.me().number_of_submitted_cards === this.state.cards_per_player
    }
  }
}
</script>

<style>
.submit-card {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.card-form {
  display: block;
  width: 40vh;
  height: 60vh;
  border: 1px solid #a3a3a3;
  border-radius: 5px;
  padding: 3vh;
  text-align: left;
}

.card-form:focus {
  outline: none;
}

.submit-btn {
  width: 40vh;
  margin-top: 20px;
  padding: 5px 0;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  border: 2px solid #000000;
  cursor: pointer;
}

.submit-btn:hover {
  color: #ffffff;
  background-color: #000000;
}
</style>
