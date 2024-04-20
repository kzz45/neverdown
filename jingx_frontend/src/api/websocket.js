import axios from "axios";
import protoRoot from "@/proto/proto.js";

const protoApi = protoRoot.github.com.kzz45.neverdown.pkg.apis.jingx.v1;
const protoRequest = protoRoot.github.com.kzz45.neverdown.pkg.jingx.proto;

let socket = null;
let reconnectTimer = false;

export const init_socket_data = (ns, kind, verb, raw) => {
  let gvk = kind.split("-");
  const msg = {
    verb: verb,
    namespace: ns,
    groupVersionKind: {
      group: gvk[0],
      version: gvk[1],
      kind: gvk[2],
    },
    raw: raw,
  };
  const request = protoRequest["Request"];
  const message = request.create(msg);
  const senddata = request.encode(message).finish();
  return senddata;
};

export const connectSocket = async (token, store) => {
  if (!token) {
    return;
  }
  if (socket && socket.readyState !== socket.CLOSED) {
    return;
  }
  if (reconnectTimer) {
    return;
  }
  reconnectTimer = true;
  setTimeout(() => {
    reconnectTimer = false;
  }, 1000);
  const cig = await axios.get("config/config.json");
  let env = cig.data;
  const wsUrl = "wss://" + String(env.VUE_APP_BASE_API) + "/api/";
  socket = new WebSocket(wsUrl + token);
  socket.binaryType = "arraybuffer";
  socket.onopen = function () {
    if (socket.readyState === 1) {
      store.dispatch("socket/socket_onopen");
    }
  };
  socket.onerror = function (err) {
    console.log("socket error: ", err);
  };
  socket.onclose = function () {
    store.dispatch("socket/socket_onclose");
    let realToken = localStorage.getItem("neverdown_jingx_token") || "";
    connectSocket(realToken, store);
  };
  socket.onmessage = (msg) => {
    const buffer = new Uint8Array(msg.data);
    store.commit("socket/SOCKET_ONMESSAGE", buffer);
  };
};

export const sendSocketMessage = async (msg, store) => {
  if (socket === null) {
    let realToken = localStorage.getItem("neverdown_jingx_token") || "";
    await connectSocket(realToken, store);
  }
  if (socket && socket.readyState === socket.OPEN) {
    // 若是ws开启状态
    socket.send(msg);
    return;
  }
  setTimeout(function () {
    sendSocketMessage(msg, store);
  }, 1000);
};

export const getProtoParam = (param, gvk) => {
  const message = protoApi[gvk.kind].create(param);
  const sendParam = protoApi[gvk.kind].encode(message).finish();
  return sendParam;
};

export const binaryToStr = (fileData) => {
  const dataString = "";
  for (let i = 0; i < fileData.length; i++) {
    dataString += String.fromCharCode(fileData[i]);
  }
  return dataString;
};
