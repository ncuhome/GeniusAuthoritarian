import { create, StoreApi, UseBoundStore } from "zustand";
import { Group } from "@api/v1/user/group";

interface AppFormState {
  name: string;
  callback: string;
  permitAll: boolean;
  permitGroups?: Group[];

  nameError: boolean;
  callbackError: boolean;
}

interface AppFormActions {
  reset: () => void;
  setState: <T extends keyof AppFormState>(
    key: T
  ) => (value: AppFormState[T]) => void;
}

const initialAppForm: AppFormState = {
  name: "",
  callback: "https://",
  permitAll: false,
  permitGroups: undefined,

  nameError: false,
  callbackError: false,
};

export type UseAppForm = UseBoundStore<StoreApi<AppFormState & AppFormActions>>;

const createAppForm = (): UseAppForm =>
  create<AppFormState & AppFormActions>()((set) => ({
    ...initialAppForm,

    reset: () => {
      set(initialAppForm);
    },
    setState: (key) => (value) => set({ [key]: value }),
  }));

export const useAppForm = createAppForm();
export const useAppModifyForm = createAppForm();
