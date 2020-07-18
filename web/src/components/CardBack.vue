<template>
  <div
    class="card-back"
    @click="$emit('click')"
  >WTC</div>
</template>

<script>
export default {
  name: 'CardBack'
}
</script>

<style scoped>
.card-back {
  display: flex;
  justify-content: center;
  align-items: center;
  position: absolute;
  width: 150px;
  height: 225px;
  color: #ffffff;
  font-size: 2em;
  background-color: #000000;
  border: 1px solid #a3a3a3;
  border-radius: 5px;
  cursor: pointer;
}
</style>
