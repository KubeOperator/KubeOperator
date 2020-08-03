import {BaseModel} from '../../../shared/class/BaseModel';

export class BackupAccount extends BaseModel {
    name: string;
    region: string;
    credentialVars: {} = {};
    status: string;
    type: string;
}
