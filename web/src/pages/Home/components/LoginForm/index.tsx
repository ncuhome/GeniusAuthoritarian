import {FC} from "react";

import feishu from '@/assets/img/login/feishu.png'

import {Stack, Box, Typography, ListItem, ListItemButton, ListItemText, ListItemIcon,
    List} from "@mui/material";

export const LoginForm:FC = ()=>{
    return <Box sx={{
        backgroundColor: '#343434',
        width: '25rem',
        maxWidth: '100%',
        overflowY: 'auto',
        padding: '2rem 3rem',
        boxShadow: '0 0 10px 0 #343434',
        borderRadius: '1.2rem',
    }}>
        <Stack sx={{
            minWidth: '100%',
            textAlign: 'center'
        }} justifyContent={'center'}>
            <Typography variant={'h4'} sx={{
                marginBottom: '2rem'
            }}>登录</Typography>

            <List>
                <ListItem disablePadding divider>
                    <ListItemButton>
                        <ListItemIcon><img style={{
                            width: '1.8rem'
                        }} src={feishu}/></ListItemIcon>
                        <ListItemText primary={`飞书`} />
                    </ListItemButton>
                </ListItem>
            </List>
        </Stack>
    </Box>
}
export default LoginForm
