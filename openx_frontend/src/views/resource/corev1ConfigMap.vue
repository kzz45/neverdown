<template>
  <el-menu
    :default-active="activeIndex"
    mode="horizontal"
    @select="handleSelect"
  >
    <el-menu-item index="1">MetaData</el-menu-item>
    <el-menu-item index="2">Data</el-menu-item>
    <el-menu-item index="3">BinaryData</el-menu-item>
  </el-menu>
  <div class="menu-item" v-show="activeIndex === '1'">
    <MetaData :metadata="itemInfo.metadata"></MetaData>
  </div>

  <div class="menu-item" v-if="activeIndex === '2'">
    <div style="width: 100%">
      <div style="display: flex; flex-wrap: wrap">
        <div style="width: 340px; display: inline-block; overflow-x: hidden">
          <el-tag
            :id="tag"
            :key="tag"
            v-for="tag in dynamicTags"
            size="medium"
            :class="selectedConfig === tag ? 'selected-tag' : 'config-tag'"
            :disable-transitions="false"
            @click="selectEditor(tag, 'data')"
          >
            <div class="tag-name">{{ tag }}</div>
            <el-icon @click.stop="handleClose(tag, 'data')"><Close /></el-icon>
            <!-- <i style="line-height: 40px; cursor: pointer" @click.stop="handleClose(tag)" class="el-icon-delete"></i> -->
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
            @click="showInput('data')"
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

  <div class="menu-item" v-if="activeIndex === '3'">
    <div style="width: 100%">
      <div style="display: flex; flex-wrap: wrap">
        <div style="width: 340px; display: inline-block; overflow-x: hidden">
          <el-tag
            :id="tag"
            :key="tag"
            v-for="tag in binaryDataTag"
            size="medium"
            :class="
              selectedConfigBinary === tag ? 'selected-tag' : 'config-tag'
            "
            :disable-transitions="false"
            @click="selectEditor(tag, 'binaryData')"
          >
            <div class="tag-name">{{ tag }}</div>
            <el-icon @click.stop="handleClose(tag, 'binary')"
              ><Close
            /></el-icon>
            <!-- <i style="line-height: 40px; cursor: pointer" @click.stop="handleClose(tag)" class="el-icon-delete"></i> -->
          </el-tag>
          <el-input
            class="input-new-tag"
            v-if="inputVisibleBinary"
            v-model="inputValueBinary"
            ref="saveTagInputBinary"
            @keyup.enter="$event.target.blur()"
            @blur="handleInputConfirmBinary"
          >
          </el-input>
          <el-button
            v-else
            class="button-new-tag"
            size="small"
            @click="showInput('binaryData')"
            >+ New Tag</el-button
          >
        </div>
        <div class="yaml-editor yaml-style" style="">
          <yaml-editor
            ref="yamlEditorBinary"
            :config="yamlbinaryData"
            @changed="yamlChangedBinary"
          />
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
import { initLabels } from "./util";
import YamlEditor from "@/components/YamlEditor.vue";
import MetaData from "./components/metadata.vue";
import { computed, nextTick, onMounted, ref } from "vue";
import { binaryToStr, strToBinary } from "@/api/socket";

const props = defineProps<{
  itemInfo?: any;
}>();

let selectedConfig = ref("");
let yamlData = ref("");

let selectedConfigBinary = ref("");
let yamlbinaryData = ref("");

let dynamicTags = computed(() => {
  let tags = [];
  for (let dataKey in props.itemInfo.data) {
    tags.push(dataKey);
  }
  return tags;
});

let binaryDataTag = computed(() => {
  let tags = [];
  for (let dataKey in props.itemInfo.binaryData) {
    tags.push(dataKey);
  }
  return tags;
});

let yamlEditor = $ref(null);
let saveTagInput = $ref(null);

let yamlEditorBinary = $ref(null);
let saveTagInputBinary = $ref(null);

initYaml();
function initYaml() {
  let tags = [];
  for (let dataKey in props.itemInfo.data) {
    tags.push(dataKey);
  }
  selectEditor(tags[0] || "", "data");

  let tagsBinary = [];
  for (let dataKey in props.itemInfo.binaryData) {
    tagsBinary.push(dataKey);
  }
  selectEditor(tagsBinary[0] || "", "binaryData");
}

// onMounted(() => {
//   let tags = []
//   for(let dataKey in props.itemInfo.data){
//     tags.push(dataKey)
//   }
//   selectEditor(tags[0] || '')
// })

function selectEditor(tag: string, objType?: string) {
  if (objType === "data") {
    selectedConfig.value = tag;
    let val = tag ? props.itemInfo.data[tag] : "";
    yamlData.value = "1";
    nextTick(() => {
      yamlData.value = val;
    });
  } else {
    selectedConfigBinary.value = tag;
    let val = tag ? props.itemInfo.binaryData[tag] : "";
    yamlbinaryData.value = "1";
    nextTick(() => {
      yamlbinaryData.value = binaryToStr(val);
    });
  }
}

let inputVisible = ref(false);
let inputValue = ref("");

let inputVisibleBinary = ref(false);
let inputValueBinary = ref("");

async function showInput(type: string) {
  if (type === "data") {
    inputVisible.value = true;
    await nextTick();
    saveTagInput.$refs.input.focus();
  } else {
    inputVisibleBinary.value = true;
    await nextTick();
    saveTagInputBinary.$refs.input.focus();
  }
}
function handleInputConfirm() {
  if (inputValue.value) {
    props.itemInfo.data[inputValue.value] = "";
    selectEditor(inputValue.value, "data");
  }

  yamlEditor.setValue("");
  inputVisible.value = false;
  inputValue.value = "";
}
function handleInputConfirmBinary() {
  if (inputValueBinary.value) {
    props.itemInfo.binaryData[inputValueBinary.value] = "";
    selectEditor(inputValueBinary.value, "binaryData");
  }

  yamlEditorBinary.setValue("");
  inputVisibleBinary.value = false;
  inputValueBinary.value = "";
}

function handleClose(tag: string, tagType: string) {
  console.log("asdas");
  let tags = [];
  if (tagType === "data") {
    delete props.itemInfo.data[tag];
    for (let dataKey in props.itemInfo.data) {
      tags.push(dataKey);
    }
  } else {
    delete props.itemInfo.binaryData[tag];
    for (let dataKey in props.itemInfo.binaryData) {
      tags.push(dataKey);
    }
  }

  let moveToTag = tags[0] || "";
  selectEditor(moveToTag, tagType);
}

function yamlChanged(yamlValue: string) {
  if (selectedConfig.value)
    props.itemInfo.data[selectedConfig.value] = yamlValue;
}
function yamlChangedBinary(yamlValue: string) {
  if (selectedConfigBinary.value)
    props.itemInfo.binaryData[selectedConfigBinary.value] =
      strToBinary(yamlValue);
}

const activeIndex = ref("1");
const handleSelect = (key: string, keyPath: string[]) => {
  activeIndex.value = keyPath[0] || "1";
};
</script>
