import { AxiosError } from "axios";

declare global {
  interface ApiError<T> extends AxiosError {
    msg: string;
    response?: {
      data?: T;
    };
  }

  interface ApiResponse<T> {
    code: number;
    data: T;
  }
}
