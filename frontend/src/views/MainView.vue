<template>
  <div class="home">
    <MainComponent msg="Welcome to SSE Vue.js 3 Demo App" />
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import MainComponent from "@/components/MainComponent.vue";

export default defineComponent({
  name: "MainView",
  components: {
    MainComponent,
  },
});
</script>
