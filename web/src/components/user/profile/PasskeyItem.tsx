import { FC, useState } from "react";

import { ListItem, ListItemText, Collapse, Divider } from "@mui/material";

interface Props {
  item: User.Passkey.Cred;
  divider: boolean;
}

export const PasskeyItem: FC<Props> = ({ item, divider }) => {
  const [deleted, setDeleted] = useState(false);

  return (
    <>
      <Collapse in={!deleted}>
        <ListItem>
          <ListItemText primary={item.name ? item.name : `Device${item.id}`} />
        </ListItem>
      </Collapse>
      {divider ? <Divider variant={"middle"} component="li" /> : undefined}
    </>
  );
};

export default PasskeyItem;
