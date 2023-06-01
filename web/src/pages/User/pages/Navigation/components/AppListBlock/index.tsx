import { FC } from "react";

import { BlockArea, NavAppCard } from "..";
import { Grid } from "@mui/material";

interface Props {
  title: string;
  apps: App.Info[];
}

export const AppListBlock: FC<Props> = ({ title, apps }) => {
  return (
    <BlockArea title={title}>
      <Grid container spacing={2}>
        {apps.map((app) => (
          <Grid key={app.id} item xs={12} sm={6} md={4}>
            <NavAppCard app={app} />
          </Grid>
        ))}
      </Grid>
    </BlockArea>
  );
};
export default AppListBlock;
