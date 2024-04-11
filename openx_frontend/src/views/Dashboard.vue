<template>
  <el-container>
    <el-header class="headerSty">
      <div class="header-ico move-sty" style="float: left">
        <!-- <img style="height: 100%;" src="../assets/k8s.png" /> -->
      </div>
      <div class="ns-title">
        {{ nsValue }}
      </div>
      <div
        class="header-back"
        style="display: inline-block"
        @click="backToNamespace"
      >
        <el-icon style="margin-right: 20px"><ArrowLeftBold /></el-icon>
        Namespace
      </div>
      <div class="header-ico">
        <Avatar></Avatar>
      </div>
      <div
        class="header-sel"
        tabindex="0"
        hidefocus="0"
        @mouseenter.stop="selFocus"
        @mouseleave="leaveNamespace()"
        @click.stop="() => {}"
      >
        {{ nsValue }}
        <transition name="selbounce">
          <div
            class="ns-opt"
            v-show="showOption"
            @mouseenter="enterNamespace()"
          >
            <el-input
              ref="nsinput"
              id="id-ns-1"
              v-model="nsInputValue"
              @click.stop="() => {}"
              style="width: calc(100% - 20px); margin: 10px 10px; height: 32px"
              placeholder="Please input Namespace"
            />
            <div
              v-for="(item, index) in optionsNsList"
              :key="item.Name"
              :id="'id-ns' + index"
              tabindex="1"
              style="outline: none"
              :class="
                index === nsValueSelected
                  ? 'ns-card-selected'
                  : 'ns-card-option'
              "
              @click.stop="changeNs(item.Name)"
            >
              {{ item.Name }}
            </div>
          </div>
        </transition>
      </div>
    </el-header>

    <el-container>
      <el-aside class="side-group">
        <div
          v-for="(gvk, gvkindex) in gvkMenu"
          :key="gvkindex"
          @mouseleave="leaveMenu()"
          :class="gvk.menuOpen ? 'group-menu-open' : 'group-menu'"
          @mouseenter.stop="openMenu(gvkindex)"
        >
          <span :style="chcekKindSelected(gvk.groupVersion)">{{
            groupFormat(gvk.groupVersion)
          }}</span>
        </div>
      </el-aside>
      <el-aside
        :class="collapsed ? 'sideColl' : 'sideHide'"
        @mouseleave.stop="backKind()"
      >
        <!-- <div v-for="(gvk, gvkindex) in gvkMenu" :key="gvkindex"
          :class="gvk.menuOpen ? 'menu-style-open' : 'menu-style'">
          <div style="padding-left: 20px" @click="gvk.menuOpen = !gvk.menuOpen">
          {{ groupFormat(gvk.groupVersion) }}
          <el-icon @click="tagLabelClose(anno.label)"><ArrowDownBold v-if="gvk.menuOpen" /><ArrowRightBold v-else /></el-icon>
          </div>
          <div v-for="(kind, index) in gvk.kinds" :key="index"
            @click="goResource(gvk.groupVersion, kind)"
            :class="chcekKind(gvk.groupVersion, kind)?'menu-item-selected':'menu-item-style'">
            {{ `${kind}` }}
          </div>
        </div> -->
        <div
          v-for="(kind, index) in findKinds(gvkMenu)"
          :key="index"
          @click="goResource(kind)"
          :class="chcekKind(kind) ? 'menu-item-selected' : 'menu-item-style'"
        >
          {{ `${kind}` }}
        </div>
      </el-aside>
      <el-main class="about-main" style="position: relative">
        <div class="close-btn" @click="changeSide">
          <el-icon
            ><ArrowLeftBold v-if="collapsed" /><ArrowRightBold v-else
          /></el-icon>
        </div>
        <router-view></router-view>
      </el-main>
    </el-container>
  </el-container>
  <div class="right-container" v-if="compareBtn">
    <div class="compare-title">对比列表</div>
    <div
      class="compare-list"
      v-for="(compareObj, index) in compareList"
      :key="index"
    >
      <div style="margin: 10px 0px; font-weight: 600">
        {{ compareObj.protoType }}
      </div>
      <div
        v-for="(item, index) in compareObj.items"
        class="compare-obj"
        :key="index"
      >
        <div class="obj-title">
          {{ item.metadata.namespace }}
          <el-icon class="obj-close" @click="removeItem(item, index)"
            ><Close
          /></el-icon>
        </div>
        <div class="name">
          {{ item.metadata.name }}
        </div>
      </div>
    </div>
    <div class="compare-footer">
      <div class="footer-compare" @click="compare">对比</div>
      <div class="footer-clear" @click="clearCompare">清空</div>
    </div>
  </div>
  <div
    class="dialog-delay"
    v-show="compareDialog"
    @click="compareDialog = !compareDialog"
  ></div>
  <transition name="fade">
    <div class="compare-dialog-body" v-if="compareDialog">
      <div style="margin: 20px">
        <div class="compare-box">
          <div class="top-name">
            <div
              class="compare-info"
              v-for="(info, indexName) in comparedItems"
              :key="indexName"
            >
              <div class="name-info">
                {{ info.metadata.name }}
                <el-tag>{{ info.metadata.protoType }}</el-tag>
              </div>
              <el-tag type="success">{{ info.metadata.namespace }}</el-tag>
            </div>
          </div>
          <div class="compare-table">
            <div
              class="table-title"
              style="border-top-left-radius: 10px; border-top: 1px solid #ddd"
            >
              <div class="pod">Pod Template</div>
            </div>
            <div class="table-container">
              <div class="compare-menu">
                <div class="menu-title" :style="fetchStyle('image')">image</div>
                <div class="menu-title" :style="fetchStyle('imagePullPolicy')">
                  imagePullPolicy
                </div>
                <div class="menu-title" :style="fetchStyle('volume')">
                  volume
                </div>
                <div class="menu-title" :style="fetchStyle('volumeMounts')">
                  volumeMounts
                </div>
                <div class="menu-title" :style="fetchStyle('command')">
                  command
                </div>
                <div class="menu-title" :style="fetchStyle('args')">args</div>
                <div class="menu-title" :style="fetchStyle('port')">port</div>
                <div class="menu-title" :style="fetchStyle('env')">env</div>
              </div>
              <div
                class="table-item"
                v-for="(info, indexName) in comparedItems"
                :key="indexName"
              >
                <div class="menu-title" :style="fetchStyle('image')">
                  {{ getPod(info)?.image }}
                </div>
                <div class="menu-title" :style="fetchStyle('imagePullPolicy')">
                  {{ getPod(info)?.imagePullPolicy }}
                </div>
                <div class="menu-title" :style="fetchStyle('volume')">
                  <div
                    v-for="volume in getPod(info)?.volumes"
                    class="volume"
                    :style="compareStyle('volume', volume)?.box"
                  >
                    <div class="volume-title">
                      <div>Name: {{ volume.name }}</div>
                    </div>
                    <div v-if="volume.volumeSource.configMap">
                      <div style="text-align: left">
                        configMap :
                        {{
                          volume.volumeSource.configMap.localObjectReference
                            .name
                        }}
                      </div>
                      <el-table
                        :data="volume.volumeSource.configMap.items"
                        size="small"
                        border
                      >
                        <el-table-column label="key" width="180">
                          <template #default="scope">
                            <span
                              :style="
                                compareStyle('volume', volume, scope.row)
                                  ?.configRow
                              "
                              >{{ scope.row.key }}</span
                            >
                          </template>
                        </el-table-column>
                        <el-table-column label="path">
                          <template #default="scope">
                            <span
                              :style="
                                compareStyle('volume', volume, scope.row)
                                  ?.configRow
                              "
                              >{{ scope.row.path }}</span
                            >
                          </template>
                        </el-table-column>
                      </el-table>
                    </div>
                    <div
                      v-if="volume.volumeSource.hostPath"
                      style="text-align: left"
                    >
                      <div :style="compareStyle('volume', volume)?.path">
                        path: {{ volume.volumeSource.hostPath.path }}
                      </div>
                      <div :style="compareStyle('volume', volume)?.type">
                        type: {{ volume.volumeSource.hostPath.type }}
                      </div>
                    </div>
                    <div
                      v-if="volume.volumeSource.secret"
                      style="text-align: left"
                    >
                      <div :style="compareStyle('volume', volume)?.secretName">
                        secretName: {{ volume.volumeSource.secret.secretName }}
                      </div>
                    </div>
                  </div>
                </div>
                <div class="menu-title" :style="fetchStyle('volumeMounts')">
                  <el-table
                    size="small"
                    :data="getPod(info)?.volumeMounts"
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
                        <span
                          :style="
                            compareStyle('volumeMounts', {}, scope.row)
                              ?.volumeMounts
                          "
                          >{{ scope.row[v.field] }}</span
                        >
                      </template>
                    </el-table-column>
                  </el-table>
                </div>
                <div class="menu-title" :style="fetchStyle('command')">
                  <div class="tag-group">
                    <div
                      class="label-tag"
                      :key="command"
                      v-for="command in getPod(info)?.command"
                      :style="compareStyle('command', {}, command)?.command"
                    >
                      {{ command }}
                    </div>
                  </div>
                </div>
                <div class="menu-title" :style="fetchStyle('args')">
                  <div class="tag-group">
                    <div
                      class="label-tag"
                      :key="arg"
                      v-for="arg in getPod(info)?.args"
                      :style="compareStyle('args', {}, arg)?.args"
                    >
                      {{ arg }}
                    </div>
                  </div>
                </div>
                <div class="menu-title" :style="fetchStyle('port')">
                  <el-table
                    size="small"
                    :data="getPod(info)?.ports"
                    border
                    style="width: 100%"
                    highlight-current-row
                  >
                    <el-table-column
                      v-for="v in containCloumn"
                      :key="v.field"
                      :label="v.title"
                      :width="v.width"
                    >
                      <template #default="scope">
                        <span
                          :style="compareStyle('port', {}, scope.row)?.port"
                          >{{ scope.row[v.field] }}</span
                        >
                      </template>
                    </el-table-column>
                  </el-table>
                </div>
                <div class="menu-title" :style="fetchStyle('env')">
                  <el-table
                    size="small"
                    :data="getPod(info)?.env"
                    border
                    style="width: 100%"
                    highlight-current-row
                  >
                    <el-table-column key="name" label="name" :width="100">
                      <template #default="scope">
                        <span
                          :style="compareStyle('env', {}, scope.row)?.env"
                          >{{ scope.row.name }}</span
                        >
                      </template>
                    </el-table-column>
                    <el-table-column key="value" label="value" :width="180">
                      <template #default="scope">
                        <span
                          :style="compareStyle('env', {}, scope.row)?.env"
                          >{{ scope.row.value }}</span
                        >
                      </template>
                    </el-table-column>
                    <el-table-column
                      key="fieldPath"
                      label="fieldPath"
                      :width="100"
                    >
                      <template #default="scope">
                        <div
                          v-if="scope.row.valueFrom"
                          :style="compareStyle('env', {}, scope.row)?.env"
                        >
                          <span>{{
                            scope.row.valueFrom.fieldRef.fieldPath
                          }}</span>
                        </div>
                      </template>
                    </el-table-column>
                  </el-table>
                </div>
              </div>
            </div>
            <div class="table-title">
              <div class="pod">Service</div>
            </div>
            <div class="table-container">
              <div class="compare-menu">
                <div
                  class="menu-title"
                  :style="fetchSerStyle('clusterIP', 'service')"
                >
                  clusterIP
                </div>
                <div
                  class="menu-title"
                  :style="fetchSerStyle('type', 'service')"
                >
                  type
                </div>
                <div
                  class="menu-title"
                  :style="fetchSerStyle('ports', 'service')"
                >
                  ports
                </div>
              </div>
              <div
                class="table-item"
                v-for="(info, indexName) in comparedItems"
                :key="indexName"
              >
                <div
                  class="menu-title"
                  :style="fetchSerStyle('clusterIP', 'service')"
                >
                  {{ getService(info, "service")?.clusterIP }}
                </div>
                <div
                  class="menu-title"
                  :style="fetchSerStyle('type', 'service')"
                >
                  {{ getService(info, "service")?.type }}
                </div>
                <div
                  class="menu-title"
                  :style="fetchSerStyle('ports', 'service')"
                >
                  <el-table
                    size="small"
                    :data="getService(info, 'service')?.ports"
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
                        <div
                          v-if="Object.keys(scope.row).length > 0"
                          :style="compareService(scope.row, 'service')"
                        >
                          <span>{{ scope.row[v.field] }}</span>
                        </div>
                      </template>
                    </el-table-column>
                    <el-table-column label="targetPort">
                      <el-table-column label="type" width="50">
                        <template #default="scope">
                          <div
                            v-if="scope.row?.targetPort"
                            :style="compareService(scope.row, 'service')"
                          >
                            <span>{{ scope.row.targetPort.type }}</span>
                          </div>
                        </template>
                      </el-table-column>
                      <el-table-column label="intVal" width="60">
                        <template #default="scope">
                          <div
                            v-if="scope.row?.targetPort"
                            :style="compareService(scope.row, 'service')"
                          >
                            <span>{{ scope.row.targetPort.intVal }}</span>
                          </div>
                        </template>
                      </el-table-column>
                      <el-table-column label="strVal" width="60">
                        <template #default="scope">
                          <div
                            v-if="scope.row?.targetPort"
                            :style="compareService(scope.row, 'service')"
                          >
                            <span>{{ scope.row.targetPort.strVal }}</span>
                          </div>
                        </template>
                      </el-table-column>
                    </el-table-column>
                  </el-table>
                </div>
              </div>
            </div>
            <div class="table-title">
              <div class="pod">extService</div>
            </div>
            <div class="table-container">
              <div class="table-container">
                <div class="compare-menu">
                  <div
                    class="menu-title"
                    :style="fetchSerStyle('clusterIP', 'extService')"
                  >
                    clusterIP
                  </div>
                  <div
                    class="menu-title"
                    :style="fetchSerStyle('type', 'extService')"
                  >
                    type
                  </div>
                  <div
                    class="menu-title"
                    :style="fetchSerStyle('ports', 'extService')"
                  >
                    ports
                  </div>
                </div>
                <div
                  class="table-item"
                  v-for="(info, indexName) in comparedItems"
                  :key="indexName"
                >
                  <div
                    class="menu-title"
                    :style="fetchSerStyle('clusterIP', 'extService')"
                  >
                    {{ getService(info, "extService")?.clusterIP }}
                  </div>
                  <div
                    class="menu-title"
                    :style="fetchSerStyle('type', 'extService')"
                  >
                    {{ getService(info, "extService")?.type }}
                  </div>
                  <div
                    class="menu-title"
                    :style="fetchSerStyle('ports', 'extService')"
                  >
                    <el-table
                      size="small"
                      :data="getService(info, 'extService')?.ports"
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
                          <div
                            v-if="Object.keys(scope.row).length > 0"
                            :style="compareService(scope.row, 'extService')"
                          >
                            <span>{{ scope.row[v.field] }}</span>
                          </div>
                        </template>
                      </el-table-column>
                      <el-table-column label="targetPort">
                        <el-table-column label="type" width="50">
                          <template #default="scope">
                            <div
                              v-if="scope.row?.targetPort"
                              :style="compareService(scope.row, 'extService')"
                            >
                              <span>{{ scope.row.targetPort.type }}</span>
                            </div>
                          </template>
                        </el-table-column>
                        <el-table-column label="intVal" width="60">
                          <template #default="scope">
                            <div
                              v-if="scope.row?.targetPort"
                              :style="compareService(scope.row, 'extService')"
                            >
                              <span>{{ scope.row.targetPort.intVal }}</span>
                            </div>
                          </template>
                        </el-table-column>
                        <el-table-column label="strVal" width="60">
                          <template #default="scope">
                            <div
                              v-if="scope.row?.targetPort"
                              :style="compareService(scope.row, 'extService')"
                            >
                              <span>{{ scope.row.targetPort.strVal }}</span>
                            </div>
                          </template>
                        </el-table-column>
                      </el-table-column>
                    </el-table>
                  </div>
                </div>
              </div>
            </div>
            <div class="table-title">
              <div class="pod">spec</div>
            </div>
            <div class="table-container">
              <div class="compare-menu">
                <div class="menu-title" :style="fetchSpecStyle('replicas')">
                  replicas
                </div>
                <div class="menu-title" :style="fetchSpecStyle('role')">
                  role
                </div>
                <div
                  class="menu-title"
                  :style="fetchSpecStyle('accessControlId')"
                >
                  accessControlId
                </div>
                <div
                  class="menu-title"
                  :style="fetchSpecStyle('loadBalancerId')"
                >
                  loadBalancerId
                </div>
                <div
                  class="menu-title"
                  :style="fetchSpecStyle('overrideListeners')"
                >
                  overrideListeners
                </div>
                <div class="menu-title" :style="fetchSpecStyle('slbStatus')">
                  slb Status
                </div>
              </div>
              <div
                class="table-item"
                v-for="(info, indexName) in comparedItems"
                :key="indexName"
              >
                <div class="menu-title" :style="fetchSpecStyle('replicas')">
                  {{ getSpec(info)?.replicas }}
                </div>
                <div class="menu-title" :style="fetchSpecStyle('role')">
                  {{ getSpec(info)?.role }}
                </div>
                <div
                  class="menu-title"
                  :style="fetchSpecStyle('accessControlId')"
                >
                  {{ getSpec(info)?.accessControlId }}
                </div>
                <div
                  class="menu-title"
                  :style="fetchSpecStyle('loadBalancerId')"
                >
                  {{ getSpec(info)?.loadBalancerId }}
                </div>
                <div
                  class="menu-title"
                  :style="fetchSpecStyle('overrideListeners')"
                >
                  {{ getSpec(info)?.overrideListeners }}
                </div>
                <div class="menu-title" :style="fetchSpecStyle('slbStatus')">
                  {{ getSpec(info)?.slbStatus }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </transition>
</template>
<script setup lang="ts">
import { Close } from "@element-plus/icons-vue";
import { useStore } from "@/store";
import {
  onMounted,
  watch,
  onBeforeUnmount,
  ref,
  computed,
  nextTick,
  provide,
} from "vue";
import { useRoute, useRouter } from "vue-router";
import Avatar from "@/components/avatar.vue";
import { ElNotification } from "element-plus";
import { goToNav } from "@/views/resource/components/tabelUtil";
import { cloneDeep, throttle } from "lodash";
import { routerToNamespace } from "./resource/util";

const route = useRoute();
const router = useRouter();
const store = useStore();
// header namespace切换
const nslist: any = ref([]);
getNamespaceList();
function getNamespaceList() {
  let rule = JSON.parse(String(localStorage.getItem("clusterRole")));
  if (!rule) {
    ElNotification({ title: "未登陆", type: "error", duration: 3000 });
    router.push({ name: "Index" });
  }
  let nsListTemp = [];
  for (let ns in rule) {
    nsListTemp.push({
      Name: ns,
    });
  }
  nslist.value = nsListTemp;
}

const nsValue = ref("default");
const showOption = ref(false);
provide("showOption", showOption);
let nsInputValue = ref("");
let optionsNsList: any = computed(() => {
  let nameList = [];
  for (let { Name } of nslist.value) {
    if (!Number(Name.slice(-2))) {
      Name = Name.slice(0, Name.length - 1) + "0" + Name.slice(-1);
    }
    nameList.push(Name);
  }
  nameList.sort();
  let sortNsList = [];
  for (let name of nameList) {
    const nameIndex = nslist.value.findIndex((n: any) => {
      let tempName = name;
      if (name.slice(-2, -1) === "0") {
        tempName = name.slice(0, -2) + name.slice(-1);
      }
      return n.Name === tempName;
    });

    sortNsList.push(nslist.value[nameIndex]);
  }
  return sortNsList.filter((ns: any) => {
    return nsValue.value != ns.Name && ns.Name.includes(nsInputValue.value);
  });
});
let nsValueSelected = ref(-1);
const namespaceInputSelectUp = throttle(
  () => {
    namespaceInputSelect("up");
  },
  150,
  { leading: true, trailing: false }
);
const namespaceInputSelectDown = throttle(
  () => {
    namespaceInputSelect("down");
  },
  150,
  { leading: true, trailing: false }
);

document.addEventListener("keydown", function (e) {
  if (!showOption.value) return;
  if (e.key === "ArrowDown") {
    namespaceInputSelectDown();
  } else if (e.key === "ArrowUp") {
    namespaceInputSelectUp();
  } else if (e.key === "Enter") {
    namespaceInputConfirm();
  }
});

function namespaceInputSelect(isup: string) {
  if (isup === "up") {
    if (nsValueSelected.value > -1) {
      nsValueSelected.value = nsValueSelected.value - 1;
      goToNav("ns" + nsValueSelected.value);
    } else {
      nsinput.value.$refs.input.focus();
    }
  } else {
    if (nsValueSelected.value + 1 < optionsNsList.value.length) {
      nsValueSelected.value = nsValueSelected.value + 1;
      document.getElementById("ns-ops-index0")?.focus();
      goToNav("ns" + (nsValueSelected.value - 1));
    }
  }
}
let leftNsTimer: any = ref(null);
function enterNamespace() {
  clearTimeout(leftNsTimer.value);
}
function leaveNamespace() {
  clearTimeout(leftNsTimer.value);
  leftNsTimer.value = setTimeout(() => {
    if (showOption.value) {
      showOption.value = false;
      nsInputValue.value = "";
      nsValueSelected.value = -1;
    }
  }, 200);
}
function namespaceInputConfirm() {
  let toID = 0;
  if (nsValueSelected.value >= 0) {
    toID = nsValueSelected.value;
  }
  if (optionsNsList.value[toID]) {
    changeNs(optionsNsList.value[toID].Name);
    nsValueSelected.value = -1;
  }
  showOption.value = false;
}
// 侧边栏控制
const collapsed = ref(window.innerWidth > 768 ? true : false);
function changeSide() {
  collapsed.value = !collapsed.value;
}
function backToNamespace() {
  router.replace({ name: "Namespace" });
}

function changeNs(ns: string) {
  // showOption.value = false
  const pathGVK = route.path.split("/")[2];
  const gvkarr = pathGVK.split("-");
  const gvkObj = {
    gv: `${gvkarr[0]}-${gvkarr[1]}`,
    kind: `${gvkarr[2]}`,
  };
  routerToNamespace(ns, gvkObj);
  nsInputValue.value = "";
}
let nsinput = ref({
  $refs: {
    input: {
      focus,
    },
  },
});
async function selFocus() {
  showOption.value = true;
  enterNamespace();
  await nextTick();
  nsinput.value.$refs.input.focus();
}
let resourceType = ref(<any>[]);
const localGvk = ref("");

let gvkList = JSON.parse(String(localStorage.getItem("gvkList")));

let gvkMenu = ref(<any>[]);

function initGvk() {
  let gvList = <any>[];
  let gvkMenuInit = [];
  const pathGVK = route.path.split("/")[2];
  if (!pathGVK) {
    return;
  }
  const gvnAscSort = [];
  for (let gvk of gvkList) {
    gvnAscSort.push(`${gvk.gv}-${gvk.kind}`);
  }
  gvnAscSort.sort();
  const sortGvkList = [];
  for (let gvk of gvnAscSort) {
    const sortIndex = gvkList.findIndex((val: any) => {
      return `${val.gv}-${val.kind}` === gvk;
    });
    if (
      sortIndex >= 0 &&
      !gvkList[sortIndex].gv.startsWith("metrics.k8s") &&
      !gvkList[sortIndex].gv.startsWith("autoscaling")
    ) {
      sortGvkList.push(gvkList[sortIndex]);
    }
  }
  for (let gvk of sortGvkList) {
    if (!gvList.includes(gvk.gv)) {
      gvList.push(gvk.gv);
      gvkMenuInit.push({
        groupVersion: gvk.gv,
        kinds: [gvk.kind],
        menuOpen: pathGVK.startsWith(gvk.gv),
      });
    } else {
      let thisIndex = gvkMenuInit.findIndex((val) => {
        return val.groupVersion === gvk.gv;
      });
      gvkMenuInit[thisIndex].kinds.push(gvk.kind);
    }
  }
  gvkMenu.value = gvkMenuInit;
}
function chcekKind(kind: string) {
  const findIndex = gvkMenu.value.findIndex((gvk: any) => {
    return gvk.menuOpen;
  });
  let gv = "";
  if (findIndex >= 0) {
    gv = gvkMenu.value[findIndex].groupVersion;
  }
  const pathGVK = route.path.split("/")[2];
  return `${gv}-${kind}` === pathGVK;
}

function goResource(kind: string) {
  const findIndex = gvkMenu.value.findIndex((gvk: any) => {
    return gvk.menuOpen;
  });
  let gv = "";
  if (findIndex >= 0) {
    gv = gvkMenu.value[findIndex].groupVersion;
  }
  router.replace({ path: `/${nsValue.value}/${gv}-${kind}` });
}

onMounted(() => {
  window.addEventListener("click", () => {
    showOption.value = false;
  });
  resourceType.value = route.params.sider;
  const pathSplit = route.path.split("/");
  nsValue.value = pathSplit[1];
  localGvk.value = pathSplit[2];
});

initGvk();
watch(
  () => route.path,
  (newPath) => {
    const pathSplit = newPath.split("/");
    gvkList = JSON.parse(String(localStorage.getItem("gvkList")));
    if (nsValue.value != pathSplit[1]) {
      nsValue.value = pathSplit[1];
      initGvk();
    }
    localGvk.value = pathSplit[2];
  }
);

function groupFormat(groupVersion: string) {
  return groupVersion.replace("-", ".");
}
let openMenuFlg: any = ref(null);
let backTimer: any = ref(null);
function leaveMenu() {
  clearTimeout(openMenuFlg.value);
}
function openMenu(index: any) {
  clearTimeout(openMenuFlg.value);
  clearTimeout(backTimer.value);
  openMenuFlg.value = setTimeout(() => {
    collapsed.value = true;
    for (let menuIndex in gvkMenu.value) {
      if (menuIndex == index) {
        gvkMenu.value[menuIndex].menuOpen = true;
      } else {
        gvkMenu.value[menuIndex].menuOpen = false;
      }
    }
  }, 200);
}
function backKind() {
  backTimer.value = setTimeout(() => {
    const pathGVK = route.path.split("/")[2].split("-");
    for (let menu of gvkMenu.value) {
      if (menu.groupVersion === `${pathGVK[0]}-${pathGVK[1]}`) {
        menu.menuOpen = true;
      } else {
        menu.menuOpen = false;
      }
    }
  }, 200);
}
function findKinds(gvkMenus: any) {
  const findIndex = gvkMenus.findIndex((gvk: any) => {
    return gvk.menuOpen;
  });
  if (findIndex >= 0) {
    return gvkMenus[findIndex].kinds;
  } else {
    return [];
  }
}
function chcekKindSelected(gv: string) {
  if (localGvk.value.startsWith(gv)) {
    return "color: #6699cc";
  } else {
    return "color: #1818196B";
  }
}

let compareBtn = ref(false);
let compareList = ref(<any>[]);
let comparedItems = computed(() => {
  let storeCompareList = [];
  for (let obj of compareList.value) {
    for (let item of obj.items) {
      storeCompareList.push(item);
    }
  }
  return storeCompareList;
});
watch(
  () => store.state.user.compareItems.length,
  (newItemLength: number) => {
    compareBtn.value = newItemLength > 0;
    let storeCompareList = [];
    for (let item of cloneDeep(store.state.user.compareItems)) {
      let itemProto = item.metadata.protoType;
      const typeIndex = storeCompareList.findIndex((storeItem: any) => {
        return storeItem.protoType === itemProto;
      });
      if (typeIndex >= 0) {
        const itemIndex = storeCompareList[typeIndex].items.findIndex(
          (objItem: any) => {
            return (
              objItem.metadata.name === item.metadata.name &&
              objItem.metadata.namespace === item.metadata.namespace
            );
          }
        );
        if (itemIndex < 0) {
          storeCompareList[typeIndex].items.push(item);
        }
      } else {
        storeCompareList.push({
          protoType: itemProto,
          items: [item],
        });
      }
    }
    compareList.value = storeCompareList;
  },
  {
    immediate: true,
  }
);
function removeItem(item: any, delIndex: number) {
  const itemProto = item.metadata.protoType;
  let changedItems = cloneDeep(compareList.value);
  const typeIndex = changedItems.findIndex((storeItem: any) => {
    return storeItem.protoType === itemProto;
  });
  if (changedItems[typeIndex].items.length === 1) {
    changedItems.splice(typeIndex, 1);
  } else {
    changedItems[typeIndex].items.splice(delIndex, 1);
  }
  let commitItems = [];
  for (let obj of changedItems) {
    for (let itemC of obj.items) {
      commitItems.push(itemC);
    }
  }
  store.commit("user/setCompareItems", commitItems);
}

let comAllPods = <any>[];
let comAllSers = <any>[];
let comAllextSers = <any>[];
let comAllSpecs = <any>[];
watch(
  () => comparedItems.value.length,
  () => {
    let pods = [],
      sers = [],
      extSers = [],
      specs = [];
    for (let item of cloneDeep(comparedItems.value)) {
      pods.push(getPod(item)), sers.push(getService(item, "service"));
      extSers.push(getService(item, "extService"));
      specs.push(getSpec(item));
    }
    comAllPods = pods;
    comAllSers = sers;
    comAllextSers = extSers;
    comAllSpecs = specs;
  },
  {
    immediate: true,
  }
);

function clearCompare() {
  store.commit("user/setCompareItems", []);
}
let compareDialog = ref(false);
function compare() {
  compareDialog.value = true;
}
function getPod(item: any) {
  if (item.metadata.protoType === "Pod") {
    return {
      image: item.spec.containers[0].image,
      imagePullPolicy: item.spec.containers[0].imagePullPolicy,
      volumes: item.spec.volumes,
      volumeMounts: item.spec.containers[0].volumeMounts,
      command: item.spec.containers[0].command,
      args: item.spec.containers[0].args,
      ports: item.spec.containers[0].ports,
      env: item.spec.containers[0].env,
    };
  }
  if (
    item.metadata.protoType === "Application" ||
    item.metadata.protoType === "MysqlConfig"
  ) {
    let podItem = item.pod;
    return {
      image: podItem.spec.containers[0].image,
      imagePullPolicy: podItem.spec.containers[0].imagePullPolicy,
      volumes: podItem.spec.volumes,
      volumeMounts: podItem.spec.containers[0].volumeMounts,
      command: podItem.spec.containers[0].command,
      args: podItem.spec.containers[0].args,
      ports: podItem.spec.containers[0].ports,
      env: podItem.spec.containers[0].env,
    };
  }

  return {
    image: "",
    imagePullPolicy: "",
    volumes: [],
    volumeMounts: [],
    ports: [],
    env: [],
    command: [],
    args: [],
  };
}
const volumeMountsCloumn = [
  { field: "name", title: "name", width: 90 },
  { field: "mountPath", title: "mountPath", width: 110 },
  { field: "readOnly", title: "readOnly", width: 50 },
  { field: "subPath", title: "subPath", width: 70 },
  { field: "subPathExpr", title: "subPathExpr", width: 60 },
];
const containCloumn = [
  { field: "name", title: "name", width: 60 },
  { field: "containerPort", title: "containerPort", width: 100 },
  { field: "hostIP", title: "hostIP", width: 70 },
  { field: "hostPort", title: "hostPort", width: 80 },
  { field: "protocol", title: "protocol", width: 70 },
];
function compareStyle(key: string, compareObj: any, row: any = {}) {
  if (key === "volume") {
    let box = "",
      type = "",
      configMap = "",
      path = "",
      configRow = "",
      secretName = "";
    box = diff("box", compareObj);
    if (!box) {
      type = diff("type", compareObj);
      path = diff("path", compareObj);
      configMap = diff("configMap", compareObj);
      configRow = diff("configRow", compareObj, row);
    }
    secretName = diff("secretName", compareObj);
    return {
      box,
      type,
      configMap,
      path,
      configRow,
      secretName,
    };
  }
  if (key === "volumeMounts") {
    return {
      volumeMounts: diffMount(key, row),
    };
  }
  if (key === "command") {
    return {
      command: diffMount(key, row),
    };
  }
  if (key === "args") {
    return {
      args: diffMount(key, row),
    };
  }
  if (key === "port") {
    return {
      port: diffMount(key, row),
    };
  }
  if (key === "env") {
    return {
      env: diffMount(key, row),
    };
  }
}
function diff(key: string, compareObj: any, row: any = {}) {
  let vols = cloneDeep(comAllPods);
  for (let volIndex = 0; volIndex < vols.length - 1; volIndex++) {
    let sameArr = vols[volIndex].volumes.filter((vol: any) => {
      let hadVol = false;
      if (key === "box") {
        let stringVol = vol.name;
        for (let nextVol of vols[volIndex + 1].volumes) {
          if (stringVol === nextVol.name) {
            hadVol = true;
          }
        }
      }
      if (key === "path") {
        let stringVol = vol?.volumeSource?.hostPath?.path;
        for (let nextVol of vols[volIndex + 1].volumes) {
          let nextPath = nextVol?.volumeSource?.hostPath?.path;
          if (stringVol === nextPath) {
            hadVol = true;
          }
        }
      }
      if (key === "type") {
        let stringVol = vol?.volumeSource?.hostPath?.type;
        for (let nextVol of vols[volIndex + 1].volumes) {
          let nextPath = nextVol?.volumeSource?.hostPath?.type;
          if (stringVol === nextPath) {
            hadVol = true;
          }
        }
      }
      if (key === "configMap") {
        let stringVol = vol.volumeSource?.configMap?.localObjectReference?.name;
        for (let nextVol of vols[volIndex + 1].volumes) {
          let nextPath =
            nextVol?.volumeSource?.configMap?.localObjectReference?.name;
          if (stringVol === nextPath) {
            hadVol = true;
          }
        }
      }
      if (key === "secretName") {
        let stringVol = vol.volumeSource?.secret?.secretName;
        for (let nextVol of vols[volIndex + 1].volumes) {
          let nextPath = nextVol?.volumeSource?.secret?.secretName;
          if (stringVol === nextPath) {
            hadVol = true;
          }
        }
      }
      if (key === "configRow") {
        let rowString = `key=${row.key};value=${row.value}`;
        for (let nextVol of vols[volIndex + 1].volumes) {
          let nextItems = nextVol?.volumeSource?.configMap?.items || [];
          for (let rowItwm of nextItems) {
            let nextString = `key=${rowItwm.key};value=${rowItwm.value}`;
            if (rowString === nextString) {
              hadVol = true;
            }
          }
        }
      }
      return hadVol;
    });
    vols[volIndex + 1].volumes = sameArr;
  }
  let findObjIndex = vols[vols.length - 1].volumes.findIndex((volsame: any) => {
    if (key === "box") {
      return volsame.name === compareObj.name;
    }
    if (key === "path") {
      return (
        volsame?.volumeSource?.hostPath?.path ===
        compareObj?.volumeSource?.hostPath?.path
      );
    }
    if (key === "type") {
      return (
        volsame?.volumeSource?.hostPath?.type ===
        compareObj?.volumeSource?.hostPath?.type
      );
    }
    if (key === "configMap") {
      return (
        volsame?.volumeSource?.configMap?.localObjectReference?.name ===
        compareObj?.volumeSource?.configMap?.localObjectReference?.name
      );
    }
    if (key === "secretName") {
      return (
        volsame?.volumeSource?.secret?.secretName ===
        compareObj?.volumeSource?.secret?.secretName
      );
    }
    if (key === "configRow") {
      let rowString = `key=${row.key};value=${row.value}`;
      let hadItem = false;
      let nextItems = volsame?.volumeSource?.configMap?.items || [];
      for (let sameItem of nextItems) {
        let nextString = `key=${sameItem.key};value=${sameItem.value}`;
        if (rowString === nextString) {
          hadItem = true;
        }
      }
      return hadItem;
    }
  });
  if (findObjIndex < 0) {
    if (key === "box") {
      return "border: 1px solid red;";
    }
    return "color: red;";
  } else {
    return "";
  }
}
function diffMount(key: string, row: any) {
  let vols = cloneDeep(comAllPods);
  let findObjIndex = -1;
  if (key === "volumeMounts") {
    for (let volIndex = 0; volIndex < vols.length - 1; volIndex++) {
      let sameArr = vols[volIndex].volumeMounts.filter((vol: any) => {
        let hadVol = false;
        let stringVol = JSON.stringify(vol);
        for (let nextVol of vols[volIndex + 1].volumeMounts) {
          if (stringVol === JSON.stringify(nextVol)) {
            hadVol = true;
          }
        }
        return hadVol;
      });
      vols[volIndex + 1].volumeMounts = sameArr;
    }
    findObjIndex = vols[vols.length - 1].volumeMounts.findIndex(
      (mount: any) => {
        let rowString = JSON.stringify(row);
        return rowString === JSON.stringify(mount);
      }
    );
  }
  if (key === "command") {
    for (let volIndex = 0; volIndex < vols.length - 1; volIndex++) {
      let sameArr = vols[volIndex].command.filter((vol: any) => {
        let hadVol = false;
        let stringVol = JSON.stringify(vol);
        for (let nextVol of vols[volIndex + 1].command) {
          if (stringVol === JSON.stringify(nextVol)) {
            hadVol = true;
          }
        }
        return hadVol;
      });
      vols[volIndex + 1].command = sameArr;
    }
    findObjIndex = vols[vols.length - 1].command.findIndex((mount: any) => {
      let rowString = JSON.stringify(row);
      return rowString === JSON.stringify(mount);
    });
  }
  if (key === "args") {
    for (let volIndex = 0; volIndex < vols.length - 1; volIndex++) {
      let sameArr = vols[volIndex].args.filter((vol: any) => {
        let hadVol = false;
        let stringVol = JSON.stringify(vol);
        for (let nextVol of vols[volIndex + 1].args) {
          if (stringVol === JSON.stringify(nextVol)) {
            hadVol = true;
          }
        }
        return hadVol;
      });
      vols[volIndex + 1].args = sameArr;
    }
    findObjIndex = vols[vols.length - 1].args.findIndex((mount: any) => {
      let rowString = JSON.stringify(row);
      return rowString === JSON.stringify(mount);
    });
  }
  if (key === "port") {
    for (let volIndex = 0; volIndex < vols.length - 1; volIndex++) {
      let sameArr = vols[volIndex].ports.filter((vol: any) => {
        let hadVol = false;
        let stringVol = JSON.stringify(vol);
        for (let nextVol of vols[volIndex + 1].ports) {
          if (stringVol === JSON.stringify(nextVol)) {
            hadVol = true;
          }
        }
        return hadVol;
      });
      vols[volIndex + 1].ports = sameArr;
    }
    findObjIndex = vols[vols.length - 1].ports.findIndex((mount: any) => {
      let rowString = JSON.stringify(row);
      return rowString === JSON.stringify(mount);
    });
  }
  if (key === "env") {
    for (let volIndex = 0; volIndex < vols.length - 1; volIndex++) {
      let sameArr = vols[volIndex].env.filter((vol: any) => {
        let hadVol = false;
        let stringVol = JSON.stringify(vol);
        for (let nextVol of vols[volIndex + 1].env) {
          if (stringVol === JSON.stringify(nextVol)) {
            hadVol = true;
          }
        }
        return hadVol;
      });
      vols[volIndex + 1].env = sameArr;
    }
    findObjIndex = vols[vols.length - 1].env.findIndex((mount: any) => {
      let rowString = JSON.stringify(row);
      return rowString === JSON.stringify(mount);
    });
  }

  if (findObjIndex < 0) {
    return "color: red;";
  } else {
    return "";
  }
}
function fetchStyle(key: string) {
  let allPod = cloneDeep(comAllPods);
  let equal = true,
    volIndex = 1,
    firstItem = allPod[0];
  for (let pod of allPod) {
    if (key === "image" && firstItem.image != pod.image) {
      equal = false;
    }
    if (
      key === "imagePullPolicy" &&
      firstItem.imagePullPolicy != pod.imagePullPolicy
    ) {
      equal = false;
    }
    if (key === "volume" && pod.volumes.length > volIndex) {
      volIndex = pod.volumes.length;
    }
    if (key === "volumeMounts" && pod.volumeMounts.length > volIndex) {
      volIndex = pod.volumeMounts.length;
    }
    if (key === "command" && pod.command.length > volIndex) {
      volIndex = pod.command.length;
    }
    if (key === "args" && pod.args.length > volIndex) {
      volIndex = pod.args.length;
    }
    if (key === "port" && pod.ports.length > volIndex) {
      volIndex = pod.ports.length;
    }
    if (key === "env" && pod.env.length > volIndex) {
      volIndex = pod.env.length;
    }
  }
  if (key === "volume") {
    let lineHeight = volIndex * 140;
    return `height: ${lineHeight}px; `;
  }
  if (key === "volumeMounts") {
    let lineHeight = volIndex * 69 + 56;
    return `height: ${lineHeight}px; `;
  }
  if (key === "command") {
    let lineHeight = volIndex * 32;
    return `height: ${lineHeight}px; `;
  }
  if (key === "args") {
    let lineHeight = volIndex * 32;
    return `height: ${lineHeight}px; `;
  }
  if (key === "port") {
    let lineHeight = volIndex * 55 + 32;
    return `height: ${lineHeight}px; `;
  }
  if (key === "env") {
    let lineHeight = volIndex * 69 + 32;
    return `height: ${lineHeight}px; `;
  }
  return `height: 60px; background-color: ${equal ? "" : "#fef0f0"}; color: ${
    equal ? "" : "red"
  }`;
}
const portsCloumn = [
  { field: "name", title: "name", width: 50 },
  { field: "nodePort", title: "nodePort", width: 50 },
  { field: "port", title: "port", width: 60 },
  { field: "protocol", title: "protocol", width: 50 },
];
function getService(item: any, serviceKey: string) {
  if (
    item.metadata.protoType === "Application" ||
    item.metadata.protoType === "MysqlConfig"
  ) {
    let serviceItem: any = {};
    if (serviceKey === "service") {
      serviceItem = item.service;
    } else {
      serviceItem = item.extensionService || {
        spec: {
          clusterIP: "",
          type: "",
          ports: [],
        },
      };
    }
    return {
      clusterIP: serviceItem.spec.clusterIP,
      type: serviceItem.spec.type,
      ports: serviceItem.spec.ports,
    };
  }
  return {
    clusterIP: "",
    type: "",
    ports: [],
  };
}
function fetchSerStyle(key: string, serviceKey: string) {
  let allSer = [];
  if (serviceKey === "service") {
    allSer = cloneDeep(comAllSers);
  } else {
    allSer = cloneDeep(comAllextSers);
  }
  let equal = true,
    firstItem = allSer[0],
    portIndex = 1;
  for (let ser of allSer) {
    if (key === "clusterIP" && firstItem.clusterIP != ser.clusterIP) {
      equal = false;
    }
    if (key === "type" && firstItem.type != ser.type) {
      equal = false;
    }
    if (key === "ports" && ser.ports.length > portIndex) {
      portIndex = ser.ports.length;
    }
  }
  if (key === "ports") {
    let lineHeight = portIndex * 55 + 60;
    return `height: ${lineHeight}px; `;
  }
  return `height: 30px; background-color: ${equal ? "" : "#fef0f0"}; color: ${
    equal ? "" : "red"
  }`;
}
function compareService(row: any = {}, serviceKey: string) {
  let vols = <any>[];
  if (serviceKey === "service") {
    vols = cloneDeep(comAllSers);
  } else {
    vols = cloneDeep(comAllextSers);
  }
  let findObjIndex = -1;
  for (let volIndex = 0; volIndex < vols.length - 1; volIndex++) {
    let sameArr = vols[volIndex].ports.filter((vol: any) => {
      let hadVol = false;
      let stringVol = JSON.stringify(vol);
      for (let nextVol of vols[volIndex + 1].ports) {
        if (stringVol === JSON.stringify(nextVol)) {
          hadVol = true;
        }
      }
      return hadVol;
    });
    vols[volIndex + 1].ports = sameArr;
  }
  findObjIndex = vols[vols.length - 1].ports.findIndex((mount: any) => {
    let rowString = JSON.stringify(row);
    return rowString === JSON.stringify(mount);
  });
  if (findObjIndex < 0) {
    return "color: red;";
  } else {
    return "";
  }
}

function getSpec(item: any) {
  if (
    item.metadata.protoType === "Application" ||
    item.metadata.protoType === "MysqlConfig"
  ) {
    return {
      replicas: item.replicas,
      role: item.role || "",
      accessControlId:
        item?.cloudNetworkConfig?.aliyunSLB?.accessControlId || "",
      loadBalancerId: item?.cloudNetworkConfig?.aliyunSLB?.loadBalancerId || "",
      overrideListeners:
        item?.cloudNetworkConfig?.aliyunSLB?.overrideListeners || "",
      slbStatus: item?.cloudNetworkConfig?.aliyunSLB?.status || "",
    };
  }
  return {
    replicas: "",
    role: "",
    accessControlId: "",
    loadBalancerId: "",
    overrideListeners: "",
    slbStatus: "",
  };
}
function fetchSpecStyle(key: string) {
  let allSer = cloneDeep(comAllSpecs);
  let Ereplicas = true,
    Erole = true,
    EaccessControlId = true,
    EloadBalancerId = true,
    EoverrideListeners = true,
    EslbStatus = true,
    firstItem = allSer[0];
  for (let ser of allSer) {
    if (firstItem.replicas != ser.replicas) {
      Ereplicas = false;
    }
    if (firstItem.role != ser.role) {
      Erole = false;
    }
    if (firstItem.accessControlId != ser.accessControlId) {
      EaccessControlId = false;
    }
    if (firstItem.loadBalancerId != ser.loadBalancerId) {
      EloadBalancerId = false;
    }
    if (firstItem.overrideListeners != ser.overrideListeners) {
      EoverrideListeners = false;
    }
    if (firstItem.slbStatus != ser.slbStatus) {
      EslbStatus = false;
    }
  }
  let equal = true;
  if (key === "replicas") equal = Ereplicas;
  if (key === "role") equal = Erole;
  if (key === "accessControlId") equal = EaccessControlId;
  if (key === "loadBalancerId") equal = EloadBalancerId;
  if (key === "overrideListeners") equal = EoverrideListeners;
  if (key === "slbStatus") equal = EslbStatus;
  return `height: 30px; background-color: ${equal ? "" : "#fef0f0"}; color: ${
    equal ? "" : "red"
  }`;
}
</script>
<style lang="scss" scoped>
@import "css/dashboard.scss";
@import "resource/css/animations.scss";
.fade-enter-active {
  -webkit-animation: scale-up-center 0.1s ease-in-out both;
  animation: scale-up-center 0.1s ease-in-out both;
}
.fade-leave-active {
  -webkit-animation: scale-out-center 0.1s cubic-bezier(0.55, 0.085, 0.68, 0.53)
    both;
  animation: scale-out-center 0.1s cubic-bezier(0.55, 0.085, 0.68, 0.53) both;
}
</style>
