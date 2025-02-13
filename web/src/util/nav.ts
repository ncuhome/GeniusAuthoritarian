import { NavigateFunction } from "react-router";

export function ThrowError(
  nav: NavigateFunction,
  title: string,
  content?: string,
  retryAppCode?: string,
) {
  nav("/error", {
    state: { title, content, retryAppCode } as ErrorState,
    replace: true,
  });
}

export function GoLogin(
  nav: NavigateFunction,
  appCode: string,
  replace?: boolean,
) {
  nav("/?appCode=" + appCode, { replace });
}
