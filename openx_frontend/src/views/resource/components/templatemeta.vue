<template>
  <div class="temp-display">
    <div class="temp-nav">
      <div
        v-for="nav in navList"
        :key="nav"
        :class="foucsEl === nav ? 'nav-item-selected' : 'nav-item'"
        @click="skipToNav(nav, foucsEl)"
      >
        {{ nav }}
      </div>
    </div>
    <div class="temp-info">
      <div class="spec-item" v-show="props.disablelabels">
        <div class="spec-value" id="id-labels">
          <div class="meta-label">labels</div>
          <div class="meta-value">
            <div class="tag-group">
              <div
                class="label-tag"
                v-for="(anno, key) in initLabels(props.poddata.metadata.labels)"
                v-bind:key="key"
              >
                {{ showLabel(anno) }}
                <el-icon @click="tagLabelClose(anno.label)"><Close /></el-icon>
              </div>
            </div>
            <el-button
              size="small"
              @click="addLabels('metadatalabels')"
              style="margin-top: 5px"
            >
              + add Matchlabel
            </el-button>
          </div>
        </div>
      </div>
      <div class="spec-item" v-show="props.disableannotation">
        <div class="spec-value" id="id-annotations">
          <div class="meta-label">annotations</div>
          <div class="meta-value">
            <div class="tag-group">
              <el-tooltip
                effect="dark"
                v-for="(anno, key) in initAnnotations(
                  props.poddata.metadata.annotations
                )"
                v-bind:key="key"
                class="box-item"
                :content="showLabel(anno)"
                placement="top-end"
              >
                <div class="label-tag">
                  {{ fetchLabel(anno) }}
                  <el-icon @click="tagAnnotationsClose(anno.label)"
                    ><Close
                  /></el-icon>
                </div>
              </el-tooltip>
              <!-- <div class="label-tag" v-for="(anno, key) in initLabels(props.poddata.metadata.annotations)"
            v-bind:key="key">
              {{showLabel(anno)}}
              <el-icon @click="tagAnnotationsClose(anno.label)"><Close /></el-icon>
            </div> -->
            </div>
            <el-button
              size="small"
              @click="addLabels('annotations')"
              style="margin-top: 5px"
            >
              + add annotation
            </el-button>
          </div>
        </div>
      </div>
      <div class="spec-item">
        <div class="spec-value" id="id-volumes">
          <div class="meta-label">
            volumes
            <el-button size="small" @click="addVolume">添加</el-button>
          </div>
          <div
            class="volume"
            v-for="(volume, volIndex) in props.poddata.spec.volumes"
            v-bind:key="volume.name"
          >
            <div class="volume-title">
              <div>Name: {{ volume.name }}</div>
              <div v-if="!volume.volumeSource.secret">
                <el-button size="small" @click="editVolume(volume, volIndex)"
                  >编辑</el-button
                >
                <el-popconfirm
                  title="确定删除?"
                  @confirm="deleteVolume(volIndex)"
                >
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
        </div>
        <div class="spec-value" id="id-Container">
          <div class="meta-label">
            Container
            <el-button
              size="small"
              v-show="containerLimit()"
              @click="addContainer"
              >添加</el-button
            >
            <el-button
              size="small"
              v-show="containerOne()"
              @click="delOneContainer"
              type="danger"
              plain
              >删除</el-button
            >
            <el-button
              size="small"
              v-show="containerCancel()"
              @click="cancelDelOne"
              type="primary"
              plain
              >撤销删除</el-button
            >
          </div>
          <div class="contain-list" v-show="!isApp">
            <div
              v-for="(contain, key) in props.poddata.spec.containers"
              v-bind:key="key"
              :class="
                conIndex === key ? 'contain-name selected-name' : 'contain-name'
              "
              @click="selectContainer(key)"
            >
              {{ contain.name }}
              <el-icon
                class="container-delete"
                @click.stop="deleteContainer(key)"
                ><Delete
              /></el-icon>
            </div>
          </div>
          <div v-if="props.poddata.spec.containers[conIndex]">
            <Container
              :containerdata="props.poddata.spec.containers[conIndex]"
              :key="conIndex"
              :appname="appname"
            />
          </div>
        </div>
        <div class="spec-value" id="id-NodeSelector">
          <div class="meta-label">NodeSelector</div>
          <div class="meta-value">
            <div class="tag-group">
              <el-tooltip
                effect="dark"
                v-for="(anno, key) in initLabels(
                  props.poddata.spec.nodeSelector
                )"
                v-bind:key="key"
                :content="showLabel(anno)"
                placement="top-end"
              >
                <div class="label-tag">
                  {{ showLabel(anno) }}
                  <el-icon @click="nodeSelectorClose(anno.label)"
                    ><Close
                  /></el-icon>
                </div>
              </el-tooltip>
            </div>
            <el-button
              size="small"
              @click="addLabels('nodeSelector')"
              style="margin-top: 5px"
            >
              + add Matchlabel
            </el-button>
          </div>
        </div>
        <div class="spec-value" id="id-serviceAccountName">
          <div class="meta-label">serviceAccountName</div>
          <div class="meta-value">
            <el-select
              v-model="props.poddata.spec.serviceAccountName"
              style="width: 100%"
            >
              <el-option
                v-for="ServiceAccount in ServiceAccountList"
                :key="ServiceAccount.metadata.name"
                :label="ServiceAccount.metadata.name"
                :value="ServiceAccount.metadata.name"
              />
            </el-select>
          </div>
        </div>
        <div class="spec-value" id="id-securityContext">
          <div
            class="meta-label"
            style="display: flex; gap: 40px; align-items: center"
          >
            securityContext :
            <el-switch v-model="openSecurity" @change="openSecurityChange" />
          </div>
          <div class="meta-value" v-if="props.poddata.spec.securityContext">
            <div class="security-item">
              supplementalGroups:
              <el-checkbox
                v-model="switchsupplementalGroups"
                @change="
                  (e) => {
                    securityChange(e, 'supplementalGroups');
                  }
                "
              />
              <div v-if="switchsupplementalGroups" class="tag-group">
                <div
                  class="label-tag"
                  v-for="(mode, index) in props.poddata.spec.securityContext
                    .supplementalGroups"
                  :key="mode"
                >
                  {{ mode }}
                  <el-popconfirm
                    title="确定删除?"
                    @confirm="
                      tagClose(
                        props.poddata.spec.securityContext.supplementalGroups,
                        index
                      )
                    "
                  >
                    <template #reference>
                      <el-icon><Close /></el-icon>
                    </template>
                  </el-popconfirm>
                </div>
                <TagInput @value-input="handleInputSupplementalGroups" />
              </div>
            </div>
            <div class="security-item">
              sysctls:
              <el-checkbox
                v-model="switchsysctls"
                @change="
                  (e) => {
                    securityChange(e, 'sysctls');
                  }
                "
              />
              <div v-if="switchsysctls" class="tag-group">
                <div
                  class="label-tag"
                  v-for="(mode, index) in props.poddata.spec.securityContext
                    .sysctls"
                  :key="mode"
                >
                  {{ mode }}
                  <el-popconfirm
                    title="确定删除?"
                    @confirm="
                      tagClose(
                        props.poddata.spec.securityContext.sysctls,
                        index
                      )
                    "
                  >
                    <template #reference>
                      <el-icon><Close /></el-icon>
                    </template>
                  </el-popconfirm>
                </div>
                <TagInput @value-input="handleInputSysctls" />
              </div>
            </div>
            <div class="security-item">
              seLinuxOptions:
              <el-checkbox
                v-model="switchseLinuxOptions"
                @change="
                  (e) => {
                    securityChange(e, 'seLinuxOptions');
                  }
                "
              />
              <div
                v-if="props.poddata.spec.securityContext.seLinuxOptions"
                style="display: flex; justify-content: flex-start; gap: 5px"
              >
                <el-input
                  v-model="
                    props.poddata.spec.securityContext.seLinuxOptions.level
                  "
                  placeholder="Please input"
                >
                  <template #prepend>level</template>
                </el-input>
                <el-input
                  v-model="
                    props.poddata.spec.securityContext.seLinuxOptions.role
                  "
                  placeholder="Please input"
                >
                  <template #prepend>role</template>
                </el-input>
                <el-input
                  v-model="
                    props.poddata.spec.securityContext.seLinuxOptions.type
                  "
                  placeholder="Please input"
                >
                  <template #prepend>type</template>
                </el-input>
                <el-input
                  v-model="
                    props.poddata.spec.securityContext.seLinuxOptions.user
                  "
                  placeholder="Please input"
                >
                  <template #prepend>user</template>
                </el-input>
              </div>
              <div
                v-else
                style="display: flex; justify-content: flex-start; gap: 5px"
              >
                <el-input disabled placeholder="level">
                  <template #prepend>level</template>
                </el-input>
                <el-input disabled placeholder="role">
                  <template #prepend>role</template>
                </el-input>
                <el-input disabled placeholder="type">
                  <template #prepend>type</template>
                </el-input>
                <el-input disabled placeholder="user">
                  <template #prepend>user</template>
                </el-input>
              </div>
            </div>
            <div class="security-item">
              seccompProfile:
              <el-checkbox
                v-model="switchseccompProfile"
                @change="
                  (e) => {
                    securityChange(e, 'seccompProfile');
                  }
                "
              />
              <div
                v-if="props.poddata.spec.securityContext.seccompProfile"
                style="display: flex; justify-content: flex-start"
              >
                <el-input
                  v-model="
                    props.poddata.spec.securityContext.seccompProfile
                      .localhostProfile
                  "
                  placeholder="Please input"
                >
                  <template #prepend>localhostProfile</template>
                </el-input>
                <div class="host-prepend" style="margin-left: 5px">type</div>
                <el-select
                  v-model="
                    props.poddata.spec.securityContext.seccompProfile.type
                  "
                  style="margin-right: 20px; margin-top: 1px"
                >
                  <el-option label="Unconfined" value="Unconfined" />
                  <el-option label="RuntimeDefault" value="RuntimeDefault" />
                  <el-option label="Localhost" value="Localhost" />
                </el-select>
              </div>
              <div v-else style="display: flex; justify-content: flex-start">
                <el-input disabled placeholder="localhostProfile">
                  <template #prepend>localhostProfile</template>
                </el-input>
                <div class="host-prepend" style="margin-left: 5px">type</div>
                <el-select disabled style="seccompProfile">
                  <el-option label="Unconfined" value="Unconfined" />
                  <el-option label="RuntimeDefault" value="RuntimeDefault" />
                  <el-option label="Localhost" value="Localhost" />
                </el-select>
              </div>
            </div>
            <!-- <div class="security-item">
            windowsOptions:
            <el-checkbox v-model="switchwindowsOptions" @change="e => {securityChange(e, 'windowsOptions')}"/>
            <div v-if="props.poddata.spec.securityContext.windowsOptions">
              <div style="display: flex; justify-content: flex-start; gap: 5px;margin-bottom: 10px">
                <el-input v-model="props.poddata.spec.securityContext.windowsOptions.gmsaCredentialSpec" placeholder="Please input">
                  <template #prepend>gmsaCredentialSpec</template>
                </el-input>
                <el-input v-model="props.poddata.spec.securityContext.windowsOptions.gmsaCredentialSpecName" placeholder="Please input">
                  <template #prepend>gmsaCredentialSpecName</template>
                </el-input>
              </div>
              <div style="display: flex; justify-content: flex-start;">
                <div class="host-prepend">hostProcess</div>
                <el-select v-model="props.poddata.spec.securityContext.windowsOptions.hostProcess" style="margin-right: 20px;margin-top: 1px">
                  <el-option label="true" :value="true" />
                  <el-option label="false" :value="false" />
                </el-select>
                <el-input v-model="props.poddata.spec.securityContext.windowsOptions.runAsUserName" placeholder="Please input">
                  <template #prepend>runAsUserName</template>
                </el-input>
              </div>
            </div>
            <div v-else>
              <div style="display: flex; justify-content: flex-start; gap: 5px;margin-bottom: 10px">
                <el-input placeholder="gmsaCredentialSpec" disabled>
                  <template #prepend>gmsaCredentialSpec</template>
                </el-input>
                <el-input placeholder="gmsaCredentialSpecName" disabled>
                  <template #prepend>gmsaCredentialSpecName</template>
                </el-input>
              </div>
              <div style="display: flex; justify-content: flex-start;">
                <div class="host-prepend">hostProcess</div>
                <el-select style="margin-right: 20px;margin-top: 1px" disabled>
                  <el-option label="true" :value="true" />
                  <el-option label="false" :value="false" />
                </el-select>
                <el-input placeholder="runAsUserName" disabled>
                  <template #prepend>runAsUserName</template>
                </el-input>
              </div>
            </div>

          </div> -->
            <div class="secur-grid">
              <div class="security-item">
                runAsGroup:
                <el-checkbox
                  v-model="switchRunAsGroup"
                  @change="
                    (e) => {
                      securityChange(e, 'runAsGroup');
                    }
                  "
                />
                <div v-if="switchRunAsGroup">
                  <el-input
                    v-model="props.poddata.spec.securityContext.runAsGroup"
                    placeholder="Please input"
                  ></el-input>
                </div>
                <div v-else>
                  <el-input placeholder="runAsGroup" disabled></el-input>
                </div>
              </div>
              <div class="security-item">
                runAsNonRoot:
                <el-checkbox
                  v-model="switchRunAsNonRoot"
                  @change="
                    (e) => {
                      securityChange(e, 'runAsNonRoot');
                    }
                  "
                />
                <div v-if="switchRunAsNonRoot">
                  <el-select
                    v-model="props.poddata.spec.securityContext.runAsNonRoot"
                    style="margin-right: 20px; margin-top: 1px"
                  >
                    <el-option label="true" :value="true" />
                    <el-option label="false" :value="false" />
                  </el-select>
                </div>
                <div v-else>
                  <el-input placeholder="runAsNonRoot" disabled></el-input>
                </div>
              </div>
              <div class="security-item">
                runAsUser:
                <el-checkbox
                  v-model="switchRunAsUser"
                  @change="
                    (e) => {
                      securityChange(e, 'runAsUser');
                    }
                  "
                />
                <div v-if="switchRunAsUser">
                  <el-input
                    v-model="props.poddata.spec.securityContext.runAsUser"
                    placeholder="Please input"
                  ></el-input>
                </div>
                <div v-else>
                  <el-input placeholder="runAsUser" disabled></el-input>
                </div>
              </div>
              <div class="security-item">
                fsGroup:
                <el-checkbox
                  v-model="switchFsGroup"
                  @change="
                    (e) => {
                      securityChange(e, 'fsGroup');
                    }
                  "
                />
                <div v-if="switchFsGroup">
                  <el-input
                    v-model="props.poddata.spec.securityContext.fsGroup"
                    placeholder="Please input"
                  ></el-input>
                </div>
                <div v-else>
                  <el-input placeholder="fsGroup" disabled></el-input>
                </div>
              </div>
              <div class="security-item">
                fsGroupChangePolicy:
                <el-checkbox
                  v-model="switchFsGroupChangePolicy"
                  @change="
                    (e) => {
                      securityChange(e, 'fsGroupChangePolicy');
                    }
                  "
                />
                <div v-if="switchFsGroupChangePolicy">
                  <el-select
                    v-model="
                      props.poddata.spec.securityContext.fsGroupChangePolicy
                    "
                  >
                    <el-option label="Always" value="Always" />
                    <el-option label="OnRootMismatch" value="OnRootMismatch" />
                  </el-select>
                </div>
                <div v-else>
                  <el-input
                    placeholder="fsGroupChangePolicy"
                    disabled
                  ></el-input>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="spec-value" id="id-imagePullSecrets">
          <div class="meta-label">imagePullSecrets</div>
          <div class="meta-value">
            <el-table
              size="small"
              :data="props.poddata.spec.imagePullSecrets"
              border
              style="width: 100%"
              highlight-current-row
            >
              <el-table-column key="name" label="Name" width="500">
                <template #default="scope">
                  <span v-if="scope.row.isSet">
                    <el-select
                      v-model="scope.row.name"
                      placeholder="Select"
                      style="width: 100%"
                    >
                      <el-option
                        v-for="secret in SecretList"
                        :key="secret.metadata.name"
                        :label="secret.metadata.name"
                        :value="secret.metadata.name"
                      />
                    </el-select>
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
                        props.poddata.spec.imagePullSecrets
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
                        props.poddata.spec.imagePullSecrets,
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
        <div class="spec-value" id="id-Affinity">
          <div class="meta-label">Affinity</div>
          <div class="meta-value">
            <el-select
              v-model="affinityName"
              clearable
              style="width: 100%"
              @change="affinityChange"
            >
              <el-option
                v-for="(Affinity, affinityIndex) in AffinityList"
                :key="affinityIndex"
                :label="Affinity.metadata.name"
                :value="Affinity.metadata.name"
              />
            </el-select>
            <Affinity
              v-if="props.poddata.spec.affinity"
              :affinity="props.poddata.spec.affinity"
            />
          </div>
        </div>
        <div class="spec-value" id="id-Tolerations">
          <div class="meta-label">
            Tolerations
            <span class="tol-tips" v-if="showAddTolTips">{{
              showAddTolTips
            }}</span>
          </div>
          <div class="meta-value">
            <el-table
              size="small"
              :data="props.poddata.spec.tolerations"
              border
              style="width: 100%"
              highlight-current-row
            >
              <el-table-column
                v-for="v in tolerationColumn"
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
                    class="el-tag el-tag--info el-tag--mini"
                    style="cursor: pointer"
                    @click="
                      tolerationEdit(
                        scope.row,
                        scope.row.isSet,
                        props.poddata.spec.tolerations
                      )
                    "
                  >
                    {{ scope.row.isSet ? "保存" : "修改" }}
                  </span>
                  <span
                    v-if="!scope.row.isSet"
                    class="el-tag el-tag--danger el-tag--mini"
                    style="cursor: pointer"
                    @click="
                      rowDelete(props.poddata.spec.tolerations, scope.$index)
                    "
                  >
                    删除
                  </span>
                  <span
                    v-else
                    class="el-tag el-tag--mini"
                    style="cursor: pointer"
                    @click="cancelTolerationEdit(scope.row)"
                  >
                    取消
                  </span>
                </template>
              </el-table-column>
            </el-table>
            <div class="tol-add">
              <el-select
                v-model="tolerationName"
                clearable
                style="width: calc(100% - 130px)"
                @change="tolerationChange"
              >
                <el-option
                  v-for="(Toleration, tolerationIndex) in TolerationList"
                  :key="tolerationIndex"
                  :label="Toleration.metadata.name"
                  :value="Toleration.metadata.name"
                />
              </el-select>
              <div
                class="el-table-add-row"
                style="width: 120px"
                @click="tolerationAdd()"
              >
                <span>+ 添加</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <el-dialog v-model="showAnnotations" width="30%" append-to-body>
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
    <el-button style="margin-top: 10px" size="small" @click="confirmAddTag">
      确定
    </el-button>
  </el-dialog>

  <el-dialog v-model="showVolume" width="60%" append-to-body>
    <div>
      Name:
      <el-input
        v-model="vlomeEditName"
        @input="
          (e) => {
            vlomeEditName = e.trim();
          }
        "
        size="small"
      />
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
                  size="mini"
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
                  configmapItemEdit(scope.row, scope.row.isSet, configMap.items)
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
            v-for="types in hostPathType"
            :key="types"
            :label="types"
            :value="types"
          />
        </el-select>
      </div>
    </div>
    <!-- <div v-if="vlomeType === 'secret'">
      <div>secret:
        <el-select v-model="secret" style="width: 100%" placeholder="Select" size="small">
          <el-option v-for="secret in SecretList"
            :key="secret.metadata.name"
            :label="secret.metadata.name"
            :value="secret.metadata.name"
          />
        </el-select>
      </div>
    </div> -->
    <el-button style="margin-top: 10px" size="small" @click="confirmAddVlome">
      确定
    </el-button>
  </el-dialog>
