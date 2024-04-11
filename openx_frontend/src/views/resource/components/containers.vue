<template>
  <!-- <div v-if="metadata" class="meta-item">
    <div class="time-input">
      log
      <div class="edit-content">
        <el-button @click="goTerm(0)" class="time-btn">all</el-button>
        <el-button @click="goTerm(100)" class="time-btn">100s</el-button>
        <el-button @click="goTerm(1000)" class="time-btn">1000s</el-button>
        <el-input v-model="logSeconds" style="width: 200px" placeholder="Please input" @keyup.enter="goTerm(logSeconds)" />
      </div>
    </div>
    <div @click="goBash" class="time-input" style="cursor: pointer">
      bash
    </div>
    <div v-if="props.containerdata.currentDownloading" class="time-input" style="color:#5865f2">
      downloading
      <i class="el-icon-loading oper-icon"></i>
    </div>
    <div v-else class="time-input">
      download current
      <div class="edit-content">
        <el-button @click="downloadLog(false,0)" class="time-btn">all</el-button>
        <el-button @click="downloadLog(false, 100)" class="time-btn">100s</el-button>
        <el-button @click="downloadLog(false, 1000)" class="time-btn">1000s</el-button>
        <el-input v-model="logSeconds" style="width: 200px" placeholder="Please input" @keyup.enter="downloadLog(false, logSeconds)" />
      </div>
    </div>
    <div v-if="props.containerdata.previousDownloading" class="time-input" style="color:#5865f2">
      downloading
      <i class="el-icon-loading oper-icon"></i>
    </div>
    <div v-else @click="downloadLog(true, 0)" class="time-input" style="cursor: pointer">
      download previous
    </div>
  </div> -->
  <div class="meta-item">
    <div class="meta-label">Name:</div>
    <div class="meta-value">
      <el-input
        v-model="props.containerdata.name"
        size="small"
        placeholder="Please input Name"
      />
    </div>
  </div>
  <div class="meta-item">
    <div class="meta-label">Image:</div>
    <div class="meta-value">
      <ImageHub
        :image="props.containerdata.image"
        @image-change="
          (image) => {
            props.containerdata.image = image;
          }
        "
      />
    </div>
  </div>
  <div class="meta-item">
    <div class="meta-label">Command:</div>
    <div class="meta-value">
      <!-- <span style="color: blue;">{{'*以(,)分割参数 示例: /server/publisher,/server/lllidan-gateway'}}</span>
      <el-input v-model="tempCommand" @change="commandChange" type="textarea"
        size="small" placeholder="Please input Command" /> -->
      <!-- <el-tag :id="command" :key="command" v-for="(command , cindex) in props.containerdata.command" size="medium" class="config-tag"
        :disable-transitions="false">
        <div class="tag-name">{{command}}</div>
        <i style="line-height: 40px; cursor: pointer" @click.stop="handleClose(props.containerdata.command, cindex)" class="el-icon-delete"></i>
      </el-tag> -->
      <div class="tag-group">
        <div
          class="label-tag"
          :key="command"
          v-for="(command, cindex) in props.containerdata.command"
        >
          {{ command }}
          <el-icon @click="handleClose(props.containerdata.command, cindex)"
            ><Close
          /></el-icon>
        </div>
      </div>
      <el-input
        class="input-new-tag"
        v-if="inputCommandVisible"
        v-model="tempCommand"
        ref="saveCommandInput"
        @keyup.enter="$event.target.blur()"
        @blur="commandChange"
      >
      </el-input>
      <el-button
        v-else
        class="button-new-tag"
        size="small"
        @click="showCommandInput"
        >+ New Command</el-button
      >
    </div>
  </div>
  <div class="meta-item">
    <div class="meta-label">Args:</div>
    <div class="meta-value">
      <!-- <span style="color: blue;">{{'*以(,)分割参数 示例: -alsologtostderr=true,-v=4,-configPath=/var/conf/conf.yaml'}}</span>
      <el-input v-model="tempArg" @change="tempArgChange" type="textarea"
        size="small" placeholder="Please input Arguments" /> -->
      <div class="tag-group">
        <div
          class="label-tag"
          :key="command"
          v-for="(command, cindex) in props.containerdata.args"
        >
          {{ command }}
          <el-icon @click="handleClose(props.containerdata.args, cindex)"
            ><Close
          /></el-icon>
        </div>
      </div>
      <el-input
        class="input-new-tag"
        v-if="inputArgVisible"
        v-model="tempArg"
        ref="saveArgInput"
        @keyup.enter="$event.target.blur()"
        @blur="tempArgChange"
      >
      </el-input>
      <el-button
        v-else
        class="button-new-tag"
        size="small"
        @click="showArgInput"
        >+ New Arg</el-button
      >
    </div>
  </div>
  <div class="meta-item">
    <div class="meta-label">Ports:</div>
    <div class="meta-value">
      <el-table
        size="small"
        :data="props.containerdata.ports"
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
                containerPortEdit(
                  scope.row,
                  scope.row.isSet,
                  props.containerdata.ports
                )
              "
            >
              {{ scope.row.isSet ? "保存" : "修改" }}
            </span>
            <span
              v-if="!scope.row.isSet"
              class="el-tag el-tag--danger el-tag--mini"
              @click="rowDelete(props.containerdata.ports, scope.$index)"
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
        @click="containerPortAdd()"
      >
        <span>+ 添加</span>
      </div>
    </div>
  </div>
  <div class="meta-item">
    <div class="meta-label">
      Env:
      <div class="td-btns">
        <el-button class="action-btn" @click="envImport">import</el-button>
        <el-button class="action-btn" @click="envExport">export</el-button>
      </div>
    </div>
    <div class="meta-value">
      <el-table
        size="small"
        :data="props.containerdata.env"
        border
        style="width: 100%"
        highlight-current-row
      >
        <el-table-column key="name" label="name" :width="200">
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
              @click="
                envEdit(scope.row, scope.row.isSet, props.containerdata.env)
              "
              >{{ scope.row.name }}</span
            >
          </template>
        </el-table-column>
        <el-table-column key="value" label="value" :width="200">
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
              @click="
                envEdit(scope.row, scope.row.isSet, props.containerdata.env)
              "
              >{{ scope.row.value }}</span
            >
          </template>
        </el-table-column>
        <el-table-column key="fieldPath" label="fieldPath" :width="260">
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
              @click="
                envEdit(scope.row, scope.row.isSet, props.containerdata.env)
              "
            >
              {{ scope.row.isSet ? "保存" : "修改" }}
            </span>
            <span
              v-if="!scope.row.isSet"
              class="el-tag el-tag--danger el-tag--mini"
              @click="rowDelete(props.containerdata.env, scope.$index)"
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
      <el-alert
        title="Env name重复，请修改"
        type="error"
        v-if="envAddError"
        show-icon
        :closable="false"
      />
      <div class="el-table-add-row" style="width: 99.2%" @click="envAdd()">
        <span>+ 添加</span>
      </div>
    </div>
  </div>
  <div class="meta-item">
    <div class="meta-label">Resources:</div>
    <div class="meta-value">
      <el-tabs type="border-card">
        <el-tab-pane key="limits" label="limits">
          <div style="justify-content: space-evenly; display: flex">
            <el-input
              v-model="props.containerdata.resources.limits.cpu.string"
              style="width: 40%"
            >
              <template #prepend>cpu</template>
            </el-input>
            <el-input
              v-model="props.containerdata.resources.limits.memory.string"
              style="width: 40%"
            >
              <template #prepend>memory</template>
            </el-input>
          </div>
        </el-tab-pane>
        <el-tab-pane key="requests" label="requests">
          <div style="justify-content: space-evenly; display: flex">
            <el-input
              v-model="props.containerdata.resources.requests.cpu.string"
              style="width: 40%"
            >
              <template #prepend>cpu</template>
            </el-input>
            <el-input
              v-model="props.containerdata.resources.requests.memory.string"
              style="width: 40%"
            >
              <template #prepend>memory</template>
            </el-input>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
  <div class="meta-item">
    <div class="meta-label">VolumeMounts:</div>
    <div class="meta-value">
      <el-table
        size="small"
        :data="props.containerdata.volumeMounts"
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
                type="text"
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
                volumeMountsEdit(
                  scope.row,
                  scope.row.isSet,
                  props.containerdata.volumeMounts
                )
              "
            >
              {{ scope.row.isSet ? "保存" : "修改" }}
            </span>
            <span
              v-if="!scope.row.isSet"
              class="el-tag el-tag--danger el-tag--mini"
              @click="rowDelete(props.containerdata.volumeMounts, scope.$index)"
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
  </div>
  <div class="meta-item">
    <div class="meta-label">ImagePullPolicy:</div>
    <div class="meta-value">
      <el-select
        v-model="props.containerdata.imagePullPolicy"
        placeholder="Select"
        style="width: 100%"
      >
        <el-option
          v-for="policy in imagePullPolicyOptions"
          :key="policy"
          :label="policy"
          :value="policy"
        />
      </el-select>
    </div>
  </div>
  <div class="meta-item">
    <div class="meta-label" style="display: flex; gap: 40px">
      SecurityContext:
      <el-switch v-model="openSecurity" @change="openSecurityChange" />
    </div>
    <div
      v-if="containerdata.securityContext"
      class="meta-value"
      style="padding-top: 25px"
    >
      <div class="security-item">
        capabilities:
        <el-checkbox
          v-model="switchcapabilities"
          @change="
            (e) => {
              securityChange(e, 'capabilities');
            }
          "
        />
        <div v-if="containerdata.securityContext.capabilities">
          <div class="host-prepend" style="width: 100px">add:</div>
          <div class="tag-group">
            <div
              class="label-tag"
              v-for="(mode, index) in containerdata.securityContext.capabilities
                .add"
              :key="mode"
            >
              {{ mode }}
              <el-popconfirm
                title="确定删除?"
                @confirm="
                  tagClose(
                    containerdata.securityContext.capabilities.add,
                    index
                  )
                "
              >
                <template #reference>
                  <el-icon><Close /></el-icon>
                </template>
              </el-popconfirm>
            </div>
            <TagInput @value-input="handleInputAdd" />
          </div>
          <div class="host-prepend" style="width: 100px">drop:</div>
          <div class="tag-group">
            <div
              class="label-tag"
              v-for="(mode, index) in containerdata.securityContext.capabilities
                .drop"
              :key="mode"
            >
              {{ mode }}
              <el-popconfirm
                title="确定删除?"
                @confirm="
                  tagClose(
                    containerdata.securityContext.capabilities.drop,
                    index
                  )
                "
              >
                <template #reference>
                  <el-icon><Close /></el-icon>
                </template>
              </el-popconfirm>
            </div>
            <TagInput @value-input="handleInputDrop" />
          </div>
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
          v-if="containerdata.securityContext.seLinuxOptions"
          style="display: flex; justify-content: flex-start; gap: 5px"
        >
          <el-input
            v-model="containerdata.securityContext.seLinuxOptions.level"
            placeholder="Please input"
          >
            <template #prepend>level</template>
          </el-input>
          <el-input
            v-model="containerdata.securityContext.seLinuxOptions.role"
            placeholder="Please input"
          >
            <template #prepend>role</template>
          </el-input>
          <el-input
            v-model="containerdata.securityContext.seLinuxOptions.type"
            placeholder="Please input"
          >
            <template #prepend>type</template>
          </el-input>
          <el-input
            v-model="containerdata.securityContext.seLinuxOptions.user"
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
          v-if="containerdata.securityContext.seccompProfile"
          style="display: flex; justify-content: flex-start"
        >
          <el-input
            v-model="
              containerdata.securityContext.seccompProfile.localhostProfile
            "
            placeholder="Please input"
          >
            <template #prepend>localhostProfile</template>
          </el-input>
          <div class="host-prepend" style="margin-left: 5px">type</div>
          <el-select
            v-model="containerdata.securityContext.seccompProfile.type"
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
        <div v-if="containerdata.securityContext.windowsOptions">
          <div style="display: flex; justify-content: flex-start; gap: 5px;margin-bottom: 10px">
            <el-input v-model="containerdata.securityContext.windowsOptions.gmsaCredentialSpec" placeholder="Please input">
              <template #prepend>gmsaCredentialSpec</template>
            </el-input>
            <el-input v-model="containerdata.securityContext.windowsOptions.gmsaCredentialSpecName" placeholder="Please input">
              <template #prepend>gmsaCredentialSpecName</template>
            </el-input>
          </div>
          <div style="display: flex; justify-content: flex-start;">
            <div class="host-prepend">hostProcess</div>
            <el-select v-model="containerdata.securityContext.windowsOptions.hostProcess" style="margin-right: 20px;margin-top: 1px">
              <el-option label="true" :value="true" />
              <el-option label="false" :value="false" />
            </el-select>
            <el-input v-model="containerdata.securityContext.windowsOptions.runAsUserName" placeholder="Please input">
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
          allowPrivilegeEscalation:
          <el-checkbox
            v-model="switchallowPrivilegeEscalation"
            @change="
              (e) => {
                securityChange(e, 'allowPrivilegeEscalation');
              }
            "
          />
          <div v-if="switchallowPrivilegeEscalation">
            <el-select
              v-model="containerdata.securityContext.allowPrivilegeEscalation"
              style="margin-right: 20px; margin-top: 1px"
            >
              <el-option label="true" :value="true" />
              <el-option label="false" :value="false" />
            </el-select>
          </div>
          <div v-else>
            <el-select disabled style="margin-right: 20px; margin-top: 1px">
              <el-option label="true" :value="true" />
              <el-option label="false" :value="false" />
            </el-select>
          </div>
        </div>
        <div class="security-item">
          privileged:
          <el-checkbox
            v-model="switchprivileged"
            @change="
              (e) => {
                securityChange(e, 'privileged');
              }
            "
          />
          <div v-if="switchprivileged">
            <el-select
              v-model="containerdata.securityContext.privileged"
              style="margin-right: 20px; margin-top: 1px"
            >
              <el-option label="true" :value="true" />
              <el-option label="false" :value="false" />
            </el-select>
          </div>
          <div v-else>
            <el-select disabled style="margin-right: 20px; margin-top: 1px">
              <el-option label="true" :value="true" />
              <el-option label="false" :value="false" />
            </el-select>
          </div>
        </div>
        <div class="security-item">
          readOnlyRootFilesystem:
          <el-checkbox
            v-model="switchreadOnlyRootFilesystem"
            @change="
              (e) => {
                securityChange(e, 'readOnlyRootFilesystem');
              }
            "
          />
          <div v-if="switchreadOnlyRootFilesystem">
            <el-select
              v-model="containerdata.securityContext.readOnlyRootFilesystem"
              style="margin-right: 20px; margin-top: 1px"
            >
              <el-option label="true" :value="true" />
              <el-option label="false" :value="false" />
            </el-select>
          </div>
          <div v-else>
            <el-select disabled style="margin-right: 20px; margin-top: 1px">
              <el-option label="true" :value="true" />
              <el-option label="false" :value="false" />
            </el-select>
          </div>
        </div>
        <div class="security-item">
          runAsGroup:
          <el-checkbox
            v-model="switchrunAsGroup"
            @change="
              (e) => {
                securityChange(e, 'runAsGroup');
              }
            "
          />
          <div v-if="switchrunAsGroup">
            <el-input
              v-model="containerdata.securityContext.runAsGroup"
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
            v-model="switchrunAsNonRoot"
            @change="
              (e) => {
                securityChange(e, 'runAsNonRoot');
              }
            "
          />
          <div v-if="switchrunAsNonRoot">
            <el-select
              v-model="containerdata.securityContext.runAsNonRoot"
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
            v-model="switchrunAsUser"
            @change="
              (e) => {
                securityChange(e, 'runAsUser');
              }
            "
          />
          <div v-if="switchrunAsUser">
            <el-input
              v-model="containerdata.securityContext.runAsUser"
              placeholder="Please input"
            ></el-input>
          </div>
          <div v-else>
            <el-input placeholder="runAsUser" disabled></el-input>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.meta-title {
  font-size: 16px;
  font-weight: bold;
  color: #303133;
  text-align: left;
}
.meta-item {
  display: flex;
  padding: 10px 0px;
  border-bottom: 1px solid #f1f3f4;
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
.tag-group {
  display: flex;
  justify-content: flex-start;
  flex-wrap: wrap;
  .label-tag {
    margin-top: 5px;
    margin-bottom: 3px;
    margin-right: 10px;
    background-color: #f1f3f4;
    border: 1px solid #ccc;
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
.time-input {
  height: 40px;
  line-height: 40px;
  width: 200px;
  border: 1px solid #ccc;
  border-radius: 5px;
  position: relative;
  margin-right: 10px;
  .edit-content {
    display: none;
    color: #7a8295;
    position: absolute;
    background-color: #f1f3f4;
    padding: 10px;
    top: 40px;
    z-index: 1;
    font-size: 20px;
    .time-btn {
      width: 100px;
      margin: 0px 10px 0px 0px;
    }
    .el-input__inner {
      width: 100px;
      height: 35px;
      line-height: 35px;
      position: relative;
      top: -1px;
      left: -10px;
    }
  }
  &:hover {
    .edit-content {
      display: flex;
      justify-content: flex-end;
    }
  }
}
.td-btns {
  display: flex;
  gap: 10px;
  justify-content: space-around;
  .action-btn {
    width: 70px;
    padding: 5px;
    border: 1px solid #ccc;
    border-radius: 5px;
    cursor: pointer;
    &:hover {
      border: 1px solid #202123;
    }
  }
}
</style>

<script setup lang="ts">
import { rowEdit, rowDelete, rowCancelEdit } from "./tabelUtil";
import { computed, inject, nextTick, onMounted, ref, toRaw, watch } from "vue";
import ImageHub from "./imagahub.vue";
import TagInput from "./taginput.vue";
import { openTerm, download } from "./../podutils";
import { ElNotification } from "element-plus";

const props = defineProps<{
  containerdata?: any;
  key?: { type: number; defalut: 0 };
  appname?: string;
}>();

let logSeconds = ref(0);
const metadata = inject("podMetadata");
let openSecurity = ref(false);

let switchcapabilities = ref(false);
let switchseLinuxOptions = ref(false);
let switchseccompProfile = ref(false);
let switchwindowsOptions = ref(false);
let switchallowPrivilegeEscalation = ref(false);
let switchprivileged = ref(false);
let switchreadOnlyRootFilesystem = ref(false);
let switchrunAsGroup = ref(false);
let switchrunAsNonRoot = ref(false);
let switchrunAsUser = ref(false);

watch(
  () => props.appname,
  () => {
    openSecurity.value = Boolean(props.containerdata.securityContext);
  },
  { immediate: true }
);
onMounted(() => {
  if (props.containerdata.securityContext) {
    openSecurity.value = true;
    let newSecur = protoApi.SecurityContext.create();
    for (let securIndex in newSecur) {
      if (!props.containerdata.securityContext.hasOwnProperty(securIndex)) {
        // props.containerdata.securityContext[securIndex] = newSecur[securIndex]
      } else {
        if (securIndex === "capabilities") switchcapabilities.value = true;
        if (securIndex === "seLinuxOptions") switchseLinuxOptions.value = true;
        if (securIndex === "seccompProfile") switchseccompProfile.value = true;
        // if(securIndex === 'windowsOptions') switchwindowsOptions.value = true
        if (securIndex === "allowPrivilegeEscalation")
          switchallowPrivilegeEscalation.value = true;

        if (securIndex === "privileged") switchprivileged.value = true;
        if (securIndex === "readOnlyRootFilesystem")
          switchreadOnlyRootFilesystem.value = true;

        if (securIndex === "runAsGroup") switchrunAsGroup.value = true;
        if (securIndex === "runAsNonRoot") switchrunAsNonRoot.value = true;
        if (securIndex === "runAsUser") switchrunAsUser.value = true;
      }
      // if(newSecur[securIndex] === null){
      //   if(securIndex === 'capabilities'){
      //     switchcapabilities.value = false
      //   }
      //   if(securIndex === 'seLinuxOptions'){
      //     switchseLinuxOptions.value = false
      //   }
      //   // if(securIndex === 'windowsOptions'){
      //   //   switchwindowsOptions.value = false
      //   // }
      //   if(securIndex === 'seccompProfile'){
      //     switchseccompProfile.value = false
      //   }
      // }
    }
  } else {
    openSecurity.value = false;
  }
  initCon();
});

let saveSecurityContext = ref({});

import proto from "../../../../proto/proto";
import { fa } from "element-plus/es/locale";
const protoApi: any = proto.k8s.io.api.core.v1;
function initCon() {
  let newSecur = protoApi.SecurityContext.create();
  (newSecur.runAsGroup = ""),
    (newSecur.runAsNonRoot = false),
    (newSecur.runAsUser = ""),
    (newSecur.allowPrivilegeEscalation = false),
    (newSecur.privileged = false),
    (newSecur.readOnlyRootFilesystem = false);
  for (let securIndex in newSecur) {
    if (newSecur[securIndex] === null) {
      if (securIndex === "capabilities") {
        newSecur[securIndex] = protoApi.Capabilities.create();
        newSecur[securIndex].add = [];
        newSecur[securIndex].drop = [];
      }
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
        newSecur[securIndex].type = "Unconfined";
        newSecur[securIndex].localhostProfile = "";
      }
    }
  }
  saveSecurityContext.value = newSecur;
}

function openSecurityChange(hasSecur: boolean) {
  if (hasSecur) {
    props.containerdata.securityContext = saveSecurityContext.value;
    switchcapabilities.value = true;
    switchseLinuxOptions.value = true;
    switchseccompProfile.value = true;
    switchwindowsOptions.value = true;
    switchallowPrivilegeEscalation.value = true;
    switchprivileged.value = true;
    switchreadOnlyRootFilesystem.value = true;
    switchrunAsGroup.value = true;
    switchrunAsNonRoot.value = true;
    switchrunAsUser.value = true;
  } else {
    delete props.containerdata.securityContext;
  }
}

function securityChange(haschild: boolean, child: string) {
  if (haschild) {
    switch (child) {
      case "capabilities":
        props.containerdata.securityContext[child] = {
          add: [],
          drop: [],
        };
        break;
      case "seLinuxOptions":
        props.containerdata.securityContext[child] = {
          level: "",
          role: "",
          type: "",
          user: "",
        };
        break;
      case "seccompProfile":
        props.containerdata.securityContext[child] = {
          localhostProfile: "",
          type: "Localhost",
        };
        break;
      // case 'windowsOptions':
      // props.containerdata.securityContext[child] = {
      //   gmsaCredentialSpec: '', gmsaCredentialSpecName: '',
      //   hostProcess: false, runAsUserName: ''
      // }
      // break
      case "allowPrivilegeEscalation":
        props.containerdata.securityContext[child] = false;
        break;
      case "privileged":
        props.containerdata.securityContext[child] = false;
        break;
      case "readOnlyRootFilesystem":
        props.containerdata.securityContext[child] = false;
        break;
      case "runAsNonRoot":
        props.containerdata.securityContext[child] = false;
        break;
      default:
        props.containerdata.securityContext[child] = "";
        break;
    }
  } else {
    delete props.containerdata.securityContext[child];
  }
}

async function goTerm(time: number) {
  await openTerm(metadata, props.containerdata, "log", time);
}
async function goBash() {
  await openTerm(metadata, props.containerdata, "bash", 0);
}
async function downloadLog(isprevious: boolean, time: number) {
  if (isprevious) {
    props.containerdata.previousDownloading = true;
  } else {
    props.containerdata.currentDownloading = true;
  }
  await download(metadata, props.containerdata, isprevious, time);
  if (isprevious) {
    props.containerdata.previousDownloading = false;
  } else {
    props.containerdata.currentDownloading = false;
  }
}

let tempCommand = ref("");
let inputCommandVisible = ref(false);
function commandChange() {
  inputCommandVisible.value = false;
  if (tempCommand.value) {
    props.containerdata.command.push(tempCommand.value);
  }
}
let saveCommandInput = ref();
function showCommandInput() {
  inputCommandVisible.value = true;
  nextTick(() => {
    saveCommandInput.value.$refs.input.focus();
  });
}

let tempArg = ref("");
let inputArgVisible = ref(false);
function tempArgChange() {
  inputArgVisible.value = false;
  if (tempArg.value) {
    props.containerdata.args.push(tempArg.value);
  }
}
let saveArgInput = ref();
function showArgInput() {
  inputArgVisible.value = true;
  nextTick(() => {
    saveArgInput.value.$refs.input.focus();
  });
}

function handleClose(proData: any, delIndex: number) {
  proData.splice(delIndex, 1);
}

const containCloumn = [
  { field: "name", title: "name", width: 150 },
  { field: "containerPort", title: "containerPort", width: 120 },
  { field: "hostIP", title: "hostIP", width: 150 },
  { field: "hostPort", title: "hostPort", width: 120 },
  { field: "protocol", title: "protocol", width: 120 },
];
// const envCloumn = [
//   { field: 'name', title: 'name', width: 300 },
//   { field: 'value', title: 'value', width: 200 }
// ]
const volumeMountsCloumn = [
  { field: "name", title: "name", width: 150 },
  { field: "mountPath", title: "mountPath", width: 150 },
  { field: "readOnly", title: "readOnly", width: 150 },
  { field: "subPath", title: "subPath", width: 150 },
  { field: "subPathExpr", title: "subPathExpr", width: 150 },
];
const imagePullPolicyOptions = ["Always", "Never", "IfNotPresent"];

function containerPortAdd() {
  const newPort = {
    containerPort: 80,
    hostIP: "",
    hostPort: 0,
    name: "",
    protocol: "TCP",
  };
  props.containerdata.ports.push(newPort);
}
let savePortsData = ref({});
function containerPortEdit(row: any, isSet: boolean, allData: any) {
  rowEdit(row, isSet, allData, savePortsData);
}
function cancelEdit(row: any) {
  rowCancelEdit(row, savePortsData);
}

function envAdd() {
  const newEnv = {
    name: "name",
    value: "value",
  };
  props.containerdata.env.push(newEnv);
}
let saveEnvData = ref({});
let envAddError = ref(false);
function envEdit(row: any, isSet: boolean, allData: any) {
  if (isSet) {
    let allEnv = toRaw(allData),
      thisROw = toRaw(row);
    const envFilter = allEnv.filter((envar: any) => {
      return envar.name === thisROw.name;
    });
    if (envFilter.length >= 2) {
      envAddError.value = true;
      return;
    }
    envAddError.value = false;
  }
  rowEdit(row, isSet, allData, saveEnvData);
}
function cancelEnvEdit(row: any) {
  envAddError.value = false;
  rowCancelEdit(row, saveEnvData);
}

function volumeMountsAdd() {
  const newVolumeMount = {
    mountPath: "/data",
    name: "name",
    readOnly: false,
    subPath: "",
    subPathExpr: "",
  };
  props.containerdata.volumeMounts.push(newVolumeMount);
}
let saveVolumeMountData = ref({});
function volumeMountsEdit(row: any, isSet: boolean, allData: any) {
  rowEdit(row, isSet, allData, saveVolumeMountData);
}
function cancelVolumeMountsEdit(row: any) {
  rowCancelEdit(row, saveVolumeMountData);
}

function envImport() {
  if (localStorage.getItem("envStorage")) {
    const envStorage = JSON.parse(String(localStorage.getItem("envStorage")));
    for (let env of envStorage) {
      const propEnvFind = props.containerdata.env.findIndex((propEnv: any) => {
        return propEnv.name === env.name;
      });
      if (propEnvFind < 0) {
        props.containerdata.env.push(env);
      } else {
        props.containerdata.env[propEnvFind] = env;
      }
    }
  }
}
function envExport() {
  localStorage.setItem("envStorage", JSON.stringify(props.containerdata.env));
  ElNotification({
    title: "导出成功",
    message: "",
    type: "success",
    duration: 1000,
  });
}

function handleInputAdd(inputStr: string) {
  props.containerdata.securityContext.capabilities.add.push(inputStr);
}
function handleInputDrop(inputStr: string) {
  props.containerdata.securityContext.capabilities.drop.push(inputStr);
}
function tagClose(obj: any, index: number) {
  obj.splice(index, 1);
}
</script>
