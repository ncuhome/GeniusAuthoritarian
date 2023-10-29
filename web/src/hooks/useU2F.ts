import { useUserApiV1 } from "@api/v1/user/hook";

import { shallow } from "zustand/shallow";
import useU2fDialog from "@store/useU2fDialog";

export const useU2F = () => {
  const dialogOpen = useU2fDialog((state) => state.open);
  const { phone, mfa, passkey } = useU2fDialog(
    (state) => ({
      phone: state.phone,
      mfa: state.mfa,
      passkey: state.passkey,
    }),
    shallow
  );

  const {} = useUserApiV1("u2f/", {
    revalidateOnFocus: false,
  });
};

export default useU2F;
