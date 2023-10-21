import { FC } from "react";

import {
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  SxProps,
} from "@mui/material";

interface Props {
  logo: string;
  text: string;
  disableDivider?: boolean;
  onClick: () => void;
  sx?: SxProps;
}

export const LoginItem: FC<Props> = ({ logo, text, disableDivider, onClick, sx }) => {
  return (
    <ListItem disablePadding divider={!disableDivider} sx={sx}>
      <ListItemButton onClick={onClick}>
        <ListItemIcon>
          <img
            style={{
              width: "1.8rem",
            }}
            src={logo}
            alt={text}
          />
        </ListItemIcon>
        <ListItemText primary={text} />
      </ListItemButton>
    </ListItem>
  );
};
export default LoginItem;
