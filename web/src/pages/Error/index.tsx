import {FC} from "react";
import './styles.css'

import {Box} from "@mui/material";
import {ClearRounded} from '@mui/icons-material';

export const Error:FC = ()=>{
    return <Box sx={{
        height: '100%',
        width: '100%',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center'
    }}>
        <Box sx={{
            display: 'flex',
            "&>svg": {
                fontSize: '25rem',
                animation: 'error-page-fork .5s ease',
                animationFillMode: 'forwards'
            }
        }}>
            <ClearRounded/>
        </Box>
    </Box>
}

export default Error
