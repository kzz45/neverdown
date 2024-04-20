import { Module } from "vuex";
import { sendSocketMessage } from "@/api/socket"
import proto from '@p/proto'
const protoApi = proto.github.com.kzz45.neverdown.pkg.openx.aggregator.proto

function serverPing() {
  let msg = {
    'verb': 'ping',
    'namespace': '',
    'groupVersionKind': { group: "apps", version: "v1", kind: '' }
  }
  const request = protoApi.Request,
    message = request.create(msg), senddata = request.encode(message).finish()
  // console.log("PING: " + new Date());
  sendSocketMessage(senddata, store)
}

interface Socket {
  maker: any,
  isConnected: boolean,
  message: string,
  reconnectError: boolean,
  heartBeatInterval: number,
  heartBeatTimer: any
}

interface StoreSock {
  socket: Socket
}

const store: Module<StoreSock, unknown> = {
  namespaced: true,
  state() {
    return {
      socket: {
        maker: null,
        // 连接状态
        isConnected: false,
        // 消息内容
        message: "",
        // 重新连接错误
        reconnectError: false,
        // 心跳消息发送时间
        heartBeatInterval: 5000,
        // 心跳定时器
        heartBeatTimer: 0
      }
    }
  },
  mutations: {
    // 连接打开
    SOCKET_ONOPEN(state) {
      state.socket.isConnected = true;
      // 连接成功时启动定时发送心跳消息，避免被服务器断开连接
      console.log("连接已建立: " + new Date());
      state.socket.heartBeatTimer = setInterval(() => {
        serverPing()
      }, state.socket.heartBeatInterval);
    },
    // 连接关闭
    SOCKET_ONCLOSE(state, event) {
      state.socket.isConnected = false;
      // 连接关闭时停掉心跳消息
      clearInterval(state.socket.heartBeatTimer);
      state.socket.heartBeatTimer = 0;
      console.log("连接已断开: " + new Date());
    },
    // 发生错误
    SOCKET_ONERROR(state, event) {
      if (!localStorage.getItem('neverdown_openx_token')) {
        clearInterval(state.socket.heartBeatTimer)
      }
      console.error(state, event);
    },
    // 收到服务端发送的消息
    SOCKET_ONMESSAGE(state, message) {
      state.socket.message = message;
    },
    // 自动重连
    SOCKET_RECONNECT(state, count) {
      console.info("消息系统重连中...", state, count);
    },
    // 重连错误
    SOCKET_RECONNECT_ERROR(state) {
      state.socket.reconnectError = true;
    }
  },
  actions: {
    onmessage(context, message) {
      context.commit('SOCKET_ONMESSAGE', message)
    },
    socket_onopen(context) {
      context.commit('SOCKET_ONOPEN')
    },
    socket_onclose(context) {
      context.commit('SOCKET_ONCLOSE')
    }
  }
}

export default store