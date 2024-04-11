<template>
  <el-skeleton :rows="3" :loading="loading" :animated="true" >
    <template #template>
      <div style="padding: 0px 80px 0px 50px;">
        <div class="table-title card-skeleton"></div>
        <div class="table-content">
          <div class="table-nodata card-skeleton" style="margin-bottom: 10px;"></div>
          <div class="table-nodata card-skeleton"></div>
        </div>
      </div>
    </template>
    <template #default>
      <div class="header-btn">
        <div class="ns-card" @click="newItem" v-show="verifyVerb('create')">
          <div class="ns-title">{{i18nt('common.add')}}</div>
          <div>
            <el-icon class="ns-to"><Plus /></el-icon>
          </div>
        </div>
        <div class="ns-card" @click="importConfig" v-show="verifyVerb('create')">
          <div class="ns-title">{{i18nt('common.import')}}</div>
          <div>
            <el-icon class="ns-to"><Upload /></el-icon>
          </div>
        </div>
        <div class="ns-card" @click="importYaml" v-show="verifyVerb('create')">
          <div class="ns-title">导入yaml文件</div>
          <div>
            <el-icon class="ns-to"><Upload /></el-icon>
          </div>
        </div>
        <div class="ns-card" @click="exportList('export')" v-show="showExportList && verifyVerb('create')">
          <div class="ns-title">批量导入至</div>
          <div>
            <el-icon class="ns-to"><DocumentCopy /></el-icon>
          </div>
        </div>
        <div class="ns-card" @click="exportList('delete')" v-show="showExportList && verifyVerb('delete')">
          <div class="ns-title">批量删除</div>
          <div>
            <el-icon class="ns-to"><DeleteFilled /></el-icon>
          </div>
        </div>
        <div class="ns-card" @click="exportList('edit')" v-show="showExportList && verifyVerb('update')">
          <div class="ns-title">批量编辑</div>
          <div>
            <el-icon class="ns-to"><EditPen /></el-icon>
          </div>
        </div>
      </div>
      <div style="padding: 0px 80px 0px 50px;position: relative;">        
        <div class="action-tips">
          *右键条目呼出菜单
        </div>
        <div class="table-title">
          <div class="table-th" style="line-height: 40px;">
            <el-checkbox v-model="checkAllPod" size="large" @change="checkAllPodChange"
              style="float: left;margin-left:20px" />
            <div style="display: flex; justify-content: space-evenly; align-items: center;">
              <span :style="orderType==='name'?'color:#8dc63f':''">Name: </span>
              <lzm-name-filter ref="nameSearch" :filterStr="filterName" @filterinput="filterChange"></lzm-name-filter>
            </div>

            <el-select v-model="searchProject" clearable v-if="showSearchProject" style="display:inline-block;width:100px">
                <el-option v-for="proName in projectOptions" :key="proName" :label="proName" :value="proName" />
              </el-select>
          </div>
          <div class="table-th" v-for="td in showInfo" :key="td">
            <div v-if="td === 'phase'" >
              phase: 
              <el-select v-model="phaseType" clearable placeholder="Select" style="display:inline-block;width:100px">
                <el-option label="Pending" value="Pending" />
                <el-option label="Running" value="Running" />
              </el-select>
            </div>
            <div v-if="td === 'nodeName'">
              nodeName: 
              <el-select v-model="nodeName" clearable placeholder="Select" style="display:inline-block;width:150px;">
                <el-option :label="nName" :value="nName" v-for="nName of nodenameOptions" />
              </el-select>
            </div>
            <div v-if="td === 'type'" >
              type:
              <el-select v-model="eventType" clearable placeholder="Select" style="display:inline-block;width:100px">
                <el-option label="Normal" value="Normal" />
                <el-option label="Warning" value="Warning" />
              </el-select>
            </div>
            <div v-if="td === 'volume.timezone'" >
              volume.timezone: 
              <el-select v-model="volumeType" placeholder="Select" style="display:inline-block;width:130px">
                <el-option :label="nName" :value="nName" v-for="nName of volumeTypeOptions" />
              </el-select>
            </div>
            <div v-if="td === 'createTime'">
              <span :style="orderType.startsWith('createTime')?'color:#8dc63f':''">createTime: </span>
              <div>
                <span style="font-size: 14px;font-weight: 400">orderBy:</span>
                <el-select v-model="orderType" placeholder="Select" style="display:inline-block;width:145px">
                  <el-option v-for="type in sortOptions" :key="type.value" :label="type.label" :value="type.value" />
                </el-select>
              </div>
            </div>
            <div v-if="td === 'cpu/memory'" style="padding-left: 20px">
              <span :style="orderType.startsWith('cpu')?'color:#8dc63f':''">cpu</span> <br />
              <span :style="orderType.startsWith('memory')?'color:#8dc63f':''">memory</span>
            </div>
            <div v-if="td === 'container'">
              <span :style="orderType.startsWith('restartCount')?'color:#8dc63f':''">container</span>
            </div>
            
            <span v-if="!tdTitle.includes(td)">{{td}}</span>
          </div>
          <!-- <div class="table-th">Action</div> -->
        </div>
        <div class="table-content">
          <div class="table-nodata" v-if="paginationArr.length === 0">
            No Data
          </div>
          <div class="table-data" v-for="pod in paginationArr" v-bind:key="pod.metadata.name">
            <div class="table-tr" :oncontextmenu="(e) => { openAction(e, pod) }" :style="statusPrompt(pod)">
              <div class="table-td" style="text-align: left;padding-left: 20px">
                <div class="td-cell" style="display: flex; align-items: center; gap: 10px;justify-content: space-between;">
                  <el-checkbox v-model="pod.isChecked" @click.stop="popClick" size="large" @change="e => { checkPod(e, pod) }" />                    
                  <span style="margin-left: 10px;cursor:pointer;text-align: left" @click="getDetail(pod)">
                    {{pod.metadata.name}}
                  </span>
                  <div>
                    <el-popover
                      placement="top-start"
                      :width="450"
                    >
                      <template #reference>
                        <el-icon style="cursor: pointer" v-show="needShowLabels(pod.metadata)"><Collection /></el-icon>
                      </template>
                      <div v-for="labelItem in getLabels(pod.metadata)" :key="labelItem.key">
                        {{labelItem.key}}:<span style="color: #8dc63f;margin-left: 10px">{{labelItem.value}}</span>
                      </div>
                    </el-popover>
                    <el-popover
                      placement="top-start"
                      :width="450"
                    >
                      <template #reference>
                        <div v-show="needShowLabelsWatchPolicy(pod.metadata)" class="label-policy-auto">
                          {{getLabelsPolicy(pod.metadata).auto}}
                        </div>
                      </template>
                      <div >
                        watch-policy.openx.neverdown.io:<span style="color: #8dc63f;margin-left: 10px">{{getLabelsPolicy(pod.metadata).watchPolicy}}</span>
                      </div>
                    </el-popover>
                    <el-popover
                      placement="top-start"
                      :width="450"
                    >
                      <template #reference>
                        <el-icon style="cursor: pointer" v-show="needShowAffinity(pod.spec)"><LocationFilled /></el-icon>
                      </template>
                      <div v-for="Affi in getAffinity(pod.spec)" :key="Affi.key">
                        key:<span style="color: #8dc63f;margin-left: 10px">{{Affi.key}}</span>
                        operator:<span style="color: #8dc63f;margin-left: 10px">{{Affi.operator}}</span>
                        <div style="display: flex; flex-wrap: wrap" >
                          <div>
                            values:<span style="color: #8dc63f;margin-left: 10px" v-for="val in Affi.values" :key="val">
                            {{val}}
                          </span>
                          </div>
                        </div>
                      </div>
                    </el-popover>
                  </div>
                </div>
              </div>
              <div class="table-td" v-for="td in showInfo" :key="td">
                <div class="td-cell" 
                v-if="typeof(fetchInfo(td, pod)) === 'string' || typeof(fetchInfo(td, pod)) === 'number'">
                  {{ fetchInfo(td, pod) }}
                <!-- <div class="pod-delete-btn" @click.stop="confirmDelete(pod)" v-show="verifyVerb('delete') && fetchInfo(td, pod) === 'Running'">删除</div> -->
                </div>
                <div class="td-cell" v-else >
                  <div class="app-group" v-if="fetchInfoType(td, pod, 'applications')">
                    <div v-for="info of fetchInfoApp(td, pod)" :key="info" @click="goApp(pod, info)" class="app-info">
                      <div class="app-name">
                        {{ info.appName }}
                      </div>
                    </div>
                  </div>
                  <div class="app-group" v-if="fetchInfoType(td, pod, 'type')">
                    <div :style="fetchInfoApp(td, pod)==='Normal'?'':'color: #deb500'">{{fetchInfoApp(td, pod)}}</div>
                  </div>
                  <div class="app-group" v-if="fetchInfoType(td, pod, 'init')">
                    <div v-for="info of fetchInfoApp(td, pod)" :key="info" @click="goApp(pod, info)" class="app-info">
                      <div class="cytag" style="align-items: center;">
                        <div class="cytag-right" style="background-color: #5f6ae9;">
                          {{ info.replicas }}
                        </div>
                        <div v-if="fetchInfoHasHpa(td, pod)" style="margin-left:10px;cursor:pointer">
                          <div class="cytag-runbtn" style="background-color: #e3604e" v-if="Number(info.replicas) > 0" @click.stop="editPod(pod, info, 0)">
                            Stop<el-icon style="font-size:1.1rem"><VideoPause /></el-icon>
                          </div>
                          <div class="cytag-runbtn" style="background-color: #5f6ae9" v-else @click.stop="editPod(pod, info, 1)">
                            Run<el-icon style="font-size:1.1rem"><CaretRight /></el-icon>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="app-group" v-if="fetchInfoType(td, pod, 'taints')">
                    <div class="taints-title taints-title-style" v-show="fetchInfoApp(td, pod).length > 0">
                      <div>effect</div>
                      <div>key</div>
                      <div>value</div>
                    </div>
                    <div v-for="info of fetchInfoApp(td, pod)" :key="info.key">
                      <div class="taints-title">
                        <div class="taints-item">{{info.effect}}</div>
                        <div class="taints-item">{{info.key}}</div>
                        <div class="taints-item">{{info.value}}</div>
                      </div>
                    </div>
                  </div>
                  <div class="app-group" v-if="fetchInfoType(td, pod, 'cpu/memory') && verifyMetricsVerb()">
                    <div><span style="color: #8dc63f;">cpu:</span>{{getCpuInfo(td, pod).cpu}} </div>
                    <div><span style="color: #8dc63f;">memory:</span>{{getCpuInfo(td, pod).memory}} </div>
                  </div>
                  <div class="app-group" v-if="fetchInfoType(td, pod, 'container')">
                    <div :class="Boolean(pod.isOpen) ? 'td-cell operator-open' :'td-cell operator'" style="font-size: 1.5rem;" @click="openRow(pod)">
                      <span class="operator-icon">
                        <el-icon :style="pod.isOpen?'color: #8dc63f':''"><Menu /></el-icon>
                        <el-badge v-if="startCount(pod.status.containerStatuses)" style="left: 20px; line-height: 1rem;" :value="startCount(pod.status.containerStatuses)" type="danger"/>
                      </span>
                    </div>
                  </div>
                  <div class="app-group" v-if="fetchInfoType(td, pod, 'conditions')">
                    <div :class="Boolean(pod.isConditions) ? 'td-cell operator-open' :'td-cell operator'" style="font-size: 1.5rem;" @click="openRowConditions('', pod)">
                      <span class="operator-icon">
                        <el-icon :style="pod.isConditions?'color: #8dc63f':''"><MoreFilled /></el-icon>
                      </span>
                    </div>
                  </div>
                  <div class="app-group" v-if="fetchInfoType(td, pod, 'DaemonSetStatus')">
                    <div v-for="info of fetchInfoDaemonSetStatus(td, pod)" :key="info" class="status-group-4">
                      <div class="cytag">
                        <div class="cytag-left">ready</div><div class="cytag-right" style="background-color: #5f6ae9">
                          {{ info.numberReady }}
                        </div>
                      </div>
                      <div class="cytag">
                        <div class="cytag-left">updated</div><div class="cytag-right" style="background-color: #0e7fc0">
                          {{ info.updatedNumberScheduled }}
                        </div>
                      </div>
                      <div class="cytag">
                        <div class="cytag-left">ava</div><div class="cytag-right" style="background-color: #0e7fc0">
                          {{ info.numberAvailable }}
                        </div>
                      </div>
                      <div class="cytag">
                        <div class="cytag-left">unava</div><div class="cytag-right" style="background-color: #e3604e">
                          {{ info.numberUnavailable }}
                        </div>
                      </div>
                      <div class="cytag">
                        <div class="cytag-left">desired</div><div class="cytag-right" style="background-color: #0e7fc0">
                          {{ info.desiredNumberScheduled }}
                        </div>
                      </div>
                      <div class="cytag">
                        <div class="cytag-left">current</div><div class="cytag-right" style="background-color: #0e7fc0">
                          {{ info.currentNumberScheduled }}
                        </div>
                      </div>
                      <div class="cytag">
                        <div class="cytag-left">observed</div><div class="cytag-right" style="background-color: #0e7fc0">
                          {{ info.observedGeneration }}
                        </div>
                      </div>
                      <div class="cytag">
                        <div class="cytag-left">mis</div><div class="cytag-right" style="background-color: #e3604e">
                          {{ info.numberMisscheduled }}
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="app-group" v-if="fetchInfoType(td, pod, 'DeploymentStatus')">
                    <div v-for="info of fetchInfoDeploymentStatus(td, pod)" :key="info" class="status-group">
                      <div class="cytag">
                        <div class="cytag-left">replicas</div><div class="cytag-right" style="background-color: #5f6ae9">
                          {{ info.replicas }}
                        </div>
                      </div>
                      <div class="cytag">
                        <div class="cytag-left">ready</div><div class="cytag-right" style="background-color: #0e7fc0">
                          {{ info.readyReplicas }}
                        </div>
                      </div>
                      <div class="cytag">
                        <div class="cytag-left">updated</div><div class="cytag-right" style="background-color: #0e7fc0">
                          {{ info.updatedReplicas }}
                        </div>
                      </div>
                      <div class="cytag">
                        <div class="cytag-left">ava</div><div class="cytag-right" style="background-color: #0e7fc0">
                          {{ info.availableReplicas }}
                        </div>
                      </div>
                      <div class="cytag">
                        <div class="cytag-left">unava</div><div class="cytag-right" style="background-color: #e3604e">
                          {{ info.unavailableReplicas }}
                        </div>
                      </div>
                      <div class="conditions-btn" @click="openRowConditions('deploymentStatus', pod)">
                        conditions
                      </div>
                    </div>
                  </div>
                  <div class="app-group" v-if="fetchInfoType(td, pod, 'HPAStatus')">
                    <div v-for="info of fetchInfoHPAStatus(td, pod)" :key="info" class="status-group">
                      <div class="cytag">
                        <div class="cytag-left">current</div><div class="cytag-right" style="background-color: #5f6ae9">
                          {{ info.currentReplicas }}
                        </div>
                      </div>
                      <div class="cytag">
                        <div class="cytag-left">desired</div><div class="cytag-right" style="background-color: #0e7fc0">
                          {{ info.desiredReplicas }}
                        </div>
                      </div>
                      <div class="conditions-btn" @click="openRowConditions('horizontalPodAutoscalerStatus', pod)">
                        conditions
                      </div>
                    </div>
                  </div>
                  <div class="app-group" v-if="fetchInfoType(td, pod, 'StatefulSetStatus')">
                    <div v-for="info of fetchInfoStatefulSetStatus(td, pod)" :key="info" class="status-group">
                      <div class="cytag">
                        <div class="cytag-left">replicas:</div><div class="cytag-right" style="background-color: #5f6ae9">
                          {{ info.replicas }}
                        </div>
                      </div>
                      <div class="cytag">
                        <div class="cytag-left">ready</div><div class="cytag-right" style="background-color: #0e7fc0">
                          {{ info.readyReplicas }}
                        </div>
                      </div>
                      <div class="cytag">
                        <div class="cytag-left">updated</div><div class="cytag-right" style="background-color: #0e7fc0">
                          {{ info.updatedReplicas }}
                        </div>
                      </div>
                      <div class="cytag">
                        <div class="cytag-left">current</div><div class="cytag-right" style="background-color: #0e7fc0">
                          {{ info.currentReplicas }}
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <!-- <div class="table-td"><div class="td-cell">
                <div class="td-btns">
                  <div class="action-btn detail" @click="getDetail(pod)">详情</div>
                  <div class="action-btn detail" @click="exportDetail(pod)">导出</div>
                  <el-checkbox v-model="pod.isChecked" label="选择" size="large" border
                     v-show="verifyVerb('create')" @change="e => { checkPod(e, pod) }" />
                  <div class="action-btn delete" @click.stop="deleteItem(pod, $event)" v-show="verifyVerb('delete')">删除</div>
                </div>
              </div></div> -->
            </div>
            <transition name="podopen">
              <div class="container" v-if="Boolean(pod.isOpen)">
                <div class="container-title">
                  <div class="container-th">Name</div>
                  <div class="container-th">RestartCount</div>
                  <div class="container-th">image</div>
                  <div class="container-th" v-if="verifyMetricsVerb()">limits</div>
                  <div class="container-th" v-if="verifyMetricsVerb()">requests</div>
                  <div class="container-th" v-if="verifyMetricsVerb()">cpu/memory</div>
                  <div class="container-th">Operations</div>
                </div>
                <div class="container-tr" v-for="con in fetchHasContainer(pod)" :key="con.name">
                  <div class="container-td"><div class="container-cell" style="padding-left: 10px;">{{ con.name }}</div></div>
                  <div class="container-td"><div class="container-cell" :style="Number(con.restartCount) > 0 ? 'color: red': ''">{{ con.restartCount }}</div></div>
                  <div class="container-td">
                    <div class="container-cell" >
                      <div>{{ fetchImage(con.name, pod.spec.containers) }}</div>
                      <el-popover
                        placement="top-start"
                        :width="650"
                        trigger="hover"
                        @before-enter="getTags"
                      >
                        <template #reference>
                         <div>
                           <span style="color: #8dc63f;margin-left: 10px">@sha256:</span>{{ fetchImageIDshort(pod, con.name) }}
                         </div>
                        </template>
                        <div>
                         <div><span style="color: #00aeef;margin: 10px">@sha256:</span>{{ fetchImageID(pod, con.name).id }}</div>
                         <div><span style="color: #00aeef;margin: 10px">author:</span>{{ fetchImageID(pod, con.name).author }}</div>
                         <div><span style="color: #00aeef;margin: 10px">lastModifiedTime:</span>{{ fetchImageID(pod, con.name).time }}</div>
                         <div><span style="color: #00aeef;margin: 10px">branch:</span>{{ fetchImageID(pod, con.name).branch }}</div>
                         <div><span style="color: #00aeef;margin: 10px">commitHash:</span>{{ fetchImageID(pod, con.name).commitHash }}</div> 
                        </div>
                      </el-popover>
                    </div>
                  </div>
                  <div class="container-td" v-if="verifyMetricsVerb()">
                    <div class="container-cell" >
                      <div><span style="color: #8dc63f;">cpu:</span>{{fetchrRsourcesInfo('limits', pod, con.name).cpu}}</div>
                      <div><span style="color: #8dc63f;">mem:</span>{{fetchrRsourcesInfo('limits', pod, con.name).memory}} </div>
                    </div>
                  </div>
                  <div class="container-td" v-if="verifyMetricsVerb()">
                    <div class="container-cell" >
                      <div><span style="color: #8dc63f;">cpu:</span>{{fetchrRsourcesInfo('requests', pod, con.name).cpu}}</div>
                      <div><span style="color: #8dc63f;">mem:</span>{{fetchrRsourcesInfo('requests', pod, con.name).memory}} </div>
                    </div>
                  </div>
                  <div class="container-td" v-if="verifyMetricsVerb()">
                    <div class="container-cell" >
                      <div><span style="color: #8dc63f;">cpu:</span>{{fetchCpuInfo('cpu/memory', pod, con.name).cpu}}</div>
                      <div><span style="color: #8dc63f;">memory:</span>{{fetchCpuInfo('cpu/memory', pod, con.name).memory}} </div>
                    </div>
                  </div>
                  <div class="container-td">
                     <div class="meta-item">
                      <div class="time-input" @mouseenter.stop="timeInputFocus = true" @mouseleave="timeInputFocus = false">
                        log
                        <div class="edit-content">
                          <el-button @click="goTerm(pod.metadata, con, 0)" class="time-btn">all</el-button>
                          <el-button @click="goTerm(pod.metadata, con, 100)" class="time-btn">100s</el-button>
                          <el-button @click="goTerm(pod.metadata, con, 1000)" class="time-btn">1000s</el-button>
                          <el-input v-model="logSeconds" style="width: 200px" placeholder="Please input" @keyup.enter="goTerm(pod.metadata, con, logSeconds)" />
                        </div>
                      </div>
                      <div @click="goBash(pod.metadata, con)" class="time-input" style="cursor: pointer">
                        bash
                      </div>
                      <div v-if="con.currentDownloading" class="time-input" style="color:#5865f2">
                        downloading
                        <el-icon class="oper-icon"><Loading /></el-icon>
                      </div>
                      <div v-else class="time-input" @mouseenter.stop="timeInputFocus = true" @mouseleave="timeInputFocus = false">
                        current <el-icon ><Download /></el-icon>
                        <div class="edit-content">
                          <el-button @click="downloadLog(pod.metadata, con, false, 0)" class="time-btn">all</el-button>
                          <el-button @click="downloadLog(pod.metadata, con, false, 100)" class="time-btn">100s</el-button>
                          <el-button @click="downloadLog(pod.metadata, con, false, 1000)" class="time-btn">1000s</el-button>
                          <el-input v-model="logSeconds" style="width: 200px" placeholder="Please input" @keyup.enter="downloadLog(pod.metadata, con, false, logSeconds)" />
                        </div>
                      </div>
                      <div v-if="con.previousDownloading" class="time-input" style="color:#5865f2">
                        downloading
                        <el-icon class="oper-icon"><Loading /></el-icon>
                      </div>
                      <div v-else @click="downloadLog(pod.metadata, con, true, 0)" class="time-input" style="cursor: pointer">
                        previous <el-icon ><Download /></el-icon>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </transition>
            <transition name="podopen">
              <div class="conditions" v-show="Boolean(pod.isConditions)">
                <div class="container-title">
                  <div class="container-th">type</div>
                  <div class="container-th">status</div>
                  <div class="container-th">message</div>
                  <div class="container-th">reason</div>
                  <div class="container-th">lastTransitionTime</div>
                </div>
                <div v-if="fetchHasConditions(pod).length === 0" style="padding: 15px 0px;font-weight: 400">
                  no data
                </div>
                <div v-if="isAppConditions(pod)">
                  <div v-for="(appCon, conditionsAppIndex) in fetchHasConditions(pod)" :key="conditionsAppIndex">
                    <div class="conditions-name">{{ appCon.appName }}</div>
                    <div class="container-tr" v-for="(con, conditionsIndex) in appCon.conditions" :key="conditionsIndex">
                      <div class="container-td"><div class="container-cell" style="padding-left: 10px;">{{ con.type }}</div></div>
                      <div class="container-td"><div class="container-cell" style="padding-left: 10px;">{{ con.status }}</div></div>
                      <div class="container-td"><div class="container-cell" style="padding-left: 10px;">{{ con.message }}</div></div>
                      <div class="container-td"><div class="container-cell" style="padding-left: 10px;">{{ con.reason }}</div></div>
                      <div class="container-td"><div class="container-cell" style="padding-left: 10px;">{{ formatTime(con.lastTransitionTime.seconds) }}</div></div>
                    </div>
                  </div>
                </div>
                <div v-else>
                  <div class="container-tr" v-for="(con, conditionsIndex) in fetchHasConditions(pod)" :key="conditionsIndex">
                    <div class="container-td"><div class="container-cell" style="padding-left: 10px;">{{ con.type }}</div></div>
                    <div class="container-td"><div class="container-cell" style="padding-left: 10px;">{{ con.status }}</div></div>
                    <div class="container-td"><div class="container-cell" style="padding-left: 10px;">{{ con.message }}</div></div>
                    <div class="container-td"><div class="container-cell" style="padding-left: 10px;">{{ con.reason }}</div></div>
                    <div class="container-td"><div class="container-cell" style="padding-left: 10px;">{{ formatTime(con.lastTransitionTime.seconds) }}</div></div>
                  </div>
                </div>                
              </div>
            </transition>
          </div>
        </div>
        <el-pagination layout="total, prev, pager, next" style="text-align: left; padding-left: 50px;"
          v-model:currentPage="currentPage"  :total="SecretArr.length" />
      </div>
    </template>
  </el-skeleton>
  <div class="dialog-delay" v-show="editDialog" @click="editDialog = !editDialog"></div>
  <transition name="fade">
    <div class="dialog-body" v-if="editDialog">
      <div>
        <div style="height: 30px;color: black;text-align: left">{{i18nt(`common.${editType}`)}}</div>
        <el-icon style="color: #909399;position: absolute; top: 20px; right: 50px; cursor: pointer;"
        @click="editDialog = !editDialog"><Close /></el-icon>
      </div>
      <div style="width: 100%; border-bottom: 1px solid #ccc; padding-bottom: 15px;">
        <component :is="currentTabForm()" :itemInfo="currentItem" :initConfig="dialogInit"></component>
      </div>
      <div class="configMapDialog">
        <div class="dialog-btn" @click="confirmEdit" v-show="verifyVerb('update')">
          {{ i18nt('common.confirm') }}
        </div>
      </div>
    </div>
  </transition>
  <div class="dialog-delay" v-show="nsSelectDialog" @click="nsSelectDialog = !nsSelectDialog"></div>
  <transition name="fade">
    <div class="dialog-ns-export" v-if="nsSelectDialog">
      <div v-if="selectDialogType === 'export'">
        导入至:
        <div class="dialog-box" style="margin-bottom: 20px;margin-top: 5px">
          <div :class="batchns === ns ?'ns-export-sel': 'ns-export'" v-for="ns in nsSelectList" :key="ns" @click="selectBatchNs(ns)">
            {{ns}}
          </div>
        </div>
      </div>
      <div v-if="selectDialogType === 'delete'">
        <el-alert title="警告！确定要批量删除吗？" style="margin-bottom: 20px" type="error" show-icon :closable="false" />
        <el-checkbox v-model="ref_batchDeleteOneByOne.isNeed" size="large" :label="`逐个删除(每间隔${ref_batchDeleteOneByOne.gapSeconds}秒删除一个)`"/>          
        <span v-if="ref_batchDeleteOneByOne.isDoing" v-text="`已删除${ref_batchDeleteOneByOne.idxDoing}个/共${ref_batchDeleteOneByOne.sum}个`"></span>      
      </div>
      已选择 <span style="color: #8dc63f; font-weight: 400">{{allChecked.length}}</span> 项:
      <el-button type="primary" @click="yamlList()" plain>批量导出为yaml</el-button>
      <div class="object-box">
        <div class="export-object" v-for="item in allChecked" :key="item.metadata.name">
          {{item.metadata.name}}
        </div>
      </div>
      <div class="cytag" v-if="selectDialogType === 'edit'" style="align-items: center;justify-content: flex-start;">
        <div style="margin-left:10px;cursor:pointer;display:flex">
          <div class="cytag-runbtn" style="background-color: #e3604e" @click.stop="editPodBatch(0)">
            all Stop<el-icon style="font-size:1.1rem"><VideoPause /></el-icon>
          </div>
          <div class="cytag-runbtn" style="background-color: #5f6ae9;margin-left:20px" @click.stop="editPodBatch(1)">
            all Run<el-icon style="font-size:1.1rem"><CaretRight /></el-icon>
          </div>
        </div>
      </div>
      <span v-if="selectDialogType === 'edit'">同时编辑至其他namespace:</span>
      <div v-if="selectDialogType === 'edit'" style="max-height: 200px; overflow-y: auto;margin-bottom:5px">
        <div class="dialog-box" style="margin-bottom: 20px;margin-top: 5px">
          <div :class="batchEditNameSel(ns) ?'ns-export-sel': 'ns-export'" v-for="ns in nsSelectList" :key="ns" @click="selectBatchEditNs(ns)">
            {{ns}}
          </div>
        </div>
      </div>
      <div v-if="selectDialogType === 'edit'" style="max-height: 40vh; overflow-y: auto;">
        <EditlistFrom :volumes="volumes4Edit" :VolumeMounts="VolumeMounts4Edit" :variables="annotations4Edit" :affList="batchListAff"
          :batchEditType="batchEditType" @input="changeBatchType" @tag-change="changeBatchTag" @policyChange="changeBatchPolicy"
           :env="Envs4Edit" :ConfigMapList="ConfigMapList" @projectChange="changeBatchProject" @affChange="changeBatchAff" /> 
      </div>
      <div v-if="selectDialogType === 'export'">
        Affinity:
        <div>
          <el-checkbox v-for="aff in batchListAff" v-model="aff.isChecked"
            :key="aff.metadata.name" :label="aff.metadata.name" size="large" />
        </div>
      </div>
      <div v-if="selectDialogType === 'export'">
        Toleration:
        <div>
          <el-checkbox v-for="tol in batchListTol" v-model="tol.isChecked"
            :key="tol.metadata.name" :label="tol.metadata.name" size="large" />
        </div>
      </div>
      <div v-if="selectDialogType === 'delete' && !ref_batchDeleteOneByOne.isDoing" class="list-delete-pop">
        <div class="delete-confirm" @click="batchDelete">确定</div>
        <div class="delete-cancel" @click="nsSelectDialog = !nsSelectDialog">取消</div>
      </div>
      <div v-if="selectDialogType === 'export'" class="list-delete-pop">
        <div class="delete-cancel" @click="importList(batchns)">确定</div>
      </div>
      <div v-if="selectDialogType === 'edit'" class="list-delete-pop">
        <span style="color: #CCC;font-size: .8rem">暂只支持openx</span>
        <div class="delete-cancel" @click="editList()">确定</div>
      </div>
    </div>
  </transition>
  <div class="dialog-delay" v-show="deleteDialog" @click="deleteDialog = !deleteDialog"></div>
  <transition name="fade">
    <div class="delete-dialog-body" v-if="deleteDialog">
      <el-alert title="警告！您正在删除 Openx，该行为影响较大，确定要删除吗？" type="error" show-icon :closable="false" />
      <div class="delete-pop">
        <div class="delete-confirm" @click="itemDelete(delItem)">确定</div>
        <div class="delete-cancel" @click="deleteDialog = !deleteDialog">取消</div>
      </div>      
    </div>
  </transition>
  <div class="dialog-delay" v-show="yamlDialog" @click="yamlDialog = !yamlDialog"></div>
  <transition name="fade">

    <div class="dialog-ns-export" v-if="yamlDialog">
      
      <el-upload
        class="upload-demo" drag action="https://"
        :http-request="uploadConfig"
        :on-change="handleChange"
        multiple
      >
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">
          Drop file here or <em>click to upload</em>
        </div>
      </el-upload>
      <div class="delete-pop">
        <div class="delete-confirm" @click="yamlImport">确定</div>
        <div class="delete-cancel" @click="yamlDialog = !yamlDialog">取消</div>
      </div>  
    </div>
  </transition>
  <transition name="fade">
    <div v-if="showPop" class="test-pop"
      :style="popStyle" @click.stop="popClick">
        <div style="text-align: left;">确定删除吗？</div>
        <div style="text-align: left; font-weight: 400;color: black">{{delItem.metadata.name}}</div>
        <div class="pop-confirm" @click="confirmDelete(delItem)">确定</div>
        <div class="popover-arrow"></div>
    </div>
  </transition>

  <transition name="fade">
    <div v-if="actionMenu" class="action-menu-pop"
      :style="actionMenuStyle">
        <div class="td-btns">
          <div class="action-btn detail" @click="getDetail(actionItem)">详情</div>
          <div class="action-btn detail" @click="addCompare(actionItem)" v-show="verifyCompare()">加入对比</div>
          <div class="action-btn detail" @click="exportDetail(actionItem)">导出</div>
          <div class="action-btn detail">导出为yaml
            <div style="display: flex">
              <el-button @click="exportDetailYamlSimp(actionItem)">纯净</el-button>
              <el-button @click="exportDetailYaml(actionItem)">完整</el-button>
            </div>
          </div>
          <!-- <el-checkbox v-model="actionItem.isChecked" label="选择" border @click.stop="popClick"
              v-show="verifyVerb('create')" @change="e => { checkPod(e, actionItem) }" /> -->
          <!-- <div class="action-btn detail" @click="exportDetail(pod)" v-show="verifyVerb('create')">选择</div> -->
          <div class="action-btn delete" @click.stop="deleteItem(actionItem, $event)" v-show="verifyVerb('delete')">删除</div>
        </div>
    </div>
  </transition>

