import {FC} from "react";
import {useQuery} from "@hooks";
import './styles.css'

import {Box, Typography} from "@mui/material";
import {ClearRounded} from '@mui/icons-material';

export const Error:FC = ()=>{
    const [title]=useQuery('title', '未知错误')
    const [content]=useQuery('content','')

    return <Box sx={{
        height: '100%',
        width: '100%',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center'
    }}>
        <Box sx={{
            display: 'flex',
            flexDirection: 'column',
            textAlign: 'center',
            justifyContent: 'center',
        }}>
            <Box sx={{
                height: '15.5rem',
                "& svg": {
                    borderStyle: 'solid',
                    borderRadius: '50%',
                    fontSize: '12rem',
                    animation: 'error-page-fork .5s ease',
                    animationFillMode: 'forwards',
                }
            }}>
                <ClearRounded/>
            </Box>
            <Typography variant={'h4'}>{title}</Typography>
            {content?<Typography sx={{
                color: '#999',
                marginTop: '0.4rem',
                wordBreak: 'break-all',
            }}>{content}</Typography>:undefined}
        </Box>
    </Box>
}

export default Error
