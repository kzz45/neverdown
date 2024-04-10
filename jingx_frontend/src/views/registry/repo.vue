镜像仓库
<template>
  <div class="app-container">
    <el-card class="box-card">
      <el-button size="small" icon="el-icon-back" @click="goback"
        >返回</el-button
      >
      <el-table
        :data="page_repo_list"
        size="small"
        empty-text="啥也没有"
        border
      >
        <el-table-column label="Repo名称" prop="name"></el-table-column>
        <el-table-column label="项目名称" prop="project_name"></el-table-column>
        <el-table-column label="Tag"></el-table-column>
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
              @click="goto_tag(scoped.row)"
            ></el-button>
            <el-button
              type="primary"
              icon="el-icon-edit"
              size="mini"
              @click="update_repo(scoped.row)"
            ></el-button>
            <el-popconfirm
              title="确定删除吗？"
              confirm-button-text="确定"
              cancel-button-text="不了"
              style="margin-left: 10px"
              @confirm="delete_repo(scoped.row)"
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
        :total="this.repo_list.length"
        layout="total, prev, pager, next"
        style="text-align: left; margin-top: 20px"
      >
      </el-pagination>
    </el-card>
  </div>
</template>

<script>
import store from "@/store";
import { mapGetters } from "vuex";
import { parseTime } from "@/utils";
import { init_socket_data, sendSocketMessage } from "@/api/websocket";
import protoRoot from "@/proto/proto.js";
const protoApi = protoRoot.github.com.kzz45.neverdown.pkg.apis.jingx.v1;
const protoRequest = protoRoot.github.com.kzz45.neverdown.pkg.jingx.proto;

export default {
  name: "RegistryRepo",
  filters: {
    parseTime(time, cFormat) {
      return parseTime(time, cFormat);
    },
  },
  data() {
    return {
      repo_list: [],
      currentPage: 1,
    };
  },
  watch: {
    message: function () {
      this.socket_onmessage(this.message);
    },
  },
  computed: {
    ...mapGetters(["message"]),
    page_repo_list: function () {
      return this.repo_list.slice(
        (this.currentPage - 1) * 10,
        this.currentPage * 10
      );
    },
  },
  created() {
    this.socket_connect();
  },
  methods: {
    goback() {
      this.$router.push({ path: "/" });
    },
    goto_tag(row) {
      localStorage.setItem("repoName", row.name);
      this.$router.push({ path: "/tag/tag" });
    },
    update_repo() {},
    delete_repo() {},
    socket_connect() {
      const send_data = init_socket_data(
        "discovery-jingx",
        "jingx-v1-Repository",
        "list"
      );
      sendSocketMessage(send_data, store);
    },
    socket_onmessage(msg) {
      const result = protoRequest.Response.decode(msg);
      if (result.verb === "list") {
        const project_list = protoApi[
          `${result.groupVersionKind.kind}List`
        ].decode(result.raw).items;

        const project_list_filter = project_list.filter((val) => {
          return (
            val.spec.repositoryMeta.projectName ===
            localStorage.getItem("projectName")
          );
        });
        project_list_filter.sort((itemA, itemB) => {
          return (
            itemA.metadata.creationTimestamp.seconds -
            itemB.metadata.creationTimestamp.seconds
          );
        });
        this.repo_list = [];
        for (let pl of project_list_filter) {
          this.repo_list.push({
            name: pl.spec.repositoryMeta.repositoryName,
            project_name: pl.spec.repositoryMeta.projectName,
            create_time: pl.metadata.creationTimestamp.seconds,
          });
        }
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
</style>