</template>

<style lang="scss" scoped>
@import "css/animations.scss";
.test-pop{
  position: fixed;
  width: 200px;
  padding: 15px 15px 30px 15px;
  color: red;
  border-radius: 5px;
  z-index: 11;
  background: #fff;
  box-shadow:  0px 0px 4px 1px #ff4d4f;
  .pop-confirm{
    color: #fff;
    border-color: #ff4d4f;
    background: #ff4d4f;
    width: 50px;
    position: absolute;
    right: 20px;
    cursor: pointer;
    &:hover{
      border-color: #ff7875;
      background: #ff7875;
    }
  }
}
.action-menu-pop{
  position: fixed;
  padding: 10px;
  color: red;
  border-radius: 5px;
  z-index: 10;
  background: #fff;
  box-shadow: 3px 3px 15px 5px #cdcfcf;
}
.fade-enter-active{
  -webkit-animation: scale-up-center 0.1s ease-in-out both;
	        animation: scale-up-center 0.1s ease-in-out both;
}
.fade-leave-active {
  -webkit-animation: scale-out-center 0.1s cubic-bezier(0.550, 0.085, 0.680, 0.530) both;
	        animation: scale-out-center 0.1s cubic-bezier(0.550, 0.085, 0.680, 0.530) both;
}
.podopen-enter-active {
  animation: bounce-in .3s cubic-bezier(0.390, 0.575, 0.565, 1.000) both;
}
.podopen-leave-active {
  animation: bounce-out .3s cubic-bezier(0.390, 0.575, 0.565, 1.000) both;
}
@import "css/list.scss";

  .delete-pop{
    display: flex; gap: 15px;
    position: absolute;
    bottom: 10px;
    right: 10px;
    .delete-confirm{
      color: #fff;
      border-radius: 5px;
      border:1px solid #ff4d4f;
      background: #ff4d4f;
      width: 50px;
      padding: 5px 10px;
      cursor: pointer;
      &:hover{
        border:1px solid #ff7875;
        background: #ff7875;
      }
    }
    .delete-cancel{
      color: #646464;
      border-radius: 5px;
      border: 1px solid #ccc;
      background: #fff;
      width: 50px;
      padding: 5px 10px;
      cursor: pointer;
      &:hover{
        border: 1px solid #646464;
      }
    }
  }
