export type ErrorsResponse = {
    errors: string[]; // Array of error messages
};

export type GenericResponse = {
    data: any; // Data returned from the API
    next_page: string; // URL for the next page of results
};
