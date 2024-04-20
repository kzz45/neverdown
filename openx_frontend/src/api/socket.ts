import proto from '@p/proto'
import router from "@/router";
import axios from 'axios'

const protoApi = proto.k8s.io.api
const protoOpenxApi = proto.github.com.kzz45.neverdown.pkg.apis
const protoNeverdown = proto.github.com.kzz45.neverdown.pkg

let socket: any = null;

let reconnectTimer = false

import { ElNotification } from 'element-plus'
export const connectSocket = async (token: string, store: any) => {
  if (!token) {
    return
  }
  if (socket && socket.readyState !== socket.CLOSED) {
    return
  }
  if (reconnectTimer) {
    return
  }
  reconnectTimer = true
  setTimeout(() => {
    reconnectTimer = false
  }, 1000)
  const cig = await axios.get('config/config.json')
  let env = cig.data
  const wsUrl = 'wss://' + String(env.VITE_BASE_URL) + '/api/'
  socket = new WebSocket(wsUrl + token);
  socket.binaryType = 'arraybuffer'
  socket.onopen = function () {
    if (socket.readyState === 1) {
      store.dispatch('socket/socket_onopen')
    }
  };
  socket.onerror = function (err: any) {
    if (router.currentRoute.value.fullPath === '/') {
      ElNotification({
        title: `Error`, message:
          '链接失败，请检查证书或网络', type: 'error', duration: 3000
      })
    } else {
      router.replace({ name: 'Index' })
    }
  };
  socket.onclose = function () {
    store.dispatch('socket/socket_onclose')
    
    if (localStorage.getItem('neverdown_openx_token') && router.currentRoute.value.fullPath != '/') {
      connectSocket(String(localStorage.getItem('neverdown_openx_token')), store)
    }
  };
  socket.onmessage = (msg: any) => {
    const buffer = new Uint8Array(msg.data)
    store.commit("socket/SOCKET_ONMESSAGE", buffer)
  };
};

export const sendSocketMessage = async (msg: any, store: any) => {
  if (socket === null) {
    await connectSocket(String(localStorage.getItem('neverdown_openx_token')), store)
  }
  if (socket && socket.readyState === socket.OPEN) {
    // 若是ws开启状态
    socket.send(msg)
    return
  }
  setTimeout(function () {
    sendSocketMessage(msg, store)
  }, 1000)
  // if (socket === null) {
  //   await connectSocket(String(localStorage.getItem('token')), store)
  //   setTimeout(function() {
  //     sendSocketMessage(msg, store)
  //   }, 100)
  // } else if (socket.readyState === socket.OPEN) {
  //   // 若是ws开启状态
  //   socket.send(msg)
  // } else if (socket.readyState === socket.CONNECTING) {
  //   // 若是 正在开启状态，则等待100ms后重新调用
  //   setTimeout(function() {
  //     sendSocketMessage(msg, store)
  //   }, 100)
  // } else if (socket.readyState === socket.CLOSED) {
  //   if(localStorage.getItem('token') && router.currentRoute.value.fullPath != '/'){
  //     await connectSocket(String(localStorage.getItem('token')), store)
  //     setTimeout(function() {
  //       sendSocketMessage(msg, store)
  //     }, 100)
  //   }
  // } else {
  //   // 若未开启 ，则等待100ms后重新调用
  //   setTimeout(function() {
  //     sendSocketMessage(msg, store)
  //   }, 100)
  // }
};

export const closeSocket = () => {
  if (socket) {
    socket.close()
    socket = null
  }
}

export const initSocketData = (goal: string = 'Request', ns: string = '', kind: any = '', verb: string = '', raw: any = null) => {

  let gvk = kind.split('-')
  let findItem: any = {
    group: gvk[0],
    version: gvk[1],
    kind: gvk[2]
  }
  let msg = {
    'verb': verb,
    'namespace': ns,
    'groupVersionKind': { group: findItem.group === 'core' ? '' : findItem.group, version: findItem.version, kind: findItem.kind },
    'raw': raw
  }
  const request = protoNeverdown.openx.aggregator.proto[goal]
  const message = request.create(msg), senddata = request.encode(message).finish()
  return senddata
}

export const getProtoParam = (param: any, gvk: any = {}) => {
  let message: any, sendParam: any
  if (gvk.group === 'openx') {
    message = protoOpenxApi[gvk.group][gvk.version][gvk.kind].create(param)
    sendParam = protoOpenxApi[gvk.group][gvk.version][gvk.kind].encode(message).finish()
  } else {
    message = protoApi[gvk.group][gvk.version][gvk.kind].create(param)
    sendParam = protoApi[gvk.group][gvk.version][gvk.kind].encode(message).finish()
  }
  return sendParam
}

export const binaryToStr = (fileData: any) => {
  let dataString = ''
  for (let i = 0; i < fileData.length; i++) {
    dataString += String.fromCharCode(fileData[i])
  }
  return dataString
}

export const strToBinary = (str: string) => {
  let arr = [];
  for (let i = 0, j = str.length; i < j; ++i) {
    arr.push(str.charCodeAt(i));
  }
  let tmpUint8Array = new Uint8Array(arr);
  return tmpUint8Array
}

