import axios from "axios";
import { BaseUrlV1, apiV1ErrorHandler } from "@api/base";

function GoLogin() {
  window.location.href = "/";
}

const apiV1User = axios.create({
  baseURL: BaseUrlV1 + "user/",
});
apiV1User.interceptors.request.use((req) => {
  const token = localStorage.getItem("token");
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
    localStorage.removeItem("token");
    GoLogin();
    return Promise.reject(err);
  }
  return Promise.reject(apiV1ErrorHandler(err));
});
