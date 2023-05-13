import { FC, useRef, useState, useEffect } from "react";
import { useLoadingToast, useMount, useInterval, useTimeout } from "@hooks";
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
  Collapse,
  Skeleton,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  OutlinedInput,
  ListItemText,
  Typography,
  Box,
} from "@mui/material";

import { GetOwnedAppList, ApplyApp } from "@api/v1/user/app";
import { ListGroups } from "@api/v1/user/group";

import { shallow } from "zustand/shallow";
import { useUser, useAppForm, useGroup } from "@store";

export const App: FC = () => {
  const apps = useUser((state) => state.apps);
  const setApps = useUser((state) => state.setState("apps"));
  const setDialog = useUser((state) => state.setDialog);

  const [onRequestApps, setOnRequestApps] = useState(true);
  const [loadAppsToast, closeAppsToast] = useLoadingToast();

  const nameInput = useRef<HTMLInputElement | null>(null);
  const callbackInput = useRef<HTMLInputElement | null>(null);

  const [name, callback, permitAll, permitGroups, nameError, callbackError] =
    useAppForm(
      (state) => [
        state.name,
        state.callback,
        state.permitAll,
        state.permitGroups,
        state.nameError,
        state.callbackError,
      ],
      shallow
    );
  const [
    setName,
    setCallback,
    setPermitAll,
    setPermitGroups,
    setNameError,
    setCallbackError,
  ] = useAppForm(
    (state) => [
      state.setState("name"),
      state.setState("callback"),
      state.setState("permitAll"),
      state.setState("permitGroups"),
      state.setState("nameError"),
      state.setState("callbackError"),
    ],
    shallow
  );
  const resetForm = useAppForm((state) => state.reset);
  const [showSelectGroups, setShowSelectGroups] = useState(!permitAll);

  const groups = useGroup((state) => state.groups);
  const setGroups = useGroup((state) => state.setState("groups"));
  const onRequestGroups = useRef(false);

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

  async function loadGroups() {
    onRequestGroups.current = true;
    try {
      const data = await ListGroups();
      setGroups(data);
    } catch ({ msg }) {
      if (msg) toast.error(msg as string);
    }
    onRequestGroups.current = false;
  }

  async function checkForm(): Promise<boolean> {
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

    if (!permitAll && (!permitGroups || permitGroups.length === 0)) {
      const yes = await setDialog({
        title: "警告",
        content: "您没有授权任何身份组使用，你确定要继续创建吗",
      });
      if (!yes) return false;
    }

    return true;
  }

  async function createApp() {
    if (!(await checkForm())) return;
    try {
      const data = await ApplyApp(name, callback, permitAll);
      setApps([data, ...apps!]);
      resetForm();
      setDialog({
        title: "密文仅在此显示一次，请妥善保管",
        content: (
          <Stack>
            <Typography>AppSecret:</Typography>
            <Box
              sx={{
                overflowY: "auto",
              }}
            >
              <pre>{data.appSecret}</pre>
            </Box>
          </Stack>
        ),
      });
    } catch ({ msg }) {
      if (msg) toast.error(msg as string);
    }
  }

  useTimeout(() => setShowSelectGroups(false), permitAll ? 300 : null);
  useEffect(() => {
    if (!permitAll) {
      setShowSelectGroups(true);
      if (onRequestGroups.current) return;
      loadGroups();
    }
  }, [permitAll]);

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
          <Grid item xs={12}>
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
          <Grid
            item
            xs={12}
            sm={6}
            sx={{
              transition: `padding .3s ease-out${
                permitAll ? " .3s" : ""
              }, opacity 0.3s ease-out${permitAll ? "" : " .3s"}`,
              py: permitAll ? "0!important" : undefined,
              opacity: permitAll ? "0" : undefined,
            }}
          >
            <Collapse in={showSelectGroups}>
              {groups ? (
                <FormControl fullWidth>
                  <InputLabel>授权身份组</InputLabel>
                  <Select
                    multiple
                    value={permitGroups || []}
                    onChange={({ target: { value } }) =>
                      setPermitGroups(
                        typeof value === "string" ? value.split(",") : value
                      )
                    }
                    input={<OutlinedInput label="授权身份组" />}
                    renderValue={(selected) => selected.join(", ")}
                  >
                    {groups.map((group) => (
                      <MenuItem key={group.id} value={group.name}>
                        <Checkbox
                          checked={
                            (permitGroups?.indexOf(group.name) ?? -2) > -1
                          }
                        />
                        <ListItemText primary={group.name} />
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              ) : (
                <Skeleton variant={"rounded"} width={"100%"} height={56} />
              )}
            </Collapse>
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
