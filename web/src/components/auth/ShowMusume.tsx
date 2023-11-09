import { FC, useMemo } from "react";

import logo from "@/assets/img/logo-white.png";
import logo_webp from "@/assets/img/logo-white.webp";

import bkg from "@/assets/img/bkg.png";
import bgk_2910230214 from "@/assets/img/bkg_2910230214.png";
import bgk_627660024 from "@/assets/img/bgk_627660024.png";
import bgk_627660022 from "@/assets/img/bgk_627660022.png";

import bgk_webp from "@/assets/img/bkg.webp";
import bgk_2910230214_webp from "@/assets/img/bkg_2910230214.webp";
import bgk_627660024_webp from "@/assets/img/bgk_627660024.webp";
import bgk_627660022_webp from "@/assets/img/bgk_627660022.webp";

import { Stack, Box, Paper } from "@mui/material";

const images = [bkg, bgk_2910230214, bgk_627660024, bgk_627660022];
const webpImages = [
  bgk_webp,
  bgk_2910230214_webp,
  bgk_627660024_webp,
  bgk_627660022_webp,
];

export const ShowMusume: FC = () => {
    const img = useMemo(() => {
      const index = Math.floor(Math.random() * images.length);
      return { src: images[index], webp: webpImages[index] };
    }, []);

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
            "&>picture>img": {
              maxHeight: "100%",
              maxWidth: "100%",
              width: "15rem",
            },
          }}
        >
          <picture>
            <source type="image/webp" srcSet={logo_webp} />
            <img src={logo} alt={"家园工作室"} />
          </picture>
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
          <picture>
            <source type="image/webp" srcSet={img.webp} />
            <img
              style={{
                width: "100%",
              }}
              src={img.src}
              alt={"看板 MuSiMie"}
              title={"“走，上工！”"}
            />
          </picture>
        </Box>
      </Stack>
    );
};
export default ShowMusume;
