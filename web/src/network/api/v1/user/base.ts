import axios, { AxiosError } from "axios";
import { BaseUrlV1, apiV1ErrorHandler } from "@api/base";

import { useUser } from "@store";

function GoLogin() {
  window.location.href = "/";
  console.log(4);
}

export function Logout() {
  console.log(1);
  useUser.getState().setAuth(null);
  console.log(2);
  GoLogin();
  console.log(3);
}

export const apiV1User = axios.create({
  baseURL: BaseUrlV1 + "user/",
});
apiV1User.interceptors.request.use((req) => {
  const token = useUser.getState().token;
  if (token) {
    req.headers["Authorization"] = token;
  } else {
    GoLogin();
    const controller = new AbortController();
    controller.abort();
    return {
      ...req,
      signal: controller.signal,
    };
  }
  return req;
}, undefined);
apiV1User.interceptors.response.use(undefined, (err: AxiosError) => {
  console.log(err);
  if (err.response?.status === 401) {
    Logout();
    return Promise.reject(err);
  }
  return Promise.reject(apiV1ErrorHandler(err));
});
