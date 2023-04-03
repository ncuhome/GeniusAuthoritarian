import {FC} from "react";

import feishuLogo from '@/assets/img/login/feishu.png'
import dingLogo from '@/assets/img/login/ding.png'

import {Stack, Box, Typography, List} from "@mui/material";
import {LoginItem} from './components'

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
                <LoginItem logo={feishuLogo} text={'飞书'}/>
                <LoginItem logo={dingLogo} text={'钉钉'} disableDivider/>
            </List>
        </Stack>
    </Box>
}
export default LoginForm
