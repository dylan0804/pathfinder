<template>
  <div 
    class="flex flex-col h-screen" 
  >
    <router-view />
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";

const num = ref(5)

</script>