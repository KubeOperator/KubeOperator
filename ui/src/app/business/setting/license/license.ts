export class License {
    message: string;
    status: string;
    license: {
        corporation: string;
        product: string;
        edition: string;
        count: number;
        expired: string;
    };
}
