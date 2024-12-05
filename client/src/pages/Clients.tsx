import { Visibility } from "@mui/icons-material";
import { CircularProgress, Container, IconButton, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from "@mui/material";
import { useEffect, useState } from "react";
import { fetchClients } from "../utils/requests/client";
import { ClientInfo } from "../utils/interfaces/client";

/**
 * Creates a table row for a client
 * @param param0 [ClientProps]
 * @rerurns [JSX.Element] TableRow
 */
const ClientRow = ({ name, contact_name, contact_email, number_of_vehicles }: ClientInfo, key: number) => {
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
    const [clients, setClients] = useState<ClientInfo[]>([])
    const [nextPage, setNextPage] = useState('')
    const [isLoading, setIsLoading] = useState(true)
    const [errorMessage, setErrorText] = useState('')

    const getClients = async () => {
        try {
            const response = await fetchClients()
            setClients(response.clients)
            setNextPage(response.next_page)
            setIsLoading(false)
        } catch (e: any) {
            setErrorText(e.message)
            console.error(e)
            setIsLoading(false)
        }
    }

    useEffect(() => {
        try {
            document.title = "Clients | Starter Project"
            getClients()

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
                            {clients?.map((client, i) => (
                                <ClientRow {...client} key={i} />
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>
                {isLoading && <CircularProgress style={{ marginTop: 50 }} />}
                {errorMessage && <p style={{color: 'red'}}>Error: {errorMessage}</p>}
            </Container>
        </div>
    )
};

export default Clients;