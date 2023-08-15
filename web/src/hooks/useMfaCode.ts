import useMfaCodeDialog from "@store/useMfaCodeDialog";

// 限 user 路由使用
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
