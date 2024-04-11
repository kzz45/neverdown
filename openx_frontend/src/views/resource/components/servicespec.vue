<template>
  <div class="spec-item">
    <div class="spec-label">clusterIP</div>
    <div class="spec-value">
      <el-input v-model="props.spec.clusterIP" size="small" />
    </div>
  </div>
  <div class="spec-item">
    <div class="spec-label">type</div>
    <div class="spec-value">
      <el-select
        v-model="props.spec.type"
        placeholder="Select"
        style="width: 100%"
        clearable
      >
        <el-option
          v-for="types in typeOptions"
          :key="types"
          :label="types"
          :value="types"
        />
      </el-select>
    </div>
  </div>
  <div class="spec-item">
    <div class="spec-label">selector</div>
    <div class="spec-value">
      <div class="tag-group">
        <div
          class="label-tag"
          v-for="(anno, key) in initLabels(props.spec.selector)"
          v-bind:key="key"
        >
          {{ fetchLabel(anno) }}
          <el-icon @click="tagLabelClose(anno.label)"><Close /></el-icon>
        </div>
      </div>
      <el-button
        size="small"
        @click="addLabels('selector')"
        style="margin-top: 5px"
      >
        + add Matchlabel
      </el-button>
    </div>
  </div>
  <div class="spec-item">
    <div class="spec-label">ports</div>
    <div class="spec-value" v-if="props.spec && props.spec.ports">
      <el-table
        size="small"
        :data="props.spec.ports"
        border
        style="width: 100%"
        highlight-current-row
      >
        <el-table-column
          v-for="v in portsCloumn"
          :key="v.field"
          :label="v.title"
          :width="v.width"
        >
          <template #default="scope">
            <div v-if="Object.keys(scope.row).length > 0">
              <span v-if="scope.row.isSet">
                <el-input
                  v-model="scope.row[v.field]"
                  size="small"
                  placeholder="请输入内容"
                />
              </span>
              <span v-else>{{ scope.row[v.field] }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="targetPort">
          <el-table-column label="type">
            <template #default="scope">
              <div v-if="scope.row?.targetPort">
                <span v-if="scope.row.isSet">
                  <el-input
                    v-model="scope.row.targetPort.type"
                    size="small"
                    placeholder="请输入内容"
                  />
                </span>
                <span v-else>{{ scope.row.targetPort.type }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="intVal">
            <template #default="scope">
              <div v-if="scope.row?.targetPort">
                <span v-if="scope.row.isSet">
                  <el-input
                    v-model="scope.row.targetPort.intVal"
                    size="small"
                    placeholder="请输入内容"
                  />
                </span>
                <span v-else>{{ scope.row.targetPort.intVal }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="strVal">
            <template #default="scope">
              <div v-if="scope.row?.targetPort">
                <span v-if="scope.row.isSet">
                  <el-input
                    v-model="scope.row.targetPort.strVal"
                    size="small"
                    placeholder="请输入内容"
                  />
                </span>
                <span v-else>{{ scope.row.targetPort.strVal }}</span>
              </div>
            </template>
          </el-table-column>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template #default="scope">
            <div v-if="Object.keys(scope.row).length > 0">
              <span
                class="el-tag el-tag--info el-tag--mini"
                style="cursor: pointer"
                @click="
                  containerPortEdit(
                    scope.row,
                    scope.row.isSet,
                    props.spec.ports
                  )
                "
              >
                {{ scope.row.isSet ? "保存" : "修改" }}
              </span>
              <span
                v-if="!scope.row.isSet"
                class="el-tag el-tag--danger el-tag--mini"
                @click="rowDelete(props.spec.ports, scope.$index)"
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
            </div>
          </template>
        </el-table-column>
      </el-table>
      <div
        class="el-table-add-row"
        style="width: 99.2%"
        @click="containerPortAdd()"
      >
        <span>+ 添加</span>
      </div>
    </div>
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
@import "./../css/spec.scss";
</style>

<script setup lang="ts">
import { cloneDeep } from "lodash";
import { initLabels, showLabel } from "../util";
import { onMounted, ref } from "vue";
import { rowEdit, rowDelete, rowCancelEdit } from "./tabelUtil";

const props = defineProps<{
  spec?: any;
}>();

onMounted(() => {
  if (props.spec && props.spec.ports) {
    for (let port of props.spec.ports) {
      if (port.appProtocol === "") {
        delete port.appProtocol;
      }
    }
  }
});

const typeOptions = ["ClusterIP", "NodePort", "LoadBalancer", "ExternalName"];
const portsCloumn = [
  { field: "name", title: "name", width: 150 },
  { field: "nodePort", title: "nodePort", width: 120 },
  { field: "port", title: "port", width: 120 },
  { field: "protocol", title: "protocol", width: 120 },
];

function fetchLabel(label: any) {
  return `${label.label} : ${label.value}`;
}
function tagLabelClose(labelKey: string) {
  delete props.spec.selector[labelKey];
}
let addTitle = ref("");
let showLabelAdd = ref(false);
let addTag = ref({
  key: "",
  value: "",
});
function addLabels(title: string) {
  addTitle.value = title;
  showLabelAdd.value = true;
}
function confirmAddTag() {
  const addKey = addTag.value.key;
  const addValue = addTag.value.value;
  if (addKey) {
    showLabelAdd.value = false;
    if (addTitle.value === "selector") {
      props.spec.selector[addKey] = addValue;
    }
  }
}
let savePortsData: any = ref({});
function containerPortEdit(row: any, isSet: boolean, allData: any) {
  rowEdit(row, isSet, allData, savePortsData);
  targetPortEdit.value = savePortsData.value.targetPort
    ? JSON.stringify(savePortsData.value.targetPort)
    : "";
}
function cancelEdit(row: any) {
  rowCancelEdit(row, savePortsData);
}
function containerPortAdd() {
  const newPort = {
    name: "",
    nodePort: 0,
    port: 3306,
    protocol: "TCP",
    targetPort: { type: 0, intVal: 3306, strVal: "" },
  };
  props.spec.ports.push(newPort);
}
let targetPortEdit = ref("");
function portChange(row: any) {
  row.targetPort = JSON.parse(targetPortEdit.value);
}
</script>
