import {FC} from "react";

import feishuLogo from '@/assets/img/login/feishu.png'
import dingLogo from '@/assets/img/login/ding.png'

import {Stack, Box, Typography, List} from "@mui/material";
import {LoginItem} from './components'

import {GetFeishuLoginUrl} from "@api/v1/login";
import toast from "react-hot-toast";

export const LoginForm:FC = ()=>{
    async function goFeishuLogin() {
        try {
            const url=await GetFeishuLoginUrl("https://dashboard.ncuos.com/")
            window.open(url,'_blank')
        } catch ({msg}) {
            if(msg)toast.error(msg as string)
        }
    }

    return <Box sx={{
        backgroundColor: '#343434',
        width: '25rem',
        maxWidth: '100%',
        overflowY: 'auto',
        padding: '2rem 3rem',
        boxShadow: '0 0 4px 0 #343434',
        borderRadius: '0.45rem',
    }}>
        <Stack sx={{
            minWidth: '100%',
            textAlign: 'center'
        }} justifyContent={'center'}>
            <Typography variant={'h4'} sx={{
                marginBottom: '2rem'
            }}>登录</Typography>

            <List>
                <LoginItem logo={feishuLogo} text={'飞书'} onClick={goFeishuLogin}/>
                <LoginItem logo={dingLogo} text={'钉钉'} onClick={()=>{}} disableDivider/>
            </List>
        </Stack>
    </Box>
}
export default LoginForm
