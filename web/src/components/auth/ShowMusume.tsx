import { FC, useMemo } from "react";

import Picture from "@/components/Picture";
import { Stack, Box, Paper } from "@mui/material";

const images = ["bkg", "bkg_2910230214", "bkg_627660024", "bkg_627660022"];

export const ShowMusume: FC = () => {
  const imgName = useMemo(() => {
    const index = Math.floor(Math.random() * images.length);
    return images[index];
  }, []);

  return (
    <Stack
      sx={{
        height: "100%",
        width: "100%",
      }}
      component={Paper}
      elevation={3}
    >
      <Box
        sx={{
          padding: "2.5rem 4rem",
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          "&>picture>img": {
            maxHeight: "100%",
            maxWidth: "100%",
            width: "15rem",
          },
        }}
      >
        <Picture name={"logo-white"} alt={"NCUHOME"} />
      </Box>
      <Box
        sx={{
          flexGrow: 1,
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          boxSizing: "border-box",
          overflow: "hidden",
        }}
      >
        <Picture
          dir={"login_bgk"}
          name={imgName}
          alt={"GoGoGo"}
          imgStyle={{
            width: "100%",
          }}
        />
      </Box>
    </Stack>
  );
};
export default ShowMusume;
