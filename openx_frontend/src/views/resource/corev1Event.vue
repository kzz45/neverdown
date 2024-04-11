<template>
  <el-menu
    :default-active="activeIndex"
    mode="horizontal"
    @select="handleSelect"
  >
    <el-menu-item index="1">MetaData</el-menu-item>
    <el-menu-item index="2">Message</el-menu-item>
  </el-menu>
  <div class="menu-item" v-show="activeIndex === '1'">
    <MetaData :metadata="itemInfo.metadata"></MetaData>
  </div>
  <div class="menu-item" v-show="activeIndex === '2'">
    <div class="spec-group">
      <div class="spec-group-item">
        <div class="spec-label">type</div>
        <div class="spec-value" style="width: 400px">
          {{ props.itemInfo.type }}
        </div>
      </div>
      <div class="spec-group-item">
        <div class="spec-label">reason</div>
        <div class="spec-value" style="width: 400px">
          {{ props.itemInfo.reason }}
        </div>
      </div>
    </div>
    <div class="spec-group">
      <div class="spec-group-item">
        <div class="spec-label">source component</div>
        <div class="spec-value" style="width: 400px">
          {{ props.itemInfo.source.component }}
        </div>
      </div>
      <div class="spec-group-item">
        <div class="spec-label">source host</div>
        <div class="spec-value" style="width: 400px">
          {{ props.itemInfo.source.host }}
        </div>
      </div>
    </div>
    <div class="involved">involvedObject:</div>
    <InvolvedObject :metadata="itemInfo.involvedObject"></InvolvedObject>
    <div class="involved">
      message:
      <div class="message-text">
        {{ props.itemInfo.message }}
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@import "./css/spec.scss";
.spec-group {
  margin-top: 20px;
  display: flex;
  justify-content: space-around;
  align-items: center;
  .spec-group-item {
    .spec-label {
      font-size: 1rem;
      font-weight: 400;
      width: 150px;
      text-align: left;
    }
    display: flex;
    gap: 20px;
    align-items: center;
  }
}
.involved {
  width: 100%;
  text-align: left;
  margin-top: 30px;
  padding-top: 10px;
  border-top: 1px solid #ccc;
}
.message-text {
  display: block;
  margin-top: 10px;
  padding: 10px;
  background: #272822;
  width: 100%;
  height: 400px;
  overflow: auto;
  color: white;
  cursor: text;
  direction: ltr;
  font-family: monospace;
  font-size: 13px;
  font-variant-ligatures: contextual;
  font-weight: 600;
}
</style>

<script setup lang="ts">
import MetaData from "./components/metadata.vue";
import InvolvedObject from "./components/objectReference.vue";
import { ref } from "vue";
const props = defineProps<{
  itemInfo?: any;
}>();

const activeIndex = ref("2");
const handleSelect = (key: string, keyPath: string[]) => {
  activeIndex.value = keyPath[0] || "1";
};
</script>
