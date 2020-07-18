<template>
  <button
    class="button"
    @click="$emit('click')"
  >{{ title }}</button>
</template>

<script>
export default {
  name: 'Button',
  props: {
    title: String
  }
}
</script>

<style scoped>
.button {
  background-color: white;
  color: black;
  border: 2px solid #555555;
  padding: 25px 0;
  margin: 10px;
  cursor: pointer;
  font-size: 1em;
  width: 250px;
}

.button:hover {
  background-color: #555555;
  color: white;
}

button:focus {
  outline: 0;
}
</style>
