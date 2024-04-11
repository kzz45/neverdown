<template>
  <div class="meta-grid">
    <div class="meta-item">
      <div class="meta-label">Namespace:</div>
      <div class="meta-value">
        {{ props.metadata.namespace }}
      </div>
    </div>
    <div class="meta-item">
      <div class="meta-label">Name:</div>
      <div class="meta-value">
        {{ props.metadata.name }}
      </div>
    </div>
    <div class="meta-item">
      <div class="meta-label">resourceVersion:</div>
      <div class="meta-value">
        {{ props.metadata.resourceVersion }}
      </div>
    </div>
    <div class="meta-item">
      <div class="meta-label">kind:</div>
      <div class="meta-value">
        {{ props.metadata.kind }}
      </div>
    </div>
    <div class="meta-item">
      <div class="meta-label">apiVersion:</div>
      <div class="meta-value">
        {{ props.metadata.apiVersion }}
      </div>
    </div>
    <div class="meta-item">
      <div class="meta-label">uid:</div>
      <div class="meta-value">
        {{ props.metadata.uid }}
      </div>
    </div>
    <div class="meta-item">
      <div class="meta-label">fieldPath:</div>
      <div class="meta-value">
        {{ props.metadata.fieldPath }}
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
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
  padding-left: 40px;
  .meta-item {
    display: flex;
    margin-bottom: 8px;
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
import { cloneDeep } from "lodash";
import { inject, onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import { formatTime } from "../util";

const route = useRoute();

const props = defineProps<{
  metadata?: any;
}>();
</script>
