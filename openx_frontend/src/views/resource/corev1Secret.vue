<template>
  <el-menu
    :default-active="activeIndex"
    mode="horizontal"
    @select="handleSelect"
  >
    <el-menu-item index="1">MetaData</el-menu-item>
    <el-menu-item index="2">Type</el-menu-item>
  </el-menu>

  <div class="menu-item" v-show="activeIndex === '1'">
    <MetaData :metadata="itemInfo.metadata"></MetaData>
  </div>
  <div class="menu-item" v-if="activeIndex === '2'">
    <div class="spec-item">
      <div class="spec-label">Type</div>
      <div class="spec-value">
        <el-select
          v-model="props.itemInfo.type"
          placeholder="Select"
          style="width: 100%"
        >
          <el-option
            v-for="types in typeOptions"
            :key="types"
            :label="types"
            :value="types"
          />
        </el-select>
      </div>
    </div>
    <div class="spec-item">
      <div class="spec-label">Data</div>
      <div class="spec-value">
        <div style="width: 100%">
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
                :id="tag"
                :key="tag"
                v-for="tag in dynamicTags"
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
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@import "./css/spec.scss";
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
</style>

<script setup lang="ts">
import { binaryToStr, strToBinary } from "@/api/socket";
import { initLabels } from "./util";
import YamlEditor from "@/components/YamlEditor.vue";
import MetaData from "./components/metadata.vue";
import { computed, nextTick, onMounted, ref } from "vue";

const typeOptions = [
  "Opaque",
  "kubernetes.io/service-account-token",
  "kubernetes.io/dockercfg",
  "kubernetes.io/dockerconfigjson",
  "kubernetes.io/basic-auth",
  "kubernetes.io/ssh-auth",
  "kubernetes.io/tls",
  "bootstrap.kubernetes.io/token",
];

const props = defineProps<{
  itemInfo?: any;
}>();

let selectedConfig = ref("");
let yamlData = ref("");

let dynamicTags = computed(() => {
  let tags = [];
  for (let dataKey in props.itemInfo.data) {
    tags.push(dataKey);
  }
  return tags;
});
onMounted(() => {
  let tags = [];
  for (let dataKey in props.itemInfo.data) {
    tags.push(dataKey);
  }
  selectEditor(tags[0] || "");
});

function selectEditor(tag: string) {
  selectedConfig.value = tag;
  let val = tag ? props.itemInfo.data[tag] : "";
  yamlData.value = "1";
  setTimeout(() => {
    yamlData.value = binaryToStr(val);
  });
}

let inputVisible = ref(false);
let inputValue = ref("");
let yamlEditor = $ref(null);
let saveTagInput = $ref(null);

async function showInput() {
  inputVisible.value = true;
  await nextTick();
  saveTagInput.$refs.input.focus();
}
function handleInputConfirm() {
  if (inputValue.value) {
    props.itemInfo.data[inputValue.value] = strToBinary("");
    selectEditor(inputValue.value);
  }

  yamlEditor.setValue("");
  inputVisible.value = false;
  inputValue.value = "";
}

function handleClose(tag: string) {
  delete props.itemInfo.data[tag];
  let tags = [];
  for (let dataKey in props.itemInfo.data) {
    tags.push(dataKey);
  }
  let moveToTag = tags[0] || "";
  selectEditor(moveToTag);
}

function yamlChanged(yamlValue: string) {
  if (selectedConfig.value)
    props.itemInfo.data[selectedConfig.value] = strToBinary(yamlValue);
}

const activeIndex = ref("1");
const handleSelect = (key: string, keyPath: string[]) => {
  activeIndex.value = keyPath[0] || "1";
};
</script>
