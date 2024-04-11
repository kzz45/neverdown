<template>
  <div class="log-style" ref="logcontainer">
    <div
      v-for="(msg, index) in items"
      class="log-line"
      :style="checkLevel(msg)"
      :id="String(index)"
      :key="index"
    >
      {{ msg }}
    </div>
    <div v-if="showBtn" class="oper-btn" @click="stopScroll">
      {{ scrollInto ? "stop" : "continue" }}
      <i
        :class="
          scrollInto
            ? 'el-icon-lock oper-icon'
            : 'el-icon-caret-right oper-icon'
        "
      ></i>
    </div>
    <div class="oper-btn" style="right: 200px" @click="clearScroll">clear</div>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick, onMounted, computed } from "vue";
import { useRouter, useRoute } from "vue-router";
import axios from "axios";

const route = useRoute();

let logMsg = ref<string[]>([]);
let logMsgTemp = ref<string[]>([]);
const token = localStorage.getItem("termToken");
// localStorage.removeItem('termToken')
function clearScroll() {
  logMsg.value = [];
  logMsgTemp.value = [];
}
let type = "";
if (route.query.type === "bash") {
  type =
    "/ssh/namespace/${route.query.namespace}/pod/${route.query.podname}/container/${route.query.containername}/pod/token/";
} else {
  type = `/log/stream/namespace/${route.query.namespace}/pod/${route.query.podname}/container/${route.query.containername}/sinceSeconds/${route.query.logSeconds}/sinceTime/nil/token/`;
}
let socket = ref();

onMounted(async () => {
  const cig = await axios.get("config/config.json");
  let env = cig.data;
  const url = "wss://" + String(env.VITE_BASE_URL) + type + token;
  socket.value = new WebSocket(url);

  init();
});

function init() {
  socket.value.onopen = openWatch;
  // 监听socket消息
  socket.value.onmessage = getMag;
  socket.value.onerror = geterror;
}

function initXtermPing() {
  setInterval(function () {
    const order = { type: "ping", input: "", rows: 0, cols: 0 };
    sendmsg(JSON.stringify(order));
  }, 1000);
}
let scrollInto = ref(true);
function stopScroll() {
  scrollInto.value = !scrollInto.value;
}
let showBtn = ref(false);
let stopBtnTimer = ref();
function getMag(event: any) {
  let reader = new FileReader();
  reader.onload = async function (e: any) {
    console.log("e.target.result", e.target.result);
    if (e.target.result.includes('{"level":')) {
      const infoArr = e.target.result.split("\n");
      for (let info of infoArr) {
        logMsg.value.push(`${info}`);
      }
      // const infoArr = e.target.result.split('{"level":')
      // for(let info of infoArr){
      //   if(info.toLowerCase().startsWith('"error"') || info.toLowerCase().startsWith('"debug"') || info.toLowerCase().startsWith('"warn"') || info.toLowerCase().startsWith('"info"') || info.toLowerCase().startsWith('"fatal"')){
      //     info = '{"level":' + info
      //   }
      //   const regex=/\{"level":(.+?)\"\}/g;
      //   let result
      //   if((result = regex.exec(info))!=null) {
      //     if(result[0]){
      //       logMsg.value.push(result[0])
      //       info = info.replace(result[0], "")
      //       info = info.replace('\n', "")
      //     }
      //     if(info){
      //       logMsg.value.push(info)
      //     }
      //   }else{
      //     logMsg.value.push(`${info}`)
      //   }
      // }
    } else {
      logMsg.value.push(e.target.result);
      // if(scrollInto.value){
      //   logMsgTemp.value.push(e.target.result)
      // }
    }

    clearTimeout(stopBtnTimer.value);
    showBtn.value = true;
    // stopBtnTimer.value= setTimeout(() => {
    //   showBtn.value = false
    // }, 500)
    await nextTick();
    if (scrollInto.value) {
      document.getElementById(String(items.value.length - 1))?.scrollIntoView({
        // behavior: "smooth",
        block: "end",
      });
    }
  };
  reader.readAsText(event.data);
}

function checkLevel(msg: string) {
  if (
    msg.toLowerCase().includes('{"level":"error"') ||
    msg.toLowerCase().includes('{"level":"fatal"')
  ) {
    return "color: #f35554";
  } else if (msg.toLowerCase().includes('{"level":"debug"')) {
    return "color: #6fdc00";
  } else if (msg.toLowerCase().includes('{"level":"warn"')) {
    return "color: #deb500";
  }
}

function openWatch() {
  initXtermPing();
}
function sendmsg(order: any) {
  socket.value.send(order);
}
function geterror() {}
function socClose() {
  if (socket.value) {
    socket.value.close();
  }
}

let toTop = ref(1);
let items = computed(() => {
  return logMsg.value;

  // if(scrollInto.value){
  //   return logMsg.value.slice(-(100 * toTop.value))
  // }else{
  //   return logMsgTemp.value.slice(-(100 * toTop.value))
  // }
});
let logcontainer: any = ref(null);
let scrollTimer: any = ref(false);
onMounted(() => {
  // document.addEventListener('scroll', async() => {
  //   let msgLength = 0
  //   if(scrollInto.value){
  //     msgLength = logMsg.value.length
  //   }else{
  //     msgLength = logMsgTemp.value.length
  //   }
  //   if(msgLength >= 100){
  //     if(document.scrollingElement){
  //       if(document.scrollingElement.scrollTop <= 10){
  //         if(scrollTimer.value) return
  //         setTimeout(() => {
  //           scrollTimer.value = false
  //         }, 1000)
  //         toTop.value += 1
  //         scrollTimer.value = true
  //         await nextTick()
  //         document.getElementById(String(101))?.scrollIntoView({
  //           // behavior: "smooth",
  //           block:    "end"
  //         });
  //       }
  //     }
  //   }
  // })
});
</script>

<style lang="scss" scoped>
// @import url(https://cdn.jsdelivr.net/npm/firacode@6.2.0/distr/fira_code.css);
.log-style {
  // font-family: "Fira Code", monospace;
  color: #bbbbbb;
  font-size: 16px;
  overflow-y: auto;
  text-align: left;
}
.log-line {
  padding: 0px 10px;
  margin-bottom: 3px;
  white-space: pre-line;
  width: inherit;
  word-break: break-all;
}
.oper-btn {
  display: inline-block;
  width: 125px;
  background-color: rgb(88 101 242 / 50%);
  color: hwb(0deg 100% 0% / 10%);
  &:hover {
    background-color: rgb(88, 101, 242);
    color: white;
  }
  text-align: center;
  line-height: 1rem;
  margin-top: 10px;
  margin-bottom: 10px;
  font-size: 1rem;
  padding: 10px;
  cursor: pointer;
  border-radius: 3px;
  position: fixed;
  bottom: 0;
  right: 10px;
}
</style>
