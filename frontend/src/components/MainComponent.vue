<template>
  <div class="hello">
    <h1>{{ msg }}</h1>
    <div id="wrapper">
      <div />
      <div />
      <div>
        <h1 class="time" :class="actualColor">{{ dateTime }}</h1>
      </div>
      <div />
      <div />
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, inject, onMounted, ref } from "vue";

export default defineComponent({
  name: "MainComponent",
  props: {
    msg: String,
  },
  setup() {
    const webSocket = inject<EventSource>("webSocket");
    const dateTime = ref("");
    const primaryColor = ref(true);
    const actualColor = ref("green");

    onMounted(async () => {
      const response = await fetch("/api/get-time", {
        method: "GET",
      });

      const json = await response.json();
      const date = new Date(json);
      dateTime.value = date.toLocaleString();

      function changeColor() {
        if (primaryColor.value) {
          actualColor.value = "green";
        } else {
          actualColor.value = "blue";
        }
        primaryColor.value = !primaryColor.value;
      }

      changeColor();

      if (webSocket === undefined) {
        throw new Error("web socket is undefined");
      }

      webSocket.onmessage = (event: MessageEvent) => {
        const json = JSON.parse(event.data);
        const date = new Date(json);
        dateTime.value = date.toLocaleString();
        changeColor();
      };
    });

    return {
      actualColor,
      primaryColor,
      dateTime,
      webSocket,
    };
  },
});
</script>

<style scoped>
#wrapper {
  display: flex;
  /*height: 200px;*/
}

#wrapper > div {
  flex-grow: 1;
}

.time {
  padding: 5px;
  transition: all 1s;
}

.blue {
  color: white;
  background: blue;
  opacity: 0.7;
}

.green {
  color: white;
  background: green;
  opacity: 0.7;
}
</style>