</template>

<style lang="scss" scoped>
.temp-display {
  display: flex;
  position: relative;

  .temp-nav {
    width: 180px;
    position: fixed;
    // top: 250px;
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
    overflow: auto;
    // max-height: calc(100vh - 279px);
    padding-left: 180px;
    //width: calc(100% - 180px);
  }
}
.meta-title {
  font-size: 16px;
  font-weight: bold;
  color: #303133;
  text-align: left;
}
.spec-item {
  .spec-label {
    width: 100%;
    text-align: left;
  }
  .spec-value {
    width: 100%;
    .meta-label {
      text-align: left;
      padding: 10px 0px;
    }
    .meta-value {
      padding-left: 100px;
      text-align: left;
      line-height: 2rem;
      font-weight: 400;
      display: flex;
      flex-direction: column;
      align-items: flex-start;
    }
  }
}
.volume {
  border: 1px solid #ccc;
  padding: 10px;
  margin: 5px;
  border-radius: 3px;
  font-size: 0.8rem;
  font-weight: 400;
}
.el-table-add-row {
  width: 100%;
  border: 1px dashed #c1c1cd;
  border-radius: 3px;
  cursor: pointer;
  justify-content: center;
  display: flex;
  line-height: 34px;
}
.tag-group {
  display: flex;
  justify-content: flex-start;
  flex-wrap: wrap;
  .label-tag {
    margin-top: 5px;
    margin-bottom: 3px;
    background-color: #f1f3f4;
    border: 1px solid #ccc;
    margin-right: 10px;
    border-radius: 4px;
    padding: 3px 3px 3px 10px;
    display: flex;
    justify-content: space-between;
    font-weight: 500;
    i {
      margin-left: 10px;
      cursor: pointer;
      &:hover {
        color: red;
      }
    }
  }
}
.contain-list {
  display: flex;
  justify-content: flex-start;
}
.contain-name {
  line-height: 1.5715;
  margin-right: 10px;
  display: inline-block;
  font-weight: 400;
  white-space: nowrap;
  text-align: center;
  border: 1px solid transparent;
  box-shadow: 0 2px #00000004;
  cursor: pointer;
  padding: 6px 8px;
  font-size: 14px;
  border-radius: 2px;
  color: #000000d9;
  border-color: #d9d9d9;
  background: #fff;
}
.container-delete {
  margin-left: 10px;
  color: white;
  &:hover {
    color: red;
  }
}
.selected-name {
  color: #fff;
  border-color: #1890ff;
  background: #1890ff;
  text-shadow: 0 -1px 0 rgb(0 0 0 / 12%);
  box-shadow: 0 2px #0000000b;
}
.volume-title {
  display: flex;
  justify-content: space-between;
}
.tol-add {
  display: flex;
  width: 100%;
}
.tol-tips {
  font-size: 0.5rem;
  margin-left: 50px;
  color: red;
}
.host-prepend {
  background-color: #f5f7fa;
  height: 30px;
  line-height: 30px;
  font-size: 14px;
  color: #909399;
  border: 1px solid #dcdfe6;
  padding: 0px 20px;
  position: relative;
  right: -1px;
  top: 1px;
  border-top-left-radius: 4px;
  border-bottom-left-radius: 4px;
}
.capability-style {
  width: 100%;
  border: 1px solid #dcdfe6;
}
.secur-grid {
  display: grid;
  grid-template-columns: 50% 50%;
  gap: 10px;
}
</style>

