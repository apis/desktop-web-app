import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";

const app = createApp(App);

const url = new URL("/api/time-event", window.location.href);
url.protocol = url.protocol.replace("http", "ws");
const webSocket = new WebSocket(url.href, ["time"]);

app.provide("webSocket", webSocket);
app.use(router).mount("#app");
