<template>
  <div class="login-bg">
    <div class="login-form">
      <div
        style="width: 100%; text-align: center; flex: 1 1 auto; display: flex"
      >
        <div class="login-input">
          <div style="font-weight: 700; font-size: 22px; margin-bottom: 20px">
            Kubernetes Dashboard Platform
          </div>
          <el-form
            label-position="top"
            style="width: 100%"
            label-width="80px"
            ref="loginForm"
            :rules="rules"
            :model="loginInfo"
          >
            <el-form-item label="账号" prop="account">
              <template #label>
                <label for="elac" class="login-label">{{
                  $t("login.account")
                }}</label>
              </template>
              <el-input
                id="elac"
                style="font-size: 16px"
                autocomplete="off"
                v-model="loginInfo.account"
                :input-style="{
                  height: '40px',
                  'border-radius': '3px',
                  'font-weight': '400px',
                }"
              ></el-input>
              <template #error="scope">
                <div
                  style="
                    text-align: left;
                    font-size: 14px;
                    color: #ed4245;
                    font-weight: 500;
                  "
                >
                  {{ scope.error }}
                </div>
              </template>
            </el-form-item>
            <el-form-item label="密码" prop="password">
              <template #label>
                <label for="elpw" class="login-label">{{
                  $t("login.password")
                }}</label>
              </template>
              <el-input
                id="elpw"
                style="font-size: 16px"
                autocomplete="off"
                v-model="loginInfo.password"
                show-password
                :input-style="{
                  height: '40px',
                  'border-radius': '3px',
                  'font-weight': '400px',
                }"
              ></el-input>
              <template #error="scope">
                <span style="float: left; font-size: 14px; color: #ed4245">
                  {{ scope.error }}
                </span>
              </template>
            </el-form-item>
          </el-form>
          <div class="loginBtn" v-if="!loginLoading" @click="login">
            {{ $t("login.submit") }}
          </div>
          <div class="loginLoading" v-else>{{ $t("login.submiting") }}</div>
        </div>
        <div class="dart"></div>
        <div class="right-icon">
          <img
            :class="loginLoading ? 'loading-icon' : ''"
            width="160"
            height="160"
            alt="k8s"
            src="../assets/k8s.png"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";
import { watch } from "vue";
import { loginApi } from "@/api/login";
import { useStore } from "@/store";
// import store from '@/store'
import router from "@/router";
import { connectSocket, closeSocket } from "@/api/socket";
import { ElNotification } from "element-plus";

localStorage.clear();

const { t } = useI18n();
let loginInfo = $ref({
  account: "",
  password: "",
});
let rules = $ref({
  account: [
    { required: true, message: t("login.accountTips"), trigger: "blur" },
    { min: 2, max: 50, message: t("login.accountLength"), trigger: "blur" },
  ],
  password: [
    { required: true, message: t("login.passwordTips"), trigger: "blur" },
  ],
});
let loginLoading = $ref(false);
let loginForm = $ref({});
const store = useStore();
closeSocket();
watch(
  () => store.state.socket.socket.isConnected,
  (connect) => {
    if (connect) {
      router.push({ name: "Namespace" });
    }
  }
);
function login() {
  loginLoading = true;
  loginForm.validate((vali: boolean) => {
    if (vali) {
      const AccountMeta: any = {
        username: loginInfo.account,
        password: loginInfo.password,
      };
      loginApi("AccountMeta", "", AccountMeta, "Context")
        .then((res: any) => {
          if (res === 403) {
            loginLoading = false;
            return;
          }
          if (!res) {
            ElNotification({
              title: "登录失败",
              message: "网络异常",
              type: "error",
              duration: 1500,
            });
            loginLoading = false;
            return;
          }
          if (!res.clusterRole) {
            ElNotification({
              title: "登录失败",
              message: "网络异常",
              type: "error",
              duration: 1500,
            });
            loginLoading = false;
            return;
          }
          if (res.clusterRole.spec.rules.length <= 0) {
            ElNotification({
              title: "登录失败",
              message: "权限不足",
              type: "error",
              duration: 1500,
            });
            loginLoading = false;
            return;
          }
          if (res.token) {
            localStorage.setItem("username", AccountMeta.username);
            localStorage.setItem("token", res.token);
            localStorage.setItem(
              "clusterRole",
              JSON.stringify(initRule(res.clusterRole.spec.rules))
            );
            connectSocket(res.token, store);
          }
        })
        .catch((err: any) => {
          if (String(err).startsWith("timeout")) {
            ElNotification({
              title: "登录失败",
              message: "网络超时",
              type: "error",
              duration: 1500,
            });
          }
          loginLoading = false;
        });
    } else {
      loginLoading = false;
    }
  });
}

function initRule(plicyRule: any) {
  let namespace: any = [];
  let routerGroup: any = {};
  for (let plicy of plicyRule) {
    if (plicy.verbs.length > 0) {
      const groupVersion = `${plicy.groupVersionKind.group}/${plicy.groupVersionKind.version}`;
      if (!namespace.includes(plicy.namespace)) {
        namespace.push(plicy.namespace);
        let gvkVerb: any = {};
        gvkVerb[groupVersion] = {
          gvk: [
            {
              kind: plicy.groupVersionKind.kind,
              verbs: plicy.verbs,
            },
          ],
        };
        routerGroup[plicy.namespace] = gvkVerb;
      } else {
        if (routerGroup[plicy.namespace][groupVersion]) {
          routerGroup[plicy.namespace][groupVersion].gvk.push({
            kind: plicy.groupVersionKind.kind,
            verbs: plicy.verbs,
          });
        } else {
          routerGroup[plicy.namespace][groupVersion] = {
            gvk: [
              {
                kind: plicy.groupVersionKind.kind,
                verbs: plicy.verbs,
              },
            ],
          };
        }
      }
    }
  }
  return routerGroup;
}
</script>
<style lang="scss">
.login-input {
  .el-form-item__label {
    line-height: 20px !important;
  }
}
@import "css/login.scss";
</style>