<script setup lang="ts">
import { cloneDeep } from "lodash";
import {
  initLabels,
  returnResource,
  showLabel,
  fetchLabel,
  initAnnotations,
} from "./../util";
import {
  rowEdit,
  rowDelete,
  rowCancelEdit,
  initLimits,
  goToNav,
} from "./tabelUtil";
import Container from "./containers.vue";
import Affinity from "./affinity.vue";
import TagInput from "./taginput.vue";
import proto from "../../../../proto/proto";
import { inject, nextTick, onMounted, ref, watch } from "vue";

const isApp = inject("isApp");
// <{ poddata?: any, disablelabels: Boolean }>
const props: any = defineProps({
  poddata: Object,
  appname: { type: String, default: "" },
  disablelabels: { type: Boolean, default: true },
  disableannotation: { type: Boolean, default: true },
});

const protoApi: any = proto.k8s.io.api.core.v1;

const hostPathType = [
  "DirectoryOrCreate",
  "Directory",
  "FileOrCreate",
  "File",
  "Socket",
  "CharDevice",
  "BlockDevice",
];

let switchRunAsGroup = ref(false);
let switchRunAsNonRoot = ref(false);
let switchRunAsUser = ref(false);
let switchFsGroup = ref(false);
let switchFsGroupChangePolicy = ref(false);

