import { rbacApi } from "@/utils/request";

export function login(AccountMeta) {
  return rbacApi("AccountMeta", "/login", AccountMeta, "Context");
}

export function getInfo(token) {
  return request({
    url: "/vue-admin-template/user/info",
    method: "get",
    params: { token },
  });
}

export function logout() {
  return request({
    url: "/logout",
    method: "post",
  });
}
