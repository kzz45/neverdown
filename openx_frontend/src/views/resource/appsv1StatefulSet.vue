<template>
  <el-menu
    :default-active="activeIndex"
    mode="horizontal"
    @select="handleSelect"
  >
    <el-menu-item index="1">MetaData</el-menu-item>
    <el-menu-item index="2">Spec</el-menu-item>
    <el-menu-item index="3">PodTemplate</el-menu-item>
  </el-menu>
  <div v-show="activeIndex === '1'">
    <MetaData :metadata="itemInfo.metadata"></MetaData>
  </div>
  <div class="spec-item" v-show="activeIndex === '2'">
    <div class="spec-title">Spec</div>
    <div class="spec-item">
      <div class="spec-label">Replicas</div>
      <div class="spec-value">
        <el-input v-model="props.itemInfo.spec.replicas" size="small" />
      </div>
    </div>
    <div class="spec-item">
      <div class="spec-label">serviceName</div>
      <div class="spec-value">
        <el-select
          v-model="props.itemInfo.spec.serviceName"
          placeholder="Select"
          style="width: 100%"
        >
          <el-option
            v-for="service in ServiceList"
            :key="service.metadata.name"
            :label="service.metadata.name"
            :value="service.metadata.name"
          />
        </el-select>
      </div>
    </div>
    <div class="spec-label">Selector</div>
    <div class="spec-value">
      <div class="meta-label">Matchlabels</div>
      <div class="meta-value">
        <div class="tag-group">
          <el-tooltip
            effect="dark"
            v-for="(anno, key) in initLabels(
              props.itemInfo.spec.selector.matchLabels
            )"
            v-bind:key="key"
            :content="showLabel(anno)"
            placement="top-end"
          >
            <div class="label-tag">
              {{ fetchLabel(anno) }}
              <el-icon @click="tagLabelClose(anno.label)"><Close /></el-icon>
            </div>
          </el-tooltip>
        </div>
        <el-button
          size="small"
          @click="addLabels('matchlabels')"
          style="margin-top: 5px"
        >
          + add Matchlabel
        </el-button>
      </div>
    </div>
  </div>
  <div class="menu-item" style="padding-top: 20px" v-show="activeIndex === '3'">
    <TemplateMeta :poddata="itemInfo.spec.template"></TemplateMeta>
  </div>

  <el-dialog v-model="showLabelAdd" width="30%" append-to-body>
    <div style="margin-bottom: 20px">
      key:
      <el-input
        v-model="addTag.key"
        size="small"
        placeholder="Please input key"
      />
    </div>
    <div style="margin-bottom: 20px">
      value:
      <el-input
        v-model="addTag.value"
        size="small"
        placeholder="Please input value"
      />
    </div>
    <el-button size="small" @click="confirmAddTag"> 确定 </el-button>
  </el-dialog>
</template>

<style lang="scss" scoped>
@import "./css/spec.scss";
</style>

<script setup lang="ts">
import { ref, watch } from "vue";
import { initLabels, returnResource } from "./util";
import MetaData from "./components/metadata.vue";
import TemplateMeta from "./components/templatemeta.vue";
const props = defineProps<{
  itemInfo?: any;
}>();
function fetchLabel(label: any) {
  return `${label.label} : ${label.value}`;
}
function tagLabelClose(labelKey: string) {
  delete props.itemInfo.spec.selector.matchLabels[labelKey];
}
let showLabelAdd = ref(false);
let addTag = ref({
  key: "",
  value: "",
});
let addTitle = ref("");
function addLabels(title: string) {
  addTitle.value = title;
  showLabelAdd.value = true;
}
function confirmAddTag() {
  const addKey = addTag.value.key;
  const addValue = addTag.value.value;
  if (addKey) {
    showLabelAdd.value = false;
    if (addTitle.value === "matchlabels") {
      props.itemInfo.spec.selector.matchLabels[addKey] = addValue;
    }
  }
}
function showLabel(anno: any) {
  return `${anno.label} : ${anno.value}`;
}

import { useRoute } from "vue-router";
import { useStore } from "@/store";
import { initSocketData, sendSocketMessage } from "@/api/socket";
const route = useRoute();
const store = useStore();

let ServiceList = ref([]);
let getList = function (gvk: string) {
  const nsGvk = route.path.split("/");
  const senddata = initSocketData("Request", nsGvk[1], gvk, "list");
  sendSocketMessage(senddata, store);
};
getList("core-v1-Service");
function initMessage(msg: any, type: string) {
  const nsGvk = route.path.split("/");
  const gvkArr = type.split("-");
  let gvkObj = {
    group: gvkArr[0],
    version: gvkArr[1],
    kind: gvkArr[2],
  };
  try {
    let resultList = returnResource(msg, nsGvk[1], gvkObj, loadOver);
    if (resultList) {
      resultList.sort((itemL: any, itemR: any) => {
        const itemLTime = itemL.metadata.creationTimestamp.seconds;
        const itemRTime = itemR.metadata.creationTimestamp.seconds;
        return itemRTime - itemLTime;
      });
      if (type === "core-v1-Service") {
        ServiceList.value = resultList;
      }
    }
  } catch (e) {
    console.log("error");
  }
}
watch(
  () => store.state.socket.socket.message,
  (msg) => {
    initMessage(msg, "core-v1-Service");
  }
);
let loadOver = function () {};

const activeIndex = ref("1");
const handleSelect = (key: string, keyPath: string[]) => {
  activeIndex.value = keyPath[0] || "1";
};
</script>
