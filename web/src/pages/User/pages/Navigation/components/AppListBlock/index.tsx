import { FC } from "react";

import { BlockArea, NavAppCard } from "..";
import { Grid } from "@mui/material";

import { App } from "@api/v1/user/app";

interface Props {
  title: string;
  apps: App[];
}

export const AppListBlock: FC<Props> = ({ title, apps }) => {
  return (
    <BlockArea title={title}>
      <Grid container spacing={2}>
        {apps.map((app) => (
          <Grid key={app.id} item xs={6} sm={4}>
            <NavAppCard app={app} />
          </Grid>
        ))}
      </Grid>
    </BlockArea>
  );
};
export default AppListBlock;
