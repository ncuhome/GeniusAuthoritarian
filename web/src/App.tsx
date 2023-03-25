import {BrowserRouter, Route, Routes} from "react-router-dom";
import {Toaster} from "react-hot-toast";

import {Home} from './pages'

export default function App() {
    return <>
        <Toaster toastOptions={{
            style: {
                borderRadius: '20px',
                background: '#2f2f2f',
                color: '#fff',
            },
        }}/>
        <BrowserRouter>
            <Routes>
                <Route index element={<Home/>} />
            </Routes>
        </BrowserRouter>
    </>
}
