<template>
  <ServiceMeta :service="itemInfo"></ServiceMeta>
</template>

<style lang="scss" scoped>
.meta-title {
  font-size: 16px;
  font-weight: bold;
  color: #303133;
  text-align: left;
}
.meta-grid {
  display: grid;
  margin-top: 10px;
  grid-template-columns: 50% 50%;
  .meta-item {
    display: flex;
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
}
</style>

<script setup lang="ts">
import ServiceMeta from "./components/servicemeta.vue";

const props = defineProps<{
  itemInfo?: any;
}>();
</script>
