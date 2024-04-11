<template>
  <div class="meta-grid">
    <div class="meta-item">
      <div class="meta-label">Name:</div>
      <div class="meta-value">
        <el-input
          v-model="props.metadata.name"
          @input="
            (e) => {
              props.metadata.name = e.trim();
            }
          "
          placeholder="Please input Name"
          style="width: calc(100% - 100px)"
        />
      </div>
    </div>
    <div class="meta-item">
      <div class="meta-label">Namespace:</div>
      <div class="meta-value">
        {{ props.metadata.namespace }}
      </div>
    </div>
    <div class="meta-item">
      <div class="meta-label">annotations:</div>
      <div class="meta-value">
        <div class="tag-group">
          <el-tooltip
            effect="dark"
            v-for="(anno, key) in initAnnotations(props.metadata.annotations)"
            v-bind:key="key"
            class="box-item"
            :content="showLabel(anno)"
            placement="top-end"
          >
            <div class="label-tag">
              {{ fetchLabel(anno) }}
              <el-icon @click="tagClose(anno.label)"><Close /></el-icon>
            </div>
          </el-tooltip>
        </div>
        <el-button
          size="small"
          @click="addAnnotations('annotations')"
          style="margin-top: 5px"
        >
          + add annotation
        </el-button>
      </div>
    </div>
    <div class="meta-item">
      <div class="meta-label">labels:</div>
      <div class="meta-value">
        <div class="tag-group">
          <div
            class="label-tag"
            v-for="(anno, key) in initAnnotations(props.metadata.labels)"
            v-bind:key="key"
          >
            {{ showLabel(anno) }}
            <el-icon @click="tagLabelClose(anno.label)"><Close /></el-icon>
          </div>
        </div>

        <el-button
          size="small"
          @click="addAnnotations('labels')"
          style="margin-top: 5px"
        >
          + add label
        </el-button>
      </div>
    </div>
    <div class="meta-item" v-if="editType === 'update'">
      <div class="meta-label">resourceVersion:</div>
      <div class="meta-value">
        {{ props.metadata.resourceVersion }}
      </div>
    </div>
    <div
      class="meta-item"
      v-if="editType === 'update' && props.metadata.creationTimestamp"
    >
      <div class="meta-label">creationTime:</div>
      <div class="meta-value">
        {{ formatTime(props.metadata.creationTimestamp.seconds) }}
      </div>
    </div>
  </div>
  <el-dialog v-model="showAnnotations" width="30%" append-to-body>
    <div style="margin-bottom: 20px">
      key:
      <el-input
        v-model="addTag.key"
        size="small"
        placeholder="Please input key"
      />
    </div>
    <div style="margin-bottom: 20px">
      value:
      <el-input
        v-model="addTag.value"
        size="small"
        placeholder="Please input value"
      />
    </div>
    <el-button size="small" @click="confirmAddTag"> 确定 </el-button>
  </el-dialog>
</template>

<style lang="scss" scoped>
.tag-group {
  display: flex;
  justify-content: flex-start;
  flex-wrap: wrap;
  .label-tag {
    margin-top: 5px;
    margin-bottom: 3px;
    background-color: #f1f3f4;
    border: 1px solid #ccc;
    margin-right: 10px;
    border-radius: 4px;
    padding: 3px 3px 3px 10px;
    display: flex;
    justify-content: space-between;
    font-weight: 500;
    i {
      margin-left: 10px;
      cursor: pointer;
      &:hover {
        color: red;
      }
    }
  }
}

.meta-title {
  font-size: 16px;
  font-weight: bold;
  color: #303133;
  text-align: left;
}
.meta-grid {
  display: grid;
  margin-top: 10px;
  grid-template-columns: 50% 50%;
  .meta-item {
    display: flex;
    margin-bottom: 8px;
    .meta-label {
      width: 150px;
      line-height: 2rem;
      text-align: left;
      font-weight: 500;
    }
    .meta-value {
      width: 80%;
      padding-right: 30px;
      padding-left: 10px;
      text-align: left;
      line-height: 2rem;
      font-weight: 400;
      display: flex;
      flex-direction: column;
      align-items: flex-start;
    }
  }
}
</style>

<script setup lang="ts">
import { watch } from "fs";
import { cloneDeep } from "lodash";
import { inject, onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import { formatTime } from "./../util";

const route = useRoute();

const props = defineProps<{
  metadata?: any;
}>();

const showInfoList = [
  "name",
  "namespace",
  "annotations",
  "generation",
  "resourceVersion",
  "uid",
];
const editType = inject("editType");

onMounted(() => {
  const nsGvk = route.path.split("/");
  const namespace: string = nsGvk[1] || "";
  props.metadata.namespace = namespace;
});

function initAnnotations(ans: any) {
  let tags = [];
  for (let key in ans) {
    tags.push({
      label: key,
      value: ans[key],
    });
  }
  return tags;
}
function tagClose(tagKey: string) {
  delete props.metadata.annotations[tagKey];
}
function tagLabelClose(tagKey: string) {
  delete props.metadata.labels[tagKey];
}
function fetchAnno(anno: any) {
  if (anno.label.length < 34) {
    const fetchKey = anno.label;
    const fetchValue =
      anno.value.length < 5 ? anno.value : anno.value.slice(0, 5) + "...";
    return `${fetchKey} : ${fetchValue}`;
  } else {
    const fetchKey = anno.label.slice(0, 26) + "...";
    const fetchValue =
      anno.value.length < 5 ? anno.value : anno.value.slice(0, 5) + "...";
    return `${fetchKey} : ${fetchValue}`;
  }
}
function fetchLabel(anno: any) {
  if (anno.label.length < 20) {
    const fetchKey = anno.label;
    const fetchValue =
      anno.value.length < 20 ? anno.value : anno.value.slice(0, 15) + "...";
    return `${fetchKey} : ${fetchValue}`;
  } else {
    const fetchKey = anno.label.slice(0, 18) + "...";
    const fetchValue =
      anno.value.length < 20 ? anno.value : anno.value.slice(0, 15) + "...";
    return `${fetchKey} : ${fetchValue}`;
  }
}
function showLabel(anno: any) {
  return `${anno.label} : ${anno.value}`;
}
let showAnnotations = ref(false);
let addTag = ref({
  key: "",
  value: "",
});
let addTitle = ref("");
function addAnnotations(title: string) {
  addTitle.value = title;
  showAnnotations.value = true;
}
function confirmAddTag() {
  const addKey = addTag.value.key;
  const addValue = addTag.value.value;
  if (addKey) {
    showAnnotations.value = false;
  }
  if (addTitle.value === "annotations") {
    props.metadata.annotations[addKey] = addValue;
  } else {
    props.metadata.labels[addKey] = addValue;
  }
}
</script>
