import { login, logout, getInfo } from "@/api/user";
import {
  getToken,
  setToken,
  removeToken,
  getTokenExpire,
  setTokenExpire,
  removeTokenExpire,
} from "@/utils/auth";
import { resetRouter } from "@/router";

const getDefaultState = () => {
  return {
    token: getToken(),
    token_expire: "",
    name: "",
    avatar: "",
  };
};

const state = getDefaultState();

const mutations = {
  RESET_STATE: (state) => {
    Object.assign(state, getDefaultState());
  },
  SET_TOKEN: (state, token) => {
    state.token = token;
  },
  SET_TOKEN_EXPIRE: (state, expireAt) => {
    state.token_expire = expireAt;
  },
  SET_NAME: (state, name) => {
    state.name = name;
  },
  SET_AVATAR: (state, avatar) => {
    state.avatar = avatar;
  },
};

const actions = {
  // user login
  login({ commit }, userInfo) {
    const { username, password } = userInfo;
    return new Promise((resolve, reject) => {
      login({ username: username.trim(), password: password })
        .then((response) => {
          // const { data } = response;
          // console.log("login resp: ===", response.token);
          commit("SET_TOKEN", response.token);
          commit("SET_TOKEN_EXPIRE", response.expireAt);
          setToken(response.token);
          setTokenExpire(response.expireAt);
          resolve();
        })
        .catch((error) => {
          reject(error);
        });
    });
  },

  // get user info
  getInfo({ commit, state }) {
    return new Promise((resolve, reject) => {
      getInfo(state.token)
        .then((response) => {
          const { data } = response;

          if (!data) {
            return reject("Verification failed, please Login again.");
          }

          const { name, avatar } = data;

          commit("SET_NAME", name);
          commit("SET_AVATAR", avatar);
          resolve(data);
        })
        .catch((error) => {
          reject(error);
        });
    });
  },

  // user logout
  logout({ commit, state }) {
    // removeToken(); // must remove  token  first
    // removeTokenExpire();
    // localStorage.clear()
    localStorage.removeItem("Token");
    localStorage.removeItem("Token_Expitation");
    localStorage.removeItem("appId");
    localStorage.removeItem("appName");
    resetRouter();
    commit("RESET_STATE");
    // resolve();
    // return new Promise((resolve, reject) => {
    // logout(state.token)
    //   .then(() => {
    //     removeToken(); // must remove  token  first
    //     resetRouter();
    //     commit("RESET_STATE");
    //     resolve();
    //   })
    //   .catch((error) => {
    //     reject(error);
    //   });
    // });
  },

  // remove token
  resetToken({ commit }) {
    return new Promise((resolve) => {
      removeToken(); // must remove  token  first
      commit("RESET_STATE");
      resolve();
    });
  },
};

export default {
  namespaced: true,
  state,
  mutations,
  actions,
};
