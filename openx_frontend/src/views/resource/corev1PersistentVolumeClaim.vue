<template>
  <el-skeleton :rows="3" :loading="loading" :animated="true">
    <template #template>
      <div class="header-btn">
        <div
          class="ns-card card-skeleton"
          style="background-color: #292b2f"
        ></div>
      </div>
      <div class="grid-box">
        <div class="configmap-card card-skeleton"></div>
      </div>
    </template>
    <template #default>
      <div class="header-btn">
        <div class="ns-card" @click="newConfig">
          <div class="ns-title">{{ i18nt("common.add") }}</div>
          <div><i class="el-icon-plus ns-to"></i></div>
        </div>
        <div class="ns-card" @click="importConfig">
          <div class="ns-title">{{ i18nt("common.import") }}</div>
          <div><i class="el-icon-upload2 ns-to"></i></div>
        </div>
        <div class="ns-search">
          <lzm-name-filter
            :filterStr="filterName"
            @filterinput="filterChange"
            :showName="false"
          ></lzm-name-filter>
        </div>
      </div>
      <div class="grid-box">
        <div
          class="configmap-card"
          v-for="config in configmapArr"
          :key="config.Name"
        >
          <div class="configmap-title">
            {{ config.Name }}
            <el-tooltip placement="bottom-end">
              <template #content>
                <div
                  class="more-act"
                  @click="configEditor(config, formatKeys(config.data)[0])"
                >
                  {{ i18nt("common.edit") }}
                </div>
                <div class="more-act" @click="configExport(config)">
                  {{ i18nt("common.export") }}
                </div>
                <div class="more-act" @click="confirmDel(config)">
                  {{ i18nt("common.delete") }}
                </div>
                <div class="more-act" @click="cardExportTo(config)">导出至</div>
              </template>
              <div class="config-close"><i class="el-icon-more"></i></div>
            </el-tooltip>
          </div>
          <div class="config-key">
            <div
              class="config-sty"
              v-for="configItem in formatKeys(config.data)"
              :key="configItem"
              @click="configEditor(config, configItem)"
            >
              <div class="itemkey">{{ configItem }}</div>
            </div>
          </div>
        </div>
      </div>
    </template>
  </el-skeleton>

  <el-dialog
    title="添加"
    v-model="dialogVisible"
    custom-class="configMapDialog"
    top="7vh"
    width="965px"
    :fullscreen="isfullscreen"
  >
    <template #title>
      <div style="height: 30px; color: black; text-align: left">
        {{ i18nt("common.add") }}
      </div>
      <i
        class="el-icon-full-screen"
        style="
          color: #909399;
          position: absolute;
          top: 20px;
          right: 50px;
          cursor: pointer;
        "
        @click="isfullscreen = !isfullscreen"
      ></i>
    </template>
    <div style="width: 100%">
      <div style="height: 50px; text-align: left">
        <el-input
          size="medium"
          v-model="configName"
          style="height: 100%; width: 50%"
          placeholder="请输入Name"
        ></el-input>
      </div>
      <div style="display: flex; flex-wrap: wrap">
        <div
          style="
            height: calc(100vh - 250px);
            width: 340px;
            display: inline-block;
            overflow-x: hidden;
          "
        >
          <el-tag
            v-for="tag in dynamicTags"
            :id="tag"
            :key="tag"
            size="medium"
            :class="selectedConfig === tag ? 'selected-tag' : 'config-tag'"
            :disable-transitions="false"
            @click="selectEditor(tag)"
          >
            <div class="tag-name">{{ tag }}</div>
            <el-icon @click.stop="handleClose(tag)"><Close /></el-icon>
          </el-tag>
          <el-input
            class="input-new-tag"
            v-if="inputVisible"
            v-model="inputValue"
            ref="saveTagInput"
            @keyup.enter="$event.target.blur()"
            @blur="handleInputConfirm"
          >
          </el-input>
          <el-button
            v-else
            class="button-new-tag"
            size="small"
            @click="showInput"
            >+ New Tag</el-button
          >
        </div>
        <div class="yaml-editor yaml-style" style="">
          <yaml-editor
            ref="yamlEditor"
            :config="yamlData"
            @changed="yamlChanged"
          />
        </div>
      </div>
    </div>
    <template #footer>
      <span class="dialog-footer">
        <div class="dialog-btn" @click="confirmEdit">
          {{ i18nt("common.confirm") }}
        </div>
      </span>
    </template>
  </el-dialog>

  <el-dialog
    v-model="exDialogVisible"
    custom-class="configMapDialog hsoDialog"
    width="1165px"
    top="7vh"
  >
    <template #title>
      <div style="height: 30px; color: black; text-align: left">导出至</div>
    </template>

    <div style="width: 100%; display: flex; flex-wrap: wrap">
      <div
        v-for="ns in nslist"
        @click="exportTo(ns.Name)"
        class="ns-crad"
        v-bind:key="ns.Name"
      >
        {{ ns.Name }}
      </div>
    </div>
    <template #footer>
      <span class="dialog-footer">
        <div class="dialog-btn" @click="exDialogVisible = false">
          {{ i18nt("common.cancel") }}
        </div>
      </span>
    </template>
  </el-dialog>
