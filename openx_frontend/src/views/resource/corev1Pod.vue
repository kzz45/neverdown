<template>
  <el-menu
    :default-active="activeIndex"
    mode="horizontal"
    @select="handleSelect"
  >
    <el-menu-item index="1">MetaData</el-menu-item>
    <el-menu-item index="2">PodTemplate</el-menu-item>
  </el-menu>

  <div class="menu-item" v-show="activeIndex === '1'">
    <MetaData :metadata="itemInfo.metadata"></MetaData>
  </div>

  <div class="menu-item" v-show="activeIndex === '2'">
    <PodTemplate
      :poddata="itemInfo"
      :disablelabels="false"
      :disableannotation="false"
    ></PodTemplate>
  </div>

  <!-- <el-button @click="goTerm">查看log</el-button>
  <el-button @click="goBash">bash</el-button>
  <el-button @click="downloadLog">download</el-button> -->
</template>

<style lang="scss" scoped>
@import "./css/spec.scss";
</style>

<script setup lang="ts">
import MetaData from "./components/metadata.vue";
import PodTemplate from "./components/templatemeta.vue";
import { openTerm } from "./podutils";
import { onMounted, provide, ref } from "vue";

const props = defineProps<{
  itemInfo?: any;
}>();

let logSeconds = ref(0);

provide("podMetadata", props.itemInfo.metadata);

async function goTerm() {
  const container = props.itemInfo.spec.containers[0];
  await openTerm(props.itemInfo.metadata, container, "log", logSeconds.value);
}
async function goBash() {
  const container = props.itemInfo.spec.containers[0];
  await openTerm(props.itemInfo.metadata, container, "bash", logSeconds.value);
}
async function downloadLog() {
  const container = props.itemInfo.spec.containers[0];
  await download(props.itemInfo.metadata, container, false);
}

let currentSeconds = ref(0);
import { saveAs } from "file-saver";
import axios from "axios";
async function download(data: any, container: any, type: boolean) {
  let sinceTime = 0;
  if (type) {
    sinceTime = 0;
    container.downloadPre = true;
  } else {
    sinceTime = currentSeconds.value;
    container.downloadCur = true;
  }
  const cig = await axios.get("config/config.json");
  let env = cig.data;
  const url = `https://${String(env.VITE_BASE_URL)}/log/download/namespace/${
    data.namespace
  }/pod/${data.name}/container/${
    container.name
  }/previous/${type}/sinceSeconds/${sinceTime}/sinceTime/nil`;
  axios({
    url: url,
    method: "get",
    timeout: 30000,
    headers: { Token: localStorage.getItem("token") },
  })
    .then((res) => {
      const str = new Blob([res.data], { type: "text/plain;charset=utf-8" });

      console.log("str::", str);

      const ispre = type ? "_previous" : "";
      type ? (container.downloadPre = false) : (container.downloadCur = false);
      saveAs(
        str,
        data.namespace + "_" + data.name + "_" + container.name + ispre + ".log"
      );
    })
    .catch((err) => {
      console.log("err ::", err);
      type ? (container.downloadPre = false) : (container.downloadCur = false);
    });
}

const activeIndex = ref("1");
const handleSelect = (key: string, keyPath: string[]) => {
  activeIndex.value = keyPath[0] || "1";
};
</script>
