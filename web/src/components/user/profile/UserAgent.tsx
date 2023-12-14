import { FC, ReactNode } from "react";
import parser from "ua-parser-js";

import { Stack, Typography } from "@mui/material";
import {
  DesktopWindowsOutlined,
  PhoneAndroidOutlined,
  TabletAndroidOutlined,
} from "@mui/icons-material";
import {
  FaChrome,
  FaEarthAmericas,
  FaEdge,
  FaFirefoxBrowser,
  FaSafari,
} from "react-icons/fa6";

interface Props {
  useragent: string;
}

const UserAgent: FC<Props> = ({ useragent }) => {
  const ua = parser(useragent);

  const renderDevice = () => {
    let name: string = "";
    if (ua.os.name) {
      name = ` ${ua.os.name}`;
      if (ua.os.version) {
        name += ua.os.version;
      }
    }
    if (ua.device.model) {
      name += ` ${ua.device.model}`;
    }
    if (ua.device.vendor) {
      name += ` ${ua.device.vendor}`;
    }

    let icon: ReactNode;
    switch (ua.device.type) {
      case "mobile":
        icon = <PhoneAndroidOutlined fontSize={"small"} />;
        break;
      case "tablet":
        icon = <TabletAndroidOutlined fontSize={"small"} />;
        break;
      default:
        icon = <DesktopWindowsOutlined fontSize={"small"} />;
    }

    return (
      <>
        {icon} <Typography mr={0}>{name}</Typography>
      </>
    );
  };

  const renderBrowser = () => {
    if (!ua.browser.name) return undefined;

    let icon: ReactNode;
    switch (ua.browser.name) {
      case "Chrome":
        icon = <FaChrome fontSize={"1.1rem"} />;
        break;
      case "Firefox":
        icon = <FaFirefoxBrowser fontSize={"1.1rem"} />;
        break;
      case "Safari":
        icon = <FaSafari fontSize={"1.1rem"} />;
        break;
      case "Edge":
        icon = <FaEdge fontSize={"1.1rem"} />;
        break;
      default:
        icon = <FaEarthAmericas fontSize={"1.1rem"} />;
    }

    return (
      <>
        {icon} <Typography>{ua.browser.name}</Typography>
      </>
    );
  };

  return (
    <Stack
      direction={"row"}
      alignItems={"center"}
      sx={{
        "&>*": {
          marginRight: 1,
        },
      }}
    >
      {renderBrowser()}
      {renderDevice()}
    </Stack>
  );
};
export default UserAgent;