</template>

<style lang="scss" scoped>
@import "css/ConfigMap.scss";
:deep(.el-tag__content) {
  display: flex;
  justify-content: space-between;
  width: 100%;
}
.yaml-style {
  width: calc(100% - 400px);
  display: inline-block;
  background: #ccc;
  margin-left: 20px;
  position: relative;
  height: 100%;
}
.tag-name {
  // width: calc(100% - 24px);
  text-align: left;
  overflow-x: auto;
  overflow-y: hidden;
}
.tag-name::-webkit-scrollbar {
  width: 8px;
  height: 1px;
}
.tag-name::-webkit-scrollbar-thumb {
  border-radius: 8px;
  background-color: #202225;
}
.tag-name ::-webkit-scrollbar-track {
  box-shadow: inset 0 0 5px rgba(0, 0, 0, 0.2);
  background: #36393f;
  border-radius: 8px;
}
.ns-crad {
  width: 260px;
  padding: 12px;
  font-size: 1rem;
  margin: 5px 20px;
  border-radius: 8px;
  background-color: #292b2f;
  color: #fff;
  display: flex;
  align-self: stretch;
  cursor: pointer;
  line-height: 25px;
  justify-content: space-between;
  &:hover {
    background-color: #3c3e42;
  }
}
</style>

<script setup lang="ts">
import { i18nt } from "@/i18n";
import { onMounted, watch, nextTick, ref, computed } from "vue";
import { useRouter, useRoute } from "vue-router";
import { initSocketData, sendSocketMessage, binaryToStr } from "@/api/socket";
import proto from "@p/proto";
const protoApi = proto.github.com.kzz45.neverdown.pkg.openx.aggregator.proto;
const protoV1 = proto.k8s.io.api.core.v1;
import { useStore } from "@/store";
const store = useStore();
const router = useRouter();
const route = useRoute();
import YamlEditor from "@/components/YamlEditor.vue";
import { cloneDeep } from "lodash";

let configmapList = $ref([]);
let dialogVisible = $ref(false);
let isfullscreen = $ref(false);
let loading = $ref(true);
let dialogType = $ref("create");
let filterName = $ref("");
const configmapArr = computed(() => {
  return configmapList.filter((conf: any) => {
    return conf.Name.includes(filterName);
  });
});
function filterChange(name: string) {
  filterName = name;
}
onMounted(() => {
  // getList()
});
watch(
  () => route.path,
  () => {
    loading = true;
    // getList()
  }
);
function getList() {
  const param = {
    nameSpace: route.path.split("/")[1],
    service: "list",
    resourceType: "ConfigMap",
  };
  const senddata = initSocketData(param, "Request");
  sendSocketMessage(senddata, store);
}
watch(
  () => store.state.socket.socket.message,
  (msg) => {
    // returnResource(msg)
  }
);
import { ElNotification } from "element-plus";
function returnResource(msg: any) {
  const result = protoApi.Response.decode(msg);
  if (
    result.param.service !== "list" &&
    result.param.service !== "harbor" &&
    result.param.service !== "ping"
  ) {
    getList();
  }
  if (result.param.resourceType === "ConfigMap") {
    switch (result.param.service) {
      case "create":
        if (result.code === 0) {
          ElNotification({
            title: "新增成功",
            message: "success",
            type: "success",
            duration: 2000,
          });
        } else {
          loading = false;
          ElNotification({
            title: "新增失败",
            message: binaryToStr(result.result),
            type: "error",
            duration: 3000,
          });
        }
        break;
      case "update":
        if (result.code === 0) {
          ElNotification({
            title: "修改成功",
            message: "success",
            type: "success",
            duration: 2000,
          });
        } else {
          loading = false;
          ElNotification({
            title: "修改失败",
            message: binaryToStr(result.result),
            type: "error",
            duration: 3000,
          });
        }
        break;
      case "delete":
        if (result.code === 0) {
          ElNotification({
            title: "删除成功",
            message: "success",
            type: "success",
            duration: 2000,
          });
        } else {
          loading = false;
          ElNotification({
            title: "删除失败",
            message: binaryToStr(result.result),
            type: "error",
            duration: 3000,
          });
        }
        break;
      case "list":
        configmapList = protoV1.ConfigMapList.decode(result.result).items;
        loading = false;
        break;
      default:
        break;
    }
  }
}

