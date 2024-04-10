const getters = {
  sidebar: (state) => state.app.sidebar,
  device: (state) => state.app.device,
  token: (state) => state.user.token,
  token_expire: (state) => state.user.token_expire,
  avatar: (state) => state.user.avatar,
  name: (state) => state.user.name,
};
export default getters;
