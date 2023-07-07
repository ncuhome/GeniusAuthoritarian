import { FC, useState } from "react";
import toast from "react-hot-toast";
import { numeral } from "@util/num";

import { Card, CardContent, Stack, Typography } from "@mui/material";
import { DataSaverOff, LinkOff } from "@mui/icons-material";

import { apiV1User } from "@api/v1/user/base";

import useMfaCodeDialog from "@store/useMfaCodeDialog";

interface Props {
  app: App.Info;
}

export const NavAppCard: FC<Props> = ({ app }) => {
  const setMfaCodeCallback = useMfaCodeDialog((state) =>
    state.setState("callback")
  );

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
    } catch ({ msg, response }) {
      console.log(response, (response as any)?.data);
      if ((response as any)?.data?.code === 21) {
        setMfaCodeCallback((code) => onLandingApp(id, code));
      } else if (msg) toast.error(msg as string);
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
          flexDirection={"row-reverse"}
          sx={{
            mt: 1,
            opacity: 0.8,
          }}
        >
          <Stack flexDirection={"row"} alignItems={"center"} width={"3.5rem"}>
            <DataSaverOff
              fontSize={"small"}
              sx={{
                mr: 0.7,
              }}
            />
            <span>{app.linkOff ? "--" : numeral(app.views, 0)}</span>
          </Stack>
        </Stack>
      </CardContent>
    </Card>
  );
};
export default NavAppCard;
