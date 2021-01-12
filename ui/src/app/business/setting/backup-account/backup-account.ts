import {BaseModel, BaseRequest} from '../../../shared/class/BaseModel';

export class BackupAccount extends BaseModel {
    name: string;
    bucket: string;
    credentialVars: {} = {};
    status: string;
    type: string;
}

export class BackupAccountCreateRequest extends BaseRequest {
    name: string;
    bucket: string;
    type: string;
    credentialVars: {} = {};
}

export class BackupAccountUpdateRequest extends BaseRequest {
    name: string;
    bucket: string;
    type: string;
    credentialVars: {} = {};
}

export class Project {
    id: string;
    name: string;
    checked: boolean;
}