.list-delete-pop{
    display: flex; gap: 15px;
    position: absolute;
    bottom: 15px;
    right: 40px;
    .delete-confirm{
      color: #fff;
      border-radius: 5px;
      border:1px solid #ff4d4f;
      background: #ff4d4f;
      width: 50px;
      padding: 5px 10px;
      cursor: pointer;
      &:hover{
        border:1px solid #ff7875;
        background: #ff7875;
      }
    }
    .delete-cancel{
      color: #646464;
      border-radius: 5px;
      border: 1px solid #ccc;
      background: #fff;
      width: 50px;
      padding: 5px 10px;
      cursor: pointer;
      &:hover{
        border: 1px solid #646464;
      }
    }
}
.dialog-body{
  position: fixed;
  border-radius: 10px;
  top: 3vh;
  width: calc(80vw - 40px);
  height: 90vh;
  margin-left: 10vw;
  bottom: 0;
  left: 0;
  z-index: 1001;
  overflow: auto;
  padding: 20px;
  background: #fff;
}
.app-group{
  display: grid;
  padding: 10px 0px;
  .app-info{
    height: 60px;padding-top: 5px;
    display: flex;justify-content: space-around;
  }
  .app-name{
    padding: 5px;
    width: 80%;
    font-size: 1rem;
    border: 1px solid #e4ebf0;
    background: #f1f3f4;
    box-shadow: 3px 3px 5px #cdcfcf,
                -3px -3px 5px #fff;
    color: #000;
    border-radius: 4px;
    display: flex;
    align-items: center;
    justify-content: space-evenly;
    align-self: stretch;
    cursor: pointer;
    &:hover{
      box-shadow: inset 3px 3px 5px #cdcfcf,
          inset  -3px -3px 5px #fff;
    }
  }
  .taints-title{
    font: .8em sans-serif;
    font-weight: 400;
    display: grid;
    grid-template-columns: 33% 33% 33%;
    gap: 10px
  }
  .taints-title-style{
    border-bottom: 1px solid #ccc;
    padding-bottom: 5px;
    position: relative;
    top: -10px;
  }
  .taints-item{
    color: #8dc63f;
  }
}
.status-group{
  padding: 4px;
  font-size: .8rem;
  color: #fff;
  display: grid;
  grid-template-columns: 33% 33% 33%;
  row-gap: 4px;
  column-gap: 4px;
  overflow: hidden;
  text-overflow: ellipsis;

  .conditions-btn{
    color: #000;
    cursor: pointer;
    border-radius: 20px;
    background: #ffffff;
    box-shadow:  2px 2px 2px #a6a6a6,
                -2px -2px 2px #ffffff;
    &:hover{
      background: linear-gradient(145deg, #e6e6e6, #ffffff);
    }
  }
}
.status-group-4{
  padding: 4px;
  font-size: .8rem;
  color: #fff;
  display: grid;
  grid-template-columns: 25% 25% 25% 25%;
  row-gap: 4px;
  column-gap: 4px;
  overflow: hidden;
  text-overflow: ellipsis;

  .conditions-btn{
    color: #000;
    cursor: pointer;
    border-radius: 20px;
    background: #ffffff;
    box-shadow:  2px 2px 2px #a6a6a6,
                -2px -2px 2px #ffffff;
    &:hover{
      background: linear-gradient(145deg, #e6e6e6, #ffffff);
    }
  }
}
.cytag{
    display: flex;
    font-size: 0.8rem;
    justify-content: center;
    .cytag-left{
      width: 80px;
      padding-left: 5px;
      color: #fff;
      padding: 3px 0px;
      background-color: #505050;
      text-align: center;
      border-bottom-left-radius: 2px;
      border-top-left-radius: 2px;
      line-height: 20px;
      height: 20px;
    }
    .cytag-right{
      text-align: center;
      color: #ffffff;
      border-bottom-right-radius: 2px;
      border-top-right-radius: 2px;
      padding: 3px 8px;
      line-height: 20px;
      height: 20px;
    }
    .cytag-runbtn{
      background-color: rgb(227, 96, 78);
      border-radius: 3px;
      display: flex;
      align-items: center;
      gap: 3px;
      color: #ffffff;
      padding: 3px 8px;
      line-height: 20px;
      height: 20px;
    }
  }
.configMapDialog :deep(.el-dialog__body){
  padding: 0px 20px !important;
}
.label-policy-auto{
  font-size: 12px;
  color: #5f6ae9;
  cursor: pointer;
}
.label-policy{
  font-size: 12px;
  color:rgb(227, 96, 78);
  cursor: pointer;
}
</style>

<script setup lang="ts">
import { i18nt } from '@/i18n';
import { onMounted, watch, nextTick, ref, computed, provide, inject, onErrorCaptured, toRaw} from "vue"
import { useRouter, useRoute } from 'vue-router'
import { initSocketData, sendSocketMessage, binaryToStr, strToBinary, connectSocket } from "@/api/socket"
import { initObject, encodeify, decodeify, returnResourceList,TimerUtil } from './util'
import { cloneDeep, debounce} from 'lodash'
import { useStore } from '@/store'
import { InfoInGvk, getInfoInGvk, formatTime } from './podutils'

import { ElNotification } from 'element-plus';
import { returnResource, deleteSocketData, updateSocketData, currentPageVue } from './util'
import proto from '../../../proto/proto'
import EditlistFrom from './components/editList.vue'

onErrorCaptured((error) => {
  console.log('getError', error)
})

const store = useStore()
const router = useRouter()
const route = useRoute()

let  showInfo = ref(infoAfterVerify())
function infoAfterVerify(){
  let titlelist = cloneDeep(InfoInGvk[route.path.split('/')[2] || '']) || []
  if(!verifyMetricsVerb()){
    const cpuIndex = titlelist.findIndex((tit :string) => {
      return tit === 'cpu/memory'
    })
    if(cpuIndex >= 0) titlelist.splice(cpuIndex, 1)
  }
  return titlelist
}
watch(() => route.path, () => {
  showInfo.value = infoAfterVerify()
  checkAllPod.value = false
  currentPage.value = 1
  orderType.value = 'createTime-ascending'
  filterName.value = ''
})

const tdTitle = ['phase','type','nodeName','volume.timezone', 'createTime','cpu/memory', 'container']

function fetchInfo(type: string, detail: any){
  const nsGvk = route.path.split('/')[2] || ''
  return getInfoInGvk(type, detail, nsGvk)
}
function fetchInfoType(type: string, detail: any, objType: string){
  const info = fetchInfo(type, detail)
  if(info){
    if(info.type === 'DeploymentStatus' || info.type === 'HPAStatus'){
      if('conditions' === objType){
        return false
      }
    }
    return info.type === objType
  }
  return false
}
let imageTags = ref(<any>[])
function getTags(){
  getImageList('jingx-v1-Tag')
}
function openRow(pod: any){
  pod.isConditions = false
  pod.isOpen = !pod.isOpen
}
let conditionsType = ref('')
function openRowConditions(type: string, detail: any){
  if(conditionsType.value === type){
    detail.isConditions = !detail.isConditions
  }else{
    detail.isConditions = true
  }
  conditionsType.value = type
  detail.isOpen = false
  // detail.isConditions = !detail.isConditions
}
function startCount( ContainerStatus: any ) {
  let result = 0
  for(const con of ContainerStatus){
    result += con.restartCount
  }
  return result
}
function fetchInfoApp(type: string, detail: any){
  const info = fetchInfo(type, detail)
  if(info){
    return info.appInfo
  }
  return []
}
function fetchInfoDeploymentStatus(type: string, detail: any){
  const info = fetchInfo(type, detail)
  if(info){
    return info.deploymentStatus
  }
  return []
}
function fetchInfoDaemonSetStatus(type: string, detail: any){
  const info = fetchInfo(type, detail)
  if(info){
    return info.daemonSetStatus
  }
  return []
}
function fetchInfoHPAStatus(type: string, detail: any){
  const info = fetchInfo(type, detail)
  if(info){
    return info.deploymentStatus
  }
  return []
}
function fetchInfoStatefulSetStatus(type: string, detail: any){
  const info = fetchInfo(type, detail)
  if(info){
    return info.statefulSetStatus
  }
  return []
}
function fetchTime(num: number){
  if(num <= 1000){
    return num + 'ns'
  }
  if(Number((num/1000).toFixed(2)) >= 1000){
    const usTime = Number((num/1000).toFixed(2))
    if(Number((usTime/1000).toFixed(2)) >= 1000){
      return Number((usTime/1000000).toFixed(2)) + 's'
    }else{
      return (usTime/1000).toFixed(2) + 'ms'
    }
  }else{
    return (num/1000).toFixed(2) + 'μs'
  }
}
function getCpuInfo(type: string, detail: any){
  const info = fetchInfo(type, detail)
  const nsGvk = route.path.split('/')
  if(nsGvk[2] === 'core-v1-Node'){
    const getNodeInfo = nodeMetrics.value.filter((pod:any) => {
      return pod.metadata.name === info.name
    })
    if(getNodeInfo && getNodeInfo.length >= 1){
      let cpuNum = 0, memory = 0
      const node = getNodeInfo[0]
      cpuNum += Number(node.usage.cpu.string.slice(0,-1)|| 0)
      const memoryNum = getmemoryKi(node.usage.memory.string)
      memory += memoryNum||0
      return {
        cpu: fetchTime(cpuNum),
        memory: bytesToSize(memory)
      }
    }
  }else{
    const getPodInfo = podMetrics.value.filter((pod:any) => {
      return pod.metadata.name === info.name
    })
    if(getPodInfo && getPodInfo.length >= 1){
      let cpuNum = 0, memory = 0
      for(let con of getPodInfo[0].containers){
        cpuNum += Number(con.usage.cpu.string.slice(0,-1)|| 0)
        const memoryNum = getmemoryKi(con.usage.memory.string)
        memory += memoryNum||0
      }
      return {
        cpu: fetchTime(cpuNum),
        memory: bytesToSize(memory)
      }
    }
  }

  return {
    cpu: '0',
    memory: '0'
  }
}
function getmemoryKi(stringM: string){
  if(!stringM.slice(0,-2)){
    return 0
  }
  if(stringM.endsWith('Ki')) return Number(stringM.slice(0,-2))
  if(stringM.endsWith('Mi')){
    return Number(stringM.slice(0,-2)) * 1024
  }
  if(stringM.endsWith('Gi')){
    return Number(stringM.slice(0,-2)) * 1024 * 1024
  }
  if(stringM.endsWith('Ti')){
    return Number(stringM.slice(0,-2)) * 1024 * 1024 * 1024
  }
  if(stringM.endsWith('Pi')){
    return Number(stringM.slice(0,-2)) * 1024 * 1024 * 1024 * 1024
  }
  if(stringM.endsWith('Ei')){
    return Number(stringM.slice(0,-2)) * 1024 * 1024 * 1024 * 1024 * 1024
  }
}
function bytesToSize(bytes:number) {
    if (bytes === 0) return '0 Ki';
    let k = 1000, // or 1024
        sizes = ['Ki', 'Mi', 'Gi', 'Ti', 'Pi', 'Ei'],
        i = Math.floor(Math.log(bytes) / Math.log(k));
 
   return (bytes / Math.pow(k, i)).toPrecision(3) + ' ' + sizes[i];
}
function fetchCpuInfo(type: string, detail: any, conName: string){
  const info = fetchInfo(type, detail)
  const getPodInfo = podMetrics.value.filter((pod:any) => {
    return pod.metadata.name === info.name
  })
  if(getPodInfo && getPodInfo.length >= 1){
    let cpuNum = 0, memory = 0
    if(getPodInfo[0].containers && getPodInfo[0].containers.length >= 1){
      let conIndex = getPodInfo[0].containers.findIndex((con:any) => {
        return con.name === conName
      })
      cpuNum += Number(getPodInfo[0].containers[conIndex].usage.cpu.string.slice(0,-1)|| 0)
      const memoryNum = getmemoryKi(getPodInfo[0].containers[conIndex].usage.memory.string)
      memory += memoryNum||0
      return {
        cpu: fetchTime(cpuNum),
        memory: bytesToSize(memory)
      }
    }else{
      return {
        cpu: '0',
        memory: '0'
      }
    }

  }
  return {
    cpu: '0',
    memory: '0'
  }
}
function fetchrRsourcesInfo(type: string, detail: any, conName: string){
  let cons = detail.spec.containers
  let conIndex = cons.findIndex((c:any) => {
    return c.name === conName
  })
  let thisCon = cons[conIndex]
  let cpu = 0, memory= 0
  if(type === 'limits'){
    if(thisCon?.resources?.limits?.cpu?.string){
      cpu = thisCon?.resources?.limits?.cpu.string
    }
    if(thisCon?.resources?.limits?.memory?.string){
      memory = thisCon?.resources?.limits?.memory.string
    }
  }
  if(type === 'requests'){
    if(thisCon?.resources?.requests?.cpu?.string){
      cpu = thisCon?.resources?.requests?.cpu.string
    }
    if(thisCon?.resources?.requests?.memory?.string){
      memory = thisCon?.resources?.requests?.memory.string
    }
  }
  return {
    cpu, memory
  }
}
let logSeconds = ref(0)
function fetchHasContainer(pod: any){
  if(pod.status){
    return pod.status.containerStatuses
  }else{
    return []
  }
}
function fetchImage(coName: string, containers:any){
  const fetchIndex = containers.findIndex((con:any) => {
    return con.name === coName
  })
  let image = ''
  if(fetchIndex >= 0){
    image = containers[fetchIndex].image
  }else{
    return ''
  }
  let arr = image.split('/')
  if(arr.length <= 0){
    return ''
  }
  if(arr.length >= 2){
    return arr.splice(-2).join('/')
  }else{
    return arr[0]
  }
}
function fetchImageID(pod:any, conName:string){
  let imageStr = ''
  for(let consta of pod.spec.containers){
    if(consta.name === conName){
      imageStr = consta.image
    }
  }
  const findImageIndex = imageTags.value.findIndex((image:any) => {
    let repositoryTag = image.spec.repositoryMeta.projectName + '/' + image.spec.repositoryMeta.repositoryName + ':' + image.spec.tag
    const isRepTag = imageStr.endsWith(repositoryTag)
    return isRepTag
  })
  if(findImageIndex >= 0){
    const tagSpec = imageTags.value[findImageIndex].spec
    return {
      image: imageStr,
      id: tagSpec.dockerImage.sha256,
      author: tagSpec.dockerImage.author,
      time: formatTime(tagSpec.dockerImage.lastModifiedTime),
      branch: tagSpec.gitReference.branch,
      commitHash: tagSpec.gitReference.commitHash
    }
  }else{
    return {
      image: '',
      id: '',
      author: '',
      time: '',
      branch: '',
      commitHash: ''
    }
  }
}
function fetchImageIDshort(pod:any, conName:string){
  let allID = fetchImageID(pod, conName).id
  return allID.slice(0, 7)
}
function isAppConditions(pod: any){
  if(pod.status){
    if(pod.status.items){
      return true
    }
    return false
  }
  return false
}
function fetchHasConditions(pod: any){
  if(pod.status){
    if(pod.status.conditions){
      return pod.status.conditions
    }
    if(pod.status.items){
      const appConditions = []
      const openxType = conditionsType.value
      if(!openxType){
        for(let app in pod.status.items){
          appConditions.push({
            appName: app,
            conditions: pod.status.items[app].conditions
          })
        }
      }else if(openxType === 'deploymentStatus'){
        for(let app in pod.status.items){
          appConditions.push({
            appName: app,
            conditions: pod.status.items[app].deploymentStatus.conditions
          })
        }
      }else if(openxType === 'horizontalPodAutoscalerStatus'){
        for(let app in pod.status.items){
          appConditions.push({
            appName: app,
            conditions: pod.status.items[app].horizontalPodAutoscalerStatus.conditions
          })
        }
      }
      return appConditions
    }
    return []
  }else{
    return []
  }
}

const SecretList:any = ref([])
const currentPage = ref(1)
const loading = ref(true)
const filterName = ref('')
let nameSearch:any = ref(null)
let showOption:any = inject('showOption')
let timeInputFocus = ref(false)

let phaseType = ref('')
let eventType = ref('')
let searchProject = ref('')
let nodeName = ref('')
let volumeType = ref('all')

onMounted(() => {
  document.addEventListener('keydown', function (e) {
    const returnKeys = ['Enter', 'Control', 'Alt']
    if(editDialog.value || nsSelectDialog.value || deleteDialog.value || showPop.value || showOption.value) return
    if(timeInputFocus.value) return
    actionMenu.value = false
    if(returnKeys.includes(e.key)){
      return
    }
    if (e.keyCode && (e.ctrlKey || e.metaKey)) {
        return
    }
    if(e.key === 'Escape'){
      filterName.value = ''
      if(nameSearch.value){
        nameSearch.value.filterName = ''
        nameSearch.value.handleInputConfirm()
      }
      return
    }
    if(nameSearch.value){
      nameSearch.value.showInput()
    }    
  })
})

let showSearchProject = computed(() => {
  const nsGvk = route.path.split('/')
  return nsGvk[2] === 'core-v1-Pod'
})

let orderType = ref('createTime-ascending')
const serlistTime = computed(() => {
  let resultList = SecretList.value
  if(orderType.value.startsWith('createTime')){
    resultList.sort((itemL:any, itemR:any) => {
      const itemLTime = itemL.metadata.creationTimestamp.seconds
      const itemRTime = itemR.metadata.creationTimestamp.seconds
      if(orderType.value === 'createTime-ascending'){
        return itemRTime - itemLTime
      }else{
        return itemLTime - itemRTime
      }    
    })
  }
  if(orderType.value.startsWith('cpu') || orderType.value.startsWith('memory')){
    resultList.sort((itemL:any, itemR:any) => {
      let cpuL = getCpuNum('cpu/memory',itemL)
      let cpuR = getCpuNum('cpu/memory',itemR)
      if(orderType.value.endsWith('ascending')){
        return Number(cpuL) - Number(cpuR)
      }else{
        return Number(cpuR) - Number(cpuL)
      }
    })
  }
  if(orderType.value === 'name'){
      let podName = []
      for(let podRes of resultList){
        podName.push(podRes.metadata.name)
      }
      podName.sort()
      let resultSort = []
      for(let name of podName){
        const podIndex = resultList.findIndex((podOne:any) => {
          return podOne.metadata.name === name
        })
        resultSort.push(resultList[podIndex])
      }
      resultList = resultSort
  }
  if(orderType.value === 'restartCount'){
    resultList.sort((itemL:any, itemR:any) => {
      let reStartL = 0, reStartR = 0
      if(itemL.status.containerStatuses){
        for(let con of itemL.status.containerStatuses){
          reStartL += con.restartCount 
        }
      }
      if(itemR.status.containerStatuses){
        for(let conr of itemR.status.containerStatuses){
          reStartR += conr.restartCount 
        }
      }
      return reStartR - reStartL
    })
  }
  resultList.sort((itemL:any, itemR:any) => {
    let statusL = statusPrompt(itemL) ? Number(statusPrompt(itemL)?.sortLevel) : 0
    let statusR = statusPrompt(itemR) ? Number(statusPrompt(itemR)?.sortLevel) : 0
    return statusR - statusL
  })
  return resultList
})

function getCpuNum(type:any, detail:any){
  const info = fetchInfo(type, detail)
  const nsGvk = route.path.split('/')
  if(nsGvk[2] === 'core-v1-Node'){
    const getNodeInfo = nodeMetrics.value.filter((pod:any) => {
      return pod.metadata.name === info.name
    })
    if(getNodeInfo && getNodeInfo.length >= 1){
      let cpuNum = 0, memory = 0
      const node = getNodeInfo[0]
      cpuNum += Number(node.usage.cpu.string.slice(0,-1)|| 0)
      const memoryNum = getmemoryKi(node.usage.memory.string)
      memory += memoryNum||0
      return orderType.value.startsWith('cpu') ? cpuNum : memory
    }
  }else if(nsGvk[2] === 'core-v1-Pod'){
    const getPodInfo = podMetrics.value.filter((pod:any) => {
      return pod.metadata.name === info.name
    })
    if(getPodInfo && getPodInfo.length >= 1){
      let cpuNum = 0, memory = 0
      for(let con of getPodInfo[0].containers){
        cpuNum += Number(con.usage.cpu.string.slice(0,-1)|| 0)
        const memoryNum = getmemoryKi(con.usage.memory.string)
        memory += memoryNum||0
      }
      return orderType.value.startsWith('cpu') ? cpuNum : memory
    }
  }else{
    return 0
  }
}
const projectOptions = computed(() => {
  let proOpt: string[] = []
  for(let pod of SecretList.value){    
   let jingxInfo = getLabels(pod.metadata)
   if(jingxInfo?.length > 0){
     let findProjectIndex = jingxInfo.findIndex((info) => {
       return info.key === 'jingx-project.openx.neverdown.io'
     })
     let proValue = jingxInfo[findProjectIndex].value
     if(!proOpt.includes(proValue)){
       proOpt.push(proValue)
     }     
   }
  }
  return proOpt
})
const nodenameOptions = computed(() => {
  let proOpt: string[] = []
  for(let pod of SecretList.value){    
   let nameNameStr = pod.spec.nodeName
   if(!proOpt.includes(nameNameStr)){
       proOpt.push(nameNameStr)
     }
  }
  return proOpt
})
const volumeTypeOptions = computed(() => {
  let proOpt: string[] = ['all', '未挂载']
  for(let pod of SecretList.value){
   let timezone = ''
   let podCon = pod.spec.containers
   if(podCon && podCon[0]?.volumeMounts){
      let vmIndex =  pod.spec.containers[0]?.volumeMounts.findIndex((mount: any) => {
        return mount.mountPath === '/etc/localtime'
      })
      if(vmIndex >= 0){
        timezone = pod.spec.containers[0].volumeMounts[vmIndex].subPath
      }
   }

   if(timezone && !proOpt.includes(timezone)){
       proOpt.push(timezone)
     }
  }
  return proOpt
})

const SecretArr = computed(() => {
  const phase = phaseType.value
  const eventTypeVal = eventType.value
  const name = filterName.value
  const nodeNameIp = nodeName.value
  const timezone = volumeType.value

  const afterFilterName =  serlistTime.value.filter( (service: any) => {
    let filterNameArr = name.split(',')
    let hasStr = false
    for(let searchName of filterNameArr){
      if(service.metadata.name?.includes(searchName)){
        hasStr = true
      }
    }
    return hasStr
  }) || []
  const afterFilterProject =  afterFilterName.filter( (service: any) => {
    const nsGvk = route.path.split('/')
    if(nsGvk[2]  !== 'core-v1-Pod'){
      return true
    }
    if(service.metadata?.labels && !service.metadata.labels['jingx-project.openx.neverdown.io']){
      if(!searchProject.value){
        return true
      }
    }
    if(service.metadata?.labels && service.metadata.labels['jingx-project.openx.neverdown.io']){
      if(!searchProject.value){
        return true
      }
      let labelPro = service.metadata.labels['jingx-project.openx.neverdown.io']
      return labelPro === searchProject.value
    }else{
      return false
    }
  }) || []
  const afterType = afterFilterProject.filter( (service: any) => {
    if(!service.type) return true
    return service.type.includes(eventTypeVal)
  }) || []

  const afterNode = afterType.filter((podItem: any) => {
    if(!nodeNameIp) return true
    return podItem.spec.nodeName.includes(nodeNameIp)
  }) || []

  const afterTimezone = afterNode.filter( (service: any) => {
    if(timezone === 'all') return true
    let serTimezone = '未挂载'
    if(service.spec.containers){
      if(service.spec.containers[0]?.volumeMounts){
        let vmIndex =  service.spec.containers[0]?.volumeMounts.findIndex((mount: any) => {
          return mount.mountPath === '/etc/localtime'
        })
        if(vmIndex >= 0){
          serTimezone = service.spec.containers[0].volumeMounts[vmIndex].subPath
        }
      }
    }
    return timezone === serTimezone
  }) || []

  return afterTimezone.filter( (service: any) => {
    if(!service.status || !service.status.phase) return true
    return service.status.phase.includes(phase)
  }) || []
})

const paginationArr = computed(() => {
  const page = currentPage.value
  return SecretArr.value.slice((page - 1) * 10, page * 10)
})

onMounted(() => {
  getList()
})
watch(
  () => route.path,() => {
    loading.value = true
    getList()
  }
)

function verifyVerb(verb: string){
  const nsGvk = route.path.split('/')[2] || ''
  const gvList = JSON.parse(String(localStorage.getItem('gvkList')))
  if(!gvList) return false
  const thisVerbsIndex = gvList.findIndex((gv:any) => {
    return `${gv.gv}-${gv.kind}` === nsGvk
  })
  if(thisVerbsIndex < 0) return false
  return gvList[thisVerbsIndex].verbs.includes(verb)
}
function verifyCompare(){
  const nsGvk = route.path.split('/')[2] || ''
  let whatToShow = [
    'openx.neverdown.io-v1-Redis', 'openx.neverdown.io-v1-Mysql',
    'openx.neverdown.io-v1-Openx', 'core-v1-Pod'
  ]
  return whatToShow.includes(nsGvk)
}
function verifyVerbOther(verb: string, gvkString:string){
  const nsGvk = route.path.split('/')[2] || ''
  const gvList = JSON.parse(String(localStorage.getItem('gvkList')))
  if(!gvList) return false
  const thisVerbsIndex = gvList.findIndex((gv:any) => {
    return `${gv.gv}-${gv.kind}` === gvkString
  })
  if(thisVerbsIndex < 0){
    return false
  }
  return gvList[thisVerbsIndex].verbs.includes(verb)
}

let getList = function() {
  const nsGvk = route.path.split('/')
  const senddata = initSocketData('Request', nsGvk[1], nsGvk[2], 'list')
  sendSocketMessage(senddata, store)
  getConfigList()
  if(nsGvk[2]  === 'core-v1-Pod' && verifyVerbOther('list', 'metrics.k8s.io-v1beta1-PodMetrics')){
    getMetricsList()
  }
  if(nsGvk[2]  === 'core-v1-Node' && verifyVerbOther('list', 'metrics.k8s.io-v1beta1-NodeMetrics')){
    getNodeMetricsList()
  }
}
function verifyMetricsVerb(){
  const nsGvk = route.path.split('/')
  if(nsGvk[2]  === 'core-v1-Pod' && verifyVerbOther('list', 'metrics.k8s.io-v1beta1-PodMetrics')){
    return true
  }
  if(nsGvk[2]  === 'core-v1-Node' && verifyVerbOther('list', 'metrics.k8s.io-v1beta1-NodeMetrics')){
    return true
  }
  return false
}
let getConfigList = function() {
  const nsGvk = route.path.split('/')
  const haveConfigMapKind = [
    'apps-v1-Deployment', 'apps-v1-StatefulSet', 'core-v1-Pod', 'openx.neverdown.io-v1-Mysql',
    'openx.neverdown.io-v1-Redis', 'openx.neverdown.io-v1-Openx'
  ]
  if( !haveConfigMapKind.includes(nsGvk[1]) ){
    return
  }
  const senddata = initSocketData('Request', nsGvk[1], 'core-v1-ConfigMap', 'list')
  sendSocketMessage(senddata, store)  
}
let getImageList = function(gvk: string) {
  const nsGvk = route.path.split('/')
  const senddata = initSocketData('Request', nsGvk[1], gvk, 'list')
  sendSocketMessage(senddata, store)
}
let getMetricsList = function() {
  const nsGvk = route.path.split('/')
  const senddata = initSocketData('Request', nsGvk[1], 'metrics.k8s.io-v1beta1-PodMetrics', 'list')
  sendSocketMessage(senddata, store)
}
let getNodeMetricsList = function() {
  const nsGvk = route.path.split('/')
  const senddata = initSocketData('Request', nsGvk[1], 'metrics.k8s.io-v1beta1-NodeMetrics', 'list')
  sendSocketMessage(senddata, store)
}

let editDialog = ref(false)
let editType = ref('create')
let isfullscreen = ref(false)

provide('editType', editType)

function newItem(){
  const nsGvk = route.path.split('/')[2]
  currentItem.value = initObject(nsGvk)
  editDialog.value = true
  editType.value = 'create'
}

let currentItem = ref()
let compareItem:any = null
let dialogInit:any = ref(null)

let actionMenu = ref(false)
let actionMenuStyle = ref('')
let actionItem = ref()

function openAction(e:any, item:any){
  e.preventDefault();
  let evt = window.event || arguments[0];
  let offsetY = 20
  if((window.innerHeight - e.y) <= 220){
    offsetY = 240
  }
  actionItem.value = item
  actionMenuStyle.value = `top: ${evt.y - offsetY}px;left: ${evt.x - 20}px`
  actionMenu.value = !actionMenu.value
}

function getDetail(item: any){
  console.log('getDetail', item)
  editType.value = 'update'
  currentItem.value = cloneDeep(item)
  compareItem = cloneDeep(item)
  editDialog.value = true
}

function addCompare(item: any){
  let compareObj = store.state.user.compareItems
  let compareList = cloneDeep(compareObj)
  if(item.spec?.applications){
    for(let app of item.spec?.applications){
      let cloneApp = cloneDeep(app)
      cloneApp.metadata = {
        protoType: 'Application',
        namespace: item.metadata.namespace,
        name: item.metadata.name + '-' + app.appName
      }
      compareList.push(cloneApp)
    }
  }else if(item.spec?.master){
    let cloneApp = cloneDeep(item.spec?.master)
    cloneApp.metadata = {
      protoType: 'MysqlConfig',
      namespace: item.metadata.namespace,
      name: item.metadata.name + '-' + 'master'
    }
    let cloneAppS = cloneDeep(item.spec?.slave)
    cloneAppS.metadata = {
      protoType: 'MysqlConfig',
      namespace: item.metadata.namespace,
      name: item.metadata.name + '-' + 'slave'
    }
    compareList.push(cloneApp, cloneAppS)
  }else{
    item.metadata.protoType = 'Pod'
    compareList.push(cloneDeep(item))
  }
  store.dispatch("user/setCompareItems", compareList)
}

function needShowLabels(metaData: any){
  if(!metaData.labels){
    return false
  }
  for(let labelkey in metaData.labels ){
    if(labelkey.startsWith('jingx-repository') || labelkey.startsWith('jingx-tag') || labelkey.startsWith('jingx-project')){
      return true
    }
  }
  return false
}
function needShowLabelsWatchPolicy(metaData: any){
  if(!metaData.labels){
    return false
  }
  for(let labelkey in metaData.labels ){
    if(labelkey.startsWith('watch-policy.openx.neverdown.io')){
      return true
    }
  }
  return false
}
function getLabels(metaData: any){
  if(!metaData.labels){
    return []
  }
  interface label {
    key: string,
    value: string
  }
  let resLables: label[] = []
  for(let labelkey in metaData.labels ){
    if(labelkey.startsWith('jingx-repository') || labelkey.startsWith('jingx-tag') || labelkey.startsWith('jingx-project')){
     resLables.push({
       key: labelkey,
       value: metaData.labels[labelkey]
     })
    }
  }
  return resLables
}
function needShowAffinity(spec: any){
  if(spec?.affinity?.nodeAffinity?.requiredDuringSchedulingIgnoredDuringExecution?.nodeSelectorTerms){
    return true
  }else{
    return false
  }
}
function getAffinity(spec: any){
  if(spec?.affinity?.nodeAffinity?.requiredDuringSchedulingIgnoredDuringExecution?.nodeSelectorTerms && spec.affinity?.nodeAffinity?.requiredDuringSchedulingIgnoredDuringExecution?.nodeSelectorTerms.length > 0){
    return spec.affinity?.nodeAffinity?.requiredDuringSchedulingIgnoredDuringExecution?.nodeSelectorTerms[0].matchExpressions
  }else{
    return []
  }
}
function getLabelsPolicy(metaData: any){
  for(let labelkey in metaData.labels ){
    if(labelkey.startsWith('watch-policy.openx.neverdown.io')){
      let auto = 'manual'
      let watchPolicy = metaData.labels[labelkey]
      if( watchPolicy === 'in-place-upgrade'){
        auto = 'auto-inplace'
      }
      if( watchPolicy === 'rolling-upgrade'){
        auto = 'auto-rolling'
      }
      return {
        auto,
        watchPolicy
      }
    }
  }
  return {
    auto: '',
    watchPolicy: ''
  }
}

function exportDetail(item: any){
  const nsGvk = route.path.split('/')[2]
  const cloneItem = cloneDeep(item)
  let initItem = initObject(nsGvk)
  if(initItem.metadata){
    for(let metaIndex in initItem.metadata){
      initItem.metadata[metaIndex] = cloneDeep(cloneItem.metadata[metaIndex])
    }
  }
  if(initItem.spec){
    for(let metaIndex in initItem.spec){
      initItem.spec[metaIndex] = cloneDeep(cloneItem.spec[metaIndex])
    }
  }else{
    initItem = cloneItem
  }
  delete initItem.metadata.creationTimestamp

  const ns2Gvk = route.path.split('/')
  const gvkArr = ns2Gvk[2].split('-')
  let gvkObj = {
    group: gvkArr[0],
    version: gvkArr[1],
    kind: gvkArr[2],
  }
  const encodeOtem = encodeify(gvkObj, initItem)
  localStorage.setItem(nsGvk, binaryToStr(encodeOtem))
  ElNotification({title: '导出成功',message: 'success',type: 'success',duration: 2000})
}
import { saveAs } from 'file-saver'
import json2yaml from 'json2yaml'
import YAML from 'js-yaml'
function exportDetailYaml(item: any){
  const nsGvk = route.path.split('/')[2]
  const gvkArr = nsGvk.split('-')
  const cloneItem = cloneDeep(item)
  let fileName = item.metadata?.name || ''
  let avk = {
    apiVersion: `${gvkArr[0]}/${gvkArr[1]}`,
    kind: `${gvkArr[2]}`
  }
  const yamlData = json2yaml.stringify(Object.assign(avk, cloneItem))
  const str = new Blob([yamlData], { type: 'text/plain;charset=utf-8' })
  saveAs(str, nsGvk +'-'+ fileName + '(whole).yaml')
}
function exportDetailYamlSimp(item: any){
  const nsGvk = route.path.split('/')[2]
  const gvkArr = nsGvk.split('-')
  const cloneItem = cloneDeep(item)
  let initItem = initObject(nsGvk)
  let fileName = ''
  if(initItem.spec){
    for(let metaIndex in initItem.spec){
      initItem.spec[metaIndex] = cloneDeep(cloneItem.spec[metaIndex])
    }
  }else{    
    for(let itemInitKey in initItem){
      if(!(itemInitKey == 'spec' || itemInitKey == 'metadata')){
        initItem[itemInitKey] = cloneItem[itemInitKey]
      }
    }
  }
  if(initItem.rules){
    for(let rule of initItem.rules){
      delete rule.resourceNames
      delete rule.nonResourceURLs
    }
  }
  if(initItem.metadata){
    fileName = item.metadata.name
    for(let metaIndex in initItem.metadata){
      initItem.metadata[metaIndex] = cloneDeep(cloneItem.metadata[metaIndex])
    }
    let simpleMeta = {
      'name': initItem.metadata.name,
      'namespace': initItem.metadata.namespace,
      'labels': initItem.metadata.labels,
      'annotations': initItem.metadata.annotations
    }
    if(gvkArr[2] === 'ClusterRole' || gvkArr[2] === 'ClusterRoleBinding' || gvkArr[2] === 'ServiceAccount'){
      delete simpleMeta.labels
      delete simpleMeta.annotations
    }
    initItem.metadata = simpleMeta
  }
  let avk = {
    apiVersion: `${gvkArr[0]}/${gvkArr[1]}`,
    kind: `${gvkArr[2]}`
  }  
  const yamlData = json2yaml.stringify(Object.assign(avk, initItem))
  const str = new Blob([yamlData], { type: 'text/plain;charset=utf-8' })
  saveAs(str, nsGvk +'-'+ fileName + '.yaml')
}
function yamlList(){
  const nsGvk = route.path.split('/')
  const gvkArr = nsGvk[2].split('-')
  let gvkObj = {
    group: gvkArr[0],
    version: gvkArr[1],
    kind: gvkArr[2],
  }
  let allCheckedTemp = []
  for(let dataIndex in allChecked.value){
    const nsGvk2 = nsGvk[2]
    let cloneItem = cloneDeep(allChecked.value[dataIndex])
    let avk = {
      apiVersion: `${gvkArr[0]}/${gvkArr[1]}`,
      kind: `${gvkArr[2]}`
    }
    let initItem = Object.assign(avk, cloneDeep(initObject(nsGvk2)))

    if(initItem.metadata){
      for(let metaIndex in initItem.metadata){
        initItem.metadata[metaIndex] = cloneItem.metadata[metaIndex]
      }
    }
    if(initItem.spec){
      for(let metaIndex in initItem.spec){
        initItem.spec[metaIndex] = cloneItem.spec[metaIndex]
      }

    }else{
      for(let initKey in initItem){
        if(initKey != 'spec' && initKey != 'metadata'){
          initItem[initKey] = cloneItem[initKey]
        }
      }
    }
    delete initItem.metadata.creationTimestamp
    allCheckedTemp.push(initItem)
  }
  const yamlData = json2yaml.stringify(allCheckedTemp)
  const str = new Blob([yamlData], { type: 'text/plain;charset=utf-8' })
  saveAs(str, nsGvk[1]+'_'+nsGvk[2] + '_list.yaml')
}

let yamlDialog = ref(false)
function importYaml(){
  const nsName = route.path.split('/')[1]
  let nsList = []
  const cluster = JSON.parse(String(localStorage.getItem('clusterRole')))
  for(let ns in cluster){
    if(ns !== nsName){
       nsList.push(ns)
    }
  }  
  nsSelectList.value = nsList
  uploadList.value = []
  jsonForUpload.value = []
  yamlDialog.value = true
}
let uploadList = ref(<any>[])
let jsonForUpload = ref(<any>[])
function uploadConfig(param: any) {
  let reader = new FileReader()
  reader.readAsText(param.file, 'utf-8')
  reader.onloadend = (evt) => {
    if(evt.target?.readyState == FileReader.DONE){
      let jsonResult = YAML.load(evt.target.result)
      if(Array.isArray(jsonResult)){        
        jsonForUpload.value.push(...jsonResult)
      }else{
        jsonForUpload.value.push(jsonResult)
      }      
    }
  }
}
function yamlImport(){
  const nsGvk = route.path.split('/')
  const gvkArr = nsGvk[2].split('-')
  const ns = nsGvk[1]
  let gvkObj = {
    group: gvkArr[0],
    version: gvkArr[1],
    kind: gvkArr[2],
  }
  let allCheckedTemp = []
  for(let dataIndex in jsonForUpload.value){
    const nsGvk2 = nsGvk[2]
    let cloneItem = cloneDeep(jsonForUpload.value[dataIndex])
    const gvkClone = `${cloneItem.apiVersion}-${cloneItem.kind}`
    const gvkIndex = `${gvkObj.group}/${gvkObj.version}-${gvkObj.kind}`
    if(gvkClone != gvkIndex) return
    let initItem = cloneDeep(initObject(nsGvk2))

    if(initItem.metadata){
      for(let metaIndex in initItem.metadata){
        initItem.metadata[metaIndex] = cloneItem.metadata[metaIndex]
      }
    }
    if(initItem.spec){
      for(let metaIndex in initItem.spec){
        initItem.spec[metaIndex] = cloneItem.spec[metaIndex]
      }
      let copyAff = {}
      let copyTol = []
      if(gvkObj.kind === 'Openx'){
        const selArr = batchListAff.value.filter((aff:any) => {
          return aff.isChecked
        })
        if(selArr.length > 0){
          copyAff = cloneDeep(selArr[0].spec.affinity)
        }
        const selTol = batchListTol.value.filter((aff:any) => {
          return aff.isChecked
        })
        if(selTol.length > 0){
          for(let col of selTol){
            copyTol.push(cloneDeep(col.spec.toleration))
          }
        }
      }
      if(initItem.spec.applications && initItem.spec.applications.length > 0){
        for(let app of initItem.spec.applications){
          app.pod.spec.affinity = copyAff
          app.pod.spec.tolerations = copyTol
          for(let vol of app.pod.spec.volumes){
            if(vol.volumeSource.hostPath){
              if(vol.volumeSource.hostPath.path){
                let pathArr = vol.volumeSource.hostPath.path.split('/')
                if(pathArr.length >= 5){
                  pathArr[3] = ns
                }
                vol.volumeSource.hostPath.path = pathArr.join('/')
              }
            }
          }
          for(let container of app.pod.spec.containers){
            if(container.env.length){
              for(let oneEnv of container.env){
                if(oneEnv.name === 'BATTLEVERIFY_DYNAMIC_NAMESPACE'){
                  oneEnv.value = ns
                }
              }
            }
          }
        }
      }
      // mysql and redis
      if(initItem.spec.master && initItem.spec.slave){
        let appList = [
          initItem.spec.master, initItem.spec.slave
        ]
        for(let app of appList){
          app.pod.spec.affinity = copyAff
          app.pod.spec.tolerations = copyTol
        }
      }
    }else{
      for(let initKey in initItem){
        if(initKey != 'spec' && initKey != 'metadata'){
          initItem[initKey] = cloneItem[initKey]
        }
      }
    }
    delete initItem.metadata.creationTimestamp
    initItem.metadata.namespace = ns
    allCheckedTemp.push(initItem)
  }
  for(let checkedItem of allCheckedTemp){
    const param = updateSocketData(gvkObj, checkedItem)
    const senddata = initSocketData('Request', ns, nsGvk[2], 'create', param)
    sendSocketMessage(senddata, store)
  }
  yamlDialog.value = false
}
function handleChange(file:any, fileList: any) {
  uploadList.value = fileList
}

function checkPod(isCheck: boolean, item: any){
  item.isChecked = isCheck
}

let showExportList = computed(() => {
  for(let dataItem of SecretArr.value){
    if(dataItem.isChecked){
      return true
    }
  }
  return false
})


let nsSelectDialog = ref(false)
let nsSelectList = ref(<any>[])
let batchns = ref('')
let allChecked = ref(<any>[])
let selectDialogType = ref('export')
function exportList(selectType: string){
  selectDialogType.value = selectType
  volumes4Edit.value = []
  VolumeMounts4Edit.value = []
  annotations4Edit.value = []
  Envs4Edit.value = []
  batchns.value = ''
  const cluster = JSON.parse(String(localStorage.getItem('clusterRole')))
  const nsName = route.path.split('/')[1]
  let nsList = []
  for(let ns in cluster){
    if(ns !== nsName){
       nsList.push(ns)
    }
  }  
  nsSelectList.value = nsList

  let allCheckedTemp:any = []
  for(let dataItem of SecretArr.value){
    if(dataItem.isChecked){
      allCheckedTemp.push(dataItem)
    }
  }
  allChecked.value = allCheckedTemp
  nsSelectDialog.value = true
}
watch(()=> nsSelectDialog.value, (newNs: any) => {
  if(newNs === false){
    batchns.value = ''
  }
})
watch(()=> selectDialogType.value, (typeNew: any) => {
  if(typeNew === 'edit'){
    const ns = route.path.split('/')[1]
    const senddataAff = initSocketData('Request', ns, 'openx.neverdown.io-v1-Affinity', 'list')
    sendSocketMessage(senddataAff, store)
  }
})
function selectBatchNs(ns:string){
  batchns.value = ns
  const nsGvk = route.path.split('/')
  const gvkArr = nsGvk[2].split('-')
  let gvkObj = {
    group: gvkArr[0],
    version: gvkArr[1],
    kind: gvkArr[2],
  }
  if(gvkObj.kind === 'Openx' || gvkObj.kind === 'Mysql' || gvkObj.kind === 'Redis'){
    const senddataAff = initSocketData('Request', ns, 'openx.neverdown.io-v1-Affinity', 'list')
    sendSocketMessage(senddataAff, store)
    const senddataTol = initSocketData('Request', ns, 'openx.neverdown.io-v1-Toleration', 'list')
    sendSocketMessage(senddataTol, store)
  }
}
let batchEditNamespace = ref(<string[]>[])
function selectBatchEditNs(ns:string){
  const nsIdnex = batchEditNamespace.value.findIndex((editNs: string) => {
    return editNs === ns
  })
  if(nsIdnex >=0){
    batchEditNamespace.value.splice(nsIdnex, 1)
  }else{
    batchEditNamespace.value.push(ns)
  }
}
function batchEditNameSel(ns:string){
  return batchEditNamespace.value.includes(ns)
}
function importList(ns: string){
  if(!ns) return
  const nsGvk = route.path.split('/')
  const gvkArr = nsGvk[2].split('-')
  let gvkObj = {
    group: gvkArr[0],
    version: gvkArr[1],
    kind: gvkArr[2],
  }
  let allCheckedTemp = []
  for(let dataIndex in allChecked.value){
    const nsGvk2 = nsGvk[2]
    let cloneItem = cloneDeep(allChecked.value[dataIndex])
    let initItem = cloneDeep(initObject(nsGvk2))

    if(initItem.metadata){
      for(let metaIndex in initItem.metadata){
        initItem.metadata[metaIndex] = cloneItem.metadata[metaIndex]
      }
    }
    if(initItem.spec){
      for(let metaIndex in initItem.spec){
        initItem.spec[metaIndex] = cloneItem.spec[metaIndex]
      }
      let copyAff = {}
      let copyTol = []
      if(gvkObj.kind === 'Openx' || gvkObj.kind === 'Mysql' || gvkObj.kind === 'Redis'){
        const selArr = batchListAff.value.filter((aff:any) => {
          return aff.isChecked
        })
        if(selArr.length > 0){
          copyAff = cloneDeep(selArr[0].spec.affinity)
        }
        const selTol = batchListTol.value.filter((aff:any) => {
          return aff.isChecked
        })
        if(selTol.length > 0){
          for(let col of selTol){
            copyTol.push(cloneDeep(col.spec.toleration))
          }
        }
      }
      if(initItem.spec.applications && initItem.spec.applications.length > 0){
        for(let app of initItem.spec.applications){
          app.pod.spec.affinity = copyAff
          app.pod.spec.tolerations = copyTol
          for(let vol of app.pod.spec.volumes){
            if(vol.volumeSource.hostPath){
              if(vol.volumeSource.hostPath.path){
                let pathArr = vol.volumeSource.hostPath.path.split('/')
                if(pathArr.length >= 5){
                  pathArr[3] = ns
                }
                vol.volumeSource.hostPath.path = pathArr.join('/')
              }
            }
          }
          for(let container of app.pod.spec.containers){
            if(container.env.length){
              for(let oneEnv of container.env){
                if(oneEnv.name === 'BATTLEVERIFY_DYNAMIC_NAMESPACE'){
                  oneEnv.value = ns
                }
              }
            }
          }
        }
      }
      // mysql and redis
      if(initItem.spec.master && initItem.spec.slave){
        let appList = [
          initItem.spec.master, initItem.spec.slave
        ]
        for(let app of appList){
          app.pod.spec.affinity = copyAff
          app.pod.spec.tolerations = copyTol
        }
      }

    }else{
      for(let initKey in initItem){
        if(initKey != 'spec' && initKey != 'metadata'){
          initItem[initKey] = cloneItem[initKey]
        }
      }
    }
    delete initItem.metadata.creationTimestamp
    initItem.metadata.namespace = ns
    allCheckedTemp.push(initItem)
  }
  for(let checkedItem of allCheckedTemp){
    const param = updateSocketData(gvkObj, checkedItem)
    const senddata = initSocketData('Request', ns, nsGvk[2], 'create', param)
    sendSocketMessage(senddata, store)
  }
  ElNotification({title: '导出完成',message: 'success',type: 'success',duration: 2000})
  nsSelectDialog.value = false
}

function importConfig(){
  const nsGvk = route.path.split('/')[2]
  const ns = route.path.split('/')[1]
  if(localStorage.getItem(nsGvk)){
    let encodeOtem = String(localStorage.getItem(nsGvk))
    const buffer = strToBinary(encodeOtem)
    const gvkArr = nsGvk.split('-')
    let gvkObj = {
      group: gvkArr[0],
      version: gvkArr[1],
      kind: gvkArr[2],
    }
    const item = decodeify(gvkObj, buffer)
    item.metadata.namespace = ns
    if(item.spec){
      if(item.spec.applications && item.spec.applications.length > 0){
        for(let app of item.spec.applications){
          for(let vol of app.pod.spec.volumes){
            if(vol.volumeSource.hostPath){
              if(vol.volumeSource.hostPath.path){
                let pathArr = vol.volumeSource.hostPath.path.split('/')
                if(pathArr.length >= 5){
                  pathArr[3] = ns
                }
                vol.volumeSource.hostPath.path = pathArr.join('/')
              }
            }
          }
          for(let container of app.pod.spec.containers){
            if(container.env.length){
              for(let oneEnv of container.env){
                if(oneEnv.name === 'BATTLEVERIFY_DYNAMIC_NAMESPACE'){
                  oneEnv.value = ns
                }
              }
            }
          }
        }
      }
    }
    delete item.metadata.creationTimestamp
    delete item.metadata.resourceVersion
    delete item.metadata.uid
    currentItem.value = cloneDeep(item)
  }else{
    currentItem.value = initObject(nsGvk)
  }
  editDialog.value = true
  editType.value = 'create'
}

function goApp(item: any, info: any){
  dialogInit.value = {
    appName: info.appName
  }
  getDetail(item)
}

let popStyle = ref('')
let showPop = ref(false)
let delItem = ref(null)
let deleteItem = function(item: any, e:any){
  let offsetY = 20
  if((window.innerHeight - e.y) <= 100){
    offsetY = -100
  }
  showPop.value = !showPop.value
  delItem.value = item
  popStyle.value = `top: ${e.y + offsetY}px;left: ${e.x - 100}px`
}
window.addEventListener('click', () => {
  showPop.value = false
  actionMenu.value = false
})
document.addEventListener('keydown', function (e) {
  if(e.key === 'Escape'){
    if(editDialog.value){
      editDialog.value = false
    }
    
  }
})
function popClick(){}
let checkAllPod = ref(false)
watch(() => {
  let checkedArr = SecretArr.value.filter((ser:any) => {
    return ser.isChecked
  })
  return checkedArr.length
}, ()=>{
  let checkedArr = SecretArr.value.filter((ser:any) => {
    return ser.isChecked
  })
  if(checkedArr.length === SecretArr.value.length){
    checkAllPod.value = true
  }else{
    checkAllPod.value = false
  }
})
function checkAllPodChange(e:boolean){
  if(e){
    for(let onePod of SecretArr.value){
      onePod.isChecked = true
    }
  }else{
    for(let onePod of SecretArr.value){
      onePod.isChecked = false
    }
  }
}
function getGvkGroup(){
  const nsGvk = route.path.split('/')
  const gvkArr = nsGvk[2].split('-')
  let gvkObj = {
    group: gvkArr[0],
    version: gvkArr[1],
    kind: gvkArr[2],
  }
  return gvkObj
}
function batchDelete(){
  if(ref_batchDeleteOneByOne.value.isNeed){
    void batchDeleteOneByOne();
  }else{
    nsSelectDialog.value = false
    for(let item of allChecked.value){
      itemDelete(item)
    }
  }
}
/* 逐个删除, 因为同时删除, 会导致卡死 */
const ref_batchDeleteOneByOne = ref<{
  gapSeconds:number, //每次删除的间隔时间
  isNeed: boolean, //是否需要 逐个批量删除
  isDoing: boolean,//是否正在删除中
  idxDoing:number,//正在删除的序号
  sum:number,//总数
}>({
  gapSeconds: 2,
  isNeed: false,
  isDoing:false,
  idxDoing:0,
  sum:0,
})

async function batchDeleteOneByOne(){
  let sendDataArr:any[] = []
  for(let item of allChecked.value){
    sendDataArr.push(parseItemSendData(item));
  }
  /*  */
  ref_batchDeleteOneByOne.value.isDoing = true;
  ref_batchDeleteOneByOne.value.sum = allChecked.value.length;
  ref_batchDeleteOneByOne.value.idxDoing = 0;
  for(let sendData of sendDataArr){
    ref_batchDeleteOneByOne.value.idxDoing++;
    console.log(`delete sendData [${ref_batchDeleteOneByOne.value.idxDoing}/${ref_batchDeleteOneByOne.value.sum}]`,sendData)
    sendSocketMessage(sendData, store)
    // if (ref_batchDeleteOneByOne.value.idxDoing < ref_batchDeleteOneByOne.value.sum) {//可以不做这个判断, 没有这个判断, 全做完后也会等几秒, 能看到最后一次数值变化, 体验好一些
    await TimerUtil.sleepSeconds(ref_batchDeleteOneByOne.value.gapSeconds);
    // }
  }
  ref_batchDeleteOneByOne.value.isDoing = false;
  nsSelectDialog.value = false
}
function parseItemSendData(item: any)
{
  const nsGvk = route.path.split('/')
  const gvkObj = getGvkGroup()
  const param = deleteSocketData(gvkObj, item)
  const senddata = initSocketData('Request', nsGvk[1], nsGvk[2], 'delete', param)
  return senddata;
}
function itemDelete(item: any){
  const senddata = parseItemSendData(item)
  sendSocketMessage(senddata, store)
  deleteDialog.value = false
}
let deleteDialog = ref(false)
function confirmDelete(item: any){
  const gvkObj = getGvkGroup()
  actionMenu.value = false
  showPop.value = false
  if(gvkObj.kind === 'Openx'){
    deleteDialog.value = true
  }else{
    itemDelete(item)
  }
}
// 防抖
const confirmEdit = debounce(sendEdit, 2000, {leading: true, trailing: false})
function sendEdit(){
  const nsGvk = route.path.split('/')
  sendEditInfo(currentItem.value, editType.value, nsGvk[1])
}
function sendEditInfo(item:any, type: string, namespace:string){
  const nsGvk = route.path.split('/')
  const gvkArr = nsGvk[2].split('-')
  let gvkObj = {
    group: gvkArr[0],
    version: gvkArr[1],
    kind: gvkArr[2],
  }
  const param = updateSocketData(gvkObj, item)
  if(!param){ // 表单校验失败
    return
  }
  const senddata = initSocketData('Request', namespace, nsGvk[2], type, param)
  sendSocketMessage(senddata, store)
  editDialog.value = false
}
function editPod(item:any, appInfo:any, rep: number){
  const ns = route.path.split('/')[1]
  if(item.spec.applications){
    let appIndex = item.spec.applications.findIndex((app:any) => {
      return app.appName === appInfo.appName
    })
    item.spec.applications[appIndex].replicas = rep
    sendEditInfo(item, 'update', ns)
  }else{
    item.spec[appInfo.appName].replicas = rep
    sendEditInfo(item, 'update', ns)
  }
}
function fetchInfoHasHpa(item:any, appInfo:any){
  if(item === 'init'){
    for(let appName in appInfo.status.items){
      const appIndex = appInfo.spec.applications.findIndex((app :any) => {
        let metadataName = appInfo.metadata.name
        return `${metadataName}-${app.appName}` === appName
      })
      if(appIndex >= 0 && appInfo.spec.applications[appIndex].horizontalPodAutoscalerSpec){          
        return false
      }
    }
  }
  return true
}
let currentTabForm = function() {
  return currentPageVue(route.path)
}

let batchListAff = ref(<any>[])
let batchListTol = ref(<any>[])
watch(() => store.state.socket.socket.isConnected, (status) => {
  if(status){
    getList()
  }
})
watch(
  () => store.state.socket.socket.message,
  (msg) => {
    // console.log("watch message: " + new Date());
    const nsGvk = route.path.split('/')
    const gvkArr = nsGvk[2].split('-')
    let gvkObj = {
      group: gvkArr[0],
      version: gvkArr[1],
      kind: gvkArr[2],
    }
    let cfgGvk = {
      group: 'core',
      version: 'v1',
      kind: 'ConfigMap'        
    }
    let getConfigMap = returnResource(msg, nsGvk[1], cfgGvk, loadOver)

    let metricsGvk = {
      group: 'metrics.k8s.io',
      version: 'v1beta1',
      kind: 'PodMetrics'
    }
    let getMetrics = returnResource(msg, nsGvk[1], metricsGvk, loadOver, updateMetricsWatch)

    let nodeMetricsGvk = {
      group: 'metrics.k8s.io',
      version: 'v1beta1',
      kind: 'NodeMetrics'
    }
    let getNodeMetrics = returnResource(msg, nsGvk[1], nodeMetricsGvk, loadOver, updateNodeMetricsWatch)

    let imageTagGvk = {
      group: 'jingx',
      version: 'v1',
      kind: 'Tag'
    }
    let getImageTags = returnResource(msg, nsGvk[1], imageTagGvk, ()=>{}, ()=>{})

    if(getConfigMap && getConfigMap.length >0){
      ConfigMapList.value = getConfigMap
    }
    if(getMetrics && getMetrics.length >0){
      podMetrics.value = getMetrics
    }
    if(getNodeMetrics && getNodeMetrics.length >0){
      nodeMetrics.value = getNodeMetrics
    }
    if(getImageTags && getImageTags.length >0){
      imageTags.value = getImageTags
    }
    if(batchns.value){
      let affGvk = {
        group: 'openx.neverdown.io',
        version: 'v1',
        kind: 'Affinity'
      }
      let tolGvk = {
        group: 'openx.neverdown.io',
        version: 'v1',
        kind: 'Toleration'
      }
      let affResultList = returnResource(msg, batchns.value, affGvk, loadOver)
      let tolResultList = returnResource(msg, batchns.value, tolGvk, loadOver)
      if(affResultList){
        affResultList.sort((itemL:any, itemR:any) => {
          const itemLTime = itemL.metadata.creationTimestamp.seconds
          const itemRTime = itemR.metadata.creationTimestamp.seconds
          return itemRTime - itemLTime
        })
        for(let aff of affResultList){
          aff.isChecked = false
        }
        batchListAff.value = affResultList
      }
      if(tolResultList){
        tolResultList.sort((itemL:any, itemR:any) => {
          const itemLTime = itemL.metadata.creationTimestamp.seconds
          const itemRTime = itemR.metadata.creationTimestamp.seconds
          return itemRTime - itemLTime
        })
        for(let tol of tolResultList){
          tol.isChecked = false
        }
        batchListTol.value = tolResultList
      }
    }
    if(nsSelectDialog.value && selectDialogType.value === 'edit'){
      let affGvk = {
        group: 'openx.neverdown.io',
        version: 'v1',
        kind: 'Affinity'
      }
      let editNs = route.path.split('/')[1]
      let affResultList = returnResource(msg, editNs, affGvk, loadOver)
      if(affResultList){
        affResultList.sort((itemL:any, itemR:any) => {
          const itemLTime = itemL.metadata.creationTimestamp.seconds
          const itemRTime = itemR.metadata.creationTimestamp.seconds
          return itemRTime - itemLTime
        })
        for(let aff of affResultList){
          aff.isChecked = false
        }
        batchListAff.value = affResultList
      }
    }

      let resultList = returnResource(msg, nsGvk[1], gvkObj, loadOver, updateForWatch, getList)
      if(resultList){
        SecretList.value = resultList
      }
  }
)

function updateForWatch(type: string, updateRaw: any){
  if(type === 'ADDED'){
    SecretList.value.unshift(updateRaw)
  }
  if(type === 'MODIFIED'){
    const modName = updateRaw.metadata.name
    const modIndex = SecretList.value.findIndex((ser: any) => {
      return ser.metadata.name === modName
    })
    if(modIndex >= 0){
      SecretList.value[modIndex] = updateRaw
    }    
  }
  if(type === 'DELETED'){
    const modName = updateRaw.metadata.name
    const modIndex = SecretList.value.findIndex((ser: any) => {
      return ser.metadata.name === modName
    })
    if(modIndex >= 0){
      SecretList.value.splice(modIndex, 1)
    } 
  }
}
let podMetrics = ref(<any>[])
let nodeMetrics = ref(<any>[])
function updateMetricsWatch(type: string, updateRaw: any){
  if(type === 'ADDED'){
    podMetrics.value.unshift(updateRaw)
  }
  if(type === 'MODIFIED'){
    const modName = updateRaw.metadata.name
    const modIndex = podMetrics.value.findIndex((ser: any) => {
      return ser.metadata.name === modName
    })
    if(modIndex >= 0){
      podMetrics.value[modIndex] = updateRaw
    }    
  }
  if(type === 'DELETED'){
    const modName = updateRaw.metadata.name
    const modIndex = podMetrics.value.findIndex((ser: any) => {
      return ser.metadata.name === modName
    })
    if(modIndex >= 0){
      podMetrics.value.splice(modIndex, 1)
    } 
  }
}
function updateNodeMetricsWatch(type: string, updateRaw: any){
  if(type === 'ADDED'){
    nodeMetrics.value.unshift(updateRaw)
  }
  if(type === 'MODIFIED'){
    const modName = updateRaw.metadata.name
    const modIndex = nodeMetrics.value.findIndex((ser: any) => {
      return ser.metadata.name === modName
    })
    if(modIndex >= 0){
      nodeMetrics.value[modIndex] = updateRaw
    }    
  }
  if(type === 'DELETED'){
    const modName = updateRaw.metadata.name
    const modIndex = nodeMetrics.value.findIndex((ser: any) => {
      return ser.metadata.name === modName
    })
    if(modIndex >= 0){
      nodeMetrics.value.splice(modIndex, 1)
    } 
  }
}

import { openTerm, download } from './podutils'
async function goTerm(metadata:any, containerdata:any, time: number){
  await openTerm(metadata, containerdata, 'log', time)
}
async function goBash(metadata:any, containerdata:any){
  await openTerm(metadata, containerdata, 'bash', 0)
}
async function downloadLog(metadata:any, containerdata:any,isprevious: boolean, time: number){
  if(isprevious){
    containerdata.previousDownloading = true
  }else{
    containerdata.currentDownloading = true
  }
  await download(metadata, containerdata, isprevious, time)
}

let loadOver = function(){
  loading.value = false
}
function filterChange(name: string){
  filterName.value = name
}

let volumes4Edit:any = ref([])
let VolumeMounts4Edit:any = ref([])
let Envs4Edit:any = ref([])
let ConfigMapList = ref([])
let batchEditType = ref('add')
let annotations4Edit = ref(<any>[])
let imageProject4Edit = ref('')
let imageTag4Edit = ref('')
let policy4Edit = ref('')
let aff4Edit = ref(<any>[])

function changeBatchType(e:any){
  batchEditType.value = e
}
function changeBatchProject(e: string){
  imageProject4Edit.value = e
}
function changeBatchTag(e: string){
  imageTag4Edit.value = e
}
function changeBatchPolicy(e: string){
  policy4Edit.value = e
}
function appUpdate4Batch(app: any){
      // volume
      for(let vol of volumes4Edit.value){
        const appVolIndex = app.pod.spec.volumes.findIndex((findvol:any) => {
          return findvol.name === vol.name
        })
        if(appVolIndex >= 0){
          app.pod.spec.volumes.splice(appVolIndex, 1)
        }          
      }
      if(batchEditType.value === 'add'){
        app.pod.spec.volumes.push(...volumes4Edit.value)
      }
      // policy
      if(policy4Edit.value){
        app.watchPolicy = toRaw(policy4Edit.value)
      }
      // annotations
      for(let anno of annotations4Edit.value){
        if(anno.key){            
          if(batchEditType.value === 'add'){
            app.pod.metadata.annotations[anno.key] = anno.value
          }else{
            delete app.pod.metadata.annotations[anno.key]
          }
        }
      }
      // aff
      if(aff4Edit.value && batchListAff.value.length > 0){
        let affIndex = batchListAff.value.findIndex((oneAff: any) => {
          return oneAff.metadata.name === aff4Edit.value
        })
        if(affIndex >= 0){
          app.pod.spec.affinity = batchListAff.value[affIndex].spec.affinity
        }        
      }
      // mount
      for(let con of app.pod.spec.containers){
        for(let volmon of VolumeMounts4Edit.value){
          const mountIndex = con.volumeMounts.findIndex((findvol:any) => {
            return findvol.mountPath === volmon.mountPath && findvol.name === volmon.name
          })
          if(mountIndex >= 0){
            con.volumeMounts.splice(mountIndex, 1)
          }
        }
        if(batchEditType.value === 'add'){
          con.volumeMounts.push(...VolumeMounts4Edit.value)
        }

        for(let envitem of Envs4Edit.value){
          const envIndex = con.env.findIndex((findenv:any) => {
            return findenv.name === envitem.name
          })
          if(envIndex >= 0){
            con.env.splice(envIndex, 1)
          }
        }
        if(batchEditType.value === 'add'){
          con.env.push(...Envs4Edit.value)
        }          
      }
      // image
      for(let con of app.pod.spec.containers){
        let imageArr = con.image.split('/')
        if(imageProject4Edit.value && imageArr[1]){
          imageArr[1] = imageProject4Edit.value
        }
        if(imageTag4Edit.value && imageArr[2]){
          let reTagArr = imageArr[2].split(':')
          if(reTagArr[1]){
            reTagArr[1] = imageTag4Edit.value
          }
          imageArr[2] = reTagArr.join(':')
        }
        con.image = imageArr.join('/')
      }
}
function changeBatchAff(aff: any){
  aff4Edit.value = aff
}

async function editList(){
  let allCheckedTemp = []
  const ns = route.path.split('/')[1]
  for(let Index in allChecked.value){
    let dataIndex = cloneDeep(allChecked.value[Index])
    if(dataIndex.spec && dataIndex.spec.applications){
      for(let app of dataIndex.spec.applications){
        appUpdate4Batch(app)
      }
    }
    if(dataIndex.spec){
      if(dataIndex.spec.master){
        appUpdate4Batch(dataIndex.spec.master)
      }
      if(dataIndex.spec.slave){
        appUpdate4Batch(dataIndex.spec.slave)
      }
    }
    allCheckedTemp.push(dataIndex)
  }
  for(let checked of allCheckedTemp){
    sendEditInfo(checked, 'update', ns)
  }
  if(batchEditNamespace.value.length <= 0){
    nsSelectDialog.value = false
    return
  }
  const allBatchNs = cloneDeep(batchEditNamespace.value)
  for(let names of allBatchNs){
    getListFromNs(names)
  }
}
async function editPodBatch(rep:number){
  let allCheckedTemp = []
  const ns = route.path.split('/')[1]
  for(let Index in allChecked.value){
    let dataIndex = cloneDeep(allChecked.value[Index])
    if(dataIndex.spec && dataIndex.spec.applications){
      for(let app of dataIndex.spec.applications){
        app.replicas = rep
      }
    }
    if(dataIndex.spec && dataIndex.spec.master){
      dataIndex.spec.master.replicas = rep
    }
    if(dataIndex.spec && dataIndex.spec.slave){
      dataIndex.spec.slave.replicas = rep
    }
    allCheckedTemp.push(dataIndex)
  }  
  for(let checked of allCheckedTemp){
    sendEditInfo(checked, 'update', ns)
  }
  nsSelectDialog.value = false
}
function sendBatchEditInfo(item:any, namespace:string){
  const nsGvk = route.path.split('/')
  const gvkArr = nsGvk[2].split('-')
  let gvkObj = {
    group: gvkArr[0],
    version: gvkArr[1],
    kind: gvkArr[2],
  }

  delete item.metadata.creationTimestamp
  delete item.metadata.resourceVersion
  delete item.metadata.uid

  item.metadata.namespace = namespace
  if(item.spec){
    if(item.spec.applications && item.spec.applications.length > 0){
      for(let app of item.spec.applications){
        for(let vol of app.pod.spec.volumes){
          if(vol.volumeSource.hostPath){
            if(vol.volumeSource.hostPath.path){
              let pathArr = vol.volumeSource.hostPath.path.split('/')
              if(pathArr.length >= 5){
                pathArr[3] = namespace
              }
              vol.volumeSource.hostPath.path = pathArr.join('/')
            }
          }
        }
        for(let container of app.pod.spec.containers){
          if(container.env.length){
            for(let oneEnv of container.env){
              if(oneEnv.name === 'BATTLEVERIFY_DYNAMIC_NAMESPACE'){
                oneEnv.value = namespace
              }
            }
          }
        }
      }
    }
  }

  const param = updateSocketData(gvkObj, item)

  if(!param){ // 表单校验失败
    return
  }
  const senddata = initSocketData('Request', namespace, nsGvk[2], 'delete', param)
  sendSocketMessage(senddata, store)
  const sendCreate = initSocketData('Request', namespace, nsGvk[2], 'create', param)
  setTimeout(() => {
    sendSocketMessage(sendCreate, store)
  }, 200)
}

function getListFromNs(ns: string){
  const nsGvk = route.path.split('/')
  const senddata = initSocketData('Request', ns, nsGvk[2], 'list')
  sendSocketMessage(senddata, store)
}

const protoRequest = proto.github.com.kzz45.neverdown.pkg.openx.aggregator.proto
watch(
  () => store.state.socket.socket.message, (msg) => {
    const result:any = protoRequest.Response.decode(msg)
    let resultGvk = result.groupVersionKind
    const gvkArr = route.path.split('/')[2].split('-')
    let gvkObj = { group: gvkArr[0], version: gvkArr[1], kind: gvkArr[2], }
    const gvkStr = `${resultGvk.group}-${resultGvk.version}-${resultGvk.kind}`
    const requestGvk = `${gvkObj.group}-${gvkObj.version}-${gvkObj.kind}`

    if(batchEditNamespace.value.includes(result.namespace) && gvkStr === requestGvk){
      const delIndex = batchEditNamespace.value.findIndex(nsb => {
        return nsb === result.namespace
      })
      batchEditNamespace.value.splice(delIndex, 1)
      let resultList = returnResourceList(msg, result.namespace, gvkObj)
      editListFromOther(result.namespace, resultList)
    }

  }
)
watch(()=> batchEditNamespace.value.length, (len) => {
  if(len === 0){
    nsSelectDialog.value = false
  }
})
function editListFromOther(ns:string, nspods: any){
  let allCheckedTemp = []
  let nsChecked = []
  for(let pod of allChecked.value){
    let findPod = nspods.findIndex((nspod:any) => {
      return nspod.metadata.name === pod.metadata.name
    })
    if(findPod >= 0) nsChecked.push(nspods[findPod])
  }
  for(let Index in nsChecked){
    let dataIndex = cloneDeep(nsChecked[Index])
    if(dataIndex.spec && dataIndex.spec.applications){
      for(let app of dataIndex.spec.applications){
        // volume
        for(let vol of volumes4Edit.value){
          const appVolIndex = app.pod.spec.volumes.findIndex((findvol:any) => {
            return findvol.name === vol.name
          })
          if(appVolIndex >= 0){
            app.pod.spec.volumes.splice(appVolIndex, 1)
          }          
        }
        if(batchEditType.value === 'add'){
          for(let volOfEdit of volumes4Edit.value){
            let copyVol = cloneDeep(volOfEdit)
            if(copyVol.volumeSource.hostPath){
              if(copyVol.volumeSource.hostPath.path){
                let pathArr = copyVol.volumeSource.hostPath.path.split('/')
                if(pathArr.length >= 5){
                  pathArr[3] = ns
                }
                copyVol.volumeSource.hostPath.path = pathArr.join('/')
              }
            }
            app.pod.spec.volumes.push(copyVol)
          }
        }
        // annotations
        for(let anno of annotations4Edit.value){
          if(anno.key){            
            if(batchEditType.value === 'add'){
              app.pod.metadata.annotations[anno.key] = anno.value
            }else{
              delete app.pod.metadata.annotations[anno.key]
            }
          }
        }
        // mount
        for(let con of app.pod.spec.containers){
          for(let volmon of VolumeMounts4Edit.value){
            const mountIndex = con.volumeMounts.findIndex((findvol:any) => {
              return findvol.mountPath === volmon.mountPath && findvol.name === volmon.name
            })
            if(mountIndex >= 0){
              con.volumeMounts.splice(mountIndex, 1)
            }
          }
          if(batchEditType.value === 'add'){
            con.volumeMounts.push(...VolumeMounts4Edit.value)
          }
          //Env
          for(let envitem of Envs4Edit.value){
            const envIndex = con.env.findIndex((findenv:any) => {
              return findenv.name === envitem.name
            })
            if(envIndex >= 0){
              con.env.splice(envIndex, 1)
            }
          }
          if(batchEditType.value === 'add'){
            con.env.push(...Envs4Edit.value)
          }
          //Image
          if(imageProject4Edit.value){
            let imageArr = con.image.split('/')
            if(imageProject4Edit.value && imageArr[1]){
              imageArr[1] = imageProject4Edit.value
            }
            if(imageTag4Edit.value && imageArr[2]){
              let reTagArr = imageArr[2].split(':')
              if(reTagArr[1]){
                reTagArr[1] = imageTag4Edit.value
              }
              imageArr[2] = reTagArr.join(':')
            }
            con.image = imageArr.join('/')
          }
        }
      }
    }
    allCheckedTemp.push(dataIndex)
  }
  for(let checked of allCheckedTemp){
    sendEditInfo(checked, 'update', ns)
  }
}

let sortOptions = computed(() => {
  const nsGvk = route.path.split('/')
  if(nsGvk[2]  === 'core-v1-Pod'){
    return [
      {label: 'name', value: 'name'},
      {label: 'restartCount', value: 'restartCount'},
      {label: 'cpu升序', value: 'cpu-ascending'},
      {label: 'cpu降序', value: 'cpu-descending'},
      {label: 'memory升序', value: 'memory-ascending'},
      {label: 'memory降序', value: 'memory-descending'},
      {label: 'createTime升序', value: 'createTime-ascending'},
      {label: 'createTime降序', value: 'createTime-descending'}
    ]
  }else{
    return [
      {label: 'name', value: 'name'},
      {label: 'createTime升序', value: 'createTime-ascending'},
      {label: 'createTime降序', value: 'createTime-descending'}
    ]
  }

})
const promptStyle = {
  'background-image':'linear-gradient(270deg,#fef0f0,#fff,#fef0f0)',
  'border': '1px solid red',
  'sortLevel': 1
}
function statusPrompt(pod: any){
  const statusItems = pod?.status?.items
  const objectName = pod?.metadata.name
  if(statusItems){
    for(let itemKey in statusItems){
      const deploymentStatus = statusItems[itemKey]?.deploymentStatus
      let appStatusName = itemKey.replace(`${objectName}-`, '')
      const replicasIndex = pod.spec.applications.findIndex((app:any) => {
        return app.appName === appStatusName
      })
      const replicas = pod.spec.applications[replicasIndex].replicas
      let equalFlg = (replicas === deploymentStatus.replicas && deploymentStatus.replicas === deploymentStatus.readyReplicas && deploymentStatus.readyReplicas === deploymentStatus.updatedReplicas && deploymentStatus.updatedReplicas === deploymentStatus.availableReplicas)
      if(pod.spec.applications[replicasIndex].horizontalPodAutoscalerSpec){
        const horizontalPodAutoscalerStatus = statusItems[itemKey]?.horizontalPodAutoscalerStatus
        let hapEqual = horizontalPodAutoscalerStatus.currentReplicas === horizontalPodAutoscalerStatus.desiredReplicas
        equalFlg = hapEqual && deploymentStatus.replicas === deploymentStatus.readyReplicas && deploymentStatus.readyReplicas === deploymentStatus.updatedReplicas && deploymentStatus.updatedReplicas === deploymentStatus.availableReplicas
      }
      if(!equalFlg){
        return promptStyle
      }
    }
  }
  const statusMaster = pod?.status?.master
  if(statusMaster){
    let masterReplicas = pod.spec.master.replicas
    let replicas = statusMaster.replicas, readyReplicas = statusMaster.readyReplicas, updatedReplicas = statusMaster.updatedReplicas, currentReplicas = statusMaster.currentReplicas
    let equalFlg = (masterReplicas === replicas && replicas === readyReplicas && readyReplicas === updatedReplicas && updatedReplicas === currentReplicas)
    if(!equalFlg){
      return promptStyle
    }
  }

  const statusSlave = pod?.status?.slave
  if(statusSlave){
    let slaveReplicas = pod.spec.slave.replicas
    let replicas = statusSlave.replicas, readyReplicas = statusSlave.readyReplicas, updatedReplicas = statusSlave.updatedReplicas, currentReplicas = statusSlave.currentReplicas
    let equalFlg = (slaveReplicas === replicas && replicas === readyReplicas && readyReplicas === updatedReplicas && updatedReplicas === currentReplicas)
    if(!equalFlg){
      return promptStyle
    }
  }

  const unavailableReplicas = pod?.status?.unavailableReplicas
  if(unavailableReplicas === 0 || unavailableReplicas){
    const deploymentStatus = pod?.status
    let equalFlg = (deploymentStatus.replicas === deploymentStatus.readyReplicas && deploymentStatus.readyReplicas === deploymentStatus.updatedReplicas && deploymentStatus.updatedReplicas === deploymentStatus.availableReplicas)
    if(!equalFlg){
      return promptStyle
    }
  }

  const currentReplicas = pod?.status?.currentReplicas
  if(currentReplicas === 0 || currentReplicas){
    const statefulSetStatus = pod?.status
    let replicas = statefulSetStatus.replicas, readyReplicas = statefulSetStatus.readyReplicas, updatedReplicas = statefulSetStatus.updatedReplicas, currentReplicas = statefulSetStatus.currentReplicas
    let equalFlg = (replicas === readyReplicas && readyReplicas === updatedReplicas && updatedReplicas === currentReplicas)
    if(!equalFlg){
      return promptStyle
    }
  }
  if(pod.status?.phase){
    if(pod.status.phase === 'Pending'){
      return {
        'background-image':'linear-gradient(270deg,#fef0f0,#fff,#fef0f0)',
        'border': '1px solid red',
        'sortLevel': 3
      }
    }
  }
  if(pod.status?.containerStatuses){
    const podStartCount = startCount(pod.status.containerStatuses)
    if(Number(podStartCount) > 0){
      return false
      // return {
      //   'background-image':'linear-gradient(270deg,#fbbc0525,#fff,#fbbc0525)',
      //   'border': '1px solid #fbbc0550',
      //   'sortLevel': 2
      // }
    }
  }
}

</script>