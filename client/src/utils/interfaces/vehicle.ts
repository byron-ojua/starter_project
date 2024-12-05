// Struct to match API ClientVehicle struct
export interface BasicVehicle {
    vin: string,
    mileage: number,
    largest_weight: number,
}

export interface ResponseGetVehicle {
    vehicle: VehicleInfo,
}

export interface VehicleInfo {
    vin: string,
    client_name: string,
    contact_name: string,
    contact_email: string,
    mileage: number,
    weights: number[],
}