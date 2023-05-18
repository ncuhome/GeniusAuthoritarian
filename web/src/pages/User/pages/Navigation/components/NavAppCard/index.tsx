import { FC, useState } from "react";

import {
  Card,
  CardContent,
  Typography,
  CardActions,
  Button,
} from "@mui/material";

import { App } from "@api/v1/user/app";

interface Props {
  app: App;
}

export const NavAppCard: FC<Props> = ({ app }) => {
  const [elevation, setElevation] = useState(5);

  return (
    <Card
      variant={"elevation"}
      elevation={elevation}
      onMouseEnter={() => setElevation(7)}
      onMouseLeave={() => setElevation(5)}
      sx={{
        transition: "box-shadow .3s ease-in-out",
      }}
    >
      <CardContent>
        <Typography
          gutterBottom
          variant="h5"
          sx={{
            margin: 0,
          }}
        >
          {app.name}
        </Typography>
      </CardContent>
    </Card>
  );
};
export default NavAppCard;
