<template>
  <el-menu
    :default-active="activeIndex"
    mode="horizontal"
    @select="handleSelect"
  >
    <el-menu-item index="1">MetaData</el-menu-item>
    <el-menu-item index="2">Finalizers</el-menu-item>
  </el-menu>

  <div class="menu-item" v-show="activeIndex === '1'">
    <MetaData :metadata="itemInfo.metadata"></MetaData>
  </div>

  <div class="menu-item" v-show="activeIndex === '2'">
    <div class="flex-box">
      <el-tag
        :id="tag"
        :key="tag"
        v-for="tag in itemInfo.spec.finalizers"
        size="default"
        class="config-tag"
        :disable-transitions="false"
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
      <el-button v-else class="button-new-tag" size="small" @click="showInput"
        >+ New Tag</el-button
      >
    </div>
  </div>
</template>

<style lang="scss" scoped>
:deep(.el-tag__content) {
  display: flex;
  justify-content: space-between;
  width: 100%;
}
.flex-box {
  margin-top: 20px;
  display: flex;
  justify-content: flex-start;
  flex-wrap: wrap;
  overflow-x: hidden;
}
.config-tag {
  width: 340px;
  cursor: pointer;
  font-size: 1rem;
  height: 40px;
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
  margin-right: 20px;
  background-color: #2f3136f0;
  color: white;
  border-radius: 5px;
}
.config-tag:hover {
  background-color: #2f3136d8;
}
.button-new-tag {
  width: 340px;
  height: 40px;
}
</style>

<script setup lang="ts">
import MetaData from "./components/metadata.vue";
import { nextTick, ref } from "vue";

const props = defineProps<{
  itemInfo?: any;
}>();

let inputVisible = ref();
let inputValue = ref("");
let saveTagInput = $ref(null);

async function showInput() {
  inputVisible.value = true;
  await nextTick();
  saveTagInput.$refs.input.focus();
}
function handleInputConfirm() {
  if (inputValue.value) {
    props.itemInfo.spec.finalizers.push(inputValue.value);
  }
  inputVisible.value = false;
  inputValue.value = "";
}
function handleClose(tag: string) {
  const closeIndex = props.itemInfo.spec.finalizers.findIndex(
    (finalizer: any) => {
      return finalizer === tag;
    }
  );
  if (closeIndex >= 0) {
    props.itemInfo.spec.finalizers.splice(closeIndex, 1);
  }
}

const activeIndex = ref("1");
const handleSelect = (key: string, keyPath: string[]) => {
  activeIndex.value = keyPath[0] || "1";
};
</script>
