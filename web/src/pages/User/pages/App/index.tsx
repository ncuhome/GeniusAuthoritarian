import { FC } from "react";

import { AppFormBlock, AppControlBlock } from "./components";
import { Container } from "@mui/material";

export const App: FC = () => {
  return (
    <Container>
      <AppFormBlock />
      <AppControlBlock />
    </Container>
  );
};
export default App;
