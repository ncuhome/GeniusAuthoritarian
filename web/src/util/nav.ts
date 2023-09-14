import { NavigateFunction } from "react-router-dom";

export function ThrowError(
  nav: NavigateFunction,
  title: string,
  content?: string,
  retryAppCode?: string
) {
  nav("/error", { state: { title, content, retryAppCode } as ErrorState });
}

export function GoLogin(nav: NavigateFunction, appCode: string) {
  nav("/?appCode=" + appCode);
}
