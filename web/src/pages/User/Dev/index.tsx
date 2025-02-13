import { FC } from "react";

import routes from "./route";
import Block from "@components/user/Block";
import { Container, Tab, Tabs, Typography } from "@mui/material";

import useDevRoute from "@store/useDevRoute";

export const Dev: FC = () => {
  const index = useDevRoute((state) => state.index);

  return (
    <Container>
      <Block disablePadding>
        <Tabs
          value={index}
          variant="scrollable"
          scrollButtons="auto"
          onChange={(_e, target: number) =>
            useDevRoute.setState({ index: target })
          }
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
        {routes[index].description}
      </Typography>

      {routes[index].element}
    </Container>
  );
};
export default Dev;
