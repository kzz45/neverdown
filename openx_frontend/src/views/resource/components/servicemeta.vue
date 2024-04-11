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
    <MetaData :metadata="service.metadata"></MetaData>
  </div>

  <div class="menu-item" v-show="activeIndex === '2'">
    <ServiveSpec :spec="service.spec" />
  </div>
</template>

<style lang="scss" scoped>
@import "./../css/spec.scss";
</style>

<script setup lang="ts">
import { cloneDeep } from "lodash";
import { ref } from "vue";
import MetaData from "./metadata.vue";
import ServiveSpec from "./servicespec.vue";

const props = defineProps<{
  service?: any;
}>();

const activeIndex = ref("1");
const handleSelect = (key: string, keyPath: string[]) => {
  activeIndex.value = keyPath[0] || "1";
};
</script>
