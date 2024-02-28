import { useState } from "react";

import { OverridableComponent } from "@mui/material/OverridableComponent";
import {
  Skeleton,
  Avatar as AvatarMui,
  AvatarTypeMap,
  AvatarProps,
} from "@mui/material";

export const Avatar: OverridableComponent<AvatarTypeMap> = ({
  src,
  ...props
}: AvatarProps) => {
  const [loaded, setLoaded] = useState(false);
  const [loadFailed, setLoadFailed] = useState(false);

  if (src === "" || loadFailed) return <AvatarMui {...props} />;

  return loaded ? (
    <AvatarMui {...props} src={src} />
  ) : (
    <Skeleton variant={"circular"}>
      <AvatarMui
        sx={props.sx}
        src={src}
        onLoad={() => src !== undefined && setLoaded(true)}
        onError={() => src !== undefined && setLoadFailed(true)}
      />
    </Skeleton>
  );
};
export default Avatar;
