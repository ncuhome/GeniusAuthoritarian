import { FC } from "react";

import Block from "@components/user/Block";
import { Container } from "@mui/material";

import useLoginData from "@hooks/data/useLoginData";

export const Admin: FC = () => {
    useLoginData("week");

  return (
    <Container>
      <Block></Block>

      <Block></Block>

      <Block></Block>
    </Container>
  );
};
export default Admin;
