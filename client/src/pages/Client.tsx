import {
    CircularProgress, Container, Paper, Table, TableBody, TableCell,
    TableContainer, TableHead, TableRow, Grid, Card, Box, CardHeader,
    CardContent, TablePagination, TableFooter,
} from "@mui/material";
import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { ClientInfo } from "../utils/interfaces/client";
import { BasicVehicle } from "../utils/interfaces/vehicle";
import { fetchClientInfo, fetchClientVehicles } from "../utils/requests/client";


function VehicleRow({ vin, mileage, largest_weight }: BasicVehicle) {
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
    const { id } = useParams<{ id?: string }>()
    const [client, setClient] = useState<ClientInfo>()
    const [isLoadingClient, setIsLoadingClient] = useState(true)
    const [vehicles, setVehicles] = useState<BasicVehicle[]>([])
    const [isLoadingVehicles, setIsLoadingVehicles] = useState(true)
    const [page, setPage] = useState(0);
    const [rowsPerPage, setRowsPerPage] = useState(5);
    const [errorText, setErrorText] = useState("")

    const getClientInfo = async () => {
        if (id) {
            try {
                const response = await fetchClientInfo(id)
                setClient(response.client)
                setIsLoadingClient(false)
            } catch (e: any) {
                console.error(e)
                setErrorText(e.message)
                setIsLoadingClient(false)
            }
        }
    }

    const getClientVehicles = async () => {
        if (id) {
            try {
                const response = await fetchClientVehicles(id)
                setVehicles(response.vehicles)
                setIsLoadingVehicles(false)
            } catch (e: any) {
                console.error(e)
                setErrorText(e.message)
                setIsLoadingVehicles(false)
            }
        }
    }

    useEffect(() => {
        try {
            document.title = id + " | Starter Project"

            getClientInfo()
            getClientVehicles()
        } catch (error) {
            console.error(error)
        }
    }, [id])

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
                <h1>{id}</h1>
                <Box sx={{ flexGrow: 1 }}>
                    <Grid container spacing={2} columns={3}>
                        <Grid item xs={1}>
                            <Card>
                                <CardHeader title="Client Info" />
                                {isLoadingClient ? (
                                    <CardContent>
                                        <CircularProgress />
                                    </CardContent>
                                ) : (
                                    <CardContent style={{ textAlign: 'left' }}>
                                        <h4>Contact Name</h4>
                                        <p>{client?.contact_name}</p>
                                        <h4>Contact Email</h4>
                                        <p>{client?.contact_email}</p>
                                    </CardContent>
                                )}
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
                                            {isLoadingVehicles &&
                                                <TableRow>
                                                    <TableCell colSpan={4} align="center">
                                                        <CircularProgress />
                                                    </TableCell>
                                                </TableRow>
                                            }
                                            {vehicles.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage).map((vehicle, i) => {
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
                                                    count={vehicles.length || 0}
                                                    rowsPerPage={rowsPerPage}
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
                {errorText && <p style={{ color: 'red' }}>Error: {errorText}</p>}
            </Container>
        </div>
    )
};

export default Client;