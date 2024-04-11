<template>
  <el-container>
    <el-header class="headerSty">
      <div class="header-ico move-sty" style="float: left">
        <!-- <img style="height: 100%;" src="../assets/k8s.png" /> -->
      </div>
      <div
        class="header-ico"
        style="
          float: right;
          display: flex;
          justify-content: flex-end;
          flex-direction: row;
          gap: 20px;
        "
      >
        <Avatar></Avatar>
      </div>
    </el-header>
    <el-main class="about-main">
      <el-skeleton :rows="3" :loading="loading" :animated="true">
        <template #template>
          <div style="padding: 14px; text-align: left">
            <div
              style="width: 80%; margin: 30px 0px 40px 40px; height: 30px"
              class="card-skeleton"
            />
            <div class="flex-box">
              <div class="ns-card card-skeleton" />
              <div class="ns-card card-skeleton" />
              <div class="ns-card card-skeleton" />
            </div>
          </div>
        </template>
        <template #default>
          <div class="dashbord-tips">Welcome, please select namespace :</div>
          <div class="flex-box">
            <div
              class="ns-card"
              v-for="ns in namespace"
              :key="ns"
              @click="goDashboard(ns)"
            >
              <div class="ns-title">{{ ns }}</div>
              <div><i class="el-icon-arrow-right ns-to"></i></div>
            </div>
          </div>
        </template>
      </el-skeleton>
    </el-main>
  </el-container>
</template>
<script setup lang="ts">
import Avatar from "@/components/avatar.vue";
import { onMounted, watch, ref } from "vue";
import { ElNotification } from "element-plus";
import { useStore } from "@/store";
import { useRouter, useRoute } from "vue-router";
import { currentPageVue } from "./resource/util";
import { routerToNamespace } from "./resource/util";
const store = useStore();

onMounted(() => {
  if (!localStorage.getItem("token")) {
    ElNotification({ title: "未登陆", type: "error", duration: 3000 });
    router.push({ name: "Index" });
  }
});

let loading = ref(false);
let namespace = ref(<any>[]);

let rule = JSON.parse(localStorage.getItem("clusterRole"));
getNamespaceList();

function getNamespaceList() {
  if (!rule) {
    ElNotification({ title: "未登陆", type: "error", duration: 3000 });
    router.push({ name: "Index" });
  }
  let nsList = [];
  for (let ns in rule) {
    nsList.push(ns);
  }
  nsList.sort();
  namespace.value = nsList;
}
// watch(
//   () => store.state.socket.socket.isConnected,
//   (connect) => {
//     if(connect){
//       getNamespaceList()
//       getResourceList()
//     }
//   }
// )

let resourceTypeList = <string[]>[];

watch(
  () => namespace,
  () => {
    console.log("namespace change");
  }
);

console.log(currentPageVue("/"));
const router = useRouter();

function goDashboard(ns: string) {
  routerToNamespace(ns);
}
</script>
<style lang="scss" scoped>
@import "css/about.scss";
</style>
