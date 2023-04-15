import { NavigateFunction } from "react-router-dom";

export function ThrowError(
  nav: NavigateFunction,
  title: string,
  content?: string
) {
  nav("/error", {state: {title, content}});
}
