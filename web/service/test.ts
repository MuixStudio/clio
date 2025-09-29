import type { Fetcher } from "swr";

import { post } from "./base";

export const Test: Fetcher<
  { name: string },
  { username: string; password: string }
> = ({ username, password }) => {
  return post<{ name: string }>("login", { query: { username, password } });
};
