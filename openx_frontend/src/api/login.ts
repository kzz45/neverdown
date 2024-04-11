import request from '@/api/request'
import proto from '../../proto/proto'
const requestproto: any = proto.github.com.kzz45.neverdown.pkg.openx.aggregator.proto
const rbacproto: any = proto.github.com.kzz45.neverdown.pkg.apis.rbac.v1

interface loginData {
  account: String,
  password: String
}
export async function loginApi(api: string, url: string, gmtParam: loginData = { account: '', password: '' }, decodeApi?: string, encodeApi?: string) {
  let gmtMsg = rbacproto[`${api}`].create(gmtParam)
  let data = rbacproto[`${api}`].encode(gmtMsg).finish()
  if (encodeApi) {
    gmtMsg = rbacproto[`${encodeApi}`].create(gmtParam)
    data = rbacproto[`${encodeApi}`].encode(gmtMsg).finish()
  }
  let gmtData = {
    raw: data
  }
  let msg = requestproto.Request.create(gmtData)
  const v = await request()
  return v({
    url: '/authority',
    method: 'post',
    data: new Blob([requestproto.Request.encode(msg).finish()]),
    responseType: 'blob',
    params: {
      serviceApi: decodeApi || api
    }
  })
}

export async function getGroupingfilter() {
  const v = await request()
  return v({
    url: `/rbac/casbin/groupingpolicy/filter `,
    method: 'get'
  })
}
