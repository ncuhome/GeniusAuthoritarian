import { FC, useState } from "react";

import routes from "./route";
import Block from "@components/user/Block";
import { Container, Tab, Tabs } from "@mui/material";

import useLoginData from "@hooks/data/useLoginData";

export const Admin: FC = () => {
  const [index, setIndex] = useState(0);

  useLoginData("week");

  return (
    <Container>
      <Block disablePadding>
        <Tabs
          value={index}
          variant="scrollable"
          scrollButtons="auto"
          onChange={(_e, target: number) => setIndex(target)}
        >
          {routes.map((route) => (
            <Tab key={route.label} label={route.label} />
          ))}
        </Tabs>
      </Block>

      {routes[index].element}
    </Container>
  );
};
export default Admin;
