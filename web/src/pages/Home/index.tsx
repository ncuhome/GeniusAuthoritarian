import {FC} from "react";

import {Box, Stack} from "@mui/material"
import {ShowBar, LoginForm} from "./components";

export const Home:FC = ()=>{
    return <Stack flexDirection={'row'} sx={{
        width: '100%',
        height: '100%',
        '&>div': {
            height: '100%',
        }
    }}>
        <Box sx={{
            width: '20rem'
        }}>
            <ShowBar/>
        </Box>
        <Box sx={{
            flexGrow: 1,
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            padding: '2rem 3rem',
            boxSizing: 'border-box',
        }}>
            <LoginForm/>
        </Box>
    </Stack>
}

export default Home
