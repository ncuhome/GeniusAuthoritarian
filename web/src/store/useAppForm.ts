import { create, StoreApi, UseBoundStore } from "zustand";

interface AppFormState {
  name: string;
  callback: string;
  permitAll: boolean;
  permitGroups?: User.Group[];

  nameError: boolean;
  callbackError: boolean;
}

interface AppFormActions {
  reset: () => void;

  setApp: (app: App.Detailed) => void;
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
    setApp: (app: App.Detailed) => {
      set({
        name: app.name,
        callback: app.callback,
        permitAll: app.permitAllGroup,
        permitGroups: app.groups,
      });
    },
  }));

export const useAppForm = createAppForm();
export const useAppModifyForm = createAppForm();
