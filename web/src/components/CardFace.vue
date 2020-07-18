<template>
  <div class="card-face">
    <div class="card-container">
      <div class="card-text-container">
        <div class="card-text">
          {{ text }}
        </div>
      </div>
      <div class="card-author">{{ author }}</div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'CardFace',
  props: {
    text: String,
    author: String
  }
}
</script>

<style scoped>
.card-face {
  position: absolute;
  width: 150px;
  height: 225px;
  background-color: #ffffff;
  border: 1px solid #a3a3a3;
  border-radius: 5px;
  cursor: pointer;
}

.card-container {
  position: relative;
  height: 100%;
  padding: 7%;
  text-align: left;
}

.card-text-container {
  height: 90%;
  overflow-y: scroll;
}

.card-text {
  word-break: break-all;
  font-size: 1em;
}

.card-author {
  position: absolute;
  bottom: 2%;
  right: 7%;
  font-size: 0.8em;
}
</style>
