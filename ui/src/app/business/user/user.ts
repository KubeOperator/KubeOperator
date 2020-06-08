import {BaseModel} from '../../shared/class/BaseModel';

export class User extends BaseModel {
    name: string;
    id: string;
    password: string;
    email: string;
    language: string;
    isActive: string;
}
