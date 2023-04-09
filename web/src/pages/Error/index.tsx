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
            "&>*": {
                marginBottom: '1rem!important'
            },
        }}>
            <Box sx={{
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                "& svg": {
                    borderStyle: 'solid',
                    borderRadius: '50%',
                    fontSize: '8rem',
                    animation: 'error-page-fork .5s ease',
                    animationFillMode: 'forwards',
                    boxSizing: 'border-box'
                }
            }}>
                <ClearRounded/>
            </Box>
            <Typography variant={'h4'} sx={{
                fontWeight: 600,
                letterSpacing: '0.25rem'
            }}>{title}</Typography>
            {content?<Typography sx={{
                color: '#999',
                wordBreak: 'break-all',
            }}>{content}</Typography>:undefined}
        </Box>
    </Box>
}

export default Error
