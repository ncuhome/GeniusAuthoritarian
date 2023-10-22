import { FC, useMemo, useState } from "react";
import { unix, duration } from "dayjs";

import {
  ListItem,
  ListItemText,
  Stack,
  TextField,
  IconButton,
  Divider,
  Button,
} from "@mui/material";
import { DeleteOutline, ModeEditOutlineOutlined } from "@mui/icons-material";

interface Props {
  item: User.Passkey.Cred;
  enableRename: boolean;
  onRename: (name: string) => Promise<void>;
  onDelete: () => void;
}

export const PasskeyItem: FC<Props> = ({
  item,
  enableRename,
  onRename,
  onDelete,
}) => {
  const name = useMemo(
    () => (item.name ? item.name : `Device${item.id}`),
    [item.name]
  );

  const [isEditing, setIsEditing] = useState(enableRename);
  const [newName, setNewName] = useState(name);

  const lastUsed = useMemo(() => {
    if (item.last_used_at === 0) return "还未使用过";
    const time = duration(new Date().getTime() - item.last_used_at * 1000);
    let word = "上次使用于 ";
    const month = time.months();
    if (month > 0) word += `${month} 月`;
    const day = time.days();
    if (day > 0) word += `${day} 天`;
    const minute = time.minutes();
    if (minute > 0) word += `${minute} 分钟`;
    else word += `${time.seconds()} 秒`;
    return word + "前";
  }, [item.last_used_at]);

  const renderItem = (item: User.Passkey.Cred, isEditing: boolean) => {
    if (isEditing)
      return (
        <ListItem
          secondaryAction={
            <Stack flexDirection={"row"}>
              <Button
                variant={"outlined"}
                sx={{
                  mr: 1,
                }}
                onClick={async () => {
                  if (newName !== name) await onRename(newName);
                  setIsEditing(false);
                }}
              >
                保存
              </Button>
              <Button onClick={() => setIsEditing(false)}>取消</Button>
            </Stack>
          }
        >
          <TextField
            value={newName}
            onChange={(e) => setNewName(e.target.value)}
            inputProps={{
              style: {
                height: "0.8rem",
                width: "10rem",
              },
            }}
          />
        </ListItem>
      );
    else
      return (
        <ListItem
          secondaryAction={
            <Stack flexDirection={"row"}>
              <IconButton
                onClick={() => {
                  setNewName(name);
                  setIsEditing(true);
                }}
              >
                <ModeEditOutlineOutlined />
              </IconButton>
              <Divider
                orientation="vertical"
                variant="middle"
                flexItem
                sx={{
                  mx: 0.5,
                  my: 1.5,
                }}
              />
              <IconButton color={"error"} onClick={onDelete}>
                <DeleteOutline />
              </IconButton>
            </Stack>
          }
        >
          <ListItemText
            primary={name}
            secondary={`创建于 ${unix(item.created_at).format(
              "YYYY/MM/DD"
            )}，${lastUsed}`}
          />
        </ListItem>
      );
  };

  return renderItem(item, isEditing);
};

export default PasskeyItem;
