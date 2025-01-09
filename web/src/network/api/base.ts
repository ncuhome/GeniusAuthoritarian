import axios from "axios";

export const BaseURL = `/api/`;
export const BaseUrlV1 = `${BaseURL}v1/`;

export const ErrNetwork = "网络异常";

const apiV1 = axios.create({
  baseURL: BaseUrlV1,
});

export function apiV1ErrorHandler(err: ApiError<any>): any {
  switch (true) {
    case !err || !err.response || !err.response.data:
      err.message = ErrNetwork;
      break;
    default:
      err.message = err.response?.data?.msg;
  }
  return err;
}

apiV1.interceptors.response.use(undefined, (err: ApiError<any>) => {
  if (err.name === "CanceledError") {
    return new Promise(() => {})
  }
  return Promise.reject(apiV1ErrorHandler(err));
});

export { apiV1 };
