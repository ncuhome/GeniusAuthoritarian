import { FC } from "react";

import { Block } from "@/pages/User/components";
import { Container } from "@mui/material";

export const App: FC = () => {
  return (
    <Container>
      <Block title={"New"}></Block>
      <Block title={"App"}></Block>
    </Container>
  );
};
export default App;
