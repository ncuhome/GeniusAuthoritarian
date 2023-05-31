import { createFetchHook } from "@/hooks/useFetch";
import { apiV1User } from "./base";

export const useUserApiV1 = createFetchHook(apiV1User);
