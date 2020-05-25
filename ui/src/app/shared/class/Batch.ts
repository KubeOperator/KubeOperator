import {BaseModel} from './BaseModel';

export class Batch<T extends BaseModel> {
    constructor(operation: string, items: T[]) {
        this.operation = operation;
        this.items = items;
    }

    operation: string;
    items: T [] = [];
}
