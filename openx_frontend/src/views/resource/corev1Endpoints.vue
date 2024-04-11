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
            size="default"
            :class="
              selectedConfig === Number(tag) - 1 ? 'selected-tag' : 'config-tag'
            "
            :disable-transitions="false"
            @click="selectEditor(tag)"
          >
            <div class="tag-name">{{ tag }}</div>
            <el-icon @click.stop="handleClose(tag)"><Close /></el-icon>
          </el-tag>
          <el-button class="button-new-tag" size="small" @click="addEndpoint"
            >+ New Tag</el-button
          >
        </div>
        <div class="yaml-style">
          <div v-if="dynamicTags.length > 0">
            addresses:
            <el-table
              size="small"
              :data="initsubset('addresses')"
              border
              style="width: 100%"
              highlight-current-row
            >
              <el-table-column
                v-for="v in addressesCloumn"
                :key="v.field"
                :label="v.title"
                :width="v.width"
              >
                <template #default="scope">
                  <span v-if="scope.row.isSet">
                    <el-input
                      v-model="scope.row[v.field]"
                      size="small"
                      placeholder="请输入内容"
                    />
                  </span>
                  <span v-else>{{ scope.row[v.field] }}</span>
                </template>
              </el-table-column>
              <el-table-column label="操作">
                <template #default="scope">
                  <span
                    class="el-tag el-tag--info el-tag--mini"
                    style="cursor: pointer"
                    @click="
                      subsetsEdit(
                        scope.row,
                        scope.row.isSet,
                        initsubset('addresses')
                      )
                    "
                  >
                    {{ scope.row.isSet ? "保存" : "修改" }}
                  </span>
                  <span
                    v-if="!scope.row.isSet"
                    class="el-tag el-tag--danger el-tag--mini"
                    @click="rowDelete(initsubset('addresses'), scope.$index)"
                    style="cursor: pointer"
                  >
                    删除
                  </span>
                  <span
                    v-else
                    class="el-tag el-tag--mini"
                    style="cursor: pointer"
                    @click="cancelSubsetsEdit(scope.row)"
                  >
                    取消
                  </span>
                </template>
              </el-table-column>
            </el-table>
            <div
              class="el-table-add-row"
              style="width: 99.2%"
              @click="subsetsAdd()"
            >
              <span>+ 添加</span>
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
.yaml-style {
  width: calc(100% - 400px);
  display: inline-block;
  background: #ccc;
  margin-left: 20px;
  position: relative;
  height: 100%;
}
</style>

<script setup lang="ts">
import { computed, nextTick, onMounted, ref } from "vue";
import MetaData from "./components/metadata.vue";
import { rowEdit, rowDelete, rowCancelEdit } from "./components/tabelUtil";

const props = defineProps<{
  itemInfo?: any;
}>();

let selectedConfig = ref(0);
let dynamicTags = computed(() => {
  let tags = [];
  for (let index in props.itemInfo.subsets) {
    tags.push(Number(index) + 1);
  }
  return tags;
});
function initsubset(EndpointSubset: string) {
  return props.itemInfo.subsets[selectedConfig.value]
    ? props.itemInfo.subsets[selectedConfig.value][EndpointSubset]
    : [];
}
function selectEditor(tag: string | number) {
  selectedConfig.value = Number(tag) - 1;
}
function handleClose(tag: number) {
  if (Number(tag) - 1 >= 0) {
    selectEditor(Number(tag) - 1);
  }
  props.itemInfo.subsets.splice(tag - 1, 1);
}
function addEndpoint() {
  props.itemInfo.subsets.push({
    addresses: [],
    notReadyAddresses: [],
    ports: [],
  });
}

const addressesCloumn = [
  { field: "hostname", title: "hostname", width: 300 },
  { field: "ip", title: "hostname", width: 200 },
  { field: "nodeName", title: "nodeName", width: 200 },
];

function subsetsAdd() {
  const subset = {
    hostname: "hostname",
    ip: "0.0.0.0",
    nodeName: "name",
  };
  props.itemInfo.subsets[selectedConfig.value].addresses.push(subset);
}

let saveEditData = ref({});
function subsetsEdit(row: any, isSet: boolean, allData: any) {
  rowEdit(row, isSet, allData, saveEditData);
}
function cancelSubsetsEdit(row: any) {
  rowCancelEdit(row, saveEditData);
}

const activeIndex = ref("1");
const handleSelect = (key: string, keyPath: string[]) => {
  activeIndex.value = keyPath[0] || "1";
};
</script>
