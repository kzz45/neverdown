import { createApp } from 'vue'
import App from './App.vue'
import router from './router';
import store from "./store";


import "./assets/css/setting.css"
import "./assets/css/global.css"
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

import { i18n } from './i18n';
import LzmNameFilter from 'lzm-namefilter'
const app = createApp(App);
app.use(router);
app.use(store);
app.use(ElementPlus)
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}

app.use(LzmNameFilter);
app.use(i18n);
app.mount('#app')

export default app;