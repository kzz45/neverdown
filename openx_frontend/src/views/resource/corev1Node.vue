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
      <div class="spec-label">externalID</div>
      <div class="spec-value">
        <el-input
          v-model="itemInfo.spec.externalID"
          size="small"
          style="width: 100%"
        />
      </div>
    </div>
    <div class="spec-item">
      <div class="spec-label">podCIDR</div>
      <div class="spec-value">
        <el-input
          v-model="itemInfo.spec.podCIDR"
          size="small"
          style="width: 100%"
        />
      </div>
    </div>
    <div class="spec-item">
      <div class="spec-label">providerID</div>
      <div class="spec-value">
        <el-input
          v-model="itemInfo.spec.providerID"
          size="small"
          style="width: 100%"
        />
      </div>
    </div>
    <div class="spec-item">
      <div class="spec-label">unschedulable</div>
      <div class="spec-value">
        <el-select
          v-model="itemInfo.spec.unschedulable"
          placeholder="Select"
          style="width: 100%"
        >
          <el-option label="false" :value="false" />
          <el-option label="true" :value="true" />
        </el-select>
      </div>
    </div>
    <div class="spec-item">
      <div class="spec-label">taints</div>
      <div class="spec-value">
        <el-table
          size="small"
          :data="itemInfo.spec.taints"
          border
          style="width: 100%"
          highlight-current-row
        >
          <el-table-column
            v-for="v in taintCloumn"
            :key="v.field"
            :label="v.title"
            :width="v.width"
          >
            <template #default="scope">
              <span v-if="scope.row.isSet">
                <el-select
                  v-if="v.title === 'effect'"
                  v-model="scope.row[v.field]"
                  placeholder="Select"
                  style="width: 100%"
                >
                  <el-option
                    v-for="types in effectOptions"
                    :key="types"
                    :label="types"
                    :value="types"
                  />
                </el-select>
                <el-input
                  v-else
                  v-model="scope.row[v.field]"
                  size="small"
                  placeholder="请输入内容"
                />
              </span>
              <span v-else>{{ scope.row[v.field] }}</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200">
            <template #default="scope">
              <span
                class="el-tag el-tag--info el-tag--mini"
                style="cursor: pointer"
                @click="
                  taintsEdit(scope.row, scope.row.isSet, itemInfo.spec.taints)
                "
              >
                {{ scope.row.isSet ? "保存" : "修改" }}
              </span>
              <span
                v-if="!scope.row.isSet"
                class="el-tag el-tag--danger el-tag--mini"
                @click="rowDelete(itemInfo.spec.taints, scope.$index)"
                style="cursor: pointer"
              >
                删除
              </span>
              <span
                v-else
                class="el-tag el-tag--mini"
                style="cursor: pointer"
                @click="cancelEdit(scope.row)"
              >
                取消
              </span>
            </template>
          </el-table-column>
        </el-table>
        <div class="el-table-add-row" style="width: 99.2%" @click="taintsAdd()">
          <span>+ 添加</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@import "./css/spec.scss";
</style>

<script setup lang="ts">
import MetaData from "./components/metadata.vue";
import { nextTick, ref } from "vue";
import { rowEdit, rowDelete, rowCancelEdit } from "./components/tabelUtil";

const effectOptions = ["NoSchedule", "PreferNoSchedule", "NoExecute"];
const taintCloumn = [
  { field: "effect", title: "effect", width: 150 },
  { field: "key", title: "key", width: 200 },
  { field: "value", title: "value" },
];

const props = defineProps<{
  itemInfo?: any;
}>();

let savePortsData: any = ref({});
function taintsEdit(row: any, isSet: boolean, allData: any) {
  rowEdit(row, isSet, allData, savePortsData);
}
function cancelEdit(row: any) {
  rowCancelEdit(row, savePortsData);
}
function taintsAdd() {
  const tempTain = {
    effect: "NoSchedule",
    key: "key",
    value: "",
  };
  props.itemInfo.spec.taints.push(tempTain);
}

const activeIndex = ref("1");
const handleSelect = (key: string, keyPath: string[]) => {
  activeIndex.value = keyPath[0] || "1";
};
</script>
