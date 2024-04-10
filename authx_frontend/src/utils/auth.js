import Cookies from "js-cookie";

const TokenKeyName = "Token";
const TokenKeyExpire = "Token_Expitation";

export function getToken() {
  // return Cookies.get(TokenKeyName);
  return localStorage.getItem(TokenKeyName);
}

export function getTokenExpire() {
  // return Cookies.get(TokenKeyExpire);
  return localStorage.getItem(TokenKeyExpire);
}

export function setToken(token) {
  // return Cookies.set(TokenKeyName, token);
  return localStorage.setItem(TokenKeyName, token);
}

export function setTokenExpire(token_expire) {
  // return Cookies.set(TokenKeyExpire, token_expire);
  return localStorage.setItem(TokenKeyExpire, token_expire);
}

export function removeToken() {
  // return Cookies.remove(TokenKeyName);
  return localStorage.removeItem(TokenKeyName);
}

export function removeTokenExpire() {
  // return Cookies.remove(TokenKeyExpire);
  return localStorage.removeItem(TokenKeyExpire);
}

export function getTimestamp() {
  const today = new Date();
  const date =
    today.getFullYear() + "-" + (today.getMonth() + 1) + "-" + today.getDate();
  const time =
    today.getHours() + ":" + today.getMinutes() + ":" + today.getSeconds();
  const dateTime = date + " " + time;
  return dateTime;
}