let switchsupplementalGroups = ref(false);
let switchsysctls = ref(false);
let switchseLinuxOptions = ref(false);
let switchwindowsOptions = ref(false);
let switchseccompProfile = ref(false);

if (props.poddata.spec.securityContext) {
  let newSecurInit = protoApi.PodSecurityContext.create();
  for (let securIndex in newSecurInit) {
    if (!props.poddata.spec.securityContext.hasOwnProperty(securIndex)) {
      // props.poddata.spec.securityContext[securIndex] = newSecurInit[securIndex]
    } else {
      if (securIndex === "runAsGroup") switchRunAsGroup.value = true;
      if (securIndex === "runAsNonRoot") switchRunAsNonRoot.value = true;
      if (securIndex === "runAsUser") switchRunAsUser.value = true;
      if (securIndex === "fsGroup") switchFsGroup.value = true;
      if (securIndex === "fsGroupChangePolicy")
        switchFsGroupChangePolicy.value = true;

      if (securIndex === "supplementalGroups")
        switchsupplementalGroups.value = true;
      if (securIndex === "sysctls") switchsysctls.value = true;

      if (securIndex === "seLinuxOptions") switchseLinuxOptions.value = true;
      // if(securIndex === 'windowsOptions') switchwindowsOptions.value = true
      if (securIndex === "seccompProfile") switchseccompProfile.value = true;
    }
    if (props.poddata.spec.securityContext[securIndex] === null) {
      if (securIndex === "seLinuxOptions") {
        switchseLinuxOptions.value = false;
      }
      // else if(securIndex === 'windowsOptions'){
      //   switchwindowsOptions.value = false
      // }
      else if (securIndex === "seccompProfile") {
        switchseccompProfile.value = false;
      }
    }
  }
}

