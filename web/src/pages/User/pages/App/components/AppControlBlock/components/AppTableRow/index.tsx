import { FC, useState } from "react";

import { TipIconButton } from "@components";
import { Divider, Stack, TableCell, TableRow } from "@mui/material";
import {
  DeleteForeverOutlined,
  DriveFileRenameOutlineOutlined,
} from "@mui/icons-material";

interface Props {
  app: App.Detailed;
  onModify: (app: App.Detailed) => void;
  onDelete: (app: App.Detailed) => Promise<void>;
}

export const AppTableRow: FC<Props> = ({ app, onDelete, onModify }) => {
  const [deletingApp, setDeletingApp] = useState(false);

    async function handleDeleteApp(app: App.Detailed) {
      setDeletingApp(true);
      await onDelete(app);
      setDeletingApp(false);
    }

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
            <TipIconButton title={"编辑"} onClick={() => onModify(app)}>
              <DriveFileRenameOutlineOutlined />
            </TipIconButton>
            <TipIconButton
              title={"删除"}
              onClick={() => handleDeleteApp(app)}
              disabled={deletingApp}
            >
              <DeleteForeverOutlined />
            </TipIconButton>
          </Stack>
        </TableCell>
      </TableRow>
    );
};
export default AppTableRow;
