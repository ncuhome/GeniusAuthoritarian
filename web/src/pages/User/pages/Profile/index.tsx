import {FC, useState} from "react";
import {useInterval} from "@hooks";
import toast from "react-hot-toast";

import {Container} from "@mui/material";

import {GetUserProfile, UserProfile} from "@api/v1/user/profile";

export const Profile: FC = () => {
  const [profile, setProfile] = useState<UserProfile | undefined>(undefined);

  async function loadProfile() {
    try {
      const data = await GetUserProfile();
      setProfile(data)
    } catch ({msg}) {
      if (msg) toast.error(msg as string);
    }
  }

  useInterval(loadProfile, profile ? null : 2000);

  return <Container></Container>;
};
export default Profile;