let openSecurity = ref(false);
let saveSecurityContext = ref({});
onMounted(() => {
  if (props.poddata.spec.securityContext) {
    openSecurity.value = true;
  } else {
    openSecurity.value = false;
  }
  initPod();
});

let affinityName = ref("");
watch(
  () => props.appname,
  () => {
    openSecurity.value = Boolean(props.poddata.spec.securityContext);
    affinityName.value = "";
  },
  { immediate: true }
);
function openSecurityChange(hasSecur: boolean) {
  if (hasSecur) {
    props.poddata.spec.securityContext = saveSecurityContext.value;
    for (let securIndex in props.poddata.spec.securityContext) {
      if (props.poddata.spec.securityContext.hasOwnProperty(securIndex)) {
        switch (securIndex) {
          case "runAsGroup":
            switchRunAsGroup.value = true;
            break;
          case "runAsNonRoot":
            switchRunAsNonRoot.value = true;
            break;
          case "runAsUser":
            switchRunAsUser.value = true;
            break;
          case "fsGroup":
            switchFsGroup.value = true;
            break;
          case "fsGroupChangePolicy":
            switchFsGroupChangePolicy.value = true;
            break;
          case "supplementalGroups":
            switchsupplementalGroups.value = true;
            break;
          case "sysctls":
            switchsysctls.value = true;
            break;
          case "seLinuxOptions":
            switchseLinuxOptions.value = true;
            break;
          case "seccompProfile":
            switchseccompProfile.value = true;
            break;
          // case 'windowsOptions': switchwindowsOptions.value = true
          // break
        }
      } else {
        switch (securIndex) {
          case "runAsGroup":
            switchRunAsGroup.value = false;
            break;
          case "runAsNonRoot":
            switchRunAsNonRoot.value = false;
            break;
          case "runAsUser":
            switchRunAsUser.value = false;
            break;
          case "fsGroup":
            switchFsGroup.value = false;
            break;
          case "fsGroupChangePolicy":
            switchFsGroupChangePolicy.value = false;
            break;
          case "supplementalGroups":
            switchsupplementalGroups.value = false;
            break;
          case "sysctls":
            switchsysctls.value = false;
            break;
          case "seLinuxOptions":
            switchseLinuxOptions.value = false;
            break;
          case "seccompProfile":
            switchseccompProfile.value = false;
            break;
          // case 'windowsOptions': switchwindowsOptions.value = false
          // break
        }
      }
    }
    nextTick(() => {
      skipToNav("securityContext");
    });
  } else {
    delete props.poddata.spec.securityContext;
  }
}
function securityChange(haschild: boolean, child: string) {
  if (haschild) {
    switch (child) {
      case "runAsNonRoot":
        props.poddata.spec.securityContext[child] = false;
        break;
      case "fsGroupChangePolicy":
        props.poddata.spec.securityContext[child] = "Always";
        break;
      case "supplementalGroups":
        props.poddata.spec.securityContext[child] = [];
        break;
      case "sysctls":
        props.poddata.spec.securityContext[child] = [];
        break;
      case "seLinuxOptions":
        props.poddata.spec.securityContext[child] = {
          level: "",
          role: "",
          type: "",
          user: "",
        };
        break;
      case "seccompProfile":
        props.poddata.spec.securityContext[child] = {
          localhostProfile: "",
          type: "Localhost",
        };
        break;
      // case 'windowsOptions':
      // props.poddata.spec.securityContext[child] = {
      //   gmsaCredentialSpec: '', gmsaCredentialSpecName: '',
      //   hostProcess: false, runAsUserName: ''
      // }
      // break
      default:
        props.poddata.spec.securityContext[child] = "";
        break;
    }
  } else {
    delete props.poddata.spec.securityContext[child];
  }
}
function showProp(parents: any, child: string) {
  return parents.hasOwnProperty(child);
}

