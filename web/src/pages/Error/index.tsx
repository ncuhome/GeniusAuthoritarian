import { FC, useMemo } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import "./styles.css";

import { Box, Typography, ButtonGroup, Button } from "@mui/material";
import { ClearRounded } from "@mui/icons-material";

import { GoLogin } from "@util/nav";

export const Error: FC = () => {
  const nav = useNavigate();
  const loc = useLocation();
  const state: ErrorState | undefined = useMemo(() => loc.state, [loc]);

  const title = useMemo(() => state?.title || "未知错误", [state]);

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
        {state?.content ? (
          <Typography
            sx={{
              color: "#999",
              wordBreak: "break-all",
            }}
          >
            {state.content}
          </Typography>
        ) : undefined}

        <Box mt={3.5}>
          <ButtonGroup
            variant="text"
            sx={{
              "&>button": {
                border: "unset!important",
              },
            }}
          >
            <Button
              onClick={() =>
                window.open(
                  "https://applink.feishu.cn/client/chat/chatter/add_by_link?link_token=250la2bb-add8-4fde-ac94-57376aee2e40",
                  "_blank"
                )
              }
            >
              反馈
            </Button>

            {state?.retryAppCode !== undefined ? (
              <Button onClick={() => GoLogin(nav, state!.retryAppCode!)}>
                重试
              </Button>
            ) : undefined}
          </ButtonGroup>
        </Box>
      </Box>
    </Box>
  );
};

export default Error;
