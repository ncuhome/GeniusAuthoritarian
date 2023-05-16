import { FC } from "react";

import { TipIconButton } from "@components";
import { Divider, Stack, TableCell, TableRow } from "@mui/material";
import {
  DeleteForeverOutlined,
  DriveFileRenameOutlineOutlined,
} from "@mui/icons-material";

import { App } from "@api/v1/user/app";

interface Props {
  app: App;
  handleModify: (app: App) => void;
  handleDelete: (app: App) => void;
}

export const AppTableRow: FC<Props> = ({ app, handleDelete, handleModify }) => {
  return (
    <TableRow hover role="checkbox" tabIndex={-1}>
      <TableCell>{app.name}</TableCell>
      <TableCell>{app.appCode}</TableCell>
      <TableCell>
        {app.permitAllGroup
          ? "ALL"
          : app.groups.length > 0
          ? app.groups.map((group) => group.name).join("，")
          : "NONE"}
      </TableCell>
      <TableCell>{app.callback}</TableCell>
      <TableCell>
        <Stack
          flexDirection={"row"}
          alignItems={"center"}
          divider={
            <Divider
              orientation="vertical"
              sx={{
                height: "15px",
                mx: "3px",
              }}
            />
          }
        >
          <TipIconButton title={"编辑"} onClick={() => handleModify(app)}>
            <DriveFileRenameOutlineOutlined />
          </TipIconButton>
          <TipIconButton title={"删除"} onClick={() => handleDelete(app)}>
            <DeleteForeverOutlined />
          </TipIconButton>
        </Stack>
      </TableCell>
    </TableRow>
  );
};
export default AppTableRow;
