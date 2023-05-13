import { create } from "zustand";

interface AppFormState {
  name: string;
  callback: string;
  permitAll: boolean;

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
  callback: "",
  permitAll: false,

  nameError: false,
  callbackError: false,
};

export const useAppForm = create<AppFormState & AppFormActions>()((set) => ({
  ...initialAppForm,

  reset: () => {
    set(initialAppForm);
  },
  setState: (key) => (value) => set({ [key]: value }),
}));
