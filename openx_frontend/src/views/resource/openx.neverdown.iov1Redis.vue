<template>
  <el-menu
    :default-active="activeIndex"
    mode="horizontal"
    @select="handleSelect"
  >
    <el-menu-item index="1">MetaData</el-menu-item>
    <el-menu-item index="2">Spec</el-menu-item>
  </el-menu>

  <div class="menu-item" v-show="activeIndex === '1'">
    <MetaData :metadata="itemInfo.metadata"></MetaData>
  </div>

  <div v-show="activeIndex === '2'">
    <div class="switch-btns">
      <div
        :class="mysqlBtn === 'master' ? 'mysql-btn btn-selected' : 'mysql-btn'"
        @click="changeTo('master')"
      >
        master
      </div>
      <div
        :class="mysqlBtn === 'slave' ? 'mysql-btn btn-selected' : 'mysql-btn'"
        @click="changeTo('slave')"
      >
        slave
      </div>
    </div>
    <MysqlConfig
      v-if="mysqlBtn === 'master'"
      :mysqlspec="itemInfo.spec.master"
    />
    <MysqlConfig v-if="mysqlBtn === 'slave'" :mysqlspec="itemInfo.spec.slave" />
  </div>
</template>

<style lang="scss" scoped>
@import "./css/spec.scss";
.switch-btns {
  display: flex;
  padding: 10px;
  .mysql-btn {
    width: 100px;
    height: 20px;
    padding: 5px;
    font-size: 1rem;
    margin: 0px 20px;
    border-radius: 4px;
    text-align: center;
    align-self: stretch;
    cursor: pointer;
    line-height: 20px;
    justify-content: space-between;
    border: 1px solid rgb(88, 101, 242);
  }
  .btn-selected {
    background-color: rgb(88, 101, 242);
    color: white;
  }
}
</style>

<script setup lang="ts">
import MetaData from "./components/metadata.vue";
import MysqlConfig from "./components/mysqlconfig.vue";
import { ref } from "vue";

const props = defineProps<{
  itemInfo?: any;
}>();

let mysqlBtn = ref("master");
function changeTo(gole: string) {
  mysqlBtn.value = gole;
}

const activeIndex = ref("1");
const handleSelect = (key: string, keyPath: string[]) => {
  activeIndex.value = keyPath[0] || "1";
};
</script>
