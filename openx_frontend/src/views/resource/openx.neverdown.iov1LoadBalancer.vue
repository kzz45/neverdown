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
    <div class="spec-item">
      <div class="spec-label">instance</div>
      <div class="spec-value">
        <div style="display: flex">
          <el-input
            v-model="props.itemInfo.spec.instance.key"
            placeholder="key"
            size="small"
          />
          <div class="middle-dash">:</div>
          <el-input
            v-model="props.itemInfo.spec.instance.value"
            placeholder="value"
            size="small"
          />
        </div>
      </div>
    </div>
    <div class="spec-item">
      <div class="spec-label">overrideListeners</div>
      <div class="spec-value">
        <div style="display: flex">
          <el-input
            v-model="props.itemInfo.spec.overrideListeners.key"
            placeholder="key"
            size="small"
          />
          <div class="middle-dash">:</div>
          <el-input
            v-model="props.itemInfo.spec.overrideListeners.value"
            placeholder="value"
            size="small"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@import "./css/spec.scss";
</style>

<script setup lang="ts">
import { ref } from "vue";
import { ElMessage } from "element-plus";
import MetaData from "./components/metadata.vue";
import { cloneDeep } from "lodash";

const props = defineProps<{
  itemInfo?: any;
}>();

const activeIndex = ref("1");
const handleSelect = (key: string, keyPath: string[]) => {
  activeIndex.value = keyPath[0] || "1";
};
</script>
