import { FC, ReactNode, useMemo, useState } from "react";

import { Container, Tab, Tabs, Typography } from "@mui/material";
import Block from "@components/user/Block";
import routes from "./route";

export const Dev: FC = () => {
  const [index, setIndex] = useState(0);
  const content = useMemo(() => routes[index], [index]);

  return (
    <Container>
      <Block
        sx={{
          padding: "unset!important",
        }}
      >
        <Tabs
          value={index}
          variant="scrollable"
          scrollButtons="auto"
          onChange={(e, target: number) => setIndex(target)}
        >
          {routes.map((route) => (
            <Tab key={route.label} label={route.label} />
          ))}
        </Tabs>
      </Block>

      <Typography
        sx={{
          opacity: "0.6",
          my: "0.6rem",
          mx: "0.2rem",
          whiteSpace: "pre-wrap",
        }}
      >
        {content.description}
      </Typography>

      {content.element}
    </Container>
  );
};
export default Dev;
