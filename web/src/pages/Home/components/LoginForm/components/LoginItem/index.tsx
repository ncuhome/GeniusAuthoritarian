import {FC} from "react";

import {ListItem, ListItemButton, ListItemIcon, ListItemText} from "@mui/material";

interface Props {
    logo:string
    text:string
    disableDivider?:boolean
    onClick: ()=>void
}

export const LoginItem:FC<Props> = ({logo,text,disableDivider, onClick})=>{
    return <ListItem disablePadding divider={!disableDivider}>
    <ListItemButton onClick={onClick}>
        <ListItemIcon><img style={{
            width: '1.8rem'
        }} src={logo} alt={text}/></ListItemIcon>
        <ListItemText primary={text} />
    </ListItemButton>
</ListItem>
}
export default LoginItem
