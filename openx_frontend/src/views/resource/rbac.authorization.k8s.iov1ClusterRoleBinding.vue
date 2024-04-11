<template>
  <el-menu
    :default-active="activeIndex"
    mode="horizontal"
    @select="handleSelect"
  >
    <el-menu-item index="1">MetaData</el-menu-item>
    <el-menu-item index="2">roleRef</el-menu-item>
    <el-menu-item index="3">subjects</el-menu-item>
  </el-menu>

  <div class="menu-item" v-show="activeIndex === '1'">
    <MetaData :metadata="itemInfo.metadata"></MetaData>
  </div>

  <div class="menu-item" v-show="activeIndex === '2'">
    <div class="spec-item">
      <div class="spec-label">apiGroup</div>
      <div class="spec-value">
        <el-input v-model="props.itemInfo.roleRef.apiGroup" size="small" />
      </div>
    </div>
    <div class="spec-item">
      <div class="spec-label">kind:</div>
      <div class="spec-value">
        <el-input
          v-model="props.itemInfo.roleRef.kind"
          size="small"
          placeholder="Please input kind"
        />
      </div>
    </div>
    <div class="spec-item">
      <div class="spec-label">name:</div>
      <div class="spec-value">
        <el-input
          v-model="props.itemInfo.roleRef.name"
          size="small"
          placeholder="Please input name"
        />
      </div>
    </div>
  </div>

  <div class="menu-item" v-show="activeIndex === '3'">
    <el-table
      size="small"
      :data="props.itemInfo.subjects"
      border
      style="width: 100%"
      highlight-current-row
    >
      <el-table-column
        v-for="v in subjectCloumn"
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
            class="el-tag el-tag--mini"
            style="cursor: pointer"
            @click="copySubsets(scope.row)"
          >
            复制
          </span>
          <span
            class="el-tag el-tag--info el-tag--mini"
            style="cursor: pointer"
            @click="
              subsetsEdit(scope.row, scope.row.isSet, props.itemInfo.subjects)
            "
          >
            {{ scope.row.isSet ? "保存" : "修改" }}
          </span>
          <span
            v-if="!scope.row.isSet"
            class="el-tag el-tag--danger el-tag--mini"
            @click="rowDelete(props.itemInfo.subjects, scope.$index)"
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
    <div class="el-table-add-row" style="width: 99.2%" @click="subsetsAdd()">
      <span>+ 添加</span>
      <span v-show="cloneRow.name" style="color: #409eff; font-weight: 400">{{
        cloneRow.name
      }}</span>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@import "./css/spec.scss";
</style>

<script setup lang="ts">
import { reactive, ref } from "vue";
import { ElMessage } from "element-plus";
import MetaData from "./components/metadata.vue";
import { rowEdit, rowDelete, rowCancelEdit } from "./components/tabelUtil";
import { cloneDeep } from "lodash";

const props = defineProps<{
  itemInfo?: any;
}>();

const subjectCloumn = [
  { field: "name", title: "name", width: 300 },
  { field: "namespace", title: "namespace", width: 150 },
  { field: "apiGroup", title: "apiGroup", width: 300 },
  { field: "kind", title: "kind", width: 150 },
];

let cloneRow = ref({
  name: "",
  namespace: "",
  apiGroup: "",
  kind: "User",
});
function copySubsets(row: any) {
  let copyItem = cloneDeep(row);
  copyItem.isSet = false;
  cloneRow.value = cloneDeep(copyItem);
}
function subsetsAdd() {
  // const subset = {
  //   name:'', namespace:'', apiGroup: '', kind: 'User'
  // }
  props.itemInfo.subjects.push(cloneDeep(cloneRow.value));
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
