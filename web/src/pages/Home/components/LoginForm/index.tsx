import {FC} from "react";

import {Stack, Box, Typography} from "@mui/material";

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
            <Typography variant={'h4'}>登录</Typography>
        </Stack>
    </Box>
}
export default LoginForm
