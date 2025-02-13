import { FC, useState } from "react";
import { numeral } from "@util/num";

import { Card, CardContent, Stack, Typography, Chip } from "@mui/material";
import { DataSaverOff, ExtensionOff } from "@mui/icons-material";
interface Props {
  app: App.Info;
}

export const NavAppCard: FC<Props> = ({ app }) => {
  const [elevation, setElevation] = useState(3);

  return (
    <Card
      variant={"elevation"}
      elevation={elevation}
      onMouseEnter={() => setElevation(5)}
      onMouseLeave={() => setElevation(3)}
      onClick={() => window.open(app.callback, "_blank")}
      sx={{
        transition: "box-shadow .3s ease-in-out",
        minHeight: "100%",
        display: "flex",
        flexDirection: "column",
        justifyContent: "space-between",
      }}
    >
      <CardContent
        sx={{
          pb: 1,
        }}
      >
        <Typography
          gutterBottom
          variant="h6"
          sx={{
            display: "flex",
            alignItems: "center",
            margin: 0,
          }}
        >
          {app.name}
        </Typography>
        <Typography variant={"subtitle2"} color={"text.secondary"}>
          {new URL(app.callback).host}
        </Typography>
      </CardContent>

      <Stack
        direction="row"
        sx={{
          justifyContent: "flex-end",
          pb: 1.5,
          px: 1.5,
          "&>.MuiChip-root": {
            paddingLeft: "3.5px",
          },
        }}
      >
        {app.linkOff ? (
          <Chip
            variant="outlined"
            size="small"
            icon={<ExtensionOff />}
            label={"未接入"}
            sx={{
              opacity: 0.7,
            }}
          />
        ) : (
          <Chip
            variant="outlined"
            size="small"
            color={"primary"}
            icon={<DataSaverOff color={"info"} />}
            label={numeral(app.views)}
          />
        )}
      </Stack>
    </Card>
  );
};
export default NavAppCard;
