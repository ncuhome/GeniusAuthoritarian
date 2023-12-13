import { FC, useMemo } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import "./style.css";

import img_png from "@/assets/img/error/img_20231213.png";
import img_webp from "@/assets/img/error/img_20231213.webp";

import { useTheme } from "@mui/material";
import { Box, Stack, Typography, ButtonGroup, Button } from "@mui/material";

import { GoLogin } from "@util/nav";

export const Error: FC = () => {
  const theme = useTheme();
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
          flexDirection: { xs: "column", sm: "row" },
          alignItems: "center",
        }}
      >
        <Box
          sx={{
            maxWidth: "40%",
            width: "22rem",
            margin: { xs: "0 0 1.2rem 0", sm: "0 5rem 0 0" },
            "&>picture": {
              height: "auto",
              width: "100%",
              "&>img": {
                width: "100%",
                height: "auto",
              },
            },
          }}
        >
          <picture>
            <source srcSet={img_webp} type="image/webp" />
            <img src={img_png} alt={"ERROR"} />
          </picture>
        </Box>

        <Stack alignItems={{ xs: "center", sm: "baseline" }}>
          <Box
            className={"ops"}
            sx={{
              display: { xs: "none", sm: "block" },
              color: "text.secondary",
              mb: 2,
              "&>h1:nth-of-type(1)": {
                WebkitTextStroke: `2px ${theme.palette.text.secondary}`,
              },
            }}
          >
            <Typography variant={"h1"}>OPS!</Typography>
            <Typography variant={"h1"}>OPS!</Typography>
          </Box>

          <Stack px={0.7}>
            <Typography
              variant={"h4"}
              fontWeight={"bold"}
              letterSpacing={"0.25rem"}
              mb={1}
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
          </Stack>

          <Box mt={1.5}>
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
                    "_blank",
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
        </Stack>
      </Box>
    </Box>
  );
};

export default Error;
