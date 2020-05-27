import {BaseModel, BaseRequest} from '../../../shared/class/BaseModel';

export class Credential extends BaseModel {
    id: string;
    name: string;
    username: string;
    password: string;
    privateKey: string;
    type = 'password';
}


export class CredentialCreateRequest extends BaseRequest {
    id: string;
    name: string;
    username: string;
    password: string;
    privateKey: string;
    type: string;
}
