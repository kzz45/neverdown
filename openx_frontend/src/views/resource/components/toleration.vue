<template>
  <div v-if="toleration">
    <div class="spec-item">
      <div class="spec-label">effect</div>
      <div class="spec-value">
        <el-select
          v-model="toleration.effect"
          placeholder="Select"
          style="width: 100%"
        >
          <el-option
            v-for="val in effectOptions"
            :key="val"
            :label="val"
            :value="val"
          />
        </el-select>
      </div>
    </div>

    <div class="spec-item">
      <div class="spec-label">key</div>
      <div class="spec-value">
        <el-input v-model="toleration.key" size="small" style="width: 100%" />
      </div>
    </div>

    <div class="spec-item">
      <div class="spec-label">operator</div>
      <div class="spec-value">
        <el-select
          v-model="toleration.operator"
          placeholder="Select"
          style="width: 100%"
        >
          <el-option
            v-for="val in operatorOptions"
            :key="val"
            :label="val"
            :value="val"
          />
        </el-select>
      </div>
    </div>

    <div class="spec-item">
      <div class="spec-label">value</div>
      <div class="spec-value">
        <el-input v-model="toleration.value" size="small" style="width: 100%" />
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@import "./../css/spec.scss";
</style>

<script setup lang="ts">
import { effect, ref } from "vue";

const props = defineProps<{
  toleration?: any;
}>();

const effectOptions = ["NoSchedule", "PreferNoSchedule", "NoExecute"];
const operatorOptions = ["Exists", "Equal"];
</script>
