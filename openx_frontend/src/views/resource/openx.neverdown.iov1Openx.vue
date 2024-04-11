<template>
  <el-menu
    :default-active="activeIndex"
    mode="horizontal"
    @select="handleSelect"
  >
    <el-menu-item index="1">MetaData</el-menu-item>
    <el-menu-item index="2">Applications</el-menu-item>
  </el-menu>

  <div class="menu-item" v-show="activeIndex === '1'">
    <MetaData :metadata="itemInfo.metadata"></MetaData>
  </div>

  <div v-show="activeIndex === '2'">
    <div class="switch-btns">
      <div
        v-for="(app, index) in itemInfo.spec.applications"
        :key="app.appName"
        :class="mysqlBtn === index ? 'mysql-btn btn-selected' : 'mysql-btn'"
        @click="changeTo(index)"
      >
        <el-icon @click.stop="appCopy(index)" class="copyBtn"
          ><DocumentCopy
        /></el-icon>
        {{ app.appName }}
        <div class="app-del" @click.stop="appDelete(index)">x</div>
      </div>
      <TagInput @value-input="handelTagInput" />
      <div class="appImport-btn" @click="appImport">导入</div>
    </div>

    <div v-if="itemInfo.spec.applications.length > 0">
      <Application :specapp="itemInfo.spec.applications[mysqlBtn]" />
    </div>
  </div>
</template>

<style lang="scss" scoped>
@import "./css/spec.scss";

.switch-btns {
  display: flex;
  padding: 10px;
  flex-wrap: wrap;
  .mysql-btn {
    height: 20px;
    padding: 5px;
    font-size: 1rem;
    margin: 0px 5px 5px 5px;
    border-radius: 4px;
    text-align: center;
    align-self: stretch;
    cursor: pointer;
    line-height: 20px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    border: 1px solid rgb(88, 101, 242);
    .copyBtn {
      margin-right: 10px;
      &:hover {
        color: #8dc63f;
      }
    }
  }
  .btn-selected {
    background-color: rgb(88, 101, 242);
    color: white;
  }
  .app-del {
    margin-left: 15px;
    font-weight: 300;
    &:hover {
      color: crimson;
    }
  }
  .appImport-btn {
    margin-top: 5px;
    margin-bottom: 3px;
    margin-right: 10px;
    background-color: #fff;
    border: 1px solid #ccc;
    border-radius: 4px;
    padding: 3px 8px;
    display: flex;
    justify-content: space-between;
    font-weight: 400;
    cursor: pointer;
    &:hover {
      color: #5865f2;
      border: 1px solid #5865f2;
    }
  }
}
</style>

<script setup lang="ts">
import MetaData from "./components/metadata.vue";
import Application from "./components/application.vue";
import TagInput from "./components/taginput.vue";
import { onMounted, provide, ref } from "vue";
import { initLimits, initPod, initCon } from "./components/tabelUtil";
import { cloneDeep } from "lodash-es";
import { ElNotification } from "element-plus";

const props = defineProps<{
  itemInfo?: any;
  initConfig?: any;
}>();

provide("podMetadata", props.itemInfo.metadata);

for (let app of props.itemInfo.spec.applications) {
  // initPod(app.pod)
  for (let appContainer of app.pod.spec.containers) {
    initLimits(appContainer);
    // initCon(appContainer)
  }
}

onMounted(() => {
  // if(props.itemInfo.spec.applications.length >= 0){
  //   for(let app of props.itemInfo.spec.applications){
  //     for(let appContainer of app.pod.spec.containers){
  //       initLimits(appContainer)
  //     }
  //   }
  // }
  if (props.initConfig) {
    const initAppName = props.itemInfo.spec.applications.findIndex(
      (app: any) => {
        return app.appName === props.initConfig.appName;
      }
    );
    if (initAppName >= 0) {
      mysqlBtn.value = initAppName;
      activeIndex.value = "2";
    }
  }
});

let mysqlBtn = ref(0);
function changeTo(gole: number) {
  mysqlBtn.value = gole;
}

function handelTagInput(appName: string) {
  const appInfo = {
    replicas: "0",
    appName,
    cloudNetworkConfig: {
      aliyunSLB: {},
    },
    persistentStorage: {
      storageVolumePath: "",
    },
    // horizontalPodAutoscalerSpec: {
    //   maxReplicas: 10, minReplicas: 1,
    //   metrics: [],
    //   behavior: {
    //     scaleDown: {
    //       policies: [], selectPolicy: '', stabilizationWindowSeconds: 0
    //     },
    //     scaleUp: {
    //       policies: [], selectPolicy: '', stabilizationWindowSeconds: 0
    //     }
    //   }
    // },
    watchPolicy: "in-place-upgrade",
    pod: {
      metadata: {},
      spec: {
        volumes: [],
        containers: [],
        nodeSelector: {},
        serviceAccountName: "",
        imagePullSecrets: [],
        affinity: [],
        tolerations: [],
      },
    },
    service: {
      metadata: {},
      spec: {
        clusterIP: "",
        type: "",
        ports: [],
      },
    },
    extensionService: {
      metadata: {},
      spec: {
        clusterIP: "",
        type: "",
        ports: [],
      },
    },
  };
  props.itemInfo.spec.applications.push(appInfo);
}

function appDelete(appIndex: number) {
  const appNums = props.itemInfo.spec.applications.length;
  if (appIndex < mysqlBtn.value) {
    changeTo(mysqlBtn.value - 1);
  }
  if (appIndex === Number(mysqlBtn.value) && appIndex === appNums - 1) {
    if (mysqlBtn.value - 1 >= 0) {
      changeTo(mysqlBtn.value - 1);
    }
  }
  props.itemInfo.spec.applications.splice(appIndex, 1);
}
function appCopy(appIndex: number) {
  const copyApp = cloneDeep(props.itemInfo.spec.applications[appIndex]);
  localStorage.setItem("applications", JSON.stringify(copyApp));
  ElNotification({
    title: "导出成功",
    message: "",
    type: "success",
    duration: 1000,
  });
}
function appImport() {
  if (localStorage.getItem("applications")) {
    const newCloneApp = JSON.parse(
      String(localStorage.getItem("applications"))
    );
    for (let vol of newCloneApp.pod.spec.volumes) {
      if (vol.volumeSource.hostPath) {
        if (vol.volumeSource.hostPath.path) {
          let pathArr = vol.volumeSource.hostPath.path.split("/");
          if (pathArr.length >= 5) {
            pathArr[3] = props.itemInfo.metadata.namespace;
          }
          vol.volumeSource.hostPath.path = pathArr.join("/");
        }
      }
    }
    for (let container of newCloneApp.pod.spec.containers) {
      if (container?.env?.length) {
        for (let oneEnv of container.env) {
          if (oneEnv.name === "BATTLEVERIFY_DYNAMIC_NAMESPACE") {
            oneEnv.value = props.itemInfo.metadata.namespace;
          }
        }
      }
    }
    props.itemInfo.spec.applications.push(newCloneApp);
    const appNums = props.itemInfo.spec.applications.length;
    changeTo(appNums - 1);
  }
}

const activeIndex = ref("1");
const handleSelect = (key: string, keyPath: string[]) => {
  activeIndex.value = keyPath[0] || "1";
};
</script>
