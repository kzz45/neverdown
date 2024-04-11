import { createRouter, createWebHashHistory } from "vue-router"
import loginPage from "../views/Login.vue"

const routes = [
  {
    path: '/',
    name: 'Index',
    component: loginPage,
    meta: { transition: 'slide-left' }
  },
  {
    path: "/Namespace",
    name: "Namespace",
    component: () => import("../views/Namespace.vue"),
    meta: { transition: 'slide-right' }
  },
  {
    path: "/dashboard/:ns",
    name: "Dashboard",
    component: () => import("../views/Dashboard.vue")
  },
  {
    path: "/term",
    name: "Term",
    component: () => import("../views/Term.vue")
  },
  {
    path: "/log",
    name: "Log",
    component: () => import("../views/Log.vue")
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

import { ElNotification } from 'element-plus'
router.beforeEach((to, from, next) => {
  if (from.fullPath === '/' && to.matched.length === 0) {
    const gvkList = JSON.parse(String(localStorage.getItem('gvkList')))
    if (gvkList && gvkList.length <= 0) {
      ElNotification({ title: '权限不足', message: '请联系管理员', type: 'error', duration: 3000 })
      next({ name: 'Namespace' })
    } else {
      const pathGvk = to.path.split('/')
      const ns = pathGvk[1]
      const gvk = pathGvk[2]
      router.addRoute({ name: 'dashboard', path: '/dashboard', component: () => import("@/views/Dashboard.vue") })
      for (let one of gvkList) {
        const gvkRule = `${one.gv}-${one.kind}`
        let componentString = gvkRule.replace(/-/g, '')
        router.addRoute('dashboard', { name: `${ns}${one.gv}${one.kind}`, path: `/${ns}/${gvkRule}`, component: () => import(`./../views/resource/list.vue`) })
      }
      next({ path: `/${ns}/${gvk}` })
    }
  } else {
    next()
  }
})

export default router;