import { FC } from "react";

import NewAppBlock from "./NewAppBlock";
import AppControlBlock from "./AppControlBlock";
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
