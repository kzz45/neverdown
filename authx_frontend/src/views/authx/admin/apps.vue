认证权限
<template>
  <div class="app-container">
    <el-card class="box-card">
      <el-button
        type="primary"
        size="small"
        icon="el-icon-circle-plus-outline"
        @click="create_app"
        >新增</el-button
      >
      <el-table :data="app_list" size="small" empty-text="啥也没有" border>
        <el-table-column prop="name" label="名称"></el-table-column>
        <el-table-column prop="desc" label="描述"></el-table-column>
        <el-table-column label="Appp ID" prop="id">
          <!-- <template slot-scope="scoped">
            <el-button type="text" @click="copy(scoped.row.id)">复制</el-button>
          </template> -->
        </el-table-column>
        <el-table-column label="App Secret" prop="secret">
          <!-- <template slot-scope="scoped">
            <el-button type="text" @click="copy(scoped.row.secret)"
              >复制</el-button
            >
          </template> -->
        </el-table-column>
        <el-table-column label="操作" width="180px;">
          <template slot-scope="scoped">
            <el-button
              type="success"
              size="mini"
              icon="el-icon-s-promotion"
              @click="goto_app(scoped.row)"
            ></el-button>
            <el-button
              type="primary"
              icon="el-icon-edit"
              size="mini"
              @click="update_app(scoped.row)"
            ></el-button>
            <el-popconfirm
              title="确定删除吗？"
              confirm-button-text="确定"
              cancel-button-text="不了"
              style="margin-left: 10px"
              @confirm="delete_app(scoped.row)"
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
        :total="app_list.length"
        layout="total, prev, pager, next"
        style="text-align: left; margin-top: 20px"
      >
      </el-pagination>
    </el-card>

    <!-- 应用管理--------------------------------------------------  -->
    <el-dialog
      :title="textMap[dialogStatus]"
      :visible.sync="app_dialog"
      width="50%"
    >
      <el-form
        ref="app_form_refs"
        :model="app_form"
        :rules="app_form_rules"
        size="small"
        label-width="100px"
      >
        <el-row>
          <el-col :span="12">
            <el-form-item label="App名称" prop="name">
              <el-input
                v-model="app_form.name"
                size="small"
              ></el-input> </el-form-item
          ></el-col>
          <el-col :span="12">
            <el-form-item label="App描述" prop="desc">
              <el-input v-model="app_form.desc" size="small"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button size="small" @click="app_dialog = false">取 消</el-button>
        <el-button type="primary" size="small" @click="submit_app"
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
        create_app: "新增App",
        update_app: "编辑App",
      },
      dialogStatus: "",
      app_list: [],
      currentPage: 1,
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
    this.get_app_list();
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
    create_app() {
      this.app_dialog = true;
      this.dialogStatus = "create_app";
      this.app_form = Object.assign({}, "");
    },
    update_app(row) {
      this.app_dialog = true;
      this.dialogStatus = "update_app";
      this.app_form = Object.assign({}, row);
    },
    delete_app(row) {
      const delete_data = {
        metadata: { name: row.name },
      };
      rbacApi("App", "/apps/delete", delete_data)
        .then((resp) => {
          this.app_dialog = false;
          this.get_app_list();
        })
        .catch((err) => {
          this.$message({
            type: "error",
            message: err.message,
          });
        });
    },
    submit_app() {
      if (this.dialogStatus === "create_app") {
        const create_data = {
          metadata: { name: this.app_form.name },
          spec: { desc: this.app_form.desc },
        };
        rbacApi("App", "/apps/create", create_data)
          .then((resp) => {
            this.app_dialog = false;
            this.get_app_list();
          })
          .catch((err) => {
            this.$message({
              type: "error",
              message: err.message,
            });
          });
      } else if (this.dialogStatus === "update_app") {
        const update_data = {
          metadata: { name: this.app_form.name },
          spec: { desc: this.app_form.desc },
        };
        rbacApi("App", "/apps/update", update_data)
          .then((resp) => {
            this.app_dialog = false;
            this.get_app_list();
          })
          .catch((err) => {
            this.$message({
              type: "error",
              message: err.message,
            });
          });
      }
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
    goto_app(row) {
      // console.log(row, "=====");
      localStorage.setItem("appId", row.id);
      localStorage.setItem("appName", row.name);
      this.$router.push({ path: "/authx" });
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
