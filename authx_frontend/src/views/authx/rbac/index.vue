认证权限
<template>
  <div class="app-container">
    <el-card class="box-card">
      <el-tag>应用名称: {{ appName }}</el-tag>
      <el-button
        size="small"
        icon="el-icon-back"
        style="margin-left: 10px"
        @click="goback"
        >返回</el-button
      >
      <el-tabs v-model="active_tab_name" @tab-click="tab_click">
        <!-- 账户管理--------------------------------------------------  -->
        <el-tab-pane label="账户" name="account_setting">
          <el-button
            type="primary"
            size="small"
            icon="el-icon-circle-plus-outline"
            @click="create_account"
            >新增</el-button
          >
          <el-table
            :data="account_list"
            size="small"
            empty-text="啥也没有"
            border
          >
            <el-table-column prop="username" label="用户"></el-table-column>
            <!-- <el-table-column prop="nickname" label="别名"></el-table-column> -->
            <el-table-column prop="desc" label="描述"></el-table-column>
            <el-table-column prop="role" label="角色"></el-table-column>
            <el-table-column label="密码">
              <template slot-scope="scoped">
                <el-button type="text" @click="copy(scoped.row.password)"
                  >复制</el-button
                >
              </template>
            </el-table-column>
            <el-table-column label="操作" width="180px;">
              <template slot-scope="scoped">
                <el-tooltip
                  class="item"
                  effect="dark"
                  content="重置密码"
                  placement="top"
                >
                  <el-button
                    type="warning"
                    size="mini"
                    icon="el-icon-unlock"
                    @click="reset_password(scoped.row)"
                  ></el-button>
                </el-tooltip>
                <el-button
                  type="primary"
                  icon="el-icon-edit"
                  size="mini"
                  @click="update_account(scoped.row)"
                ></el-button>
                <el-popconfirm
                  title="确定删除吗？"
                  confirm-button-text="确定"
                  cancel-button-text="不了"
                  style="margin-left: 10px"
                  @confirm="delete_account(scoped.row)"
                  @cancel="cancel_delete"
                >
                  <el-button
                    slot="reference"
                    type="danger"
                    icon="el-icon-delete"
                    size="mini"
                  ></el-button>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
          <el-pagination
            background
            :page-size="10"
            :current-page.sync="currentPage"
            :total="account_list.length"
            layout="total, prev, pager, next"
            style="text-align: left; margin-top: 20px"
          >
          </el-pagination>
        </el-tab-pane>
        <!-- 角色管理--------------------------------------------------  -->
        <el-tab-pane label="角色" name="role_setting">
          <el-button
            type="primary"
            size="small"
            icon="el-icon-circle-plus-outline"
            @click="create_role"
            >新增</el-button
          >
          <el-table :data="role_list" size="small" empty-text="啥也没有" border>
            <el-table-column prop="name" label="名称"></el-table-column>
            <el-table-column prop="desc" label="描述"></el-table-column>
            <el-table-column label="操作" width="120px;">
              <template slot-scope="scoped">
                <el-button
                  type="primary"
                  icon="el-icon-edit"
                  size="mini"
                  @click="update_role(scoped.row)"
                ></el-button>
                <el-popconfirm
                  title="确定删除吗？"
                  confirm-button-text="确定"
                  cancel-button-text="不了"
                  style="margin-left: 10px"
                  @confirm="delete_role(scoped.row)"
                  @cancel="cancel_delete"
                >
                  <el-button
                    slot="reference"
                    type="danger"
                    icon="el-icon-delete"
                    size="mini"
                  ></el-button>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
          <el-pagination
            background
            :page-size="10"
            :current-page.sync="currentPage"
            :total="role_list.length"
            layout="total, prev, pager, next"
            style="text-align: left; margin-top: 20px"
          >
          </el-pagination>
        </el-tab-pane>
        <!-- GVK管理--------------------------------------------------  -->
        <el-tab-pane label="GVK" name="gvk_setting">
          <el-button
            type="primary"
            size="small"
            icon="el-icon-circle-plus-outline"
            @click="create_gvk"
            disabled
            >新增</el-button
          >
          <el-table
            :data="page_gvk_list"
            size="small"
            empty-text="啥也没有"
            border
          >
            <el-table-column prop="name" label="GVK"></el-table-column>
            <el-table-column prop="verbs" label="verbs">
              <template slot-scope="scoped">
                <el-tag
                  v-for="(item, index) in scoped.row.verbs"
                  :key="index"
                  :value="item"
                  style="margin-right: 10px"
                  >{{ item }}</el-tag
                >
              </template>
            </el-table-column>
            <el-table-column label="操作" width="120px;">
              <template slot-scope="scoped">
                <el-button
                  type="primary"
                  icon="el-icon-edit"
                  size="mini"
                  @click="update_gvk(scoped.row)"
                  disabled
                ></el-button>
                <el-popconfirm
                  title="确定删除吗？"
                  confirm-button-text="确定"
                  cancel-button-text="不了"
                  style="margin-left: 10px"
                  @confirm="delete_gvk(scoped.row)"
                  @cancel="cancel_delete"
                >
                  <el-button
                    slot="reference"
                    type="danger"
                    icon="el-icon-delete"
                    size="mini"
                    disabled
                  ></el-button>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
          <el-pagination
            background
            :page-size="10"
            :current-page.sync="currentPage"
            :total="gvk_list.length"
            layout="total, prev, pager, next"
            style="text-align: left; margin-top: 20px"
          >
          </el-pagination>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- 账户管理--------------------------------------------------  -->
    <el-dialog
      :title="textMap[dialogStatus]"
      :visible.sync="account_dialog"
      width="50%"
    >
      <el-form
        ref="account_form_refs"
        :model="account_form"
        :rules="account_form_rules"
        size="small"
        label-width="100px"
      >
        <el-row>
          <el-col :span="12">
            <el-form-item label="名称" prop="username">
              <el-input
                v-model="account_form.username"
                placeholder=""
              ></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="别名" prop="nickname">
              <el-input
                v-model="account_form.nickname"
                placeholder=""
              ></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="描述" prop="desc">
              <el-input v-model="account_form.desc" placeholder=""></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="角色" prop="role">
              <el-select
                v-model="account_form.role"
                placeholder=""
                filterable
                clearable
              >
                <el-option
                  v-for="(item, index) in role_list"
                  :key="index"
                  :label="item.name"
                  :value="item.name"
                ></el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button size="small" @click="account_dialog = false"
          >取 消</el-button
        >
        <el-button type="primary" size="small" @click="submit_account"
          >确 定</el-button
        >
      </span>
    </el-dialog>

    <!-- 角色管理--------------------------------------------------  -->
    <el-dialog
      :title="textMap[dialogStatus]"
      :visible.sync="role_dialog"
      width="70%"
    >
      <el-form
        ref="role_form_refs"
        :model="role_form"
        :rules="role_form_rules"
        size="small"
        label-width="100px"
      >
        <el-row>
          <el-col :span="12">
            <el-form-item label="名称" prop="name">
              <el-input v-model="role_form.name" placeholder=""></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="描述" prop="desc">
              <el-input v-model="role_form.desc" placeholder=""></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="命名空间" prop="namespace">
              <el-input
                v-model="role_form.namespace"
                style="width: 130px"
              ></el-input>
              <el-button
                size="small"
                style="margin-left: 10px"
                @click="add_namespace()"
                >添加</el-button
              >
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-tabs
            v-model="default_namespace"
            @tab-remove="remove_namespace"
            style="margin-left: 100px"
            closable
          >
            <el-tab-pane
              v-for="item in role_form.ns_list"
              :key="item.namespace"
              :name="item.namespace || '默认'"
              :label="item.namespace || '默认'"
            >
              <el-row>
                <el-col>
                  <el-checkbox
                    v-model="checkedAll"
                    label="全选"
                    @change="
                      (e) => {
                        handleSelecteAll(e, item);
                      }
                    "
                    style="margin-bottom: 20px"
                  ></el-checkbox>
                  <el-button
                    size="mini"
                    style="margin-left: 20px; vertical-align: middle"
                    @click="copy_verbs"
                    >复制</el-button
                  >
                  <el-button
                    v-show="can_paste"
                    size="mini"
                    style="margin-left: 20px; vertical-align: middle"
                    >粘贴</el-button
                  >
                </el-col>
              </el-row>
              <el-row>
                <el-col
                  v-for="(tempRole, index) in gvk_list"
                  :span="12"
                  :key="index"
                  style="margin-bottom: 20px"
                >
                  <el-checkbox
                    v-model="item.type[index].checked"
                    style="margin-bottom: 10px"
                    :indeterminate="
                      isIndeterminate(
                        item.type[index].verbs,
                        gvk_list[index].verbs
                      )
                    "
                    @change="
                      (e) => {
                        handleCheckAllChange(
                          e,
                          item.type[index].verbs,
                          gvk_list[index].verbs
                        );
                      }
                    "
                    >{{ tempRole.name }}
                  </el-checkbox>

                  <el-checkbox-group
                    v-model="item.type[index].verbs"
                    size="small"
                    style="margin-left: 20px"
                  >
                    <el-checkbox-button
                      v-for="verb in tempRole.verbs"
                      :key="verb"
                      :label="verb"
                      >{{ verb }}
                    </el-checkbox-button>
                  </el-checkbox-group>
                </el-col>
              </el-row>
            </el-tab-pane>
          </el-tabs>
        </el-row>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button size="small" @click="role_dialog = false">取 消</el-button>
        <el-button type="primary" size="small" @click="submit_role"
          >确 定</el-button
        >
      </span>
    </el-dialog>
    <!-- GVK管理--------------------------------------------------  -->
    <el-dialog
      :title="textMap[dialogStatus]"
      :visible.sync="gvk_dialog"
      width="50%"
    >
      <el-form
        ref="gvk_form_refs"
        :model="gvk_form"
        :rules="gvk_form_rules"
        size="small"
        label-width="100px"
      >
        <el-row>
          <el-col :span="12">
            <el-form-item label="GVK" prop="gvk">
              <el-input v-model="gvk_form.gvk" placeholder=""></el-input>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button size="small" @click="gvk_dialog = false">取 消</el-button>
        <el-button type="primary" size="small" @click="submit_gvk"
          >确 定</el-button
        >
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { cloneDeep } from "lodash";
import { rbacApi } from "@/utils/request";
export default {
  name: "AuthXRBAC",
  data() {
    return {
      textMap: {
        create_account: "新增账户",
        update_account: "编辑账户",
        create_role: "新增角色",
        update_role: "编辑角色",
        create_gvk: "新增GVK",
        update_gvk: "编辑GVK",
      },
      appId: localStorage.getItem("appId"),
      appName: localStorage.getItem("appName"),
      currentPage: 1,
      dialogStatus: "",
      active_tab_name: "account_setting",
      account_list: [],
      account_dialog: false,
      account_form: {
        id: null,
        username: null,
        nickname: null,
        password: null,
        role: null,
        desc: null,
      },
      account_form_rules: {},
      // 角色管理--------------------------------------------------
      itemSel: [],
      default_namespace: "默认",
      can_paste: false,
      clone_verbs: [],
      role_list: [],
      ns_list: [],
      role_dialog: false,
      role_form: {
        id: null,
        name: null,
        desc: null,
        ns_list: [],
      },
      role_form_rules: {},
      checkedAll: false,
      // gvk管理--------------------------------------------------
      gvk_list: [],
      gvk_dialog: false,
      gvk_form: {},
      gvk_form_rules: {},
    };
  },
  computed: {
    page_gvk_list: function () {
      return this.gvk_list.slice(
        (this.currentPage - 1) * 10,
        this.currentPage * 10
      );
    },
  },
  watch: {
    default_namespace: function (newVal, oldVal) {
      let nIndex = this.role_form.ns_list.findIndex((fni) => {
        return fni.namespace === newVal;
      });
      this.checkedAll = false;
      this.itemSel = this.role_form.ns_list[nIndex];
    },
    immediate: true,
  },
  mounted() {
    this.get_account_list();
  },
  methods: {
    goback() {
      localStorage.removeItem("appId");
      localStorage.removeItem("appName");
      this.$router.push({ path: "/" });
    },
    copy(row) {
      navigator.clipboard.writeText(row).then(() => {
        this.$message({
          type: "success",
          message: "复制成功",
        });
      });
    },
    tab_click(tab) {
      if (tab.name === "gvk_setting") {
        this.get_gvk_list();
      } else if (tab.name === "role_setting") {
        this.get_role_list();
        this.get_gvk_list();
      } else if (tab.name === "account_setting") {
        this.get_account_list();
      }
    },
    create_account() {
      this.account_dialog = true;
      this.dialogStatus = "create_account";
      this.account_form = Object.assign({}, "");
      this.get_role_list();
    },
    update_account(row) {
      this.account_dialog = true;
      this.dialogStatus = "update_account";
      this.account_form = Object.assign({}, row);
      this.get_role_list();
    },
    reset_password(row) {
      const delete_data = {
        metadata: {
          namespace: localStorage.getItem("appId"),
          name: row.username,
        },
      };
      const create_data = {
        metadata: {
          namespace: localStorage.getItem("appId"),
        },
        spec: {
          accountMeta: {
            username: row.username,
          },
          roleRef: {
            clusterRoleName: row.role,
          },
          desc: row.desc,
        },
      };
      rbacApi("AppServiceAccount", "/appaccount/delete", delete_data)
        .then((resp) => {
          rbacApi("AppServiceAccount", "/appaccount/create", create_data)
            .then((resp) => {
              this.$message({
                type: "success",
                message: "重置成功",
              });
              this.get_account_list();
            })
            .catch((err) => {
              this.$message({
                type: "error",
                message: err.message,
              });
            });
        })
        .catch((err) => {
          this.$message({
            type: "error",
            message: err.message,
          });
        });
    },
    delete_account(row) {
      const delete_data = {
        metadata: {
          namespace: localStorage.getItem("appId"),
          name: row.username,
        },
      };
      rbacApi("AppServiceAccount", "/appaccount/delete", delete_data).then(
        (resp) => {
          this.get_account_list();
        }
      );
    },
    get_account_list() {
      const get_data = {
        metadata: {
          namespace: localStorage.getItem("appId"),
          name: "123",
        },
      };
      rbacApi(
        "AppServiceAccountList",
        "/appaccount/list",
        get_data,
        "",
        "AppServiceAccount"
      ).then((resp) => {
        const accounts = [];
        resp.items.sort((acA, acB) => {
          return (
            acA.metadata.creationTimestamp.seconds -
            acB.metadata.creationTimestamp.seconds
          );
        });
        for (const ac of resp.items) {
          accounts.push({
            username: ac.spec.accountMeta.username,
            nickname: ac.spec.accountMeta.nickname,
            password: ac.spec.accountMeta.password,
            role: ac.spec.roleRef.clusterRoleName,
            desc: ac.spec.desc,
          });
        }
        this.account_list = accounts;
      });
    },
    submit_account() {
      if (this.dialogStatus === "create_account") {
        const create_data = {
          metadata: {
            namespace: localStorage.getItem("appId"),
          },
          spec: {
            accountMeta: {
              username: this.account_form.username,
              nickname: this.account_form.nickname,
            },
            roleRef: {
              clusterRoleName: this.account_form.role,
            },
            desc: this.account_form.desc,
          },
        };
        rbacApi("AppServiceAccount", "/appaccount/create", create_data).then(
          (resp) => {
            this.account_dialog = false;
            this.get_account_list();
          }
        );
      } else if (this.dialogStatus === "update_account") {
        const update_data = {
          metadata: {
            namespace: localStorage.getItem("appId"),
          },
          spec: {
            accountMeta: {
              username: this.account_form.username,
              nickname: this.account_form.nickname,
            },
            roleRef: {
              clusterRoleName: this.account_form.role,
            },
            desc: this.account_form.desc,
          },
        };
        rbacApi("AppServiceAccount", "/appaccount/update", update_data).then(
          (resp) => {
            this.account_dialog = false;
            this.get_account_list();
          }
        );
      }
    },
    // 角色管理--------------------------------------------------
    create_role() {
      this.role_dialog = true;
      this.dialogStatus = "create_role";
      this.role_form = Object.assign({}, "");
      this.default_namespace = "默认";
      const typeList = [];
      for (const kind of this.gvk_list) {
        typeList.push({
          name: kind.name,
          verbs: [],
        });
      }
      this.role_form.ns_list = [
        {
          namespace: "",
          type: typeList,
        },
      ];
    },
    update_role(row) {
      this.role_dialog = true;
      this.dialogStatus = "update_role";
      this.role_form = Object.assign({}, row);
      let vbs = [];
      for (let roleNs of this.role_form.ns_list) {
        let npType = [];
        for (let kind of this.gvk_list) {
          npType.push({
            name: kind.name,
            verbs: [],
          });
        }
        vbs.push({
          namespace: roleNs.namespace,
          type: npType,
        });
      }
      for (let nsRole of vbs) {
        let editIndex = this.role_form.ns_list.findIndex((val) => {
          return val.namespace === nsRole.namespace;
        });
        for (let oneType of nsRole.type) {
          for (let editRes of this.role_form.ns_list[editIndex].type) {
            if (oneType.name === editRes.gvk) {
              oneType.verbs = editRes.verbs;
            }
          }
        }
      }

      this.role_form.ns_list = vbs;
      if (vbs.length > 0) {
        this.default_namespace = vbs[0].namespace || "默认";
      }
      for (let ns of vbs) {
        for (let vbIndex in ns.type) {
          if (
            this.verbCompare(
              ns.type[vbIndex].verbs,
              this.gvk_list[Number(vbIndex)].verbs
            )
          ) {
            ns.type[vbIndex].checked = true;
          } else {
            ns.type[vbIndex].checked = false;
          }
        }
      }
      // console.log(this.role_form, "====");
    },
    verbCompare(vera, verb) {
      for (let ver of verb) {
        if (!vera.includes(ver)) {
          return false;
        }
      }
      return true;
    },
    delete_role(row) {
      const roleInfo = {
        metadata: {
          namespace: localStorage.getItem("appId"),
          name: row.name,
        },
      };
      rbacApi(
        "ClusterRoleList",
        "/clusterrole/delete",
        roleInfo,
        "",
        "ClusterRole"
      )
        .then((resp) => {
          this.get_role_list();
        })
        .catch((err) => {
          this.$message({
            type: "error",
            message: err.message,
          });
        });
    },
    get_role_list() {
      const get_data = {
        metadata: {
          namespace: localStorage.getItem("appId"),
        },
      };
      rbacApi(
        "ClusterRoleList",
        "/clusterrole/list",
        get_data,
        "",
        "ClusterRole"
      )
        .then((resp) => {
          resp.items.sort((acA, acB) => {
            return (
              acA.metadata.creationTimestamp.seconds -
              acB.metadata.creationTimestamp.seconds
            );
          });
          this.role_list = [];
          for (const role of resp.items) {
            const ns_list = [];
            for (const policy of role.spec.rules) {
              const policyGvk = policy.groupVersionKind;
              let findI = ns_list.findIndex((val) => {
                return val.namespace === policy.namespace;
              });
              if (policy.verbs.length > 0) {
                if (findI >= 0) {
                  ns_list[findI].type.push({
                    gvk: `${policyGvk.group ? policyGvk.group + "." : ""}${
                      policyGvk.version ? policyGvk.version + "." : ""
                    }${policyGvk.kind}`,
                    verbs: policy.verbs,
                  });
                } else {
                  ns_list.push({
                    namespace: policy.namespace,
                    type: [
                      {
                        gvk: `${policyGvk.group ? policyGvk.group + "." : ""}${
                          policyGvk.version ? policyGvk.version + "." : ""
                        }${policyGvk.kind}`,
                        verbs: policy.verbs,
                      },
                    ],
                  });
                }
              }
            }
            this.role_list.push({
              name: role.metadata.name,
              desc: role.spec.desc,
              ns_list,
            });
          }
        })
        .catch((err) => {
          this.$message({
            type: "error",
            message: err.message,
          });
        });
    },
    copy_verbs() {
      this.can_paste = true;
    },
    add_namespace() {
      if (
        this.role_form.namespace === undefined ||
        this.role_form.namespace === ""
      ) {
        return;
      }
      let typeArr = [];
      const newIdnex = this.role_form.ns_list.findIndex((ns) => {
        return ns.namespace === this.role_form.namespace;
      });
      if (newIdnex >= 0) return;
      for (let kind of this.gvk_list) {
        typeArr.push({
          name: kind.name,
          verbs: [],
        });
      }
      this.role_form.ns_list.push({
        namespace: this.role_form.namespace,
        type: typeArr,
      });
      this.default_namespace = this.role_form.namespace;
      this.role_form.namespace = "";
    },
    remove_namespace(tab) {
      const removeIndex = this.role_form.ns_list.findIndex((val) => {
        return val.namespace === tab;
      });
      if (removeIndex === -1) {
        this.role_form.ns_list.shift();
      } else {
        this.role_form.ns_list.splice(removeIndex, 1);
      }
    },
    handleSelecteAll(checked, item) {
      if (checked) {
        item.type = cloneDeep(this.gvk_list);
        for (let ver of item.type) {
          ver.checked = true;
        }
      } else {
        for (let everyItem of item.type) {
          everyItem.verbs = [];
          everyItem.checked = false;
        }
      }
    },
    isIndeterminate(verbs, kinds) {
      return verbs.length > 0 && verbs.length < kinds.length;
    },
    handleCheckAllChange(val, verbs, kinds) {
      if (val) {
        for (let k of kinds) {
          if (!verbs.includes(k)) verbs.push(k);
        }
      } else {
        verbs.splice(0, verbs.length);
      }
    },
    submit_role() {
      if (this.dialogStatus === "create_role") {
        const gvkVerbs = this.role_form.ns_list;
        let rules = [];
        for (let gvkverb of gvkVerbs) {
          for (let allverb of gvkverb.type) {
            const gvk = allverb.name.split(".");
            let g = "",
              v = "",
              k = "";
            if (gvk.length === 3) {
              (g = gvk[0]), (v = gvk[1]), (k = gvk[2]);
            } else if (gvk.length > 3) {
              (k = gvk[gvk.length - 1]), (v = gvk[gvk.length - 2]);
              let group = "";
              for (let sar = 0; sar < gvk.length - 2; sar++) {
                group = group + gvk[sar] + ".";
              }
              g = group.substring(0, group.length - 1);
            } else if (gvk.length === 2) {
              (g = ""), (v = gvk[0]), (k = gvk[1]);
            } else {
              k = gvk[0];
            }

            if (gvkVerbs.length === 1 || gvkverb.namespace) {
              rules.push({
                namespace: gvkverb.namespace,
                groupVersionKind: { group: g, version: v, kind: k },
                verbs: allverb.verbs,
              });
            }
          }
        }
        const create_role_data = {
          metadata: {
            namespace: this.appId,
            name: this.role_form.name,
          },
          spec: {
            desc: this.role_form.desc,
            rules,
          },
        };
        rbacApi(
          "ClusterRoleList",
          "/clusterrole/create",
          create_role_data,
          "",
          "ClusterRole"
        )
          .then((resp) => {
            this.role_dialog = false;
            this.get_role_list();
          })
          .catch((err) => {
            this.$message({
              type: "error",
              message: err.message,
            });
          });
      } else if (this.dialogStatus === "update_role") {
        const gvkVerbs = this.role_form.ns_list;
        let rules = [];
        for (let gvkverb of gvkVerbs) {
          for (let allverb of gvkverb.type) {
            const gvk = allverb.name.split(".");
            let g = "",
              v = "",
              k = "";
            if (gvk.length === 3) {
              (g = gvk[0]), (v = gvk[1]), (k = gvk[2]);
            } else if (gvk.length > 3) {
              (k = gvk[gvk.length - 1]), (v = gvk[gvk.length - 2]);
              let group = "";
              for (let sar = 0; sar < gvk.length - 2; sar++) {
                group = group + gvk[sar] + ".";
              }
              g = group.substring(0, group.length - 1);
            } else if (gvk.length === 2) {
              (g = ""), (v = gvk[0]), (k = gvk[1]);
            } else {
              k = gvk[0];
            }

            if (gvkVerbs.length === 1 || gvkverb.namespace) {
              rules.push({
                namespace: gvkverb.namespace,
                groupVersionKind: { group: g, version: v, kind: k },
                verbs: allverb.verbs,
              });
            }
          }
        }

        const update_role_data = {
          metadata: {
            namespace: this.appId,
            name: this.role_form.name,
          },
          spec: {
            desc: this.role_form.desc,
            rules,
          },
        };
        // console.log(update_role_data, "=====");
        rbacApi(
          "ClusterRoleList",
          "/clusterrole/update",
          update_role_data,
          "",
          "ClusterRole"
        )
          .then((resp) => {
            this.role_dialog = false;
            this.get_role_list();
          })
          .catch((err) => {
            this.$message({
              type: "error",
              message: err.message,
            });
          });
      }
    },
    // gvk管理--------------------------------------------------
    create_gvk() {
      this.gvk_dialog = true;
      this.dialogStatus = "create_gvk";
    },
    update_gvk(row) {
      this.gvk_dialog = true;
      this.dialogStatus = "update_gvk";
    },
    delete_gvk(row) {},
    get_gvk_list() {
      const get_data = {
        metadata: {
          namespace: localStorage.getItem("appId"),
        },
      };
      rbacApi(
        "GroupVersionKindRuleList",
        "/kind/list",
        get_data,
        "",
        "GroupVersionKindRule"
      ).then((resp) => {
        this.gvk_list = [];
        const kinds = [];
        const nameLabels = [];
        for (const kind of resp.items) {
          const gvk = kind.spec.groupVersionKind;
          nameLabels.push(
            `${gvk.group ? gvk.group + "." : ""}${
              gvk.version ? gvk.version + "." : ""
            }${gvk.kind}`
          );
          kinds.push({
            name: `${gvk.group ? gvk.group + "." : ""}${
              gvk.version ? gvk.version + "." : ""
            }${gvk.kind}`,
            verbs: kind.spec.verbs,
          });
        }
        nameLabels.sort();
        const sortKinds = [];
        for (let name of nameLabels) {
          const fIndex = kinds.findIndex((k) => {
            return k.name === name;
          });
          sortKinds.push(kinds[fIndex]);
        }
        this.gvk_list = sortKinds;
      });
    },
    submit_gvk() {
      if (this.dialogStatus === "create_gvk") {
        //
      } else if (this.dialogStatus === "update_gvk") {
        //
      }
    },
    cancel_delete() {
      this.$message({
        type: "warning",
        message: "你考虑的很全面",
      });
    },
  },
};
</script>

<style scoped>
.el-input {
  width: 200px;
}

.el-select {
  width: 200px;
}

.el-table {
  width: 100%;
  margin-top: 10px;
}

.el-button {
  vertical-align: top;
}
.verb-label {
  width: 75px;
  font-weight: 600;
  text-align: left;
  display: flex;
  padding-right: 15px;
}
</style>
