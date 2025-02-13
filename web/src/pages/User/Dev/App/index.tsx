import { FC } from "react";

import NewAppBlock from "./NewAppBlock";
import AppControlBlock from "./AppControlBlock";

export const App: FC = () => {
  return (
    <>
      <NewAppBlock />
      <AppControlBlock />
    </>
  );
};
export default App;
