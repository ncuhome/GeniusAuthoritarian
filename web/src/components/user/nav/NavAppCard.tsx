import { FC, useState } from "react";
import toast from "react-hot-toast";

import { Card, CardContent, Typography } from "@mui/material";

import { apiV1User } from "@api/v1/user/base";

interface Props {
  app: App.Info;
}

export const NavAppCard: FC<Props> = ({ app }) => {
  const [elevation, setElevation] = useState(5);

  async function onLandingApp(id: number) {
    try {
      const {
        data: {
          data: { url },
        },
      } = await apiV1User.get("app/landing", {
        params: {
          id,
        },
      });
      window.open(url, "_blank");
    } catch ({ msg }) {
      if (msg) toast.error(msg as string);
    }
  }

  return (
    <Card
      variant={"elevation"}
      elevation={elevation}
      onMouseEnter={() => setElevation(7)}
      onMouseLeave={() => setElevation(5)}
      onClick={() => onLandingApp(app.id)}
      sx={{
        transition: "box-shadow .3s ease-in-out",
      }}
    >
      <CardContent>
        <Typography
          gutterBottom
          variant="subtitle1"
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