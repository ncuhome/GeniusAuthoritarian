import { FC } from "react";
import "./styles.css";

import LoginForm from "@components/auth/LoginForm";
import ShowMusume from "@components/auth/ShowMusume";
import { Box, Stack } from "@mui/material";

export const Home: FC = () => {
  return (
    <Stack
      flexDirection={"row"}
      sx={{
        width: "100%",
        height: "100%",
        "&>div": {
          height: "100%",
        },
      }}
    >
      <Box
        className={"show-bar"}
        sx={{
          width: "20rem",
        }}
      >
        <ShowMusume />
      </Box>
      <Box
        className={"login-form"}
        sx={{
          flexGrow: 1,
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          padding: "2rem 3rem",
          boxSizing: "border-box",
        }}
      >
        <LoginForm />
      </Box>
    </Stack>
  );
};

export default Home;
