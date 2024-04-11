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
      <div class="spec-label">nodeSelector</div>
      <div class="spec-value">
        <!-- <el-tag v-for="(anno, key) in initLabels(props.itemInfo.spec.nodeSelector)" style="margin-top: 5px"
          v-bind:key="key" @close="tagLabelClose(anno.label)" closable>
          {{fetchLabel(anno)}}
        </el-tag> -->
        <div class="tag-group">
          <div
            class="label-tag"
            v-for="(anno, key) in initLabels(props.itemInfo.spec.nodeSelector)"
            v-bind:key="key"
          >
            {{ fetchLabel(anno) }}
            <el-icon @click="tagLabelClose(anno.label)"><Close /></el-icon>
          </div>
        </div>

        <el-button
          size="small"
          @click="addLabels('nodeSelector')"
          style="margin-top: 5px"
        >
          + add Matchlabel
        </el-button>
      </div>
    </div>
  </div>

  <el-dialog v-model="showLabelAdd" width="30%" append-to-body>
    <div>
      key:
      <el-input
        v-model="addTag.key"
        size="small"
        placeholder="Please input key"
      />
    </div>
    <div>
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
import { ref } from "vue";
import { ElMessage } from "element-plus";
import MetaData from "./components/metadata.vue";
import { initLabels } from "./util";

const props = defineProps<{
  itemInfo?: any;
}>();

function tagLabelClose(labelKey: string) {
  delete props.itemInfo.spec.nodeSelector[labelKey];
}
function fetchLabel(label: any) {
  return `${label.label} : ${label.value}`;
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
  if (addTitle.value === "nodeSelector") {
    props.itemInfo.spec.nodeSelector[addKey] = addValue;
  }
}

const activeIndex = ref("1");
const handleSelect = (key: string, keyPath: string[]) => {
  activeIndex.value = keyPath[0] || "1";
};
</script>
