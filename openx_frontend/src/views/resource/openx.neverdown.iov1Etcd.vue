<template>
  <el-menu
    :default-active="activeIndex"
    mode="horizontal"
    @select="handleSelect"
  >
    <el-menu-item index="1">MetaData</el-menu-item>
    <el-menu-item index="2">Spec</el-menu-item>
    <el-menu-item index="3">Pod</el-menu-item>
  </el-menu>

  <div class="menu-item" v-show="activeIndex === '1'">
    <MetaData :metadata="itemInfo.metadata"></MetaData>
  </div>

  <div class="menu-item" v-show="activeIndex === '2'">
    <div class="spec-item">
      <div class="spec-label">Replicas</div>
      <div class="spec-value">
        <el-input v-model="props.itemInfo.spec.replicas" size="small" />
      </div>
    </div>
    <div class="spec-item">
      <div class="spec-label">persistentStorage:</div>
      <div class="spec-value">
        <el-input
          v-model="props.itemInfo.spec.persistentStorage.storageVolumePath"
          size="small"
          placeholder="Please input Name"
        />
      </div>
    </div>
  </div>

  <div class="menu-item" style="padding-top: 20px" v-show="activeIndex === '3'">
    <TemplateMeta :poddata="props.itemInfo.spec.pod"></TemplateMeta>
  </div>
</template>

<style lang="scss" scoped>
@import "./css/spec.scss";
</style>

<script setup lang="ts">
import { ref } from "vue";
import { ElMessage } from "element-plus";
import MetaData from "./components/metadata.vue";
import TemplateMeta from "./components/templatemeta.vue";
import { cloneDeep } from "lodash";

const props = defineProps<{
  itemInfo?: any;
}>();

const activeIndex = ref("1");
const handleSelect = (key: string, keyPath: string[]) => {
  activeIndex.value = keyPath[0] || "1";
};
</script>
