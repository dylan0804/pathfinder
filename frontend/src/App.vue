<template>
  <div 
    class="flex flex-col h-screen" 
  >
    <router-view />
  </div>

  <Toast />
</template>

<script setup lang="ts">
    import Toast from 'primevue/toast';
</script>