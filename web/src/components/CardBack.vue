<template>
  <div
    class="card-back"
    @click="$emit('click')"
  ></div>
</template>

<script>
export default {
  name: 'CardBack'
}
</script>

<style>
.card-back {
  display: block;
  position: absolute;
  width: 100px;
  height: 150px;
  background-color: #e44e4e;
  border: 1px solid #a3a3a3;
  border-radius: 5px;
  cursor: pointer;
}
</style>
