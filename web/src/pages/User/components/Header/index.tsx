import {FC} from "react";
import {useNavigate} from "react-router-dom";
import logo from "@/assets/img/logo-lg.png";

import {Box, Stack, Paper} from "@mui/material";

export const Header: FC = () => {
    const nav = useNavigate();

    const handleGoHome = () => nav("/user/");

    return (
    <Stack
      sx={{
        px: "3rem",
        height: "inherit",
      }}
      component={Paper}
      elevation={6}
    >
      <Box
        sx={{
          height: "100%",
          display: "flex",
          alignItems: "center",
          "&>img": {
            height: "60%",
          },
        }}
        onClick={handleGoHome}
      >
        <img src={logo} alt={""} />
      </Box>
    </Stack>
  );
};
export default Header;
