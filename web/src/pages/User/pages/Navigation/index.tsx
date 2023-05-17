import { FC, useState } from "react";
import { useMount, useLoadingToast, useInterval } from "@hooks";

import { Stack, Typography } from "@mui/material";

import { GetAccessibleAppList } from "@api/v1/user/app";

import { useUser } from "@store";

export const Navigation: FC = () => {
  const [loadAccessible, closeLoadAccessible] = useLoadingToast();
  const [onRequestAccessible, setOnRequestAccessible] = useState(false);

  const accessibleApps = useUser((state) => state.accessibleApps);
  const setAccessibleApps = useUser((state) =>
    state.setState("accessibleApps")
  );

  async function loadAccessibleApps() {
    setOnRequestAccessible(true);
    try {
      const data = await GetAccessibleAppList();
      setAccessibleApps(data);
      closeLoadAccessible();
    } catch ({ msg }) {
      if (msg) loadAccessible(msg as string);
    }
    setOnRequestAccessible(false);
  }

  useInterval(
    loadAccessibleApps,
    !accessibleApps && !onRequestAccessible ? 2000 : null
  );
  useMount(() => {
    if (!accessibleApps) loadAccessibleApps();
  });

  return (
    <Stack justifyContent={"center"} alignItems={"center"} height={"100%"}>
      <Typography variant={"h5"} fontWeight={"bold"} sx={{ opacity: 0.5 }}>
        别急，马上写
      </Typography>
    </Stack>
  );
};
export default Navigation;
