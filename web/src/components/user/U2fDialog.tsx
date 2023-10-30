import { FC } from "react";

import { Dialog } from "@mui/material";

import useU2fDialog from "@store/useU2fDialog";

const U2fDialog: FC = () => {
  const open = useU2fDialog((state) => state.open);

  return <Dialog open={open}></Dialog>;
};
export default U2fDialog;
