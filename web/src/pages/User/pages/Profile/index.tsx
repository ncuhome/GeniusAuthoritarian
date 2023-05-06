import { FC } from "react";
import { useInterval } from "@hooks";
import toast from "react-hot-toast";

import { Container, Box } from "@mui/material";

import { GetUserProfile } from "@api/v1/user/profile";

import { shallow } from "zustand/shallow";
import { useUser } from "@store";

export const Profile: FC = () => {
  const [profile] = useUser((state) => [state.profile], shallow);
  const [setProfile] = useUser((state) => [state.setState("profile")], shallow);

  async function loadProfile() {
    try {
      const data = await GetUserProfile();
      setProfile(data);
    } catch ({ msg }) {
      if (msg) toast.error(msg as string);
    }
  }

  useInterval(loadProfile, profile ? null : 2000);

  return (
    <Container>
      <Box>{profile ? profile.user.name : "Loading..."}</Box>
    </Container>
  );
};
export default Profile;
