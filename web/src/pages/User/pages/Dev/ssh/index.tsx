import { FC } from "react";

import { Stack } from "@mui/material";
import { LoadingButton } from "@mui/lab";
import { LockPerson, Link } from "@mui/icons-material";
import Block from "@components/user/Block";

const Ssh: FC = () => {
  return (
    <Block>
      <Stack alignItems={"center"} justifyContent={"center"} height={"10rem"}>
        <Stack
          justifyContent={"center"}
          alignItems={"center"}
          flexDirection={"row"}
          mb={"1rem"}
        >
          <Link color={"disabled"} />
          <LockPerson
            sx={{
              fontSize: "3rem",
              paddingBottom: "0.55rem",
              paddingX: "0.5rem",
            }}
          />
          <Link color={"disabled"} />
        </Stack>

        <LoadingButton variant={"outlined"}>解锁</LoadingButton>
      </Stack>
    </Block>
  );
};
export default Ssh;
