import { FC, useState } from "react";
import toast from "react-hot-toast";
import { numeral } from "@util/num";

import {
  Card,
  CardContent,
  Stack,
  Typography,
  Chip,
  CardActions,
} from "@mui/material";
import { DataSaverOff, LinkOff } from "@mui/icons-material";

import { apiV1User } from "@api/v1/user/base";

interface Props {
  app: App.Info;
}

export const NavAppCard: FC<Props> = ({ app }) => {
  const [elevation, setElevation] = useState(5);

  async function onLandingApp(id: number, code?: string) {
    try {
      const {
        data: {
          data: { url },
        },
      } = await apiV1User.get("app/landing", {
        params: {
          id,
          code,
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
      <CardContent
        sx={{
          pb: 1,
        }}
      >
        <Typography
          gutterBottom
          variant="subtitle1"
          sx={{
            display: "flex",
            alignItems: "center",
            margin: 0,
          }}
        >
          {app.name}
          {app.linkOff ? (
            <LinkOff
              fontSize={"small"}
              color={"warning"}
              style={{ display: "inline", marginLeft: "0.5rem", opacity: 0.5 }}
            />
          ) : undefined}
        </Typography>

        <Stack
          direction="row"
          justifyContent={"flex-end"}
          sx={{
            mt: 1,
          }}
        ></Stack>
      </CardContent>
      <CardActions
        sx={{
          justifyContent: "flex-end",
          pt: 0,
          pb: 1.5,
          px: 1.5,
        }}
      >
        <Chip
          variant="outlined"
          size="small"
          color={"primary"}
          icon={<DataSaverOff color={"info"} />}
          label={numeral(app.views)}
          sx={{
            opacity: app.linkOff ? 0 : 1,
          }}
        />
      </CardActions>
    </Card>
  );
};
export default NavAppCard;
