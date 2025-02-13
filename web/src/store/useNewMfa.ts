import { create } from "zustand";

interface NewMfaForm {
  step: number;
  mfaCode: string;
}

interface NewMfaAction {
  reset: () => void;
}

const initialMfaForm: NewMfaForm = {
  step: 0,
  mfaCode: "",
};

const useNewMfaForm = create<NewMfaForm & NewMfaAction>()((set) => ({
  ...initialMfaForm,

  reset: () => {
    set(initialMfaForm);
  },
}));
export default useNewMfaForm;
