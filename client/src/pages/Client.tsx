/* 
* @file Client.tsx
* @author Byron Ojua-Nice
* @version 1.0
* 
* @section DESCRIPTION
* 
* This file contains the code for the Client page. This page displays info about the client and their vehicles.
*/

import {
    CircularProgress, Container, Paper, Table, TableBody, TableCell,
    TableContainer, TableHead, TableRow, Grid, Card, Box, CardHeader,
    CardContent, TablePagination, TableFooter,
} from "@mui/material";
import axios, { AxiosResponse } from "axios";
import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

// Struct to match API ClientWithVehicles struct
type ClientProps = {
    name: string,
    contact_name: string,
    contact_email: string
    number_of_vehicles: number
}

// Struct to match API ClientVehicle struct
type VehicleProps = {
    vin: string,
    mileage: number,
    largest_weight: number,
}

// Struct to match API ClientVehicles struct
type Vehicles = {
    name: string,
    vehicles: VehicleProps[]
}

/**
 * Creates a table row for a vehicle
 * @param param0 [VehicleProps]
 * @rerurns [JSX.Element] TableRow
 */
function VehicleRow({ vin, mileage, largest_weight }: VehicleProps) {
    return (
        <TableRow hover>
            <TableCell>{vin}</TableCell>
            <TableCell>{mileage}</TableCell>
            <TableCell>{largest_weight}</TableCell>
            <TableCell align="right">
                <a href={"/vehicles/" + vin}>View</a>
            </TableCell>
        </TableRow>
    )
}

/**
 * Page that displays info about the client and their vehicles.
 * Uses id from the URL /clients/:id
 * @returns [JSX.Element] Client
 */
const Client = () => {
    const params = useParams()
    const [client, setClient] = useState<ClientProps>()
    const [vehicles, setVehicles] = useState<Vehicles>()
    const [is_loading_vehicles, setIsLoadingVehicles] = useState(true)
    const [page, setPage] = useState(0);
    const [rows_per_page, setRowsPerPage] = useState(5);

    useEffect(() => {
        try {
            document.title = params.id + " | Starter Project"
            axios.get('http://localhost:8080/clients/' + params.id)
                .then((res: AxiosResponse<ClientProps>) => {
                    setClient(res.data)
                }).catch((error) => {
                    console.error(error)
                });
            axios.get('http://localhost:8080/clients/' + params.id + '/vehicles')
                .then((res: AxiosResponse<Vehicles>) => {
                    setVehicles({
                        name: res.data.name,
                        vehicles: res.data.vehicles.sort(
                            (a: VehicleProps, b: VehicleProps) => a.vin.localeCompare(b.vin))
                    })
                    setIsLoadingVehicles(false)
                }).catch((error) => {
                    console.error(error)
                });
        } catch (error) {
            console.error(error)
        }
    }, [params.id])

    // Pagination functions
    const handleChangePage = (
        event: React.MouseEvent<HTMLButtonElement> | null,
        new_page: number,
    ) => {
        setPage(new_page);
    };

    const handleChangeRowsPerPage = (
        event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
    ) => {
        setRowsPerPage(parseInt(event.target.value, 10));
        setPage(0);
    };

    return (
        <div className="App">
            <Container>
                <h1>{client?.name}</h1>
                <Box sx={{ flexGrow: 1 }}>
                    <Grid container spacing={2} columns={3}>
                        <Grid item xs={1}>
                            <Card>
                                <CardHeader title="Client Info" />
                                <CardContent style={{ textAlign: 'left' }}>
                                    <h4>Contact Name</h4>
                                    <p>{client?.contact_name}</p>
                                    <h4>Contact Email</h4>
                                    <p>{client?.contact_email}</p>
                                </CardContent>
                            </Card>
                        </Grid>
                        <Grid item xs={2}>
                            <Card sx={{ justifyContent: 'start' }}>
                                <CardHeader title="Vehicles" />
                                <TableContainer component={Paper}>
                                    <Table aria-label="Vehicles table">
                                        <TableHead>
                                            <TableRow>
                                                <TableCell style={{ fontWeight: 'bold' }}>VIN</TableCell>
                                                <TableCell style={{ fontWeight: 'bold' }}>Mileage</TableCell>
                                                <TableCell style={{ fontWeight: 'bold' }}>Largest Weight</TableCell>
                                                <TableCell></TableCell>
                                            </TableRow>
                                        </TableHead>
                                        <TableBody>
                                            {is_loading_vehicles &&
                                                <TableRow>
                                                    <TableCell colSpan={4} align="center">
                                                        <CircularProgress />
                                                    </TableCell>
                                                </TableRow>
                                            }
                                            {vehicles?.vehicles.slice(page * rows_per_page, page * rows_per_page + rows_per_page).map((vehicle, i) => {
                                                return (
                                                    <VehicleRow {...vehicle} key={i} />
                                                )
                                            })}
                                        </TableBody>
                                        <TableFooter>
                                            <TableRow>
                                                <TablePagination
                                                    rowsPerPageOptions={[3, 5, 10, 25]}
                                                    colSpan={3}
                                                    count={vehicles?.vehicles.length || 0}
                                                    rowsPerPage={rows_per_page}
                                                    page={page}
                                                    onPageChange={handleChangePage}
                                                    onRowsPerPageChange={handleChangeRowsPerPage}
                                                />
                                            </TableRow>
                                        </TableFooter>
                                    </Table>
                                </TableContainer>
                            </Card>
                        </Grid>
                    </Grid>
                </Box>
            </Container>
        </div>
    )
};

export default Client;