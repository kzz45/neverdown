<template>
  <div class="mysql-menu">
    <el-menu
      active-text-color="#409eff"
      mode="horizontal"
      background-color="#fff"
      @select="handleSelect"
      class="el-menu-vertical-demo"
      default-active="1"
      text-color="#333"
    >
      <el-menu-item index="1">
        <span>spec</span>
      </el-menu-item>
      <el-menu-item index="2">
        <span>pod</span>
      </el-menu-item>
      <el-menu-item index="3">
        <span>service</span>
      </el-menu-item>
      <el-menu-item index="4">
        <span>extService</span>
      </el-menu-item>
    </el-menu>
  </div>
  <div class="mysql-spec">
    <div class="mysql-content">
      <div v-show="muneIndex === '1'" style="padding-bottom: 10px">
        <div class="temp-display">
          <div class="temp-nav">
            <div
              v-for="nav in navList"
              :key="nav"
              :class="foucsEl === nav ? 'nav-item-selected' : 'nav-item'"
              @click="skipToNav(nav)"
            >
              {{ nav }}
            </div>
          </div>
          <div class="temp-info">
            <div class="meta-item" id="id-replicas">
              <div class="meta-label">replicas:</div>
              <div class="meta-value">
                <el-input-number
                  v-model="props.specapp.replicas"
                  :min="0"
                  :max="300"
                  :disabled="hasHpa && editType === 'update'"
                />
              </div>
            </div>

            <div class="meta-item" id="id-appName">
              <div class="meta-label">appName:</div>
              <div class="meta-value">
                <el-input
                  v-model="props.specapp.appName"
                  size="small"
                  placeholder="Please input Name"
                />
              </div>
            </div>

            <div class="meta-item" id="id-cloudNetworkConfig">
              <div class="meta-label">cloudNetworkConfig:</div>
              <div class="meta-value">
                <AliyunSLB
                  :aliyunSLB="props.specapp.cloudNetworkConfig.aliyunSLB"
                />
              </div>
            </div>

            <div class="meta-item" id="id-persistentStorage">
              <div class="meta-label">persistentStorage:</div>
              <div class="meta-value">
                <el-input
                  v-model="props.specapp.persistentStorage.storageVolumePath"
                  size="small"
                  placeholder="Please input Name"
                />
              </div>
            </div>

            <div class="meta-item" id="id-watchPolicy">
              <div class="meta-label">watchPolicy:</div>
              <div class="meta-value">
                <el-select
                  v-model="props.specapp.watchPolicy"
                  placeholder="Select"
                  style="width: 100%"
                  @change="policyChange"
                >
                  <el-option
                    v-for="Policy in watchPolicyList"
                    :key="Policy"
                    :label="Policy"
                    :value="Policy"
                  />
                </el-select>
              </div>
            </div>
            <div class="switch-item" id="id-hpa">
              HPA:
              <el-switch
                style="margin-left: 20px"
                v-model="openHpa"
                @change="openHpaChange"
              />
            </div>
            <div
              class="meta-item"
              v-if="props.specapp.horizontalPodAutoscalerSpec"
            >
              <div class="meta-label" style="margin-top: 10px">
                horizontalPodAutoscalerSpec:
              </div>
              <div class="meta-value">
                <div>
                  <div>replicas</div>
                  <div>
                    min:
                    <el-input-number
                      v-model="
                        props.specapp.horizontalPodAutoscalerSpec.minReplicas
                      "
                      :min="0"
                      :max="
                        props.specapp.horizontalPodAutoscalerSpec.maxReplicas
                      "
                    />
                    max:
                    <el-input-number
                      v-model="
                        props.specapp.horizontalPodAutoscalerSpec.maxReplicas
                      "
                      :min="
                        props.specapp.horizontalPodAutoscalerSpec.minReplicas
                      "
                    />
                  </div>
                  <div>metrics</div>
                  <div>
                    <el-table
                      size="small"
                      :data="props.specapp.horizontalPodAutoscalerSpec.metrics"
                      border
                      style="width: 100%"
                      highlight-current-row
                    >
                      <el-table-column label="MetricSourceType" width="150">
                        <template #default="scope">
                          <span v-if="scope.row.isSet">
                            <el-select
                              v-model="scope.row.type"
                              style="width: 100%"
                              size="small"
                            >
                              <el-option
                                v-for="types in metricTypeList"
                                :key="types"
                                :label="types"
                                :value="types"
                              ></el-option>
                            </el-select>
                          </span>
                          <span v-else>{{ scope.row.type }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column label="resource">
                        <el-table-column label="name" width="150">
                          <template #default="scope">
                            <span v-if="scope.row.isSet">
                              <el-select
                                v-model="scope.row.resource.name"
                                size="small"
                                style="width: 100%"
                              >
                                <el-option
                                  v-for="types in resourceTypeList"
                                  :key="types"
                                  :label="types"
                                  :value="types"
                                ></el-option>
                              </el-select>
                            </span>
                            <span v-else>{{ scope.row.resource.name }}</span>
                          </template>
                        </el-table-column>
                        <el-table-column label="target" width="400">
                          <el-table-column label="value" width="150">
                            <template #default="scope">
                              <div
                                v-if="
                                  scope.row.resource.target.value ||
                                  scope.row.resource.target.value === ''
                                "
                              >
                                <span v-if="scope.row.isSet">
                                  <el-input
                                    v-model="scope.row.resource.target.value"
                                    size="small"
                                    placeholder="请输入内容"
                                  />
                                </span>
                                <span v-else>{{
                                  scope.row.resource.target.value
                                }}</span>
                              </div>
                            </template>
                          </el-table-column>
                          <el-table-column label="averageValue " width="200">
                            <template #default="scope">
                              <div
                                v-if="scope.row.resource.target.averageValue"
                              >
                                <span v-if="scope.row.isSet">
                                  <el-input
                                    v-model="
                                      scope.row.resource.target.averageValue
                                        .string
                                    "
                                    size="small"
                                    placeholder="请输入内容"
                                  />
                                </span>
                                <span v-else>{{
                                  scope.row.resource.target.averageValue.string
                                }}</span>
                              </div>
                            </template>
                          </el-table-column>
                          <el-table-column
                            label="averageUtilization"
                            width="150"
                          >
                            <template #default="scope">
                              <div
                                v-if="
                                  scope.row.resource.target
                                    .averageUtilization ||
                                  scope.row.resource.target
                                    .averageUtilization === ''
                                "
                              >
                                <span v-if="scope.row.isSet">
                                  <el-input
                                    v-model="
                                      scope.row.resource.target
                                        .averageUtilization
                                    "
                                    size="small"
                                    placeholder="请输入内容"
                                  />
                                </span>
                                <span v-else>{{
                                  scope.row.resource.target.averageUtilization
                                }}</span>
                              </div>
                            </template>
                          </el-table-column>
                          <el-table-column label="MetricTargetType" width="200">
                            <template #default="scope">
                              <span v-if="scope.row.isSet">
                                <el-select
                                  v-model="scope.row.resource.target.type"
                                  size="small"
                                  style="width: 100%"
                                  @change="
                                    (e) => {
                                      targetTypeChange(
                                        e,
                                        scope.row.resource.target
                                      );
                                    }
                                  "
                                >
                                  <el-option
                                    v-for="types in targetTypeList"
                                    :key="types"
                                    :label="types"
                                    :value="types"
                                  ></el-option>
                                </el-select>
                              </span>
                              <span v-else>{{
                                scope.row.resource.target.type
                              }}</span>
                            </template>
                          </el-table-column>
                        </el-table-column>
                      </el-table-column>
                      <el-table-column label="操作">
                        <template #default="scope">
                          <span
                            class="el-tag el-tag--info el-tag--mini"
                            style="cursor: pointer"
                            @click="
                              metricsEdit(
                                scope.row,
                                scope.row.isSet,
                                props.specapp.horizontalPodAutoscalerSpec
                                  .metrics
                              )
                            "
                          >
                            {{ scope.row.isSet ? "保存" : "修改" }}
                          </span>
                          <span
                            v-if="!scope.row.isSet"
                            class="el-tag el-tag--danger el-tag--mini"
                            @click="
                              rowDelete(
                                props.specapp.horizontalPodAutoscalerSpec
                                  .metrics,
                                scope.$index
                              )
                            "
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
                    <div
                      class="el-table-add-row"
                      style="width: 99.2%"
                      @click="metricsAdd()"
                    >
                      <span>+ 添加</span>
                    </div>
                  </div>
                  <div>behavior</div>
                  <div class="scale-box">
                    <div class="behavior-scale">
                      <div style="width: 150px">scaleUp:</div>
                      <el-select
                        v-model="
                          props.specapp.horizontalPodAutoscalerSpec.behavior
                            .scaleUp.selectPolicy
                        "
                      >
                        <template #prefix>selectPolicy:</template>
                        <el-option label=" Min" value="Min"></el-option>
                        <el-option label=" Max" value="Max"></el-option>
                      </el-select>
                      <el-input
                        v-model="
                          props.specapp.horizontalPodAutoscalerSpec.behavior
                            .scaleUp.stabilizationWindowSeconds
                        "
                        placeholder="Please input"
                      >
                        <template #prepend>stabilizationWindowSeconds</template>
                      </el-input>
                    </div>
                    <el-table
                      size="small"
                      :data="
                        props.specapp.horizontalPodAutoscalerSpec.behavior
                          .scaleUp.policies
                      "
                      border
                      style="width: 100%"
                      highlight-current-row
                    >
                      <el-table-column label="type" width="150">
                        <template #default="scope">
                          <span v-if="scope.row.isSet">
                            <el-select
                              v-model="scope.row.type"
                              style="width: 100%"
                              size="small"
                            >
                              <el-option
                                v-for="types in HPAScalingPolicyTypeList"
                                :key="types"
                                :label="types"
                                :value="types"
                              ></el-option>
                            </el-select>
                          </span>
                          <span v-else>{{ scope.row.type }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column label="periodSeconds" width="150">
                        <template #default="scope">
                          <span v-if="scope.row.isSet">
                            <el-input
                              v-model="scope.row.periodSeconds"
                              size="small"
                            />
                          </span>
                          <span v-else>{{ scope.row.periodSeconds }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column label="value" width="150">
                        <template #default="scope">
                          <span v-if="scope.row.isSet">
                            <el-input v-model="scope.row.value" size="small" />
                          </span>
                          <span v-else>{{ scope.row.value }}</span>
                        </template>
                      </el-table-column>

                      <el-table-column label="操作">
                        <template #default="scope">
                          <span
                            class="el-tag el-tag--info el-tag--mini"
                            style="cursor: pointer"
                            @click="
                              policiesEdit(
                                scope.row,
                                scope.row.isSet,
                                props.specapp.horizontalPodAutoscalerSpec
                                  .behavior.scaleUp.policies
                              )
                            "
                          >
                            {{ scope.row.isSet ? "保存" : "修改" }}
                          </span>
                          <span
                            v-if="!scope.row.isSet"
                            class="el-tag el-tag--danger el-tag--mini"
                            @click="
                              rowDelete(
                                props.specapp.horizontalPodAutoscalerSpec
                                  .behavior.scaleUp.policies,
                                scope.$index
                              )
                            "
                            style="cursor: pointer"
                          >
                            删除
                          </span>
                          <span
                            v-else
                            class="el-tag el-tag--mini"
                            style="cursor: pointer"
                            @click="cancelPoliciesEdit(scope.row)"
                          >
                            取消
                          </span>
                        </template>
                      </el-table-column>
                    </el-table>
                    <div
                      class="el-table-add-row"
                      style="width: 99.2%"
                      @click="policiesAdd('up')"
                    >
                      <span>+ 添加</span>
                    </div>
                  </div>
                  <div class="scale-box">
                    <div class="behavior-scale">
                      <div style="width: 150px">scaleDown:</div>
                      <el-select
                        v-model="
                          props.specapp.horizontalPodAutoscalerSpec.behavior
                            .scaleDown.selectPolicy
                        "
                      >
                        <template #prefix>selectPolicy:</template>
                        <el-option label=" Min" value="Min"></el-option>
                        <el-option label=" Max" value="Max"></el-option>
                      </el-select>
                      <el-input
                        v-model="
                          props.specapp.horizontalPodAutoscalerSpec.behavior
                            .scaleDown.stabilizationWindowSeconds
                        "
                        placeholder="Please input"
                      >
                        <template #prepend>stabilizationWindowSeconds</template>
                      </el-input>
                    </div>
                    <el-table
                      size="small"
                      :data="
                        props.specapp.horizontalPodAutoscalerSpec.behavior
                          .scaleDown.policies
                      "
                      border
                      style="width: 100%"
                      highlight-current-row
                    >
                      <el-table-column label="type" width="150">
                        <template #default="scope">
                          <span v-if="scope.row.isSet">
                            <el-select
                              v-model="scope.row.type"
                              style="width: 100%"
                              size="small"
                            >
                              <el-option
                                v-for="types in HPAScalingPolicyTypeList"
                                :key="types"
                                :label="types"
                                :value="types"
                              ></el-option>
                            </el-select>
                          </span>
                          <span v-else>{{ scope.row.type }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column label="periodSeconds" width="150">
                        <template #default="scope">
                          <span v-if="scope.row.isSet">
                            <el-input
                              v-model="scope.row.periodSeconds"
                              size="small"
                            />
                          </span>
                          <span v-else>{{ scope.row.periodSeconds }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column label="value" width="150">
                        <template #default="scope">
                          <span v-if="scope.row.isSet">
                            <el-input v-model="scope.row.value" size="small" />
                          </span>
                          <span v-else>{{ scope.row.value }}</span>
                        </template>
                      </el-table-column>

                      <el-table-column label="操作">
                        <template #default="scope">
                          <span
                            class="el-tag el-tag--info el-tag--mini"
                            style="cursor: pointer"
                            @click="
                              policiesDownEdit(
                                scope.row,
                                scope.row.isSet,
                                props.specapp.horizontalPodAutoscalerSpec
                                  .behavior.scaleDown.policies
                              )
                            "
                          >
                            {{ scope.row.isSet ? "保存" : "修改" }}
                          </span>
                          <span
                            v-if="!scope.row.isSet"
                            class="el-tag el-tag--danger el-tag--mini"
                            @click="
                              rowDelete(
                                props.specapp.horizontalPodAutoscalerSpec
                                  .behavior.scaleDown.policies,
                                scope.$index
                              )
                            "
                            style="cursor: pointer"
                          >
                            删除
                          </span>
                          <span
                            v-else
                            class="el-tag el-tag--mini"
                            style="cursor: pointer"
                            @click="cancelPoliciesDownEdit(scope.row)"
                          >
                            取消
                          </span>
                        </template>
                      </el-table-column>
                    </el-table>
                    <div
                      class="el-table-add-row"
                      style="width: 99.2%"
                      @click="policiesAdd('down')"
                    >
                      <span>+ 添加</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div v-show="muneIndex === '2'">
        <TemplateMeta
          :poddata="props.specapp.pod"
          :appname="props.specapp.appName"
          :disablelabels="true"
        ></TemplateMeta>
      </div>
      <div v-show="muneIndex === '3'">
        <ServiveSpec :spec="props.specapp.service.spec" />
      </div>
      <div v-show="muneIndex === '4'">
        <ServiveSpec :spec="props.specapp.extensionService.spec" />
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@import "../css/mysql.scss";
.el-table-add-row {
  width: 100%;
  border: 1px dashed #c1c1cd;
  border-radius: 3px;
  cursor: pointer;
  justify-content: center;
  display: flex;
  line-height: 34px;
}
.switch-item {
  margin-top: 20px;
  text-align: left;
  font-weight: 500;
  font-size: 0.9rem;
}
.scale-box {
  border: 1px solid #f1f2f4;
  padding: 10px;
  margin-bottom: 10px;
  .behavior-scale {
    display: flex;
    gap: 40px;
    margin-bottom: 10px;
  }
}
.temp-display {
  display: flex;
  position: relative;
  overflow: auto;
  max-height: calc(100vh - 400px);
  .temp-nav {
    width: 180px;
    position: sticky;
    top: 0;
    .nav-item {
      font-size: 0.8rem;
      font-weight: 400;
      margin-top: 10px;
      padding-left: 10px;
      text-align: left;
      color: #909399;
      letter-spacing: 1px;
      cursor: pointer;
      border-left: 3px solid #fff;
    }
    .nav-item-selected {
      font-size: 0.8rem;
      font-weight: 500;
      margin-top: 10px;
      padding-left: 10px;
      text-align: left;
      color: #409eff;
      letter-spacing: 1px;
      cursor: pointer;
      border-left: 3px solid #409eff;
    }
  }
  .temp-info {
    width: calc(100% - 180px);
  }
}
</style>

<script setup lang="ts">
import { cloneDeep } from "lodash";
import TemplateMeta from "./templatemeta.vue";
import AliyunSLB from "./aliyunslb.vue";
import ServiveSpec from "./servicespec.vue";
import { inject, nextTick, onMounted, provide, ref, watch } from "vue";

import { rowEdit, rowDelete, rowCancelEdit, goToNav } from "./tabelUtil";

const navList = [
  "replicas",
  "appName",
  "cloudNetworkConfig",
  "persistentStorage",
  "watchPolicy",
  "hpa",
];
let foucsEl = ref("replicas");
function skipToNav(nav: string) {
  foucsEl.value = nav;
  goToNav(nav);
}

const props = defineProps<{
  specapp?: any;
}>();

if (!props.specapp.extensionService) {
  props.specapp.extensionService = {
    metadata: {},
    spec: {
      clusterIP: "",
      type: "",
      ports: [],
    },
  };
}

let hasHpa = ref(false);
const editType = inject("editType");
let openHpa = ref(false);

provide("isApp", true);

let saveHap: any = ref({
  maxReplicas: 10,
  minReplicas: 1,
  metrics: [
    {
      resource: {
        name: "cpu",
        target: { averageUtilization: 50, type: "Utilization" },
      },
      type: "Resource",
    },
  ],
  behavior: {
    scaleDown: {
      policies: [],
      selectPolicy: "",
      stabilizationWindowSeconds: 0,
    },
    scaleUp: {
      policies: [],
      selectPolicy: "",
      stabilizationWindowSeconds: 0,
    },
  },
});

function targetTypeChange(type: any, target: any) {
  if (type === "Value") {
    target.value = "";
    delete target.averageValue;
    delete target.averageUtilization;
  }
  if (type === "AverageValue") {
    target.averageValue = { string: "" };
    delete target.value;
    delete target.averageUtilization;
  }
  if (type === "Utilization") {
    target.averageUtilization = "";
    delete target.averageValue;
    delete target.value;
  }
}

watch(
  () => props.specapp.appName,
  () => {
    openHpa.value = false;
    hasHpa.value = false;
    nextTick(() => {
      skipToNav("replicas");
    });
    if (!props.specapp.extensionService) {
      props.specapp.extensionService = {
        metadata: {},
        spec: {
          clusterIP: "",
          type: "",
          ports: [],
        },
      };
    }
    if (props.specapp.horizontalPodAutoscalerSpec) {
      openHpa.value = true;
      hasHpa.value = true;
      if (
        props.specapp.horizontalPodAutoscalerSpec.metrics &&
        props.specapp.horizontalPodAutoscalerSpec.metrics.length > 0
      ) {
        for (let metri of props.specapp.horizontalPodAutoscalerSpec.metrics) {
          if (!metri.resource.target.averageValue) {
            metri.resource.target.averageValue = {};
          }
        }
      }
      if (!props.specapp.horizontalPodAutoscalerSpec.behavior) {
        props.specapp.horizontalPodAutoscalerSpec.behavior = {
          scaleDown: {
            policies: [],
            selectPolicy: "",
            stabilizationWindowSeconds: 0,
          },
          scaleUp: {
            policies: [],
            selectPolicy: "",
            stabilizationWindowSeconds: 0,
          },
        };
      }
      saveHap.value = props.specapp.horizontalPodAutoscalerSpec;
    }
  },
  { immediate: true }
);
function openHpaChange(hapOpen: any) {
  if (hapOpen) {
    props.specapp.horizontalPodAutoscalerSpec = cloneDeep(saveHap.value);
    nextTick(() => {
      skipToNav("hpa");
    });
  } else {
    delete props.specapp.horizontalPodAutoscalerSpec;
  }
}

if (!props.specapp.cloudNetworkConfig.aliyunSLB) {
  props.specapp.cloudNetworkConfig = {
    aliyunSLB: {
      loadBalancerId: "",
      accessControlId: "",
      overrideListeners: true,
      status: "off",
    },
  };
}

const watchPolicyList = ["manual", "in-place-upgrade", "rolling-upgrade"];

let muneIndex = ref("1");
function handleSelect(gole: any) {
  muneIndex.value = gole;
}

let saveEditData = ref({});
function metricsEdit(row: any, isSet: boolean, allData: any) {
  rowEdit(row, isSet, allData, saveEditData);
}
function cancelEdit(row: any) {
  rowCancelEdit(row, saveEditData);
}
let savePolicyData = ref({});
function policiesEdit(row: any, isSet: boolean, allData: any) {
  rowEdit(row, isSet, allData, savePolicyData);
}
function cancelPoliciesEdit(row: any) {
  rowCancelEdit(row, savePolicyData);
}
let savePolicyDownData = ref({});
function policiesDownEdit(row: any, isSet: boolean, allData: any) {
  rowEdit(row, isSet, allData, savePolicyDownData);
}
function cancelPoliciesDownEdit(row: any) {
  rowCancelEdit(row, savePolicyDownData);
}

const metricTypeList = ["Resource"];
const targetTypeList = ["Utilization", "Value", "AverageValue"];
const resourceTypeList = ["cpu", "memory", "storage", "ephemeral-storage"];
const HPAScalingPolicyTypeList = ["Percent", "Pods"];
function metricsAdd() {
  const mertic = {
    type: "Resource",
    resource: {
      name: "cpu",
      target: {
        averageUtilization: 50,
        type: "Utilization",
      },
    },
  };
  props.specapp.horizontalPodAutoscalerSpec.metrics.push(mertic);
}
function policiesAdd(type: string) {
  const policy = {
    periodSeconds: 15,
    type: "Percent",
    value: 10,
  };
  if (type === "up") {
    props.specapp.horizontalPodAutoscalerSpec.behavior.scaleUp.policies.push(
      policy
    );
  } else {
    props.specapp.horizontalPodAutoscalerSpec.behavior.scaleDown.policies.push(
      policy
    );
  }
}
function policyChange(watchPolicy: string) {
  // if( watchPolicy === 'in-place-upgrade'){
  //   delete props.specapp.horizontalPodAutoscalerSpec
  // }else{
  //   const mertic = {
  //     type: "Resource",
  //     resource: {
  //       name: "cpu",
  //       target:{
  //         averageUtilization: 50,
  //         type: "Utilization"
  //       }
  //     }
  //   }
  //   props.specapp.horizontalPodAutoscalerSpec = {
  //     metrics : [mertic],
  //     maxReplicas: 10,
  //     minReplicas: 0
  //   }
  // }
}
</script>
