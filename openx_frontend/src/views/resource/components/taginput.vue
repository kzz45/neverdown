<template>
  <el-input
    class="input-new-tag"
    v-if="inputVisible"
    v-model="inputValue"
    ref="taginputref"
    @keyup.enter="$event.target.blur()"
    @blur="handleInputConfirm"
  >
  </el-input>
  <div v-else class="add-label-btn" @click="showInput">添加</div>
</template>

<style lang="scss" scoped>
@import "./../css/spec.scss";
</style>

<script lang="ts" setup>
import { nextTick, ref } from "vue";
const props = defineProps<{
  blank?: { type: boolean; defalut: false; required: false };
}>();

let inputVisible = ref(false);
let inputValue = ref("");
let taginputref = ref({
  $refs: {
    input: {
      focus,
    },
  },
});

const emit = defineEmits(["value-input"]);

async function showInput() {
  inputVisible.value = true;
  await nextTick();
  taginputref.value.$refs.input.focus();
}
function handleInputConfirm() {
  console.log("blank", props.blank);
  if (inputValue.value) {
    emit("value-input", inputValue.value);
  } else if (props.blank) {
    emit("value-input", inputValue.value);
  }
  inputVisible.value = false;
  inputValue.value = "";
}
</script>
