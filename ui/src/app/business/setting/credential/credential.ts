import {BaseModel} from '../../../shared/class/BaseModel';

export class Credential extends BaseModel {
    id: string;
    name: string;
    username: string;
    password: string;
    privateKey: string;
    type = 'password';
}
