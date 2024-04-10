import protoRoot from "@/proto/proto.js";
import { sendSocketMessage } from "@/api/websocket";
import store from "@/store";
const protoApi = protoRoot.github.com.kzz45.neverdown.pkg.jingx.proto;

function serverPing() {
  let msg = {
    verb: "ping",
    namespace: "",
    groupVersionKind: { group: "jingx", version: "v1", kind: "" },
  };
  const request = protoApi.Request;
  const message = request.create(msg);
  const senddata = request.encode(message).finish();
  sendSocketMessage(senddata, store);
}

const getDefaultState = () => {
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
      heartBeatTimer: 0,
    },
  };
};

const state = getDefaultState();

const mutations = {
  SOCKET_ONOPEN: (state) => {
    state.socket.isConnected = true;
    // 连接成功时启动定时发送心跳消息，避免被服务器断开连接
    console.log("onopen");
    state.socket.heartBeatTimer = setInterval(() => {
      // console.log("socket heartbeat ping:", new Date());
      serverPing();
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
    if (!localStorage.getItem("neverdown_authx_token")) {
      clearInterval(state.socket.heartBeatTimer);
    }
    console.error(state, event);
  },
  // 收到服务端发送的消息
  SOCKET_ONMESSAGE(state, message) {
    state.socket.message = message;
  },
  // 自动重连
  SOCKET_RECONNECT(state, count) {
    console.info("重连中...", state, count);
  },
  // 重连错误
  SOCKET_RECONNECT_ERROR(state) {
    state.socket.reconnectError = true;
  },
};

const actions = {
  onmessage({ commit }, message) {
    commit("SOCKET_ONMESSAGE", message);
  },
  socket_onopen({ commit }) {
    commit("SOCKET_ONOPEN");
  },
  socket_onclose({ commit }) {
    commit("SOCKET_ONCLOSE");
  },
};

export default {
  namespaced: true,
  state,
  mutations,
  actions,
};
