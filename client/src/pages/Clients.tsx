/* 
* @file Clients.tsx
* @author Byron Ojua-Nice
* @version 1.0
* 
* @section DESCRIPTION
* 
* This file contains the code for the Clients page. This page displays info about the clients in the database.
*/

import { Visibility } from "@mui/icons-material";
import { CircularProgress, Container, IconButton, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from "@mui/material";
import axios, { AxiosResponse } from "axios";
import React, { useEffect, useState } from "react";

// Struct to match API ClientWithVehicles struct
type ClientProps = {
    name: string,
    contact_name: string,
    contact_email: string
    number_of_vehicles: number
}

/**
 * Creates a table row for a client
 * @param param0 [ClientProps]
 * @rerurns [JSX.Element] TableRow
 */
const ClientRow = ({ name, contact_name, contact_email, number_of_vehicles }: ClientProps, key: number) => {
    var client_url = '/clients/' + name

    return (
        <TableRow>
            <TableCell>{name}</TableCell>
            <TableCell>{contact_name}</TableCell>
            <TableCell>{contact_email}</TableCell>
            <TableCell>{number_of_vehicles}</TableCell>
            <TableCell>
                <IconButton href={client_url} size="small">
                    <Visibility />
                </IconButton>
            </TableCell>
        </TableRow>
    )
}

/**
 * Pages that displays clients in the database
 * @returns [JSX.Element] Clients
 */
const Clients = () => {
    const [clients, setClients] = useState<ClientProps[]>([])
    const [is_loading, setIsLoading] = useState(true)

    useEffect(() => {
        try {
            document.title = "Clients | Starter Project"
            axios.get('http://localhost:8080/clients')
                .then((res: AxiosResponse<ClientProps[]>) => {
                    setClients(res.data.sort((a, b) => a.name.localeCompare(b.name)))
                    setIsLoading(false)
                }).catch((error) => {
                    console.error(error)
                });

        } catch (error) {
            console.error(error)
        }
    }, [])

    return (
        <div className="App">
            <h1>Clients</h1>
            <Container>
                <TableContainer component={Paper}>
                    <Table aria-label="Clients table" sx={{ minWidth: 800 }}>
                        <TableHead>
                            <TableRow>
                                <TableCell>Client Name</TableCell>
                                <TableCell>Contact Name</TableCell>
                                <TableCell>Contact Email</TableCell>
                                <TableCell>Number of Vehicles</TableCell>
                                <TableCell></TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {clients?.map((client, i) => {
                                return (
                                    <ClientRow {...client} key={i} />
                                )
                            })}
                        </TableBody>
                    </Table>
                </TableContainer>
                {is_loading && <CircularProgress style={{ marginTop: 50 }} />}
            </Container>
        </div>
    )
};

export default Clients;