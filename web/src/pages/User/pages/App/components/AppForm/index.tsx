import { FC, useEffect, useRef, useState } from "react";
import { useTimeout, useLoadingToast, useInterval } from "@hooks";
import toast from "react-hot-toast";

import { SelectPermitGroup } from "@/pages/User/pages/App/components";
import {
  Button,
  Checkbox,
  Collapse,
  FormControlLabel,
  Grid,
  Skeleton,
  Stack,
  TextField,
} from "@mui/material";

import { ListGroups } from "@api/v1/user/group";

import { shallow } from "zustand/shallow";
import { useAppForm, useGroup, useUser } from "@store";

interface Props {
  submitText: string;
  onSubmit: () => void;
  cancelText: string;
  onCancel: () => void;
}

export const AppForm: FC<Props> = ({
  submitText,
  onSubmit,
  cancelText,
  onCancel,
}) => {
  const setDialog = useUser((state) => state.setDialog);

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
  const [showSelectGroups, setShowSelectGroups] = useState(!permitAll);

  const groups = useGroup((state) => state.groups);
  const setGroups = useGroup((state) => state.setState("groups"));
  const [onRequestGroups, setOnRequestGroups] = useState(false);
  const [loadGroupToast, closeLoadGroupToast] = useLoadingToast();

  async function loadGroups() {
    setOnRequestGroups(true);
    try {
      const data = await ListGroups();
      setGroups(data);
      closeLoadGroupToast();
    } catch ({ msg }) {
      if (msg) loadGroupToast(msg as string);
    }
    setOnRequestGroups(false);
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

    if (callback.length <= 8) {
      setCallbackError(true);
      toast.error("回调地址不能为空");
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

  async function handleSubmit() {
    if (!(await checkForm())) return;
    onSubmit();
  }

    useTimeout(() => setShowSelectGroups(false), permitAll ? 300 : null);
    useInterval(
      loadGroups,
      !groups && !onRequestGroups && !permitAll ? 2000 : null
    );
    useEffect(() => {
      if (!permitAll) {
        setShowSelectGroups(true);
        if (onRequestGroups || groups) return;
        loadGroups();
      } else {
        closeLoadGroupToast();
      }
    }, [permitAll]);
    return (
      <>
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
              onChange={(e) => {
                if (e.target.value.indexOf("https://") !== 0) return;
                setCallback(e.target.value);
              }}
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
                <SelectPermitGroup
                  groups={groups}
                  permitGroups={permitGroups}
                  setPermitGroups={setPermitGroups}
                  fullWidth
                />
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
          <Button variant={"outlined"} onClick={onCancel}>
            {cancelText}
          </Button>
          <Button variant={"contained"} onClick={handleSubmit}>
            {submitText}
          </Button>
        </Stack>
      </>
    );
};
export default AppForm;