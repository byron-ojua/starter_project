import { ResponseGetClient, ResponseGetClients, ResponseGetClientVehicles } from "../interfaces/client";
import { ErrorsResponse, GenericResponse } from "../interfaces/generic";
import { buildBaseUrl } from "./build-base-url";

/**
 * Fetch clients in the database
 * @returns [Promise<ResponseGetClients>] The clients information
 */
export async function fetchClients(): Promise<ResponseGetClients> {
    const response = await fetch(buildBaseUrl() + 'clients');
    if (!response.ok) {
        const error: ErrorsResponse = await response.json();
        throw new Error(error.errors.join(", ") || "An unknown error occurred");
    }

    const res: GenericResponse = await response.json();
    const out: ResponseGetClients = {
        clients: res.data,
        next_page: res.next_page
    }
    
    return out;
}


/**
 * Fetch client information
 * @param clientId [string] The client's ID
 * @returns [Promise<ClientInfo>] The client's information
 */
export async function fetchClientInfo(clientId: string): Promise<ResponseGetClient> {
    const response = await fetch(buildBaseUrl() + `clients/${clientId}`);
    if (!response.ok) {
        const error: ErrorsResponse = await response.json();
        throw new Error(error.errors.join(", ") || "An unknown error occurred");
    }

    const res: GenericResponse = await response.json();
    const out: ResponseGetClient = {
        client: res.data,
        next_page: res.next_page
    }

    return out;
}


/**
 * Fetch vehicles for a specific client
 * @param clientId [string] The client's ID
 * @returns [Promise<Vehicles>] The vehicles information
 */
export async function fetchClientVehicles(clientId: string): Promise<ResponseGetClientVehicles> {
    const response = await fetch(buildBaseUrl() + `clients/${clientId}/vehicles`);
    if (!response.ok) {
        const error: ErrorsResponse = await response.json();
        throw new Error(error.errors.join(", ") || "An unknown error occurred");
    }
    
    const res: GenericResponse = await response.json();
    const out: ResponseGetClientVehicles = {
        vehicles: res.data,
        next_page: res.next_page
    }

    return out;
}