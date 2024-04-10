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
        component: () => import("@/views/registry/ns.vue"),
        meta: { title: "镜像服务", icon: "jingx" },
      },
    ],
  },

  {
    hidden: true,
    path: "/repo",
    component: Layout,
    redirect: "/repo",
    children: [
      {
        path: "repo",
        name: "Repo",
        component: () => import("@/views/registry/repo.vue"),
        meta: { title: "镜像仓库", icon: "repo" },
      },
    ],
  },

  {
    hidden: true,
    path: "/tag",
    component: Layout,
    redirect: "/tag",
    children: [
      {
        path: "tag",
        name: "Tag",
        component: () => import("@/views/registry/tag.vue"),
        meta: { title: "镜像版本", icon: "repo" },
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
