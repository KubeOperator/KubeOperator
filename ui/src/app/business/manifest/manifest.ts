export class Manifest {
    name: string;
    category: Category[] = [];
}

export class Category {
    name: string;
    items: Item[] = [];
}

export class Item {
    name: string;
    version: string;
}



