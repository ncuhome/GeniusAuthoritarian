import { FC, useMemo } from "react";
import logo from "@/assets/img/logo-white.png";
import bkg from "@/assets/img/bkg.png";
import bgk_2910230214 from "@/assets/img/bkg_2910230214.png";
import bgk_627660024 from "@/assets/img/bgk_627660024.png";
import bgk_627660022 from "@/assets/img/bgk_627660022.png";

import { Stack, Box, Paper } from "@mui/material";

const images = [bkg, bgk_2910230214, bgk_627660024, bgk_627660022];

export const ShowMusume: FC = () => {
  const img = useMemo(
    () => images[Math.floor(Math.random() * images.length)],
    []
  );

  return (
    <Stack
      sx={{
        height: "100%",
        width: "100%",
      }}
      component={Paper}
      elevation={5}
    >
      <Box
        sx={{
          padding: "2.5rem 4rem",
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          "&>img": {
            maxHeight: "100%",
            maxWidth: "100%",
            width: "15rem",
          },
        }}
      >
        <img src={logo} alt={"家园工作室"} />
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
        <img
          style={{
            width: "100%",
          }}
          src={img}
          alt={"看板 MuSiMie"}
          title={"“走，上工！”"}
        />
      </Box>
    </Stack>
  );
};
export default ShowMusume;
