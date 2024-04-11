<template>
  <el-menu
    :default-active="activeIndex"
    mode="horizontal"
    @select="handleSelect"
  >
    <el-menu-item index="1">MetaData</el-menu-item>
    <el-menu-item index="2">ClusterRole</el-menu-item>
  </el-menu>

  <div class="menu-item" v-show="activeIndex === '1'">
    <MetaData :metadata="itemInfo.metadata"></MetaData>
  </div>

  <div class="menu-item" v-show="activeIndex === '2'">
    <el-table
      size="small"
      :data="props.itemInfo.rules"
      border
      style="width: 100%"
      highlight-current-row
    >
      <el-table-column
        v-for="v in roleCloumn"
        :key="v.field"
        :label="v.title"
        :width="v.width"
      >
        <template #default="scope">
          <span v-if="scope.row.isSet">
            <div class="tag-group">
              <div
                class="label-tag"
                v-for="(mode, index) in scope.row[v.field]"
                :key="mode"
              >
                {{ mode }}
                <el-popconfirm
                  title="确定删除?"
                  @confirm="tagClose(scope.row[v.field], index)"
                >
                  <template #reference>
                    <el-icon><Close /></el-icon>
                  </template>
                </el-popconfirm>
              </div>
              <TagInput
                @value-input="
                  (input) => handleInputModeConfirm(input, scope.row[v.field])
                "
                :blank="true"
              />
            </div>
            <!-- <el-input v-model="scope.row[v.field]" size="mini" placeholder="请输入内容" /> -->
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
              subsetsEdit(scope.row, scope.row.isSet, props.itemInfo.rules)
            "
          >
            {{ scope.row.isSet ? "保存" : "修改" }}
          </span>
          <span
            v-if="!scope.row.isSet"
            class="el-tag el-tag--danger el-tag--mini"
            @click="rowDelete(props.itemInfo.rules, scope.$index)"
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
import TagInput from "./components/taginput.vue";
import { rowEdit, rowDelete, rowCancelEdit } from "./components/tabelUtil";
import { cloneDeep } from "lodash";

const props = defineProps<{
  itemInfo?: any;
}>();

const roleCloumn = [
  { field: "apiGroups", title: "apiGroups", width: 300 },
  { field: "resources", title: "resources", width: 300 },
  { field: "verbs", title: "verbs", width: 300 },
];
function subsetsAdd() {
  const subset = {
    apiGroups: [],
    resources: [],
    verbs: ["*"],
  };
  props.itemInfo.rules.push(subset);
}

let saveEditData = ref({});
function subsetsEdit(row: any, isSet: boolean, allData: any) {
  rowEdit(row, isSet, allData, saveEditData);
}
function cancelSubsetsEdit(row: any) {
  rowCancelEdit(row, saveEditData);
}

function tagClose(obj: any, index: number) {
  obj.splice(index, 1);
}
function handleInputModeConfirm(inputStr: string, obj: any) {
  obj.push(inputStr);
}

const activeIndex = ref("1");
const handleSelect = (key: string, keyPath: string[]) => {
  activeIndex.value = keyPath[0] || "1";
};
</script>
