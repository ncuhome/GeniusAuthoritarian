import { NavigateFunction } from "react-router-dom";
import { ErrorState } from "@/typings/error";

export function ThrowError(
  nav: NavigateFunction,
  title: string,
  content?: string,
  retryAppCode?: string
) {
  nav("/error", { state: { title, content, retryAppCode } as ErrorState });
}
