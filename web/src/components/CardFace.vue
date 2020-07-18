<template>
  <div class="card-face">
    {{ text }}
  </div>
</template>

<script>
export default {
  name: 'CardFace',
  props: {
    text: String
  }
}
</script>

<style>
.card-face {
  display: block;
  position: absolute;
  width: 100px;
  height: 150px;
  background-color: #ffffff;
  border: 1px solid #a3a3a3;
  border-radius: 5px;
  cursor: pointer;
}
</style>
