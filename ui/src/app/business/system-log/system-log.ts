export class SystemLog {
    name: string;
    operationUnit: string;
    operation: string;
    completeOperation: string;
    isExceed: boolean;
}

export class Page<T> {
    total: number;
    items: T[] = [];
}
