import { FC } from "react";

import { NewAppBlock, AppControlBlock } from "./components";
import { Container } from "@mui/material";

export const App: FC = () => {
  return (
    <Container>
      <NewAppBlock />
      <AppControlBlock />
    </Container>
  );
};
export default App;
