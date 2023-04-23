import {FC, useMemo} from "react";
import {useLocation} from "react-router-dom";
import "./styles.css";

import { Box, Typography } from "@mui/material";
import { ClearRounded } from "@mui/icons-material";

export const Error: FC = () => {
    const loc = useLocation()
    const title = useMemo(() => (loc.state && loc.state.title) ? loc.state.title : "未知错误", [loc.state])
    const content = useMemo(() => (loc.state && loc.state.content) ? loc.state.content : "", [loc.state])

    return (
        <Box
            sx={{
                height: "100%",
                width: "100%",
                display: "flex",
                justifyContent: "center",
                alignItems: "center",
            }}
    >
      <Box
        sx={{
          display: "flex",
          flexDirection: "column",
          textAlign: "center",
          justifyContent: "center",
          "&>*": {
            marginBottom: "1rem!important",
          },
        }}
      >
        <Box
          sx={{
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            "& svg": {
              borderStyle: "solid",
              borderRadius: "50%",
              fontSize: "8rem",
              animation: "error-page-fork .5s ease",
              animationFillMode: "forwards",
              boxSizing: "border-box",
            },
          }}
        >
          <ClearRounded />
        </Box>
        <Typography
          variant={"h4"}
          sx={{
            fontWeight: 600,
            letterSpacing: "0.25rem",
          }}
        >
          {title}
        </Typography>
        {content ? (
          <Typography
            sx={{
              color: "#999",
              wordBreak: "break-all",
            }}
          >
            {content}
          </Typography>
        ) : undefined}
      </Box>
    </Box>
  );
};

export default Error;
