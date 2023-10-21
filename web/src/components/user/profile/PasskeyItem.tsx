import { FC, useEffect, useMemo, useState } from "react";
import { unix, duration } from "dayjs";
import useMfaCode from "@hooks/useMfaCode";
import toast from "react-hot-toast";

import {
  ListItem,
  ListItemText,
  Stack,
  TextField,
  IconButton,
  Divider,
} from "@mui/material";
import { DeleteOutline, ModeEditOutlineOutlined } from "@mui/icons-material";

import { apiV1User } from "@api/v1/user/base";

interface Props {
  item: User.Passkey.Cred;
  onDelete: () => void;
}

export const PasskeyItem: FC<Props> = ({ item: itemProp, onDelete }) => {
  const [item, setItem] = useState(itemProp);
  const [isEditing, setIsEditing] = useState(false);

  const openMfaDialog = useMfaCode();

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

  const onEdit = () => {};

  useEffect(() => {
    setItem(itemProp);
  }, [itemProp]);

  const renderItem = (item: User.Passkey.Cred, isEditing: boolean) => {
    if (isEditing)
      return (
        <ListItem>
          <Stack flexDirection={"row"}>
            <TextField value={item.name} />
          </Stack>
        </ListItem>
      );
    else
      return (
        <ListItem
          secondaryAction={
            <Stack flexDirection={"row"}>
              <IconButton>
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
            primary={item.name ? item.name : `Device${item.id}`}
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
