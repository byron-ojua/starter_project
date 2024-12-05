import { BasicVehicle } from "./vehicle"

export interface ResponseGetClients {
    clients: ClientInfo[] 
    next_page: string
}

// Struct to match API ClientWithVehicles struct
export interface ClientInfo {
    name: string,
    contact_name: string,
    contact_email: string
    number_of_vehicles: number
}

export interface ResponseGetClient {
    client: ClientInfo
    next_page: string
}

export interface ResponseGetClientVehicles {
    vehicles: BasicVehicle[]
    next_page: string
}