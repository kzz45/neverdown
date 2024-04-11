<template>
  <div class="spec-item">
    <div class="spec-item">
      <div class="spec-value">
        nodeSelectorTerms
        <el-tabs
          v-model="editableTabsValue"
          v-if="affinity.nodeAffinity"
          type="card"
          closable
          @tab-remove="removeTab"
        >
          <el-tab-pane
            v-for="(item, index) in affinity.nodeAffinity
              .requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms"
            :key="index"
            :label="`Term${index}`"
            :name="`Term${index}`"
          >
            <el-table
              size="small"
              :data="item.matchExpressions"
              border
              style="width: 100%"
              highlight-current-row
            >
              <el-table-column
                v-for="v in affinityColumn"
                :key="v.field"
                :label="v.title"
                :width="v.width"
              >
                <template #default="scope">
                  <span v-if="scope.row.isSet">
                    <el-input
                      v-if="v.field === 'key'"
                      v-model="scope.row['key']"
                      size="small"
                      placeholder="请输入内容"
                    />
                    <el-input
                      v-else-if="v.field === 'values'"
                      v-model="scope.row[v.field]"
                      @change="
                        (e) => {
                          inputChange(e, scope.row);
                        }
                      "
                      size="small"
                      placeholder="请输入内容"
                    />
                    <el-select
                      v-else-if="v.field === 'operator'"
                      v-model="scope.row['operator']"
                    >
                      <el-option
                        v-for="affVal in affinityList"
                        :key="affVal"
                        :label="affVal"
                        :value="affVal"
                      />
                    </el-select>
                  </span>
                  <span v-else>{{ scope.row[v.field] }}</span>
                </template>
              </el-table-column>

              <!-- <el-table-column label="操作" :width="160">
              <template #default="scope">
                <span class="el-tag el-tag--info el-tag--mini" style="cursor: pointer;" @click="affEdit(scope.row, scope.row.isSet, item.matchExpressions)">
                  {{ scope.row.isSet?'保存':"修改" }}
                </span>
                <span v-if="!scope.row.isSet" class="el-tag el-tag--danger el-tag--mini" style="cursor: pointer;" @click="rowDelete(item.matchExpressions, scope.$index)">
                  删除
                </span>
                <span v-else class="el-tag  el-tag--mini" style="cursor: pointer;" @click="cancelAffEdit(scope.row)">
                  取消
                </span>
              </template>
            </el-table-column> -->
            </el-table>
            <!-- <div class="el-table-add-row" style="width: 99.2%;" @click="affAdd(item.matchExpressions)"><span>+ 添加</span></div> -->
          </el-tab-pane>
        </el-tabs>
      </div>
    </div>
    <!-- <div class="spec-item">
      <div class="spec-label">podAffinity</div>
      <div class="spec-value">
        podAffinity
      </div>
    </div>
    <div class="spec-item">
      <div class="spec-label">podAntiAffinity</div>
      <div class="spec-value">
        podAntiAffinity
      </div>
    </div> -->
  </div>
</template>

<style lang="scss" scoped>
@import "./../css/spec.scss";
</style>

<script setup lang="ts">
import { rowEdit, rowDelete, rowCancelEdit } from "./tabelUtil";
import { ElMessage } from "element-plus";
import { cloneDeep } from "lodash";
import { ref } from "vue";

const props = defineProps<{
  affinity?: any;
}>();

const affinityList = ["In", "NotIn", "Exists", "DoesNotExist", "Gt", "Lt"];

const affinityColumn = [
  { field: "key", title: "key", width: 160 },
  { field: "values", title: "values (,)分割", width: 250 },
  { field: "operator", title: "operator", width: 200 },
];
let editableTabsValue = ref("Term0");
let cacheAffData = ref();
function removeTab(termIndex: string) {
  const removeIndex = Number(termIndex.replace("Term", ""));
  props.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms.splice(
    removeIndex,
    1
  );
}

let saveAffinityData = ref({});
function affEdit(row: any, isSet: boolean, allData: any) {
  rowEdit(row, isSet, allData, saveAffinityData);
}
function cancelAffEdit(row: any) {
  rowCancelEdit(row, saveAffinityData);
}

function inputChange(e: any, row: any) {
  row.values = e.split(",");
}

function affAdd(exp: any) {
  for (let port of exp) {
    if (port.isSet) {
      ElMessage.error("请先保存");
      return;
    }
  }
  let tempPort = { key: "key", values: "value", operator: "In" };
  exp.push(tempPort);
}
</script>
