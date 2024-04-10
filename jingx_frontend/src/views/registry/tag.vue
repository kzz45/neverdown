<template>
  <div class="app-container">
    <el-card class="box-card">
      <el-button size="small" icon="el-icon-back" @click="goback"
        >返回</el-button
      >
      <el-table :data="page_tag_list" size="small" empty-text="啥也没有" border>
        <el-table-column label="名称" prop="name"></el-table-column>
        <el-table-column label="Tag" prop="tag"></el-table-column>
        <el-table-column
          label="Author"
          prop="dockerImage.author"
        ></el-table-column>
        <el-table-column label="Sha256" prop="dockerImage.sha256">
        </el-table-column>
        <el-table-column label="branch" prop="gitRef.branch"></el-table-column>
        <el-table-column label="更新时间" prop="">
          <template slot-scope="scoped">
            {{
              scoped.row.dockerImage.lastModifiedTime
                | parseTime("{y}-{m}-{d} {h}:{i}:{s}")
            }}
          </template>
        </el-table-column>
        <el-table-column label="拉取命令" width="120px;">
          <template slot-scope="scoped">
            <el-tag @click="copy_public(scoped.row)">外网</el-tag>
            <el-tag @click="copy_private(scoped.row)">内网</el-tag>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        background
        :page-size="10"
        :current-page.sync="currentPage"
        :total="this.tag_list.length"
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
  name: "RegistryTag",
  filters: {
    parseTime(time, cFormat) {
      return parseTime(time, cFormat);
    },
  },
  data() {
    return {
      tag_list: [],
      currentPage: 1,
      projectName: localStorage.getItem("projectName"),
      repoName: localStorage.getItem("repoName"),
    };
  },
  watch: {
    message: function () {
      this.socket_onmessage(this.message);
    },
  },
  computed: {
    ...mapGetters(["message"]),
    page_tag_list: function () {
      return this.tag_list.slice(
        (this.currentPage - 1) * 10,
        this.currentPage * 10
      );
    },
  },
  mounted() {
    this.socket_connect();
  },
  created() {},
  methods: {
    goback() {
      this.$router.push({ path: "/repo/repo" });
    },
    socket_connect() {
      const send_data = init_socket_data(
        "discovery-jingx",
        "jingx-v1-Tag",
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
        // console.log(project_list);
        const project_list_filter1 = project_list.filter((val) => {
          return (
            val.spec.repositoryMeta.projectName ===
            localStorage.getItem("projectName")
          );
        });
        const project_list_filter2 = project_list_filter1.filter((val) => {
          return (
            val.spec.repositoryMeta.repositoryName ===
            localStorage.getItem("repoName")
          );
        });
        project_list_filter2.sort((itemA, itemB) => {
          return (
            itemA.metadata.creationTimestamp.seconds -
            itemB.metadata.creationTimestamp.seconds
          );
        });
        // console.log(project_list_filter2);
        this.tag_list = [];
        for (let pl of project_list_filter2) {
          this.tag_list.push({
            name: pl.metadata.name,
            tag: pl.spec.tag,
            dockerImage: pl.spec.dockerImage,
            gitRef: pl.spec.gitReference,
            create_time: pl.metadata.creationTimestamp.seconds,
          });
        }
      }
    },
    copy_public(row) {
      // console.log(row);
      const public_addr =
        this.projectName +
        "/" +
        this.repoName +
        "@sha256:" +
        row.dockerImage.sha256;
      navigator.clipboard.writeText(public_addr).then(() => {
        this.$message({
          type: "success",
          message: "复制成功",
        });
      });
    },
    copy_private(row) {
      const private_addr =
        this.projectName +
        "/" +
        this.repoName +
        "@sha256:" +
        row.dockerImage.sha256;
      navigator.clipboard.writeText(private_addr).then(() => {
        this.$message({
          type: "success",
          message: "复制成功",
        });
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
