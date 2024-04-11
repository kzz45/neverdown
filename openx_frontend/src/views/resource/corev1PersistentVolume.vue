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

  <div class="menu-item" v-show="activeIndex === '2'">
    <div class="spec-item">
      <div class="spec-label">accessModes</div>
      <div class="spec-value">
        <div class="tag-group">
          <div
            class="label-tag"
            v-for="(mode, index) in itemInfo.spec.accessModes"
            :key="mode"
          >
            {{ mode }}
            <el-popconfirm
              title="确定删除?"
              @confirm="tagClose(itemInfo.spec.accessModes, index)"
            >
              <template #reference>
                <el-icon><Close /></el-icon>
              </template>
            </el-popconfirm>
          </div>
          <TagInput @value-input="handleInputModeConfirm" />
        </div>
      </div>
    </div>
    <div class="spec-item">
      <div class="spec-label">capacity</div>
      <div class="spec-value">
        <el-input
          v-model="itemInfo.spec.capacity.storage.string"
          size="small"
          style="width: 100%"
        />
      </div>
    </div>
    <div class="spec-item">
      <div class="spec-label">claimRef</div>
      <div class="spec-value">
        <el-select v-model="claimName" placeholder="Select" style="width: 100%">
          <el-option
            v-for="claim in claimList"
            :key="claim"
            :label="claim"
            :value="claim"
          />
        </el-select>
      </div>
    </div>
    <div class="spec-item">
      <div class="spec-label">mountOptions</div>
      <div class="spec-value">
        <div class="tag-group">
          <div
            class="label-tag"
            v-for="(mode, index) in itemInfo.spec.mountOptions"
            :key="mode"
          >
            {{ mode }}
            <el-popconfirm
              title="确定删除?"
              @confirm="tagClose(itemInfo.spec.mountOptions, index)"
            >
              <template #reference>
                <el-icon><Close /></el-icon>
              </template>
            </el-popconfirm>
          </div>
          <TagInput @value-input="handleInputMountConfirm" />
        </div>
      </div>
    </div>

    <div class="spec-item">
      <div class="spec-label">persistentVolumeReclaimPolicy</div>
      <div class="spec-value">
        <el-input
          v-model="itemInfo.spec.persistentVolumeReclaimPolicy"
          size="small"
          style="width: 100%"
        />
      </div>
    </div>

    <div class="spec-item">
      <div class="spec-label">storageClassName</div>
      <div class="spec-value">
        <el-input
          v-model="itemInfo.spec.storageClassName"
          size="small"
          style="width: 100%"
        />
      </div>
    </div>

    <div class="spec-item">
      <div class="spec-label">volumeMode</div>
      <div class="spec-value">
        <el-input
          v-model="itemInfo.spec.volumeMode"
          size="small"
          style="width: 100%"
        />
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@import "./css/spec.scss";
</style>

<script setup lang="ts">
import MetaData from "./components/metadata.vue";
import TagInput from "./components/taginput.vue";
import { nextTick, ref } from "vue";

const props = defineProps<{
  itemInfo?: any;
}>();

function tagClose(obj: any, index: number) {
  obj.splice(index, 1);
}
function handleInputModeConfirm(inputStr: string) {
  props.itemInfo.spec.accessModes.push(inputStr);
}

const claimList = ref(["test"]);
const claimName = ref("");

function handleInputMountConfirm(inputStr: string) {
  props.itemInfo.spec.mountOptions.push(inputStr);
}

const activeIndex = ref("1");
const handleSelect = (key: string, keyPath: string[]) => {
  activeIndex.value = keyPath[0] || "1";
};
</script>
