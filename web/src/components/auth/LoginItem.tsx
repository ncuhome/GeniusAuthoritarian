import { FC, ReactNode } from "react";

import {
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  SxProps,
} from "@mui/material";

interface Props {
  logo: string;
  webpLogo?: string;
  text: string;
  disableDivider?: boolean;
  onClick: () => void;
  sx?: SxProps;
}

export const LoginItem: FC<Props> = ({
  logo,
  webpLogo,
  text,
  disableDivider,
  onClick,
  sx,
}) => {
  return (
    <ListItem disablePadding divider={!disableDivider} sx={sx}>
      <ListItemButton onClick={onClick}>
        <ListItemIcon>
          <picture
            style={{
              display: "inline-flex",
              alignItems: "center",
            }}
          >
            {webpLogo ? (
              <source type="image/webp" srcSet={webpLogo} />
            ) : undefined}
            <img
              style={{
                width: "1.8rem",
              }}
              src={logo}
              alt={text}
            />
          </picture>
        </ListItemIcon>
        <ListItemText primary={text} />
      </ListItemButton>
    </ListItem>
  );
};
export default LoginItem;
