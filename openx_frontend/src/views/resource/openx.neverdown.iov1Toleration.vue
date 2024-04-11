<template>
  <el-menu
    :default-active="activeIndex"
    mode="horizontal"
    @select="handleSelect"
  >
    <el-menu-item index="1">MetaData</el-menu-item>
    <el-menu-item index="2">Spec</el-menu-item>
  </el-menu>
  <div class="menu-item" v-show="activeIndex === '1'">
    <MetaData :metadata="itemInfo.metadata"></MetaData>
  </div>
  <div class="menu-item" v-show="activeIndex === '2'">
    <Toleration :toleration="props.itemInfo.spec.toleration" />
  </div>
</template>

<style lang="scss" scoped>
@import "./css/spec.scss";
</style>

<script setup lang="ts">
import { ref } from "vue";
import { ElMessage } from "element-plus";
import MetaData from "./components/metadata.vue";
import Toleration from "./components/toleration.vue";
import { cloneDeep } from "lodash";

const props = defineProps<{
  itemInfo?: any;
}>();

const activeIndex = ref("1");
const handleSelect = (key: string, keyPath: string[]) => {
  activeIndex.value = keyPath[0] || "1";
};
</script>
