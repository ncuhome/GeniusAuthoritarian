import { FC } from "react";

import {
  Box,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
} from "@mui/material";
import { unix } from "dayjs";
import UserAgent from "@components/user/profile/UserAgent";

interface Props {
  records: User.LoginRecord[];
}

export const OnlineDevice: FC<Props> = ({ records }) => {
  return (
    <Box
      sx={{
        marginTop: "0.5rem",
        width: "100%",
        overflowY: "auto",
        whiteSpace: "nowrap",
      }}
    >
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>上线时间</TableCell>
            <TableCell>应用</TableCell>
            <TableCell>设备</TableCell>
            <TableCell>操作</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {records.map((record) => (
            <TableRow key={record.id}>
              <TableCell>
                {unix(record.createdAt).format("MM/DD HH:mm")}
              </TableCell>
              <TableCell>{record.target}</TableCell>
              <TableCell>
                {record.useragent ? (
                  <UserAgent useragent={record.useragent} />
                ) : undefined}
              </TableCell>
              <TableCell></TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </Box>
  );
};
export default OnlineDevice;
