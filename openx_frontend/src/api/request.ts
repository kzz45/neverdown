import axios from 'axios'
import proto from '../../proto/proto'
import { ElNotification } from 'element-plus'
import router from "@/router";

const requestproto: any = proto.github.com.kzz45.neverdown.pkg.openx.aggregator.proto
const rbacproto: any = proto.github.com.kzz45.neverdown.pkg.apis.openx.v1

export default async () => {
  const cig = await axios.get('config/config.json')
  let env = cig.data
  let config: any = {
    baseURL: 'https://' + String(env.VITE_BASE_URL),
    timeout: 5000
  }
  let service = axios.create(config)
  service.interceptors.request.use(
    config => {
      if (localStorage.getItem('neverdown_openx_token')) {
        config.headers['Token'] = localStorage.getItem('neverdown_openx_token')
      }
      return config
    },
    error => {
      return Promise.reject(error)
    }
  )
  service.interceptors.response.use(
    async response => {
      const ab = await response.data.arrayBuffer()
      const buffer = new Uint8Array(ab)
      let resp = requestproto.Response.decode(buffer)
      if (resp.code === 1 || resp.code === 401) {
        router.replace({ path: `/login` })
        ElNotification({
          title: 'Error',
          message: resp.message,
          type: 'error',
        })
      }
      if (resp.code === 403) {
        ElNotification({
          title: 'Error',
          message: '账号或密码错误',
          type: 'error',
        })
        return 403
      }
      let gmtData = null
      if (resp.raw) {
        gmtData = requestproto[`${response.config.params.serviceApi}`].decode(resp.raw)
      } else {
        gmtData = resp
      }
      return gmtData
    },
    error => {
      return Promise.reject(error)
    }
  )
  return service
};