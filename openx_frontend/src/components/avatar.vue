<template>
  <div class="mail-style">
    <a href="https://github.com/kzz45/neverdown" style="text-decoration: none" target="_blank"
      >If you like OpenX, give it a star on GitHub! </a
    >
  </div>
  <el-dropdown trigger="click" size="medium" style="padding: 10px">
    <img style="height: 100%; cursor: pointer" src="../assets/k8s.png" />
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item>{{ userName }}</el-dropdown-item>
        <el-dropdown-item @click="logout">{{
          $t("nav.logout")
        }}</el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<style scoped lang="scss">
.pop-text {
  padding-top: 40px;
  cursor: pointer;
  color: #1818196b;
  font-size: 1rem;
  &:hover {
    color: #000;
  }
}
@keyframes text-flicker-out-glow {
  0%,
  100% {
    -webkit-transform: rotate(0deg);
    transform: rotate(0deg);
    -webkit-transform-origin: 50% 100%;
    transform-origin: 50% 100%;
  }
  10%,
  20%,
  40%,
  60%,
  30%,
  50%,
  70% {
    -webkit-transform: rotate(0deg);
    transform: rotate(0deg);
  }
  75% {
    -webkit-transform: rotate(2deg);
    transform: rotate(2deg);
  }
  80% {
    -webkit-transform: rotate(-4deg);
    transform: rotate(-4deg);
  }
  85% {
    -webkit-transform: rotate(4deg);
    transform: rotate(4deg);
  }
  90% {
    -webkit-transform: rotate(-2deg);
    transform: rotate(-2deg);
  }
  95% {
    -webkit-transform: rotate(2deg);
    transform: rotate(2deg);
  }
}

.mail-to {
  background-image: linear-gradient(to right, #00aeef, #8dc63f);
  animation: text-flicker-out-glow 5s cubic-bezier(0.455, 0.03, 0.515, 0.955)
    infinite both;
  -webkit-background-clip: text;
  color: transparent;
}
.mail-style {
  display: flex;
  flex-direction: column;
  justify-content: center;
  font-size: 14px;
  font-weight: 400;
}
</style>

<script setup lang="ts">
import { i18n, setLanguage } from "@/i18n";
import { onMounted, computed, ref } from "vue";
import { useStore } from "@/store";
const st = useStore();
const languageOptions = [
  { lang: "中文", value: "zh-cn" },
  { lang: "English", value: "en" },
];
const langTrans = {
  "zh-cn": "中文",
  en: "English",
};
let languages = ref("zh-cn");
const userName = localStorage.getItem("username") || "-";
// const languages = computed(()=>{return i18n.global.locale})
function translateChange(lang: string) {
  languages.value = lang;
  setLanguage(lang);
}
import { useRouter } from "vue-router";
const router = useRouter();
function logout() {
  localStorage.removeItem("token");
  localStorage.removeItem("nsList");
  localStorage.removeItem("rtList");
  localStorage.removeItem("configMapStorage");
  localStorage.removeItem("hsoCard");
  localStorage.removeItem("mysqlCard");
  localStorage.removeItem("redisCard");
  router.push({ name: "Index" });
}
</script>
