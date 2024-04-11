<template>
  <div>
    Image :
    <div class="edit-image">
      Project:
      <el-select
        v-model="imageProject"
        class="m-2"
        placeholder="Select"
        @change="projectChange"
      >
        <el-option
          v-for="item in ProjectList"
          :key="item.metadata.name"
          :label="item.metadata.name"
          :value="item.metadata.name"
        />
      </el-select>
      Tag:
      <el-input
        v-model="imageTag"
        @change="tagChange"
        style="width: 150px"
      ></el-input>
    </div>
    WatchPolicy: <span class="edit-tips">{{ policyTips }}</span>
    <div class="edit-image">
      <el-select
        v-model="watchPolicy"
        class="m-2"
        placeholder="Select"
        @change="policyChange"
        :clearable="true"
      >
        <el-option
          v-for="Policy in watchPolicyList"
          :key="Policy"
          :label="Policy"
          :value="Policy"
        />
      </el-select>
    </div>
    <div>
      <el-radio-group
        v-model="batchEditTypeEdit"
        size="small"
        @change="changeEditType"
      >
        <el-radio label="add" border>添加</el-radio>
        <el-radio label="delete" border>删除</el-radio>
      </el-radio-group>
    </div>
    Pod Annotations:<span class="edit-tips">{{ annotationTips }}</span>
    <div
      v-for="dom in props.variables"
      v-bind:key="dom.index"
      class="variables-style"
    >
      <input
        type="text"
        v-model="dom.key"
        class="login-input"
        placeholder=""
      />:
      <input
        type="text"
        v-model="dom.value"
        class="login-input"
        placeholder=""
        style="width: 212px"
      />
      <el-button v-if="plusFlg(dom)" :icon="Plus" @click="addDom" circle />
      <el-button v-else :icon="Minus" @click="delDom(dom.index)" circle />
    </div>
    <div class="vol-line"></div>
    Volumes:<span class="edit-tips">{{ volumeTips }}</span>
    <el-button size="small" @click="addVolume">添加</el-button>
    <div
      class="volume"
      v-for="(volume, volIndex) in volumes"
      v-bind:key="volume.name"
    >
      <div class="volume-title">
        <div>Name: {{ volume.name }}</div>
        <div v-if="!volume.volumeSource.secret">
          <el-button size="small" @click="editVolume(volume, volIndex)"
            >编辑</el-button
          >
          <el-popconfirm title="确定删除?" @confirm="deleteVolume(volIndex)">
            <template #reference>
              <el-button size="small" type="danger">删除</el-button>
            </template>
          </el-popconfirm>
        </div>
      </div>
      <div v-if="volume.volumeSource.configMap">
        <div style="text-align: left">
          configMap :
          {{ volume.volumeSource.configMap.localObjectReference.name }}
        </div>
        <el-table
          :data="volume.volumeSource.configMap.items"
          size="small"
          border
          style="width: 100%"
        >
          <el-table-column prop="key" label="key" width="180" />
          <el-table-column prop="path" label="path" />
          <el-table-column label="action" width="180"> </el-table-column>
        </el-table>
      </div>
      <div v-if="volume.volumeSource.hostPath" style="text-align: left">
        <div>path: {{ volume.volumeSource.hostPath.path }}</div>
        <div>type: {{ volume.volumeSource.hostPath.type }}</div>
      </div>
      <div v-if="volume.volumeSource.secret" style="text-align: left">
        <div>secretName: {{ volume.volumeSource.secret.secretName }}</div>
      </div>
    </div>
    <div class="vol-line"></div>
    VolumeMounts:<span class="edit-tips">{{ volumeMountsTips }}</span>
    <div>
      <el-table
        size="small"
        :data="props.VolumeMounts"
        border
        style="width: 100%"
        highlight-current-row
      >
        <el-table-column
          v-for="v in volumeMountsCloumn"
          :key="v.field"
          :label="v.title"
          :width="v.width"
        >
          <template #default="scope">
            <span v-if="scope.row.isSet">
              <el-input
                v-model="scope.row[v.field]"
                size="small"
                :formatter="(value: string) => `${value}`.replace(/\s/g, '')"
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
                volumeMountsEdit(scope.row, scope.row.isSet, props.VolumeMounts)
              "
            >
              {{ scope.row.isSet ? "保存" : "修改" }}
            </span>
            <span
              v-if="!scope.row.isSet"
              class="el-tag el-tag--danger el-tag--mini"
              @click="rowDelete(props.VolumeMounts, scope.$index)"
              style="cursor: pointer"
            >
              删除
            </span>
            <span
              v-else
              class="el-tag el-tag--mini"
              style="cursor: pointer"
              @click="cancelVolumeMountsEdit(scope.row)"
            >
              取消
            </span>
          </template>
        </el-table-column>
      </el-table>
      <div
        class="el-table-add-row"
        style="width: 99.2%"
        @click="volumeMountsAdd()"
      >
        <span>+ 添加</span>
      </div>
    </div>
    <div class="vol-line"></div>
    Env:<span class="edit-tips">{{ volumeTips }}</span
    ><el-button size="small" @click="importEnv">导入</el-button>
    <div>
      <el-table
        size="small"
        :data="props.env"
        border
        style="width: 100%"
        highlight-current-row
      >
        <el-table-column key="name" label="name" :width="240">
          <template #default="scope">
            <span v-if="scope.row.isSet">
              <el-input
                v-model="scope.row.name"
                size="small"
                :formatter="(value: string) => `${value}`.replace(/\s/g, '')"
                placeholder="请输入内容"
              />
            </span>
            <span
              v-else
              @click="envEdit(scope.row, scope.row.isSet, props.env)"
              >{{ scope.row.name }}</span
            >
          </template>
        </el-table-column>
        <el-table-column key="value" label="value" :width="230">
          <template #default="scope">
            <span v-if="scope.row.isSet">
              <el-input
                v-model="scope.row.value"
                size="small"
                :formatter="(value: string) => `${value}`.replace(/\s/g, '')"
                placeholder="请输入内容"
              />
            </span>
            <span
              v-else
              @click="envEdit(scope.row, scope.row.isSet, props.env)"
              >{{ scope.row.value }}</span
            >
          </template>
        </el-table-column>
        <el-table-column key="fieldPath" label="fieldPath" :width="230">
          <template #default="scope">
            <div v-if="scope.row.valueFrom">
              <!-- <span v-if="scope.row.isSet">
                <el-input v-model="scope.row.valueFrom.fieldRef.fieldPath" size="mini" placeholder="请输入内容" />
              </span> -->
              <span>{{ scope.row.valueFrom.fieldRef.fieldPath }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <span
              class="el-tag el-tag--info el-tag--mini"
              style="cursor: pointer"
              @click="envEdit(scope.row, scope.row.isSet, props.env)"
            >
              {{ scope.row.isSet ? "保存" : "修改" }}
            </span>
            <span
              v-if="!scope.row.isSet"
              class="el-tag el-tag--danger el-tag--mini"
              @click="rowDelete(props.env, scope.$index)"
              style="cursor: pointer"
            >
              删除
            </span>
            <span
              v-else
              class="el-tag el-tag--mini"
              style="cursor: pointer"
              @click="cancelEnvEdit(scope.row)"
            >
              取消
            </span>
          </template>
        </el-table-column>
      </el-table>
      <div class="el-table-add-row" style="width: 99.2%" @click="envAdd()">
        <span>+ 添加</span>
      </div>
      <div class="vol-line"></div>
      <div>
        Affinity: <span class="edit-tips">替换，未选择时不修改</span>
        <div>
          <el-select
            v-model="affSelected"
            clearable
            style="width: 80"
            @change="affChange"
          >
            <el-option
              v-for="aff in props.affList"
              :key="aff.metadata.name"
              :label="aff.metadata.name"
              :value="aff.metadata.name"
            />
          </el-select>
        </div>
      </div>
    </div>

    <el-dialog v-model="showVolume" width="60%" append-to-body destroy-on-close>
      <div>
        Name:
        <el-input v-model="vlomeEditName" size="small" />
      </div>
      <el-radio-group
        v-model="vlomeType"
        class="ml-4"
        :disabled="vlomeEditType === 'edit'"
      >
        <el-radio label="configmap">configMap</el-radio>
        <el-radio label="hostpath">hostpath</el-radio>
      </el-radio-group>
      <div v-if="vlomeType === 'configmap'">
        <div>
          configMap:
          <el-select
            v-model="configMap.name"
            style="width: 100%"
            placeholder="Select"
            size="small"
          >
            <el-option
              v-for="cfgMap in ConfigMapList"
              :key="cfgMap.metadata.name"
              :label="cfgMap.metadata.name"
              :value="cfgMap.metadata.name"
            />
          </el-select>
          <!-- <el-input v-model="configMap.name" size="small" /> -->
        </div>
        <div>
          items:
          <el-table
            :data="configMap.items"
            size="small"
            border
            style="width: 100%"
          >
            <el-table-column
              v-for="v in configMapItemsCloumn"
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
            <el-table-column label="action">
              <template #default="scope">
                <span
                  class="el-tag el-tag--info el-tag--mini"
                  style="cursor: pointer"
                  @click="
                    configmapItemEdit(
                      scope.row,
                      scope.row.isSet,
                      configMap.items
                    )
                  "
                >
                  {{ scope.row.isSet ? "保存" : "修改" }}
                </span>
                <span
                  v-if="!scope.row.isSet"
                  class="el-tag el-tag--danger el-tag--mini"
                  @click="rowDelete(configMap.items, scope.$index)"
                  style="cursor: pointer"
                >
                  删除
                </span>
                <span
                  v-else
                  class="el-tag el-tag--mini"
                  style="cursor: pointer"
                  @click="cancelconfigmapItemEdit(scope.row)"
                >
                  取消
                </span>
              </template>
            </el-table-column>
          </el-table>
          <div
            class="el-table-add-row"
            style="width: 99.2%"
            @click="vlomeConfigmapAdd()"
          >
            <span>+ 添加</span>
          </div>
        </div>
      </div>
      <div v-if="vlomeType === 'hostpath'">
        <div>
          path:
          <el-input v-model="hostPath.path" size="small" />
        </div>
        <div>
          type:
          <el-select v-model="hostPath.type" style="margin-top: 10px">
            <el-option
              v-for="type in hostPathType"
              :key="type"
              :label="type"
              :value="type"
            />
          </el-select>
        </div>
      </div>
      <el-button style="margin-top: 10px" size="small" @click="confirmAddVlome">
        确定
      </el-button>
    </el-dialog>
  </div>
</template>
<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { rowEdit, rowDelete, rowCancelEdit } from "./tabelUtil";

/**
Annotations
 */
import { Plus, Minus } from "@element-plus/icons-vue";

import { initSocketData, sendSocketMessage } from "@/api/socket";
import { useStore } from "@/store";
import { useRouter, useRoute } from "vue-router";
import { returnResource } from "./../util";

const hostPathType = [
  "DirectoryOrCreate",
  "Directory",
  "FileOrCreate",
  "File",
  "Socket",
  "CharDevice",
  "BlockDevice",
];

const watchPolicyList = ["manual", "in-place-upgrade", "rolling-upgrade"];

const props: any = defineProps({
  variables: { type: Array, default: [] },
  volumes: { type: Array, default: [] },
  VolumeMounts: { type: Array, default: [] },
  env: { type: Array, default: [] },
  batchEditType: { type: String, default: "add" },
  ConfigMapList: { type: Array, default: [] },
  affList: { type: Array, default: [] },
});

let imageProject = ref("");
let watchPolicy = ref("");
let imageTag = ref("");
let vlomeEditType = ref("create");
let editVlomeIndex = ref(0);

function deleteVolume(volIndex: number) {
  props.volumes.splice(volIndex, 1);
}
function editVolume(volume: any, volIndex: number) {
  editVlomeIndex.value = volIndex;
  if (volume.volumeSource.configMap) {
    editVolumeCfgMap(volume);
  }
  if (volume.volumeSource.hostPath) {
    editVolumeHostPath(volume);
  }
  if (volume.volumeSource.secret) {
    editVolumeSecret(volume);
  }
}
const configMapItemsCloumn = [
  { field: "key", title: "key", width: 180 },
  { field: "path", title: "path", width: 180 },
];
let saveEditData = ref({});
function configmapItemEdit(row: any, isSet: boolean, allData: any) {
  rowEdit(row, isSet, allData, saveEditData);
}
function cancelconfigmapItemEdit(row: any) {
  rowCancelEdit(row, saveEditData);
}
let showVolume = ref(false);
let vlomeType = ref("configmap");
let vlomeEditName = ref("");
let hostPath = ref({
  path: "",
  type: "",
});
interface cmKeyPath {
  key: string;
  path: string;
}
interface configMapType {
  name: string;
  items: cmKeyPath[];
}
let configMap = ref(<configMapType>{
  name: "",
  items: [],
});
let secret = ref("");
function editVolumeCfgMap(volume: any) {
  vlomeEditType.value = "edit";
  vlomeType.value = "configmap";
  vlomeEditName.value = volume.name;
  configMap.value.name =
    volume.volumeSource.configMap.localObjectReference.name;
  configMap.value.items = volume.volumeSource.configMap.items;
  showVolume.value = true;
}
function editVolumeHostPath(volume: any) {
  vlomeEditType.value = "edit";
  vlomeType.value = "hostpath";
  vlomeEditName.value = volume.name;
  hostPath.value.path = volume.volumeSource.hostPath.path;
  hostPath.value.type = volume.volumeSource.hostPath.type;
  showVolume.value = true;
}
function editVolumeSecret(volume: any) {
  vlomeEditType.value = "edit";
  vlomeType.value = "secret";
  vlomeEditName.value = volume.name;
  secret.value = volume.volumeSource.secret.secretName;
  showVolume.value = true;
}
function addVolume() {
  vlomeEditType.value = "create";
  vlomeType.value = "configmap";
  vlomeEditName.value = "";
  configMap.value.name = "";
  configMap.value.items = [];
  showVolume.value = true;
}
function vlomeConfigmapAdd() {
  const oneConfigMapItem: cmKeyPath = {
    key: "key",
    path: "path",
  };
  configMap.value.items.push(oneConfigMapItem);
}
function confirmAddVlome() {
  if (!vlomeEditName.value) {
    return;
  }
  let vloItem: any = {};
  if (vlomeType.value === "configmap") {
    vloItem = {
      name: vlomeEditName.value,
      volumeSource: {
        configMap: {
          items: configMap.value.items,
          localObjectReference: {
            name: configMap.value.name,
          },
        },
      },
    };
    if (vlomeEditType.value === "edit") {
      props.volumes[editVlomeIndex.value] = vloItem;
    } else {
      props.volumes.push(vloItem);
    }
  } else if (vlomeType.value === "hostpath") {
    vloItem = {
      name: vlomeEditName.value,
      volumeSource: {
        hostPath: {
          path: hostPath.value.path,
          type: hostPath.value.type,
        },
      },
    };
    if (vlomeEditType.value === "edit") {
      props.volumes[editVlomeIndex.value] = vloItem;
    } else {
      props.volumes.push(vloItem);
    }
  } else {
    vloItem = {
      name: vlomeEditName.value,
      volumeSource: {
        secret: {
          secretName: secret.value,
        },
      },
    };
    if (vlomeEditType.value === "edit") {
      props.volumes[editVlomeIndex.value] = vloItem;
    } else {
      props.volumes.push(vloItem);
    }
  }
  showVolume.value = false;
}

let volumeMountsTips = computed(() => {
  if (props.batchEditType === "add") {
    return "追加，mountPath重复时覆盖";
  }
  if (props.batchEditType === "delete") {
    return "寻找以下mountPath并删除";
  }
});
/***
env
 */
function importEnv() {
  if (localStorage.getItem("envStorage")) {
    const envStorage = JSON.parse(String(localStorage.getItem("envStorage")));
    for (let env of envStorage) {
      const propEnvFind = props.env.findIndex((propEnv: any) => {
        return propEnv.name === env.name;
      });
      if (propEnvFind < 0) {
        props.env.push(env);
      } else {
        props.env[propEnvFind] = env;
      }
    }
  }
}
function envAdd() {
  const newEnv = {
    name: "name",
    value: "value",
  };
  props.env.push(newEnv);
}
let saveEnvData = ref({});
function envEdit(row: any, isSet: boolean, allData: any) {
  rowEdit(row, isSet, allData, saveEnvData);
}
function cancelEnvEdit(row: any) {
  rowCancelEdit(row, saveEnvData);
}

/***
volumeMounts
 */

const volumeMountsCloumn = [
  { field: "name", title: "name", width: 150 },
  { field: "mountPath", title: "mountPath", width: 150 },
  { field: "readOnly", title: "readOnly", width: 150 },
  { field: "subPath", title: "subPath", width: 150 },
  { field: "subPathExpr", title: "subPathExpr", width: 150 },
];
function volumeMountsAdd() {
  const newVolumeMount = {
    mountPath: "/data",
    name: "name",
    readOnly: false,
    subPath: "",
    subPathExpr: "",
  };
  props.VolumeMounts.push(newVolumeMount);
}
let saveVolumeMountData = ref({});
function volumeMountsEdit(row: any, isSet: boolean, allData: any) {
  rowEdit(row, isSet, allData, saveVolumeMountData);
}
function cancelVolumeMountsEdit(row: any) {
  rowCancelEdit(row, saveVolumeMountData);
}
let volumeTips = computed(() => {
  if (props.batchEditType === "add") {
    return "追加，Name重复时覆盖";
  }
  if (props.batchEditType === "delete") {
    return "寻找以下Name并删除";
  }
});

let batchEditTypeEdit = ref("add");
// watch(() => props.batchEditType, (type) => {
//   conso
//   batchEditTypeEdit.value = type
// })
const emit = defineEmits([
  "input",
  "projectChange",
  "tagChange",
  "policyChange",
  "affChange",
]);
function changeEditType(e: any) {
  emit("input", e);
}
let policyTips = computed(() => {
  return "不选时不修改";
});
let annotationTips = computed(() => {
  if (props.batchEditType === "add") {
    return "追加，Key重复时覆盖";
  }
  if (props.batchEditType === "delete") {
    return "寻找以下Key并删除";
  }
});
onMounted(() => {
  const prometheus1 = {
    index: 1,
    key: "prometheus.io/path",
    value: "/metrics",
  };
  const prometheus2 = {
    index: 2,
    key: "prometheus.io/port",
    value: "8081",
  };
  const prometheus3 = {
    index: 3,
    key: "prometheus.io/scrape",
    value: "true",
  };
  props.variables.push(prometheus1, prometheus2, prometheus3);
});
function addDom() {
  const domLength = props.variables.length;
  props.variables.push({
    index: domLength + 1,
    key: "",
    value: "",
  });
}
function plusFlg(dom: any) {
  const endFlg = dom.index === props.variables.length;
  const oneFlg = props.variables.length === 1;
  return endFlg || oneFlg;
}
function delDom(key: any) {
  for (let vari of props.variables) {
    if (vari.index > key) {
      vari.index = Number(vari.index) - 1;
    }
  }
  props.variables.splice(Number(key) - 1, 1);
}

let optionsProject = ref([{ value: "games", label: "games" }]);

const store = useStore();
const router = useRouter();
const route = useRoute();

let ProjectList = ref(<any>[]);
let getList = function (gvk: string) {
  const nsGvk = route.path.split("/");
  const senddata = initSocketData("Request", nsGvk[1], gvk, "list");
  sendSocketMessage(senddata, store);
};
getList("jingx-v1-Project");
function initMessage(msg: any, type: string) {
  const nsGvk = route.path.split("/");
  const gvkArr = type.split("-");
  let gvkObj = {
    group: gvkArr[0],
    version: gvkArr[1],
    kind: gvkArr[2],
  };
  try {
    let resultList = returnResource(msg, nsGvk[1], gvkObj, () => {});
    if (resultList) {
      resultList.sort((itemL: any, itemR: any) => {
        const itemLTime = itemL.metadata.creationTimestamp.seconds;
        const itemRTime = itemR.metadata.creationTimestamp.seconds;
        return itemRTime - itemLTime;
      });
      if (type === "jingx-v1-Project") {
        ProjectList.value = resultList;
      }
      // if( type === 'jingx-v1-Repository'){
      //   RepositoryList.value = resultList
      // }
      // if( type === 'jingx-v1-Tag'){
      //   TagList.value = resultList
      // }
    }
  } catch (e) {
    console.log("error");
  }
}
watch(
  () => store.state.socket.socket.message,
  (msg) => {
    const requestList = [
      "jingx-v1-Project",
      "jingx-v1-Repository",
      "jingx-v1-Tag",
    ];
    for (let requestType of requestList) {
      initMessage(msg, requestType);
    }
  }
);
function projectChange() {
  emit("projectChange", imageProject.value);
}
function policyChange() {
  emit("policyChange", watchPolicy.value);
}
function tagChange() {
  emit("tagChange", imageTag.value);
}

let affSelected = ref("");
function affChange() {
  emit("affChange", affSelected.value);
}
</script>

<style lang="scss" scoped>
.volume {
  border: 1px solid #ccc;
  padding: 10px;
  margin: 5px;
  border-radius: 3px;
  font-size: 0.8rem;
  font-weight: 400;
}
.volume-title {
  display: flex;
  justify-content: space-between;
}
.edit-tips {
  font-size: 0.8rem;
  font-weight: 400;
  color: red;
}
.vol-line {
  width: 100%;
  height: 1px;
  padding-top: 10px;
  margin-bottom: 5px;
  border-bottom: 1px solid #ccc;
}
.el-table-add-row {
  margin-top: 10px;
  width: 100%;
  height: 34px;
  border: 1px dashed #c1c1cd;
  border-radius: 3px;
  cursor: pointer;
  justify-content: center;
  display: flex;
  line-height: 34px;
}
.variables-style {
  display: flex;
  justify-content: flex-start;
  gap: 10px;
  .login-input {
    width: calc(50% - 100px);
    height: 30px;
    border-radius: 6px;
    border: 1px solid #d6d6d6;
    font-size: 14px;
    line-height: 1.29412;
    font-weight: 400;
    letter-spacing: -0.021em;
    font-family: SF Pro Text, SF Pro Icons, Helvetica Neue, Helvetica, Arial,
      sans-serif;
    display: inline-block;
    box-sizing: border-box;
    vertical-align: top;
    margin-bottom: 12px;
    padding-left: 15px;
    padding-right: 15px;
    color: #333;
    text-align: left;
    background: #fff;
    background-clip: padding-box;
    &::placeholder {
      color: #a8abb2;
    }
  }
}
.edit-image {
  padding: 10px 0px;
  display: flex;
  gap: 10px;
  align-items: center;
  font-weight: 400;
  font-size: 0.9rem;
}
</style>
