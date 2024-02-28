import { FC } from "react";

import Block from "@components/user/Block";
import { Container } from "@mui/material";

import useAdminData from "@store/useAdminData";

import { useUserApiV1 } from "@api/v1/user/hook";

export const Admin: FC = () => {
  const setLoginData = useAdminData((state) => state.setState("login"));

  useUserApiV1<User.LoginRecordAdminView[]>("admin/data/login", {
    enableLoading: true,
    onSuccess: (data) => setLoginData(data),
  });

  return (
    <Container>
      <Block></Block>

      <Block></Block>

      <Block></Block>
    </Container>
  );
};
export default Admin;