// initPod()
function initPod() {
  let newSecur = protoApi.PodSecurityContext.create();
  (newSecur.runAsGroup = ""),
    (newSecur.runAsNonRoot = false),
    (newSecur.runAsUser = ""),
    (newSecur.fsGroup = ""),
    (newSecur.fsGroupChangePolicy = "");
  for (let securIndex in newSecur) {
    if (newSecur[securIndex] === null) {
      if (securIndex === "seLinuxOptions") {
        const linux = protoApi.SELinuxOptions.create();
        (linux.level = ""),
          (linux.role = ""),
          (linux.type = ""),
          (linux.user = "");
        newSecur[securIndex] = linux;
      }
      // if(securIndex === 'windowsOptions'){
      //   let win = protoApi.WindowsSecurityContextOptions.create()
      //   win.gmsaCredentialSpec = '',  win.gmsaCredentialSpecName = '',  win.hostProcess = false,  win.runAsUserName = ''
      //   newSecur[securIndex] = win
      // }
      if (securIndex === "seccompProfile") {
        newSecur[securIndex] = protoApi.SeccompProfile.create();
        newSecur[securIndex].type = "Localhost";
        newSecur[securIndex].localhostProfile = "";
      }
    }
  }
  saveSecurityContext.value = newSecur;
}

