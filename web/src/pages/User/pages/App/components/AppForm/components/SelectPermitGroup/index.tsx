import { FC } from "react";

import {
  Checkbox,
  FormControl,
  InputLabel,
  ListItemText,
  MenuItem,
  OutlinedInput,
  Select,
  FormControlProps,
} from "@mui/material";

import { Group } from "@api/v1/user/group";

interface Props extends FormControlProps {
  groups: Group[];
  permitGroups?: Group[];
  setPermitGroups: (data?: Group[]) => void;
}

// 需要额外手动加载 group 列表
export const SelectPermitGroup: FC<Props> = ({
  groups,
  permitGroups,
  setPermitGroups,
  ...rest
}) => {
  return (
    <FormControl {...rest}>
      <InputLabel>授权身份组</InputLabel>
      <Select
        multiple
        value={permitGroups?.map((group) => group.name) || []}
        input={<OutlinedInput label="授权身份组" />}
        renderValue={(selected) => selected.join(", ")}
      >
        {groups.map((group) => {
          const checked = (permitGroups?.map(group=>group.id).indexOf(group.id) ?? -2) > -1;
          return (
            <MenuItem
              key={group.id}
              value={group.name}
              onClick={() => {
                if (checked) {
                  setPermitGroups(
                    permitGroups?.filter((g) => g.id !== group.id)
                  );
                } else {
                  setPermitGroups([...(permitGroups ?? []), group]);
                }
              }}
            >
              <Checkbox checked={checked} />
              <ListItemText primary={group.name} />
            </MenuItem>
          );
        })}
      </Select>
    </FormControl>
  );
};
export default SelectPermitGroup;
