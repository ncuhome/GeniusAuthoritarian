import { FC, useState } from "react";
import toast from "react-hot-toast";

import Block from "@components/user/Block";
import AppForm from "@components/user/dev/app/AppForm";
import { Box, Stack, Typography } from "@mui/material";

import { apiV1User } from "@api/v1/user/base";

import { useShallow } from "zustand/react/shallow";
import { useAppForm } from "@store/useAppForm";
import useUser from "@store/useUser";

export const NewAppBlock: FC = () => {
  const apps = useUser((state) => state.apps);
  const setApps = useUser((state) => state.setApps);
  const setDialog = useUser((state) => state.setDialog);

  const [name, callback, permitAll, permitGroups] = useAppForm(
    useShallow((state) => [
      state.name,
      state.callback,
      state.permitAll,
      state.permitGroups,
    ]),
  );
  const resetForm = useAppForm((state) => state.reset);

  const [onCreateApp, setOnCreateApp] = useState(false);

  async function createApp() {
    setOnCreateApp(true);
    try {
      const {
        data: { data },
      } = await apiV1User.post("dev/app/", {
        name,
        callback,
        permitAll,
        permitGroups: permitGroups?.map((group) => group.id),
      });
      setApps([data, ...apps!]);
      resetForm();
      setDialog({
        title: "AppSecret 仅显示一次，请妥善保管",
        content: (
          <Stack>
            <Typography>AppCode: </Typography>
            <pre>{data.appCode}</pre>
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
    } catch (err) {
      if (err instanceof Error) toast.error(err.message);
    }
    setOnCreateApp(false);
  }

  return (
    <Block title={"New"}>
      <AppForm
        useForm={useAppForm}
        submitText={"创建应用"}
        onSubmit={createApp}
        cancelText={"重置"}
        onCancel={resetForm}
        loading={onCreateApp}
      />
    </Block>
  );
};
export default NewAppBlock;
