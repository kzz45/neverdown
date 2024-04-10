import axios from "axios";
import { MessageBox, Message } from "element-ui";
import store from "@/store";
import { getToken } from "@/utils/auth";
import protoRoot from "@/proto/proto.js";
import router from "@/router";

const requestproto = protoRoot.github.com.kzz45.neverdown.pkg.authx.http.proto;
const rbacproto = protoRoot.github.com.kzz45.neverdown.pkg.apis.rbac.v1;

async function getService() {
  const service = axios.create({
    // baseURL: "https://127.0.0.1:8087",
    baseURL: process.env.VUE_APP_BASE_API,
    timeout: 8000,
  });

  service.interceptors.request.use(
    (config) => {
      if (store.getters.token) {
        config.headers["Token"] = getToken();
      }
      return config;
    },
    (error) => {
      console.log(error);
      return Promise.reject(error);
    }
  );

  service.interceptors.response.use(
    async (response) => {
      const ab = await response.data.arrayBuffer();
      const buffer = new Uint8Array(ab);
      let resp = requestproto.Response.decode(buffer);
      // console.log(resp.code, "======");
      if (resp.code === 401) {
        // router.push({ path: `/login` });
        // console.log("401", resp);
        Message({
          message: resp.message,
          type: "error",
          duration: 5 * 1000,
        });
      }
      if (resp.code === 1) {
        // console.log("1", resp);
        Message({
          message: resp.message,
          type: "error",
          duration: 5 * 1000,
        });
      }

      let gmtData = null;
      if (response.config.params.serviceApi === "Context") {
        gmtData = requestproto[`${response.config.params.serviceApi}`].decode(
          resp.data
        );
      } else {
        gmtData = rbacproto[`${response.config.params.serviceApi}`].decode(
          resp.data
        );
      }
      return gmtData;
    },
    (error) => {
      console.log("err" + error);
      Message({
        message: error.message,
        type: "error",
        duration: 5 * 1000,
      });
      return Promise.reject(error);
    }
  );
  return service;
}

export async function rbacApi(api, url, gmtParam, decodeApi, encodeApi) {
  let gmtMsg = rbacproto[`${api}`].create(gmtParam);
  let data = rbacproto[`${api}`].encode(gmtMsg).finish();
  if (encodeApi) {
    gmtMsg = rbacproto[`${encodeApi}`].create(gmtParam);
    data = rbacproto[`${encodeApi}`].encode(gmtMsg).finish();
  }

  let gmtData = {
    serviceRoute: url,
    data,
  };

  let msg = requestproto.Request.create(gmtData);
  // const service = await getService();
  const service = await getService();
  return service({
    url: "authz",
    method: "post",
    data: new Blob([requestproto.Request.encode(msg).finish()]),
    responseType: "blob",
    params: {
      serviceApi: decodeApi || api,
    },
  });
}
