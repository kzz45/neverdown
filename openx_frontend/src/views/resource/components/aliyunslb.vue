<template>
  <div class="alb-item">
    <div class="alb-label">loadBalancerId</div>
    <div class="alb-value">
      <el-select
        v-model="aliyunSLB.loadBalancerId"
        placeholder="Select"
        style="width: 100%"
        clearable
      >
        <el-option
          v-for="val in albList"
          :key="val.metadata.name"
          :label="val.metadata.name"
          :value="val.metadata.name"
        />
      </el-select>
    </div>
  </div>

  <div class="alb-status">
    <div class="alb-item">
      <div class="alb-label">overrideListeners</div>
      <div class="alb-value">
        <el-select
          v-model="aliyunSLB.overrideListeners"
          placeholder="Select"
          style="width: 100%"
          clearable
        >
          <el-option label="true" :value="true" />
          <el-option label="false" :value="false" />
        </el-select>
      </div>
    </div>
  </div>
  <div class="alb-item alb-inline">
    <div class="alb-label">
      accessControlId
      <div class="alb-value">
        <el-select
          v-model="aliyunSLB.accessControlId"
          placeholder="Select"
          style="width: 100%"
          clearable
        >
          <el-option
            v-for="val in aacList"
            :key="val.metadata.name"
            :label="val.metadata.name"
            :value="val.metadata.name"
          />
        </el-select>
      </div>
    </div>
    <div class="alb-label">
      status
      <div class="alb-value">
        <el-select
          v-model="aliyunSLB.status"
          placeholder="Select"
          style="width: 100%"
          clearable
        >
          <el-option label="on" value="on" />
          <el-option label="off" value="off" />
        </el-select>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.alb-item {
  width: 100%;
}
.alb-inline {
  display: flex;
  justify-content: flex-start;
  gap: 30px;
}
</style>

<script setup lang="ts">
import { ref, watch } from "vue";
import { useStore } from "@/store";
import { useRoute } from "vue-router";
import { initSocketData, sendSocketMessage } from "@/api/socket";
import { returnResource } from "../util";
const route = useRoute();

const store = useStore();
const props = defineProps<{
  aliyunSLB?: any;
}>();

let getList = function (gvk: string) {
  const nsGvk = route.path.split("/");
  const senddata = initSocketData("Request", nsGvk[1], gvk, "list");
  sendSocketMessage(senddata, store);
};

getList("openx.neverdown.io-v1-AliyunLoadBalancer");
getList("openx.neverdown.io-v1-AliyunAccessControl");

watch(
  () => store.state.socket.socket.message,
  (msg) => {
    const requestList = [
      "openx.neverdown.io-v1-AliyunLoadBalancer",
      "openx.neverdown.io-v1-AliyunAccessControl",
    ];
    for (let requestType of requestList) {
      initMessage(msg, requestType);
    }
  }
);

let albList = ref([]);
let aacList = ref([]);

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
      if (type === "openx.neverdown.io-v1-AliyunLoadBalancer") {
        albList.value = resultList;
      }
      if (type === "openx.neverdown.io-v1-AliyunAccessControl") {
        aacList.value = resultList;
      }
    }
  } catch (e) {
    console.log("error");
  }
}

let loadOver = () => {};
</script>
