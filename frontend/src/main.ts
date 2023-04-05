import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";

const app = createApp(App);

const url = new URL("/api/time-event", window.location.href);
url.protocol = url.protocol.replace("http", "ws");
const webSocket = new WebSocket(url.href, ["time"]);

// // Disable common Chrome shortcuts like Ctrl+N, Ctrl+T, Ctrl+F etc.
// window.addEventListener("keydown", (event: KeyboardEvent) => {
//   if (event.ctrlKey) {
//     event.stopPropagation();
//     event.preventDefault();
//   }
//   // return true;
// });
//
// // Disable Chrome context menu
// window.addEventListener("contextmenu", (event: MouseEvent) => {
//   event.stopPropagation();
//   event.preventDefault();
// });

app.provide("webSocket", webSocket);
app.use(router).mount("#app");
