export class Manifest {
    name: string;
    coreVars: NameVersion[] = [];
    networkVars: NameVersion[] = [];
    toolVars: NameVersion[] = [];
    otherVars: NameVersion[] = [];
    version: string;
    isActive: boolean;
}

export class Category {
    name: string;
    items: Item[] = [];
}

export class Item {
    name: string;
    version: string;
}

export class NameVersion {
    name: string;
    version: string;
}



