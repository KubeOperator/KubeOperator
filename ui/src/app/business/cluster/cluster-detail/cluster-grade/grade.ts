import {BaseModel} from '../../../../shared/class/BaseModel';

export class Grade extends BaseModel {
    grade: string;
    score: number;
    totalSum: Summary;
    listSum: Summary[] = [];
    result: [] = [];
}

export class Summary {
    success: number;
    warning: number;
    danger: number;
}