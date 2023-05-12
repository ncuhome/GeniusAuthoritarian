import { FC, useState } from "react";
import { useLoadingToast, useMount, useInterval } from "@hooks";

import { Block } from "@/pages/User/components";
import { Container } from "@mui/material";

import { GetOwnedAppList } from "@api/v1/user/app";

import { useUser } from "@store";

export const App: FC = () => {
  const apps = useUser((state) => state.apps);
  const setApps = useUser((state) => state.setState("apps"));

  const [onRequest, setOnRequest] = useState(true);
  const [loadAppsToast, closeAppsToast] = useLoadingToast();

  async function loadApps() {
    try {
      const data = await GetOwnedAppList();
      setApps(data);
      closeAppsToast();
    } catch ({ msg }) {
      if (msg) loadAppsToast(msg as string);
    }
  }

  useInterval(loadApps, !apps && !onRequest ? 2000 : null);
  useMount(() => {
    if (!apps) loadApps();
    else setOnRequest(false);
  });

  return (
    <Container>
      <Block title={"New"}></Block>
      <Block title={"App"}></Block>
    </Container>
  );
};
export default App;
