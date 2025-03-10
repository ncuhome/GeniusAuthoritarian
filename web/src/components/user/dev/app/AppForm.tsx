import { FC, useEffect, useMemo, useRef, useState } from "react";
import useTimeout from "@hooks/useTimeout";
import toast from "react-hot-toast";

import SelectPermitGroup from "@components/user/dev/app/SelectPermitGroup";
import { LoadingButton } from "@mui/lab";
import {
  Button,
  Checkbox,
  Collapse,
  FormControlLabel,
  Grid,
  Skeleton,
  Stack,
  TextField,
  Alert,
  Fade,
} from "@mui/material";

import { useUserApiV1 } from "@api/v1/user/hook";

import { useShallow } from "zustand/react/shallow";
import useGroup from "@store/useGroup";
import useUser from "@store/useUser";
import { UseAppForm } from "@store/useAppForm";

interface Props {
  useForm: UseAppForm;

  submitText: string;
  onSubmit: () => void;
  cancelText: string;
  onCancel: () => void;

  loading?: boolean;
}

const hasValidScheme = (input: string) => {
  const isLocalhost = input.includes("://localhost");
  return (
    input.indexOf("https://") === 0 ||
    (isLocalhost && input.indexOf("http://") === 0)
  );
};

export const AppForm: FC<Props> = ({
  useForm,
  submitText,
  onSubmit,
  cancelText,
  onCancel,
  loading,
}) => {
  const setDialog = useUser((state) => state.setDialog);

  const nameInput = useRef<HTMLInputElement | null>(null);
  const callbackInput = useRef<HTMLInputElement | null>(null);

  const [name, callback, permitAll, permitGroups, nameError, callbackError] =
    useForm(
      useShallow((state) => [
        state.name,
        state.callback,
        state.permitAll,
        state.permitGroups,
        state.nameError,
        state.callbackError,
      ]),
    );
  const [showSelectGroups, setShowSelectGroups] = useState(!permitAll);

  const groups = useGroup((state) => state.groups);

  const notOnlyCenterSelected = useMemo(
    () =>
      !permitAll &&
      permitGroups &&
      permitGroups.length !== 1 &&
      permitGroups.findIndex((group) => group.name === "中心") !== -1,
    [permitAll, permitGroups],
  );

  useUserApiV1<User.Group[]>(!permitAll ? "group/list" : null, {
    immutable: true,
    enableLoading: true,
    onSuccess: (data) => useGroup.setState({ groups: data }),
  });

  async function checkForm(): Promise<boolean> {
    if (!name) {
      useForm.setState({ nameError: true });
      toast.error("应用名称不能为空");
      nameInput.current!.focus();
      return false;
    } else {
      useForm.setState({ nameError: false });
    }

    if (callback.length <= 8) {
      useForm.setState({ callbackError: true });
      toast.error("回调地址不能为空");
      callbackInput.current!.focus();
      return false;
    } else {
      useForm.setState({ callbackError: false });
    }

    if (!permitAll && (!permitGroups || permitGroups.length === 0)) {
      const yes = await setDialog({
        title: "警告",
        content: "您没有授权任何身份组使用，你确定要继续吗",
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

  useEffect(() => {
    if (!permitAll) setShowSelectGroups(true);
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
            onChange={(e) => useForm.setState({ name: e.target.value })}
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
              const text = e.target.value;
              if (!hasValidScheme(text)) return;
              useForm.setState({ callback: text });
            }}
          />
        </Grid>
        <Grid item xs={12}>
          <Stack
            flexDirection={"row"}
            alignItems={"center"}
            justifyContent={"space-between"}
          >
            <FormControlLabel
              control={
                <Checkbox
                  checked={permitAll}
                  onChange={(e) =>
                    useForm.setState({ permitAll: e.target.checked })
                  }
                />
              }
              label="允许所有成员使用"
            />
            <Fade in={notOnlyCenterSelected}>
              <Alert severity={"info"} sx={{ flexGrow: 1 }}>
                仅供中心组使用才需要选择中心组
              </Alert>
            </Fade>
          </Stack>
        </Grid>
        <Grid
          item
          xs={12}
          mb={2}
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
                setPermitGroups={(data) =>
                  useForm.setState({ permitGroups: data })
                }
                fullWidth
                sx={{ pb: "1rem" }}
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
          "&>button": {
            marginLeft: "0.8rem",
          },
        }}
      >
        <Button variant={"outlined"} onClick={onCancel}>
          {cancelText}
        </Button>
        <LoadingButton
          variant={"contained"}
          loading={loading}
          onClick={handleSubmit}
        >
          {submitText}
        </LoadingButton>
      </Stack>
    </>
  );
};
export default AppForm;
