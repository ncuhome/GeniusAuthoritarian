import { FC } from "react";

import { Header } from "./components";
import { Box } from "@mui/material";

export const User: FC = () => {
  return (
    <Box>
      <Box
        sx={{
          width: "100%",
          position: "sticky",
          height: "3.5rem",
        }}
      >
        <Header />
      </Box>
    </Box>
  );
};
export default User;
