import ReactDOM from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import App from './App.tsx';
import './index.css';
import 'flowbite';
import {PelamarProvider} from "./pages/Pelamar/Pelamar.tsx";

ReactDOM.createRoot(document.getElementById('root')!).render(
    <BrowserRouter>
        <PelamarProvider>
            <App />
        </PelamarProvider>
    </BrowserRouter>
);
