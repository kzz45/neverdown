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
      <div class="spec-label">imagePullSecrets</div>
      <div class="spec-value">
        <el-table
          size="small"
          :data="itemInfo.imagePullSecrets"
          border
          style="width: 100%"
          highlight-current-row
        >
          <el-table-column key="name" label="Name" width="200">
            <template #default="scope">
              <span v-if="scope.row.isSet">
                <el-input
                  v-model="scope.row.name"
                  size="small"
                  placeholder="请输入内容"
                />
              </span>
              <span v-else>{{ scope.row.name }}</span>
            </template>
          </el-table-column>
          <el-table-column label="操作">
            <template #default="scope">
              <span
                class="el-tag el-tag--info el-tag--mini"
                style="cursor: pointer"
                @click="
                  imagePullSecretEdit(
                    scope.row,
                    scope.row.isSet,
                    itemInfo.imagePullSecrets
                  )
                "
              >
                {{ scope.row.isSet ? "保存" : "修改" }}
              </span>
              <span
                v-if="!scope.row.isSet"
                class="el-tag el-tag--danger el-tag--mini"
                @click="rowDelete(itemInfo.imagePullSecrets, scope.$index)"
                style="cursor: pointer"
              >
                删除
              </span>
              <span
                v-else
                class="el-tag el-tag--mini"
                style="cursor: pointer"
                @click="cancelimagePullSecretEdit(scope.row)"
              >
                取消
              </span>
            </template>
          </el-table-column>
        </el-table>
        <div
          class="el-table-add-row"
          style="width: 99.2%"
          @click="imagePullSecretAdd()"
        >
          <span>+ 添加</span>
        </div>
      </div>
    </div>
    <div class="spec-item">
      <div class="spec-label">secrets</div>
      <div class="spec-value">
        <el-table
          size="small"
          :data="itemInfo.secrets"
          border
          style="width: 100%"
          highlight-current-row
        >
          <el-table-column key="name" label="Name" width="220">
            <template #default="scope">
              <span v-if="scope.row.isSet">
                <el-input
                  v-model="scope.row.name"
                  size="small"
                  placeholder="请输入内容"
                />
              </span>
              <span v-else>{{ scope.row.name }}</span>
            </template>
          </el-table-column>
          <el-table-column label="操作">
            <template #default="scope">
              <span
                class="el-tag el-tag--info el-tag--mini"
                style="cursor: pointer"
                @click="
                  imageSecretEdit(scope.row, scope.row.isSet, itemInfo.secrets)
                "
              >
                {{ scope.row.isSet ? "保存" : "修改" }}
              </span>
              <span
                v-if="!scope.row.isSet"
                class="el-tag el-tag--danger el-tag--mini"
                @click="rowDelete(itemInfo.secrets, scope.$index)"
                style="cursor: pointer"
              >
                删除
              </span>
              <span
                v-else
                class="el-tag el-tag--mini"
                style="cursor: pointer"
                @click="cancelSecretEdit(scope.row)"
              >
                取消
              </span>
            </template>
          </el-table-column>
        </el-table>
        <div
          class="el-table-add-row"
          style="width: 99.2%"
          @click="imageSecretAdd()"
        >
          <span>+ 添加</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@import "./css/spec.scss";
.meta-title {
  font-size: 16px;
  font-weight: bold;
  color: #303133;
  text-align: left;
}
.meta-item {
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
</style>

<script setup lang="ts">
import { cloneDeep } from "lodash";
import { ref } from "vue";
import MetaData from "./components/metadata.vue";
import { rowEdit, rowDelete, rowCancelEdit } from "./components/tabelUtil";

const props = defineProps<{
  itemInfo?: any;
}>();

let saveEditImageData = ref({});
let saveEditSecretData = ref({});
function imagePullSecretEdit(row: any, isSet: boolean, allData: any) {
  rowEdit(row, isSet, allData, saveEditImageData);
}
function imageSecretEdit(row: any, isSet: boolean, allData: any) {
  rowEdit(row, isSet, allData, saveEditSecretData);
}
function cancelimagePullSecretEdit(row: any) {
  rowCancelEdit(row, saveEditImageData);
}
function cancelSecretEdit(row: any) {
  rowCancelEdit(row, saveEditSecretData);
}

function imagePullSecretAdd() {
  const newPullSecret = { name: "name" };
  props.itemInfo.imagePullSecrets.push(newPullSecret);
}
function imageSecretAdd() {
  const newPullSecret = { name: "name" };
  props.itemInfo.secrets.push(newPullSecret);
}

const activeIndex = ref("1");
const handleSelect = (key: string, keyPath: string[]) => {
  activeIndex.value = keyPath[0] || "1";
};
</script>
