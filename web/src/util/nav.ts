import { NavigateFunction } from "react-router-dom";

export function ThrowError(
  nav: NavigateFunction,
  title: string,
  content?: string
) {
  let target = `/error?title=${encodeURI(title)}`;
  if (content) target += "&content=" + encodeURI(content);
  nav(target);
}
