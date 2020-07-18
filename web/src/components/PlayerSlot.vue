<template>
  <div class="player-slot">
    {{ name }}
  </div>
</template>

<script>
export default {
  name: 'PlayerSlot',
  props: {
    name: String
  }
}
</script>

<style scoped>
.player-slot {
  width: 80%;
  margin: 10px 0;
  padding: 5px 0;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  border: 1px solid #555555;
}
</style>
