import { rbacApi } from "@/utils/request";

export function RbacAccountList() {
  return rbacApi("RbacServiceAccountList", "/rbacaccount/list");
}

export function RbacAccountCreate(data) {
  rbacApi("RbacServiceAccount", "/rbacaccount/create", data);
}

export function RbacAccountUpdate(data) {
  rbacApi("RbacServiceAccount", "/rbacaccount/update", data);
}

export function RbacAccountDelete(data) {
  rbacApi("RbacServiceAccount", "/rbacaccount/delete", data);
}

// --------------------------------------------------

export function RbacAppList() {
  return rbacApi("AppList", "/apps/list");
}

export function RbacAppCreate(data) {
  return rbacApi("App", "/apps/create", data);
}
export function RbacAppUpdate(data) {
  return rbacApi("App", "/apps/update", data);
}
export function RbacAppDelete(data) {
  return rbacApi("App", "/apps/delete", data);
}

// --------------------------------------------------

export function AppAccountList(data) {
  return rbacApi(
    "AppServiceAccountList",
    "/appaccount/list",
    data,
    "",
    "AppServiceAccount"
  );
}

export function AppAccountCreate(data) {
  return rbacApi("AppServiceAccount", "/appaccount/create", data);
}

export function AppAccountUpdate(data) {
  return rbacApi("AppServiceAccount", "/appaccount/update", data);
}

export function AppAccountDelete(data) {
  return rbacApi("AppServiceAccount", "/appaccount/delete", data);
}

// --------------------------------------------------

export function AppRoleList(data) {
  return rbacApi(
    "ClusterRoleList",
    "/clusterrole/list",
    data,
    "",
    "ClusterRole"
  );
}

// --------------------------------------------------

export function AppGVKList(data) {
  return rbacApi(
    "GroupVersionKindRuleList",
    "/kind/list",
    data,
    "",
    "GroupVersionKindRule"
  );
}
