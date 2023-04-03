import {FC} from "react";

import {Box, Stack} from "@mui/material"
import {ShowBar} from "./components";

export const Home:FC = ()=>{
    return <Stack flexDirection={'row'} sx={{
        width: '100%',
        height: '100%',
        '&>div': {
            height: '100%',
        }
    }}>
        <Box sx={{
            width: '40%'
        }}>
            <ShowBar/>
        </Box>
        <Box sx={{
            flexGrow: 1
        }}>

        </Box>
    </Stack>
}

export default Home
