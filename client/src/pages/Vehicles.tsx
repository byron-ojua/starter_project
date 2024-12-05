import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import {
    CircularProgress, Container, Paper, Table, TableBody, TableCell,
    TableContainer, TableHead, TableRow, Grid, Card, Box, CardHeader,
    CardContent, TablePagination, TableFooter,
} from "@mui/material";
import { fetchVehicle } from "../utils/requests/vehicle";
import { VehicleInfo } from "../utils/interfaces/vehicle";

/**
 * Creates a table row for a weight
 * @param weight [number]
 * @returns [JSX.Element] TableRow
 */
function WeightRow({ weight }: { weight: number }) {
    return (
        <TableRow hover>
            <TableCell>{weight}</TableCell>
        </TableRow>
    )
}

/**
 * Page that displays info about a vehicle.
 * Uses id from the URL /vehicles/:id
 * @returns [JSX.Element] Vehicles
 */
const Vehicles = () => {
    const { id } = useParams<{ id?: string }>()
    const [vehicle, setVehicle] = useState<VehicleInfo>()
    const [is_loading, setIsLoading] = useState(true)
    const [page, setPage] = useState(0);
    const [rows_per_page, setRowsPerPage] = useState(5);
    const [error_text, setErrorText] = useState("")

    const getVehicle = async () => {
        if (id) {
            try {
                const response = await fetchVehicle(id)
                setVehicle(response.vehicle)
                setIsLoading(false)
            } catch (e: any) {
                console.error(e)
                setErrorText(e.message)
                setIsLoading(false)
            }
        }
    }

    useEffect(() => {
        document.title = id + " | Vehicles | Starter Project"
        getVehicle()
    }, [id])

    // Pagination functions
    const handleChangePage = (
        event: React.MouseEvent<HTMLButtonElement> | null,
        newPage: number,
    ) => {
        setPage(newPage);
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
                                <CardHeader title="Vehicle Info" />
                                <CardContent style={{ textAlign: 'left' }}>
                                    <h4>VIN</h4>
                                    <p>{id}</p>
                                    <h4>Client Name</h4>
                                    <p>{vehicle?.client_name}</p>
                                    <h4>Contact Name</h4>
                                    <p>{vehicle?.contact_name}</p>
                                    <h4>Contact Email</h4>
                                    <p>{vehicle?.contact_email}</p>
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
                                                <TableCell style={{ fontWeight: 'bold' }}>Weight</TableCell>
                                            </TableRow>
                                        </TableHead>
                                        <TableBody>
                                            {is_loading &&
                                                <TableRow>
                                                    <TableCell colSpan={2} align="center">
                                                        <CircularProgress />
                                                    </TableCell>
                                                </TableRow>
                                            }
                                            {vehicle?.weights.slice(page * rows_per_page, page * rows_per_page + rows_per_page).map((weight, i) => {
                                                return (
                                                    <WeightRow key={i} weight={weight} />
                                                )
                                            })}
                                        </TableBody>
                                        <TableFooter>
                                            <TableRow>
                                                <TablePagination
                                                    rowsPerPageOptions={[3, 5, 10, 25]}
                                                    colSpan={3}
                                                    count={vehicle?.weights.length || 0}
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
                {error_text && <p style={{color: 'red'}}>Error: {error_text}</p>}
            </Container>
        </div>
    )
};

export default Vehicles;