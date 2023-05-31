import { FC } from "react";

import { AppListBlock } from "./components";
import { Container } from "@mui/material";

import { useUserApiV1 } from "@api/v1/user/hook";

import { useUser } from "@store";

export const Navigation: FC = () => {
  const accessibleApps = useUser((state) => state.accessibleApps);
  const setAccessibleApps = useUser((state) =>
    state.setState("accessibleApps")
  );

  useUserApiV1<App.Accessible>("app/accessible", {
    immutable: true,
    enableLoading: true,
    onSuccess: (data) => {
      setAccessibleApps(data);
    },
  });

  return (
    <Container>
      {accessibleApps ? (
        <>
          {accessibleApps.accessible.map((item) => (
            <AppListBlock
              key={item.group.id}
              title={item.group.name}
              apps={item.app}
            />
          ))}
          <AppListBlock title={"全站"} apps={accessibleApps.permitAll} />
        </>
      ) : null}
    </Container>
  );
};
export default Navigation;
