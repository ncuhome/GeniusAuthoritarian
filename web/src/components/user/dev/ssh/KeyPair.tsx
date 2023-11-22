import { FC } from "react";

import {
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  Stack,
  TextField,
  StackProps,
} from "@mui/material";

interface Props extends StackProps {
  keyMode: User.SSH.KeyMode;
  keys: User.SSH.Keys;

  onSetKeyMode: (mode: User.SSH.KeyMode) => void;
}

export const KeyPair: FC<Props> = ({
  keyMode,
  onSetKeyMode,
  keys,
  ...rest
}) => {
  const renderTextField = (
    label: string,
    value: string,
    multiline?: boolean,
  ) => {
    return (
      <TextField
        label={label}
        fullWidth
        multiline={multiline}
        value={value}
        InputProps={{
          readOnly: true,
        }}
        onClick={(e: any) => e.target.select()}
      />
    );
  };

  return (
    <Stack {...rest}>
      <FormControl variant={"outlined"} fullWidth>
        <InputLabel id={"key-mode-select"}>格式</InputLabel>
        <Select
          labelId={"key-mode-select"}
          label={"格式"}
          value={keyMode}
          defaultValue={"ssh"}
          onChange={(e) => onSetKeyMode(e.target.value as "pem" | "ssh")}
        >
          <MenuItem value={"ssh"}>SSH</MenuItem>
          <MenuItem value={"pem"}>PEM</MenuItem>
        </Select>
      </FormControl>

      {renderTextField("用户名", keys.username)}
      {renderTextField("公钥", keys[keyMode].public.trimEnd(), true)}
      {renderTextField("密钥", keys[keyMode].private.trimEnd(), true)}
    </Stack>
  );
};

export default KeyPair;
