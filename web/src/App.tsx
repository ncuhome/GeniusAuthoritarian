import {BrowserRouter, Route, Routes} from "react-router-dom";
import {Toaster} from "react-hot-toast";

import {Home} from './pages'

import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
const darkTheme = createTheme({
    palette: {
        mode: 'dark',
    },
});

export default function App() {
    return <>
        <Toaster toastOptions={{
            style: {
                borderRadius: '20px',
                background: '#2f2f2f',
                color: '#fff',
            },
        }}/>
        <ThemeProvider theme={darkTheme}>
            <BrowserRouter>
                <Routes>
                    <Route index element={<Home/>} />
                </Routes>
            </BrowserRouter>
        </ThemeProvider>
    </>
}
