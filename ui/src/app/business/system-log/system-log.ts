export class SystemLog {
    name: string;
    operationUnit: string;
    operation: string;
    requestPath: string;
}

export class Page<T> {
    total: number;
    items: T[] = [];
}
