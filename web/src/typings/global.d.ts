import {AxiosError, AxiosResponse} from "axios";

declare global {
  interface ApiError<T> extends AxiosError {
    msg: string;
    response?: {
      data?: T;
    } & AxiosResponse;
  }

  interface ApiResponse<T> {
    code: number;
    data: T;
  }
}
