import useMfaCodeDialog from "@store/useMfaCodeDialog";

export const useMfaCode = (): ((desc?: string) => Promise<string>) => {
  const setDialog = useMfaCodeDialog((state) => state.setDialog);
  const resetDialog = useMfaCodeDialog((state) => state.resetDialog);

  return (desc) => {
    return new Promise((resolve, reject) => {
      setDialog((code) => {
        if (code) resolve(code);
        else reject();
        resetDialog();
      }, desc);
    });
  };
};
export default useMfaCode;
