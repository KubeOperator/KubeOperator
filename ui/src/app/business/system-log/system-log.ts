export class SystemLog {
    name: string;
    operation_unit: string;
    operation: string;
}

export class Page<T> {
    total: number;
    items: T[] = [];
}
