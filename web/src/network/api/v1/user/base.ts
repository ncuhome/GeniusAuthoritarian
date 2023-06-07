import axios, { AxiosError } from "axios";
import { BaseUrlV1, apiV1ErrorHandler } from "@api/base";

import useUser from "@store/useUser";

function GoLogin() {
  window.location.href = "/";
}

export function Logout() {
  useUser.getState().setAuth(null);
  GoLogin();
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
