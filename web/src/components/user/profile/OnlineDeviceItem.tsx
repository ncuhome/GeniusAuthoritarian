import { memo, useState } from "react";
import { unix } from "dayjs";
import useU2F from "@hooks/data/useU2F";
import toast from "react-hot-toast";

import UserAgent from "@components/user/profile/UserAgent";
import { TableCell, TableRow } from "@mui/material";
import { LoadingButton } from "@mui/lab";

import useUser from "@store/useUser";
import { apiV1User } from "@api/v1/user/base";

interface Props extends User.LoginRecordOnline {}

export const OnlineDeviceItem = memo<Props>(({ ...record }) => {
  const [isLoading, setIsLoading] = useState(false);

  const setDeviceOffline = useUser((state) => state.setDeviceOffline);
  const setDialog = useUser((state) => state.setDialog);

  const { refreshToken } = useU2F();

  const onDeviceOffline = async () => {
    const yes = await setDialog({
      title: "下线此设备",
      content: record.isMe ? "请注意，你正在下线你自己" : undefined,
    });
    if (!yes) return;

    const token = await refreshToken();

    setIsLoading(true);
    try {
      await apiV1User.patch(
        "logout",
        {
          id: record.id,
        },
        {
          headers: {
            Authorization: token,
          },
        },
      );
      setDeviceOffline(record.id);
    } catch (err) {
      if (err instanceof Error) toast.error(err.message);
    }
    setIsLoading(false);
  };

  return (
    <TableRow>
      <TableCell>{`${unix(record.createdAt).format("MM/DD")}~${unix(
        record.validBefore,
      ).format("MM/DD HH:mm")}`}</TableCell>
      <TableCell>{record.target}</TableCell>
      <TableCell>
        {record.useragent ? (
          <UserAgent useragent={record.useragent} />
        ) : undefined}
      </TableCell>
      <TableCell>
        <LoadingButton
          size={"small"}
          variant={"outlined"}
          color={record.isMe ? undefined : "warning"}
          onClick={onDeviceOffline}
          loading={isLoading}
        >
          下线
        </LoadingButton>
      </TableCell>
    </TableRow>
  );
});
export default OnlineDeviceItem;
