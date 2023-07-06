import { FC, useState } from "react";
import toast from "react-hot-toast";
import { numeral } from "@util/num";

import { Card, CardContent, Stack, Typography } from "@mui/material";
import { DataSaverOff } from "@mui/icons-material";

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
            margin: 0,
          }}
        >
          {app.name}
        </Typography>

        <Stack
          flexDirection={"row-reverse"}
          sx={{
            mt: 1,
            opacity: 0.8,
          }}
        >
          <Stack flexDirection={"row"} alignItems={"center"} width={"3.8rem"}>
            <DataSaverOff
              fontSize={"small"}
              sx={{
                mr: 1,
              }}
            />
            <span>{numeral(app.views)}</span>
          </Stack>
        </Stack>
      </CardContent>
    </Card>
  );
};
export default NavAppCard;
