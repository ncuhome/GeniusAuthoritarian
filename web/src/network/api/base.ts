import axios from "axios";

export const BaseURL = `/api/`;
export const BaseUrlV1 = `${BaseURL}v1/`;

export const ErrNetwork = "网络异常";

const apiV1 = axios.create({
  baseURL: BaseUrlV1,
});

export function apiV1ErrorHandler(err: any): any {
  switch (true) {
    case err.name === "CanceledError":
      break;
    case !err || !err.response || !err.response.data:
      err.msg = ErrNetwork;
      break;
    default:
      err.msg = err.response.data.msg;
  }
  return err;
}
apiV1.interceptors.response.use(undefined, (err: any) => {
  return Promise.reject(apiV1ErrorHandler(err));
});

export { apiV1 };
