import {BaseModel} from './BaseModel';

export class Batch<T extends BaseModel> {
    constructor(method: string, items: T[]) {
        this.method = method;
        this.items = items;
    }

    method: string;
    items: T [] = [];
}
