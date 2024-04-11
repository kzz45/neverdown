<template>
  <el-menu
    :default-active="activeIndex"
    mode="horizontal"
    @select="handleSelect"
  >
    <el-menu-item index="1">MetaData</el-menu-item>
    <el-menu-item index="2">Spec</el-menu-item>
    <el-menu-item index="3">PodTemplate</el-menu-item>
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
      <div class="spec-label">Selector</div>
      <div class="spec-value">
        <div class="meta-label">Matchlabels</div>
        <div class="meta-value">
          <div class="tag-group">
            <div
              class="label-tag"
              v-for="(anno, key) in initLabels(
                props.itemInfo.spec.selector.matchLabels
              )"
              v-bind:key="key"
            >
              {{ fetchLabel(anno) }}
              <el-icon @click="tagLabelClose(anno.label)"><Close /></el-icon>
            </div>
          </div>
          <el-button
            size="small"
            @click="addLabels('matchlabels')"
            style="margin-top: 5px"
          >
            + add Matchlabel
          </el-button>
        </div>
      </div>
    </div>
  </div>
  <div class="menu-item" style="padding-top: 20px" v-show="activeIndex === '3'">
    <TemplateMeta :poddata="itemInfo.spec.template"></TemplateMeta>
  </div>
  <el-dialog v-model="showLabelAdd" width="30%" append-to-body>
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
@import "./css/spec.scss";
</style>

<script setup lang="ts">
import { initLabels, showLabel } from "./util";
import MetaData from "./components/metadata.vue";
import TemplateMeta from "./components/templatemeta.vue";
import { ref } from "vue";
const props = defineProps<{
  itemInfo?: any;
}>();

function fetchLabel(label: any) {
  return `${label.label} : ${label.value}`;
}
function tagLabelClose(labelKey: string) {
  delete props.itemInfo.spec.selector.matchLabels[labelKey];
}

let showLabelAdd = ref(false);
let addTag = ref({
  key: "",
  value: "",
});
let addTitle = ref("");
function addLabels(title: string) {
  addTitle.value = title;
  showLabelAdd.value = true;
}
function confirmAddTag() {
  const addKey = addTag.value.key;
  const addValue = addTag.value.value;
  if (addKey) {
    showLabelAdd.value = false;
    if (addTitle.value === "matchlabels") {
      props.itemInfo.spec.selector.matchLabels[addKey] = addValue;
    }
  }
}

const activeIndex = ref("1");
const handleSelect = (key: string, keyPath: string[]) => {
  activeIndex.value = keyPath[0] || "1";
};
</script>
