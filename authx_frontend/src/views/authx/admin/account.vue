认证权限
<template>
  <div class="app-container">
    <el-card class="box-card">
      <!-- 账户管理--------------------------------------------------  -->
      <el-button
        type="primary"
        size="small"
        icon="el-icon-circle-plus-outline"
        @click="craete_account"
        >新增</el-button
      >
      <el-table :data="account_list" size="small" empty-text="啥也没有" border>
        <el-table-column prop="username" label="账户"></el-table-column>
        <el-table-column label="密码">
          <template slot-scope="scoped">
            <el-button type="text" @click="copy(scoped.row.password)"
              >复制</el-button
            >
          </template>
        </el-table-column>
        <el-table-column label="App">
          <template slot-scope="scoped">
            <el-tag
              v-for="item in scoped.row.apps"
              :key="item"
              :item="item"
              type="primary"
              size="small"
              style="margin-right: 10px; margin-bottom: 10px"
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
            <el-form-item label="APPS" prop="apps">
              <el-select
                v-model="account_form.apps"
                placeholder=""
                multiple
                filterable
              >
                <el-option
                  v-for="(item, index) in app_list"
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
  </div>
</template>

<script>
import { rbacApi } from "@/utils/request";
import { sortByName } from "@/utils/utils.js";

export default {
  name: "AuthxAdmin",
  data() {
    return {
      textMap: {
        create_account: "新增账户",
        update_account: "编辑账户",
        create_app: "新增App",
        update_app: "编辑App",
      },
      dialogStatus: "",
      active_tab_name: "apps_setting",
      account_list: [],
      account_dialog: false,
      account_form: {
        id: null,
        username: null,
        password: null,
        apps: [],
      },
      account_form_rules: {},
      app_list: [],
      app_dialog: false,
      app_form: {
        id: null,
        name: null,
        desc: null,
        secret: null,
      },
      app_form_rules: {},
    };
  },
  mounted() {
    this.get_account_list();
  },
  methods: {
    copy(row) {
      navigator.clipboard.writeText(row).then(() => {
        this.$message({
          type: "success",
          message: "复制成功",
        });
      });
    },
    craete_account() {
      this.account_dialog = true;
      this.dialogStatus = "create_account";
      this.account_form = Object.assign({}, "");
      this.get_app_list();
    },
    update_account(row) {
      this.account_dialog = true;
      this.dialogStatus = "update_account";
      this.account_form = Object.assign({}, row);
      this.get_app_list();
    },
    delete_account(row) {
      const delete_data = {
        metadata: {
          name: row.username,
        },
        accountMeta: {
          username: row.username,
        },
      };
      rbacApi("RbacServiceAccount", "/rbacaccount/delete", delete_data)
        .then((resp) => {
          this.account_dialog = false;
          this.get_account_list();
        })
        .catch((err) => {
          this.$message({
            type: "error",
            message: err.message,
          });
        });
    },
    submit_account() {
      if (this.dialogStatus === "create_account") {
        const create_data = {
          metadata: {
            name: this.account_form.username,
          },
          accountMeta: {
            username: this.account_form.username,
          },
          spec: {
            apps: this.account_form.apps,
          },
        };
        rbacApi("RbacServiceAccount", "/rbacaccount/create", create_data)
          .then((resp) => {
            this.account_dialog = false;
            this.get_account_list();
          })
          .catch((err) => {
            this.$message({
              type: "error",
              message: err.message,
            });
          });
      } else if (this.dialogStatus === "update_account") {
        const update_data = {
          metadata: {
            name: this.account_form.username,
          },
          accountMeta: {
            username: this.account_form.username,
          },
          spec: {
            apps: this.account_form.apps,
          },
        };
        rbacApi("RbacServiceAccount", "/rbacaccount/update", update_data)
          .then((resp) => {
            this.account_dialog = false;
            this.get_account_list();
          })
          .catch((err) => {
            this.$message({
              type: "error",
              message: err.message,
            });
          });
      }
    },
    get_account_list() {
      rbacApi("RbacServiceAccountList", "/rbacaccount/list")
        .then((resp) => {
          this.account_list = [];
          for (const ac of resp.items) {
            this.account_list.push({
              username: ac.spec.accountMeta.username,
              password: ac.spec.accountMeta.password,
              apps: ac.spec.apps,
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
    get_app_list() {
      rbacApi("AppList", "/apps/list")
        .then((resp) => {
          this.app_list = [];
          const temp_list = [];
          for (const app of resp.items) {
            temp_list.push({
              id: app.spec.id,
              name: app.metadata.name,
              desc: app.spec.desc,
              secret: app.spec.secret,
            });
          }
          this.app_list = sortByName(temp_list);
        })
        .catch((err) => {
          this.$message({
            type: "error",
            message: err.message,
          });
        });
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
</style>
