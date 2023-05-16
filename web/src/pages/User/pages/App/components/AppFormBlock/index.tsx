import { FC } from "react";
import toast from "react-hot-toast";

import { AppForm } from "@/pages/User/pages/App/components";
import { Block } from "@/pages/User/components";
import { Box, Stack, Typography } from "@mui/material";

import { ApplyApp } from "@api/v1/user/app";

import { shallow } from "zustand/shallow";
import { useUser, useAppForm, useGroup } from "@store";

export const AppFormBlock: FC = () => {
  const apps = useUser((state) => state.apps);
  const setApps = useUser((state) => state.setState("apps"));
  const setDialog = useUser((state) => state.setDialog);

  const [name, callback, permitAll, permitGroups] = useAppForm(
    (state) => [
      state.name,
      state.callback,
      state.permitAll,
      state.permitGroups,
    ],
    shallow
  );
  const resetForm = useAppForm((state) => state.reset);

  async function createApp() {
    try {
      const data = await ApplyApp(
        name,
        callback,
        permitAll,
        permitGroups?.map((group) => group.id)
      );
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

  return (
    <Block title={"New"}>
      <AppForm
        submitText={"创建应用"}
        onSubmit={createApp}
        cancelText={"重置"}
        onCancel={resetForm}
      />
    </Block>
  );
};
export default AppFormBlock;
