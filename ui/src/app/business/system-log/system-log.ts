export class SystemLog {
    name: string;
    operationUnit: string;
    operation: string;
    completeOperation: string;
    isExceed: boolean;
    ip: string;
}

export class Page<T> {
    total: number;
    items: T[] = [];
}
