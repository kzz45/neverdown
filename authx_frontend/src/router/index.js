import Vue from "vue";
import Router from "vue-router";

Vue.use(Router);

import Layout from "@/layout";

export const constantRoutes = [
  {
    path: "/login",
    component: () => import("@/views/login/index"),
    hidden: true,
  },

  {
    path: "/404",
    component: () => import("@/views/404"),
    hidden: true,
  },

  {
    path: "/",
    component: Layout,
    redirect: "/dashboard",
    children: [
      {
        path: "dashboard",
        name: "Dashboard",
        component: () => import("@/views/authx/admin/apps"),
        meta: { title: "应用管理", icon: "home" },
      },
    ],
  },
  {
    path: "/account",
    component: Layout,
    children: [
      {
        path: "account",
        name: "Account",
        component: () => import("@/views/authx/admin/account"),
        meta: { title: "账户管理", icon: "auth" },
      },
    ],
  },
  {
    hidden: true,
    path: "/authx",
    component: Layout,
    redirect: "/authx/rbac",
    children: [
      {
        path: "rbac",
        name: "RBAC",
        component: () => import("@/views/authx/rbac"),
        meta: { title: "认证授权", icon: "auth" },
      },
    ],
  },
];

export const asyncRouter = [

  // 404 page must be placed at the end !!!
  { path: "*", redirect: "/404", hidden: true },
];

const createRouter = () =>
  new Router({
    // mode: "history", // require service support
    scrollBehavior: () => ({ y: 0 }),
    routes: constantRoutes.concat(asyncRouter),
    // routes: constantRoutes,
  });

const router = createRouter();

export function resetRouter() {
  const newRouter = createRouter();
  router.matcher = newRouter.matcher; // reset router
}

export default router;
