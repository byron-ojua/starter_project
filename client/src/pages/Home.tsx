/* 
* @file Home.tsx
* @author Byron Ojua-Nice
* @version 1.0
* 
* @section DESCRIPTION
* 
* This file contains the code for the Home page. This page is the landing page for the application.
*/

import { Business, LocalShipping } from "@mui/icons-material";
import { Card, CardActionArea, Container } from "@mui/material";
import React, { useEffect } from "react";

/**
 * Home page for the application
 * @returns [JSX.Element] Home
 */
const Home = () => {
    function selectVehicle() {
        var vehicle = prompt("Enter vehicle identification number (VIN):");
        if (vehicle != null) {
            window.location.href = "/vehicles/" + vehicle
        }
    }

    useEffect(() => {
        document.title = "Home | Starter Project"
    }, [])

    return (
        <div className="App">
            <Container>
                <h1>Home</h1>
                <div style={{ display: 'flex', flexDirection: 'row', gap: 20, justifyContent: 'center' }}>
                    <Card sx={{ width: 200, height: 200}}>
                        <CardActionArea href="/clients" style={{ width: '100%', height: '100%', alignContent: 'center' }}>
                            <Business sx={{ fontSize: 100 }} />
                            <h2>Clients</h2>
                        </CardActionArea>
                    </Card>
                    <Card sx={{ width: 200, height: 200}}>
                        <CardActionArea onClick={selectVehicle} style={{ width: '100%', height: '100%', alignContent: 'center' }}>
                            <LocalShipping sx={{ fontSize: 100 }} />
                            <h2>Vehicles</h2>
                        </CardActionArea>
                    </Card>
                </div>
            </Container>
        </div>
    )
};

export default Home;