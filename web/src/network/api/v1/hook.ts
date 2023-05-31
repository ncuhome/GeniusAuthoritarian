import { createFetchHook } from "@/hooks/useFetch";
import { apiV1 } from "@api/base";

export const useApiV1 = createFetchHook(apiV1);
