import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.scss';
import reportWebVitals from './reportWebVitals';
import { Route, Routes } from "react-router";
import { BrowserRouter } from "react-router-dom"
import { Home, Auth } from "./components";

const root = ReactDOM.createRoot(
  document.getElementById("root") as HTMLElement
);
root.render(
    <React.StrictMode>
        <BrowserRouter>
            <Routes>
                <Route path="/" element={ <Home /> } />
                <Route path="/auth" element={ <Auth /> } />
            </Routes>
        </BrowserRouter>

    </React.StrictMode>
);
reportWebVitals();