let showAnnotations = ref(false);
let addTag = ref({
  key: "",
  value: "",
});
let addTitle = ref("");
function addLabels(title: string) {
  addTitle.value = title;
  showAnnotations.value = true;
}
function tagLabelClose(tagKey: string) {
  delete props.poddata.metadata.labels[tagKey];
}
function tagAnnotationsClose(tagKey: string) {
  delete props.poddata.metadata.annotations[tagKey];
}
function nodeSelectorClose(tagKey: string) {
  delete props.poddata.spec.nodeSelector[tagKey];
}
function confirmAddTag() {
  const addKey = addTag.value.key;
  const addValue = addTag.value.value;
  if (!addKey) {
    return;
  }
  if (addTitle.value === "metadatalabels") {
    props.poddata.metadata.labels[addKey] = addValue;
  } else if (addTitle.value === "nodeSelector") {
    props.poddata.spec.nodeSelector[addKey] = addValue;
  } else if (addTitle.value === "annotations") {
    props.poddata.metadata.annotations[addKey] = addValue;
  }
  showAnnotations.value = false;
}

function affinityChange(affinity: any) {
  const foundAffnity: any = AffinityList.value.findIndex((val: any) => {
    return val.metadata.name === affinity;
  });
  if (foundAffnity >= 0) {
    props.poddata.spec.affinity =
      AffinityList.value[foundAffnity].spec.affinity;
  } else {
    props.poddata.spec.affinity = {};
  }
}

let TolerationList: any = ref([]);
let tolerationName = ref("");
let selectedToleration = ref({});
function tolerationChange(name: any) {
  const foundToleration: any = TolerationList.value.findIndex((val: any) => {
    return val.metadata.name === name;
  });
  if (foundToleration >= 0) {
    selectedToleration.value =
      TolerationList.value[foundToleration].spec.toleration;
  }
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

function vlomeConfigmapAdd() {
  const oneConfigMapItem: cmKeyPath = {
    key: "key",
    path: "path",
  };
  configMap.value.items.push(oneConfigMapItem);
}
let editVlomeIndex = ref(0);
function confirmAddVlome() {
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
      props.poddata.spec.volumes[editVlomeIndex.value] = vloItem;
    } else {
      props.poddata.spec.volumes.push(vloItem);
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
      props.poddata.spec.volumes[editVlomeIndex.value] = vloItem;
    } else {
      props.poddata.spec.volumes.push(vloItem);
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
      props.poddata.spec.volumes[editVlomeIndex.value] = vloItem;
    } else {
      props.poddata.spec.volumes.push(vloItem);
    }
  }
  showVolume.value = false;
}

let conIndex = ref(0);
function selectContainer(selectIndex: number) {
  conIndex.value = selectIndex;
}
function deleteContainer(selectIndex: number) {
  props.poddata.spec.containers.splice(selectIndex, 1);
}
function containerLimit() {
  const hasOneContainer =
    props.poddata.spec.containers && props.poddata.spec.containers.length === 1;
  if (isApp && hasOneContainer) {
    return false;
  }
  return true;
}
function containerOne() {
  const hasOneContainer =
    props.poddata.spec.containers && props.poddata.spec.containers.length === 1;
  if (isApp && hasOneContainer) {
    return true;
  }
  return false;
}
let saveContainers = <any>[];
function delOneContainer() {
  saveContainers = cloneDeep(props.poddata.spec.containers);
  props.poddata.spec.containers = [];
}
function containerCancel() {
  const hasSavedContainer = saveContainers.length >= 1;
  const hasDelOneContainer =
    props.poddata.spec.containers && props.poddata.spec.containers.length === 0;
  return hasSavedContainer && hasDelOneContainer;
}
function cancelDelOne() {
  props.poddata.spec.containers = cloneDeep(saveContainers);
  saveContainers = [];
}
function addContainer() {
  const container = {
    name: "new Container",
    image: "",
    command: [],
    args: [],
    ports: [],
    env: [],
    resources: {
      limits: {
        cpu: { srting: "0" },
        memory: { srting: "0" },
      },
      requests: {
        cpu: { srting: "0" },
        memory: { srting: "0" },
      },
    },
    volumeMounts: [],
    imagePullPolicy: "",
  };

  props.poddata.spec.containers.push(container);
}