function configExport(config: any) {
  localStorage.setItem("configMapStorage", JSON.stringify(config));
  ElNotification({
    title: "导出到缓存成功",
    message: "",
    type: "success",
    duration: 1000,
  });
}
function importConfig() {
  if (localStorage.getItem("configMapStorage")) {
    configInfo = JSON.parse(String(localStorage.getItem("configMapStorage")));
    configName = configInfo.Name;
    const keys = Object.keys(configInfo.data);
    if (keys.length > 0) {
      selectEditor(keys[0]);
    }
    dynamicTags = keys;
    dialogType = "create";
    dialogVisible = true;
  } else {
    newConfig();
  }
}
import { ElMessageBox } from "element-plus";
function confirmDel(config: any) {
  ElMessageBox.confirm("确认要删除吗?", "确认", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(() => {
    configDelete(config);
  });
}
function configDelete(config: any) {
  loading = true;
  const keys = Object.keys(config.data).join(",");
  config["keys"] = keys;
  const configData = initSocketData(config, "ConfigMap");
  const param = {
    nameSpace: route.path.split("/")[1],
    service: "delete",
    resourceType: "ConfigMap",
  };
  const senddata = initSocketData(param, "Request", configData);
  sendSocketMessage(senddata, store);
}
function formatKeys(data: object) {
  return Object.keys(data);
}

let dynamicTags = $ref(<any[]>[]);
let inputVisible = $ref(false);
let inputValue = $ref("");
let saveTagInput = $ref(null);
let selectedConfig = $ref("");
let configInfo = $ref({
  data: {},
});
let configName = $ref("");
let yamlData = $ref("");
let yamlEditor = $ref(null);
function newConfig() {
  dialogType = "create";
  dialogVisible = true;
  dynamicTags = [];
  configInfo = {
    data: {},
  };
  configName = "";
  selectedConfig = "";
  yamlData = "";
}
function selectEditor(cfg: string) {
  selectedConfig = cloneDeep(cfg);
  let val = configInfo.data[cfg] ? configInfo.data[cfg] : "";
  yamlData = "1";
  setTimeout(() => {
    yamlData = val;
  });
}

function configEditor(config: any, configKey: string) {
  configInfo = cloneDeep(config);
  configName = configInfo.Name;
  const keys = Object.keys(configInfo.data);
  if (keys.length > 0) {
    if (configKey) {
      setTimeout(async () => {
        selectEditor(configKey);
        await nextTick();
        document.getElementById(configKey).scrollIntoView({
          behavior: "smooth",
          block: "start",
        });
      }, 0);
    } else {
      selectEditor(keys[0]);
    }
  }
  dynamicTags = keys;
  dialogType = "update";
  dialogVisible = true;
}
function yamlChanged(yamlValue: string) {
  configInfo.data[selectedConfig] = yamlValue;
}
function handleClose(tag: string) {
  delete configInfo.data[tag];
  dynamicTags.splice(dynamicTags.indexOf(tag), 1);
  selectEditor(dynamicTags[0]);
}
async function showInput() {
  inputVisible = true;
  await nextTick();
  saveTagInput.$refs.input.focus();
}
function handleInputConfirm() {
  if (inputValue) {
    dynamicTags.push(inputValue);
    configInfo.data[inputValue] = "";
    selectEditor(inputValue);
  }

  yamlEditor.setValue("");
  inputVisible = false;
  inputValue = "";
}
function confirmEdit() {
  const queryEdit = {
    Name: configName,
    data: configInfo.data,
  };
  const configData = initSocketData(queryEdit, "ConfigMap");
  const param = {
    nameSpace: route.path.split("/")[1],
    service: dialogType,
    resourceType: "ConfigMap",
  };
  const senddata = initSocketData(param, "Request", configData);
  sendSocketMessage(senddata, store);
  loading = true;
  dialogVisible = false;
}

let cloneCard = ref();
let exDialogVisible = ref(false);
let nslist = $ref(
  localStorage.getItem("nsList")
    ? JSON.parse(String(localStorage.getItem("nsList")))
    : []
);

function cardExportTo(card: any) {
  cloneCard.value = cloneDeep(card);
  exDialogVisible.value = true;
}
function exportTo(ns: string) {
  cloneCard.value.resourceVersion = "";
  const configData = initSocketData(cloneCard.value, "ConfigMap");
  const param = {
    nameSpace: ns,
    service: "create",
    resourceType: "ConfigMap",
  };
  const senddata = initSocketData(param, "Request", configData);
  sendSocketMessage(senddata, store);
  exDialogVisible.value = false;
}
</script>
