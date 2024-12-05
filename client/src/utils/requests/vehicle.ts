import { ErrorsResponse, GenericResponse } from "../interfaces/generic";
import { ResponseGetVehicle } from "../interfaces/vehicle";
import { buildBaseUrl } from "./build-base-url";

/**
 * Fetch vehicle information
 * @param vehicleId [string] The client's ID
 * @returns [Promise<ResponseGetVehicle>] The vehicle's information
 */
export async function fetchVehicle(vehicleId: string): Promise<ResponseGetVehicle> {
    const response = await fetch(buildBaseUrl() + `vehicles/${vehicleId}`);
    if (!response.ok) {
        const error: ErrorsResponse = await response.json();
        throw new Error(error.errors.join(", ") || "An unknown error occurred");
    }
    
    const res: GenericResponse = await response.json();
    const out: ResponseGetVehicle = {
        vehicle: res.data
    }

    return out;
}