镜像仓库
<template>
  <div class="app-container">
    <el-card class="box-card">
      <el-button
        type="primary"
        size="small"
        icon="el-icon-circle-plus-outline"
        @click="create_ns"
        >新增</el-button
      >
      <el-table :data="ns_list" size="small" empty-text="啥也没有" border>
        <el-table-column label="项目" prop="name"></el-table-column>
        <el-table-column label="公网" prop="publicAddr"></el-table-column>
        <el-table-column label="内网" prop="privateAddr"></el-table-column>
        <el-table-column label="创建时间">
          <template slot-scope="scoped">
            {{ scoped.row.create_time | parseTime("{y}-{m}-{d} {h}:{i}:{s}") }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180px;">
          <template slot-scope="scoped">
            <el-button
              type="success"
              size="mini"
              icon="el-icon-s-promotion"
              @click="goto_repo(scoped.row)"
            ></el-button>
            <el-button
              type="primary"
              icon="el-icon-edit"
              size="mini"
              @click="update_ns(scoped.row)"
            ></el-button>
            <el-popconfirm
              title="确定删除吗？"
              confirm-button-text="确定"
              cancel-button-text="不了"
              style="margin-left: 10px"
              @confirm="delete_ns(scoped.row)"
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
        :total="ns_list.length"
        layout="total, prev, pager, next"
        style="text-align: left; margin-top: 20px"
      >
      </el-pagination>
    </el-card>

    <el-dialog
      :title="textMap[dialogStatus]"
      :visible.sync="ns_dialog"
      :show-close="true"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      width="50%"
    >
      <el-form
        ref="ns_form_refs"
        :model="ns_form"
        :rules="ns_form_rules"
        size="small"
        label-width="100px"
      >
        <el-row>
          <el-col :span="12">
            <el-form-item label="名称" prop="name">
              <el-input v-model="ns_form.name" placeholder=""></el-input>
            </el-form-item>
          </el-col>
          <!-- <el-col :span="12">
            <el-form-item label="描述" prop="desc">
              <el-input v-model="ns_form.desc" placeholder=""></el-input>
            </el-form-item>
          </el-col> -->
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="地址1" prop="addr1">
              <el-input v-model="ns_form.publicAddr" placeholder=""></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="地址2" prop="addr2">
              <el-input v-model="ns_form.privateAddr" placeholder=""></el-input>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button size="small" @click="ns_dialog = false">取 消</el-button>
        <el-button type="primary" size="small" @click="submit_ns"
          >确 定</el-button
        >
      </span>
    </el-dialog>
  </div>
</template>

<script>
import store from "@/store";
import { parseTime } from "@/utils";
import { mapGetters } from "vuex";
import {
  init_socket_data,
  sendSocketMessage,
  getProtoParam,
} from "@/api/websocket";
import protoRoot from "@/proto/proto.js";
const protoApi = protoRoot.github.com.kzz45.discovery.pkg.apis.jingx.v1;
const protoRequest = protoRoot.github.com.kzz45.discovery.pkg.jingx.proto;

export default {
  name: "RegistryNS",
  filters: {
    parseTime(time, cFormat) {
      return parseTime(time, cFormat);
    },
  },
  data() {
    return {
      textMap: {
        create_ns: "新增项目",
        update_ns: "编辑项目",
      },
      dialogStatus: "",
      ns_list: [],
      currentPage: 1,
      ns_dialog: false,
      ns_form: {
        id: null,
        name: null,
        desc: null,
        publicAddr: null,
        privateAddr: null,
      },
      ns_form_rules: {},
    };
  },
  mounted() {
    this.get_project_list();
  },
  watch: {
    message: function () {
      this.socket_onmessage(this.message);
    },
  },
  computed: {
    ...mapGetters(["message"]),
  },
  created() {},
  methods: {
    goto_repo(row) {
      localStorage.setItem("projectName", row.name);
      this.$router.push({ path: "/repo/repo" });
    },
    create_ns() {
      this.ns_dialog = true;
      this.dialogStatus = "create_ns";
      this.ns_form = Object.assign({}, "");
    },
    update_ns(row) {
      this.ns_dialog = true;
      this.dialogStatus = "update_ns";
      this.ns_form = Object.assign({}, row);
    },
    delete_ns(row) {
      const projectObj = new protoApi.Project();
      projectObj.metadata = { name: row.name };
      const gvk = {
        group: "jingx",
        version: "v1",
        kind: "Project",
      };
      const senddata = getProtoParam(projectObj, gvk);
      const new_delete_data = init_socket_data(
        "discovery-jingx",
        "jingx-v1-Project",
        "delete",
        senddata
      );
      sendSocketMessage(new_delete_data, store);
      this.get_project_list();
    },
    submit_ns() {
      if (this.dialogStatus === "create_ns") {
        const projectObj = new protoApi.Project();
        projectObj.metadata = { name: this.ns_form.name };
        projectObj.spec = {
          id: "",
          publicAddr: this.ns_form.publicAddr,
          privateAddr: this.ns_form.privateAddr,
        };
        const gvk = {
          group: "jingx",
          version: "v1",
          kind: "Project",
        };
        const senddata = getProtoParam(projectObj, gvk);
        const new_create_data = init_socket_data(
          "discovery-jingx",
          "jingx-v1-Project",
          "create",
          senddata
        );
        sendSocketMessage(new_create_data, store);
        this.ns_dialog = false;
        this.get_project_list();
      } else if (this.dialogStatus === "update_ns") {
        const projectObj = new protoApi.Project();
        projectObj.metadata = { name: this.ns_form.name };
        projectObj.spec = {
          id: "",
          publicAddr: this.ns_form.publicAddr,
          privateAddr: this.ns_form.privateAddr,
        };
        const gvk = {
          group: "jingx",
          version: "v1",
          kind: "Project",
        };
        const senddata = getProtoParam(projectObj, gvk);
        const new_update_data = init_socket_data(
          "discovery-jingx",
          "jingx-v1-Project",
          "update",
          senddata
        );
        sendSocketMessage(new_update_data, store);
        this.ns_dialog = false;
        this.get_project_list();
      }
    },
    get_project_list() {
      const send_data = init_socket_data(
        "discovery-jingx",
        "jingx-v1-Project",
        "list"
      );
      sendSocketMessage(send_data, store);
    },
    socket_onmessage(msg) {
      const result = protoRequest.Response.decode(msg);
      // console.log(result, "-----------");
      if (result.code === 1) {
        const err_msg = String.fromCharCode.apply(null, result.raw);
        this.$message({
          type: "error",
          message: err_msg,
        });
      }
      if (result.verb === "list") {
        const project_list = protoApi[
          `${result.groupVersionKind.kind}List`
        ].decode(result.raw).items;

        project_list.sort((itemA, itemB) => {
          return (
            itemA.metadata.creationTimestamp.seconds -
            itemB.metadata.creationTimestamp.seconds
          );
        });
        this.ns_list = [];
        for (let pl of project_list) {
          this.ns_list.push({
            name: pl.metadata.name,
            publicAddr: pl.spec.publicAddr,
            privateAddr: pl.spec.privateAddr,
            create_time: pl.metadata.creationTimestamp.seconds,
          });
        }
      }
    },
    copy(row) {
      navigator.clipboard.writeText(row).then(() => {
        this.$message({
          type: "success",
          message: "复制成功",
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
