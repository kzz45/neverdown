<template>
  <div class="mysql-menu">
    <el-menu
      active-text-color="#409eff"
      mode="horizontal"
      background-color="#fff"
      @select="handleSelect"
      class="el-menu-vertical-demo"
      default-active="1"
      text-color="#333"
    >
      <el-menu-item index="1">
        <span>spec</span>
      </el-menu-item>
      <el-menu-item index="2">
        <span>Template</span>
      </el-menu-item>
      <el-menu-item index="3">
        <span>service</span>
      </el-menu-item>
    </el-menu>
  </div>
  <div class="mysql-spec">
    <div class="mysql-content">
      <div v-show="muneIndex === '1'">
        <div class="meta-item">
          <div class="meta-label">replicas:</div>
          <div class="meta-value">
            <el-input-number
              v-model="props.mysqlspec.replicas"
              :min="0"
              :max="30"
            />
          </div>
        </div>

        <div class="meta-item">
          <div class="meta-label">role:</div>
          <div class="meta-value">
            <el-input
              v-model="props.mysqlspec.role"
              size="small"
              placeholder="Please input Name"
            />
          </div>
        </div>

        <div class="meta-item">
          <div class="meta-label">cloudNetworkConfig:</div>
          <div class="meta-value">
            <AliyunSLB
              :aliyunSLB="props.mysqlspec.cloudNetworkConfig.aliyunSLB"
            />
          </div>
        </div>

        <div class="meta-item">
          <div class="meta-label">persistentStorage:</div>
          <div class="meta-value">
            <el-input
              v-model="props.mysqlspec.persistentStorage.storageVolumePath"
              size="small"
              placeholder="Please input Name"
            />
          </div>
        </div>
      </div>
      <div v-show="muneIndex === '2'" style="padding-top: 20px">
        <TemplateMeta :poddata="props.mysqlspec.pod"></TemplateMeta>
      </div>
      <div v-show="muneIndex === '3'">
        <ServiveSpec :spec="props.mysqlspec.service.spec" />
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@import "../css/mysql.scss";
</style>

<script setup lang="ts">
import { cloneDeep } from "lodash";
import TemplateMeta from "./templatemeta.vue";
import ServiveSpec from "./servicespec.vue";
import AliyunSLB from "./aliyunslb.vue";
import { ref, watch } from "vue";
import { formatTime } from "./../util";

const props = defineProps<{
  mysqlspec?: any;
}>();

if (!props.mysqlspec.cloudNetworkConfig.aliyunSLB) {
  props.mysqlspec.cloudNetworkConfig = {
    aliyunSLB: {
      loadBalancerId: "",
      accessControlId: "",
      overrideListeners: true,
      status: "off",
    },
  };
}

let muneIndex = ref("1");
function handleSelect(gole: any) {
  muneIndex.value = gole;
}
</script>