function handleInputSupplementalGroups(inputStr: string) {
  props.poddata.spec.securityContext.supplementalGroups.push(inputStr);
}
function handleInputSysctls(inputStr: string) {
  props.poddata.spec.securityContext.sysctls.push(inputStr);
}
function tagClose(obj: any, index: number) {
  obj.splice(index, 1);
}
let saveImageData = ref({});
function imagePullSecretEdit(row: any, isSet: boolean, allData: any) {
  rowEdit(row, isSet, allData, saveImageData);
}
function cancelimagePullSecretEdit(row: any) {
  rowCancelEdit(row, saveImageData);
}
function imagePullSecretAdd() {
  const newImagePullSecret = {
    name: "LocalObjectReference",
  };
  props.poddata.spec.imagePullSecrets.push(newImagePullSecret);
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

const tolerationColumn = [
  { field: "effect", title: "effect", width: 180 },
  { field: "key", title: "key" },
  { field: "operator", title: "operator" },
  { field: "value", title: "value", width: 180 },
];
let saveTolerationData = ref({});
function tolerationEdit(row: any, isSet: boolean, allData: any) {
  rowEdit(row, isSet, allData, saveTolerationData);
}
function cancelTolerationEdit(row: any) {
  rowCancelEdit(row, saveTolerationData);
}

let showAddTolTips = ref("");
let tipsTimer = ref();
function tolerationAdd() {
  const fundIndex = props.poddata.spec.tolerations.findIndex((tol: any) => {
    const selectedStr = JSON.stringify(selectedToleration.value);
    const tolStr = JSON.stringify(tol);
    return selectedStr === tolStr;
  });
  if (JSON.stringify(selectedToleration.value) === "{}") {
    showAddTolTips.value = "不能添加空配置";
    tipsTimer.value = setTimeout(() => {
      showAddTolTips.value = "";
    }, 1500);
    return;
  }
  if (fundIndex >= 0) {
    showAddTolTips.value = "已存在的配置，请勿重复添加";
    tipsTimer.value = setTimeout(() => {
      showAddTolTips.value = "";
    }, 1500);
  } else {
    props.poddata.spec.tolerations.push(selectedToleration.value);
  }
}

import { useRoute } from "vue-router";
import { useStore } from "@/store";
import { initSocketData, sendSocketMessage } from "@/api/socket";
const route = useRoute();
const store = useStore();

let getList = function (gvk: string) {
  const nsGvk = route.path.split("/");
  const senddata = initSocketData("Request", nsGvk[1], gvk, "list");
  sendSocketMessage(senddata, store);
};
getList("core-v1-ConfigMap");
getList("core-v1-Secret");
getList("core-v1-ServiceAccount");
getList("openx.neverdown.io-v1-Affinity");
getList("openx.neverdown.io-v1-Toleration");

let vlomeEditType = ref("create");
let ConfigMapList = ref([]);
let SecretList = ref([]);
let ServiceAccountList = ref([]);
let AffinityList: any = ref([]);

function deleteVolume(volIndex: number) {
  props.poddata.spec.volumes.splice(volIndex, 1);
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

function initMessage(msg: any, type: string) {
  const nsGvk = route.path.split("/");
  const gvkArr = type.split("-");
  let gvkObj = {
    group: gvkArr[0],
    version: gvkArr[1],
    kind: gvkArr[2],
  };
  try {
    let resultList = returnResource(msg, nsGvk[1], gvkObj, loadOver);
    if (resultList) {
      resultList.sort((itemL: any, itemR: any) => {
        const itemLTime = itemL.metadata.creationTimestamp.seconds;
        const itemRTime = itemR.metadata.creationTimestamp.seconds;
        return itemRTime - itemLTime;
      });
      if (type === "core-v1-ConfigMap") {
        ConfigMapList.value = resultList;
      }
      if (type === "core-v1-Secret") {
        SecretList.value = resultList;
      }
      if (type === "core-v1-ServiceAccount") {
        ServiceAccountList.value = resultList;
      }
      if (type === "openx.neverdown.io-v1-Affinity") {
        AffinityList.value = resultList;
      }
      if (type === "openx.neverdown.io-v1-Toleration") {
        TolerationList.value = resultList;
      }
    }
  } catch (e) {
    console.log("error");
  }
}
watch(
  () => store.state.socket.socket.message,
  (msg) => {
    const watchObj = [
      "core-v1-ConfigMap",
      "core-v1-Secret",
      "core-v1-ServiceAccount",
      "openx.neverdown.io-v1-Affinity",
      "openx.neverdown.io-v1-Toleration",
    ];
    for (let obj of watchObj) {
      initMessage(msg, obj);
    }
  }
);
let loadOver = function () {};

for (let containerOne of props.poddata.spec.containers) {
  initLimits(containerOne);
}

let foucsEl = ref("volumes");
const navList = [
  "volumes",
  "Container",
  "NodeSelector",
  "serviceAccountName",
  "securityContext",
  "imagePullSecrets",
  "Affinity",
  "Tolerations",
];
function skipToNav(nav: string) {
  foucsEl.value = nav;
  goToNav(nav);
}
</script>
