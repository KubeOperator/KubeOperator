import {BaseModel} from '../../../../shared/class/BaseModel';

export class Grade extends BaseModel {
    grade: string;
    score: number;
    totalSum: Summary;
    listSum: Summary[] = [];
    results: NamespaceResult[] = [];
}

export class Summary {
    success: number;
    warning: number;
    danger: number;
}

export class NamespaceResult {
    namespace: string;
    results: NamespaceResultDetail[] = [];
}

export class NamespaceResultDetail {
    kind: string;
    name: string;
    podResults: PodResult[] = [];
}

export class PodResult {
    category: string;
    id: string;
    message: string;
    severity: string;
    success: boolean;
}