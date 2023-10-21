import { FC, useEffect, useMemo, useState } from "react";

import { ListItem, ListItemText, Collapse, Divider } from "@mui/material";
import { unix, duration } from "dayjs";

interface Props {
  item: User.Passkey.Cred;
  divider: boolean;
}

export const PasskeyItem: FC<Props> = ({ item: itemProp, divider }) => {
  const [deleted, setDeleted] = useState(false);
  const [item, setItem] = useState(itemProp);

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

  useEffect(() => {
    setItem(itemProp);
  }, [itemProp]);

  return (
    <>
      <Collapse in={!deleted}>
        <ListItem>
          <ListItemText
            primary={item.name ? item.name : `Device${item.id}`}
            secondary={`创建于 ${unix(item.created_at).format(
              "YYYY/MM/DD"
            )}，${lastUsed}`}
          />
        </ListItem>
      </Collapse>
      {divider ? <Divider variant={"middle"} component="li" /> : undefined}
    </>
  );
};

export default PasskeyItem;
