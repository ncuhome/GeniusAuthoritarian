import { FC } from "react";

import { Stack, Typography } from "@mui/material";

export const PageNotFound: FC = () => {
  return (
    <Stack
      sx={{
        width: "100%",
        height: "100%",
        textAlign: "center",
        "&>*": {
          fontWeight: "bold",
          letterSpacing: "3px",
        },
      }}
      justifyContent={"center"}
      alignContent={"center"}
    >
      <Typography variant={"h3"}>404</Typography>
      <Typography variant={"h5"}>Not Found</Typography>
    </Stack>
  );
};
export default PageNotFound;
