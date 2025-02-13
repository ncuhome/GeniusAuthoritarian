import { create, StoreApi, UseBoundStore } from "zustand";

interface NewMfaForm {
  step: number;
  mfaCode: string;
}

interface NewMfaAction {
  reset: () => void;

  setStep: (step: number) => void;
  setMfaCode: (mfaCode: string) => void;
}

const initialMfaForm: NewMfaForm = {
  step: 0,
  mfaCode: "",
};

export type UseAppForm = UseBoundStore<StoreApi<NewMfaForm & NewMfaAction>>;

const useNewMfaForm = create<NewMfaForm & NewMfaAction>()((set) => ({
  ...initialMfaForm,

  reset: () => {
    set(initialMfaForm);
  },
  setStep: (step: number) => {
    set({ step });
  },
  setMfaCode: (mfaCode: string) => {
    set({ mfaCode });
  },
}));
export default useNewMfaForm;
