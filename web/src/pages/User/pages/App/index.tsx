import { FC, useRef, useState } from "react";
import { useLoadingToast, useMount, useInterval } from "@hooks";
import toast from "react-hot-toast";

import { Block } from "@/pages/User/components";
import {
  Container,
  TextField,
  Grid,
  Checkbox,
  FormControlLabel,
  Stack,
  Button,
} from "@mui/material";

import { GetOwnedAppList, ApplyApp } from "@api/v1/user/app";

import { shallow } from "zustand/shallow";
import { useUser, useAppForm } from "@store";

export const App: FC = () => {
  const apps = useUser((state) => state.apps);
  const setApps = useUser((state) => state.setState("apps"));

  const [onRequestApps, setOnRequestApps] = useState(true);
  const [loadAppsToast, closeAppsToast] = useLoadingToast();

  const nameInput = useRef<HTMLInputElement | null>(null);
  const callbackInput = useRef<HTMLInputElement | null>(null);

  const [name, callback, permitAll, nameError, callbackError] = useAppForm(
    (state) => [
      state.name,
      state.callback,
      state.permitAll,
      state.nameError,
      state.callbackError,
    ],
    shallow
  );
  const [setName, setCallback, setPermitAll, setNameError, setCallbackError] =
    useAppForm(
      (state) => [
        state.setState("name"),
        state.setState("callback"),
        state.setState("permitAll"),
        state.setState("nameError"),
        state.setState("callbackError"),
      ],
      shallow
    );
  const [resetForm] = useAppForm((state) => [state.reset], shallow);

  async function loadApps() {
    setOnRequestApps(true);
    try {
      const data = await GetOwnedAppList();
      setApps(data);
      closeAppsToast();
    } catch ({ msg }) {
      if (msg) loadAppsToast(msg as string);
    }
    setOnRequestApps(false);
  }

  function checkForm(): boolean {
    if (!name) {
      setNameError(true);
      toast.error("应用名称不能为空");
      nameInput.current!.focus();
      return false;
    } else {
      setNameError(false);
    }

    if (!callback || callback.indexOf("https://") !== 0) {
      setCallbackError(true);
      if (!callback) toast.error("回调地址不能为空");
      else toast.error("回调地址仅支持 https 协议");
      callbackInput.current!.focus();
      return false;
    } else {
      setCallbackError(false);
    }

    return true;
  }

  async function createApp() {
    if (!checkForm()) return;
    try {
      const data = await ApplyApp(name, callback, permitAll);
      setApps([data, ...apps!]);
      toast.success("创建成功");
      resetForm();
    } catch ({ msg }) {
      if (msg) toast.error(msg as string);
    }
  }

  useInterval(loadApps, !apps && !onRequestApps ? 2000 : null);
  useMount(() => {
    if (!apps) loadApps();
    else setOnRequestApps(false);
  });

  return (
    <Container>
      <Block title={"New"}>
        <Grid container spacing={2} marginTop={0}>
          <Grid item xs={12} sm={6}>
            <TextField
              label={"应用名称"}
              fullWidth
              inputRef={nameInput}
              color={nameError ? "error" : "primary"}
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
          </Grid>
          <Grid item xs={12} sm={6}>
            <TextField
              label={"回调地址"}
              fullWidth
              inputRef={callbackInput}
              color={callbackError ? "error" : "primary"}
              value={callback}
              onChange={(e) => setCallback(e.target.value)}
            />
          </Grid>
          <Grid item xs={12} sm={6}>
            <FormControlLabel
              control={
                <Checkbox
                  checked={permitAll}
                  onChange={(e) => setPermitAll(e.target.checked)}
                />
              }
              label="允许所有成员使用"
            />
          </Grid>
        </Grid>

        <Stack
          flexDirection={"row"}
          justifyContent={"flex-end"}
          flexWrap={"wrap"}
          sx={{
            marginTop: "1rem",
            "&>button": {
              marginLeft: "0.8rem",
            },
          }}
        >
          <Button variant={"outlined"} onClick={resetForm}>
            重置
          </Button>
          <Button variant={"contained"} onClick={createApp}>
            创建应用
          </Button>
        </Stack>
      </Block>
      <Block title={"App"}></Block>
    </Container>
  );
};
export default App;
