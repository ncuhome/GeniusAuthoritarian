import axios from "axios";
import { BaseUrlV1, apiV1ErrorHandler } from "@api/base";

import { useUser } from "@store";

function GoLogin() {
  window.location.href = "/";
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
apiV1User.interceptors.response.use(undefined, (err: any) => {
  if (err.statusCode && err.statusCode === 401) {
    useUser.getState().setToken(null);
    GoLogin();
    return Promise.reject(err);
  }
  return Promise.reject(apiV1ErrorHandler(err));
});
