import { FC, useState } from "react";
import { useMount, useLoadingToast, useInterval } from "@hooks";

import { AppListBlock } from "./components";
import { Container, Grid } from "@mui/material";

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
    <Container>
      {accessibleApps ? (
        <AppListBlock title={"全站"} apps={accessibleApps.permitAll} />
      ) : null}
    </Container>
  );
};
export default Navigation;
