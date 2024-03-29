import {create, StoreApi, UseBoundStore} from "zustand";

interface NewMfaForm {
    step: number
    mfaCode: string
}

interface NewMfaAction {
    reset: () => void;
    setState: <T extends keyof NewMfaForm>(
        key: T
    ) => (value: NewMfaForm[T]) => void;
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
    setState: (key) => (value) => set({ [key]: value }),
}));
export default useNewMfaForm
