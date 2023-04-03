import {FC} from "react";
import logo from '@/assets/img/logo-lg.png'

import {Stack, Box} from "@mui/material";

export const ShowBar: FC = () => {
    return <Stack sx={{
        height: '100%',
        width: '100%',
        boxShadow: '0 0 15px 0 #343434',
        backgroundColor: '#343434'
    }}>
        <Box sx={{
            padding: '2.5rem 4rem',
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            '&>img': {
                maxHeight: '100%',
                maxWidth: '100%',
                width: '15rem'
            }
        }}><img src={logo} alt={"家园工作室"}/></Box>
        <Box sx={{
            flexGrow: 1,
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            boxSizing: 'border-box',
            padding: '4rem'
        }}>
            {/*不知道放什么*/}
        </Box>
    </Stack>
}
export default ShowBar
