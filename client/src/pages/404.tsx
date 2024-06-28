/* 
* @file 404.tsx
* @author Byron Ojua-Nice
* @version 1.0
* 
* @section DESCRIPTION
* 
* This file contains the code for the 404 page. This page is displayed when a user navigates to a page that does not exist.
*/

import { Card, CardActionArea, Container } from "@mui/material";
import React, { useEffect } from "react";

/**
 * Home page for the application
 * @returns [JSX.Element] Home
 */
const Error404 = () => {
    useEffect(() => {
        document.title = "404 | Starter Project"
    }, [])

    return (
        <div className="App">
            <Container>
                <h1>404</h1>
                <h2>Page not found</h2>
                <p>The page you are looking for does not exist.</p>
                
            </Container>
        </div>
    )
};

export default Error404;