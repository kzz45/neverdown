import axios from "axios";
import store from "@/store";
import { getToken } from "@/utils/auth";
import protoRoot from "@/proto/proto.js";
import { Message } from "element-ui";

const requestproto = protoRoot.github.com.kzz45.neverdown.pkg.jingx.proto;
const rbacproto = protoRoot.github.com.kzz45.neverdown.pkg.apis.rbac.v1;

async function getService() {
  const cig = await axios.get("config/config.json");
  const env = cig.data;
  let config = {
    baseURL: "https://" + String(env.VUE_APP_BASE_API),
    timeout: 5000,
  };
  // console.log("config", config);
  const service = axios.create(config);

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

      if (resp.code === 401) {
        // console.log("code: 401", resp);
        Message({
          message: resp.message,
          type: "error",
          duration: 5 * 1000,
        });
      }
      if (resp.code === 403) {
        // console.log("code: 403", resp);
        Message({
          message: "账户或者密码错误",
          type: "error",
          duration: 5 * 1000,
        });
      }
      // console.log("request resp: ", resp);
      let gmtData = null;
      if (response.config.params.serviceApi === "Context") {
        // gmtData = requestproto[`${response.config.params.serviceApi}`].decode(
        //   resp.data
        // );
        gmtData = requestproto.Context.decode(resp.raw);
      } else {
        gmtData = rbacproto[`${response.config.params.serviceApi}`].decode(
          resp.data
        );
      }
      return gmtData;
    },
    (error) => {
      // console.log("err" + error);
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

// export async function rbacApi(api, url, gmtParam, decodeApi, encodeApi) {
//   let gmtMsg = rbacproto[`${api}`].create(gmtParam);
//   let data = rbacproto[`${api}`].encode(gmtMsg).finish();
//   if (encodeApi) {
//     gmtMsg = rbacproto[`${encodeApi}`].create(gmtParam);
//     data = rbacproto[`${encodeApi}`].encode(gmtMsg).finish();
//   }

//   let gmtData = {
//     serviceRoute: url,
//     data,
//   };

//   let msg = requestproto.Request.create(gmtData);
//   const service = await getService();
//   return service({
//     url: "authority",
//     method: "post",
//     data: new Blob([requestproto.Request.encode(msg).finish()]),
//     responseType: "blob",
//     params: {
//       serviceApi: decodeApi || api,
//     },
//   });
// }

export async function rbacApi(api, url, gmtParam, decodeApi, encodeApi) {
  const gmtMsg = rbacproto[`${api}`].create(gmtParam);
  // console.log("gmtMsg", gmtMsg);
  const data = rbacproto[`${api}`].encode(gmtMsg).finish();
  if (encodeApi) {
    gmtMsg = rbacproto[`${encodeApi}`].create(gmtParam);
    data = rbacproto[`${encodeApi}`].encode(gmtMsg).finish();
  }
  const gmtData = {
    raw: data,
  };
  const msg = requestproto.Request.create(gmtData);
  const service = await getService();
  return service({
    url: "/authz",
    method: "post",
    data: new Blob([requestproto.Request.encode(msg).finish()]),
    responseType: "blob",
    params: {
      serviceApi: decodeApi || api,
    },
  });
}
