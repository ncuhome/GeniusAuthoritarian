import { FC } from "react";
import { Outlet } from "react-router";

import { Toaster } from "react-hot-toast";

export const Authoritarian: FC = () => {
  return (
    <>
      <Toaster
        toastOptions={{
          style: {
            borderRadius: "20px",
            background: "#2f2f2f",
            color: "#fff",
          },
        }}
      />
      <Outlet />
    </>
  );
};
export default Authoritarian;
