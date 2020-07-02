import {BaseModel} from '../../../../../shared/class/BaseModel';

export class StorageProvisioner extends BaseModel {
    name: string;
    type: string;
    vars: string;
}

export class CreateStorageProvisionerRequest {
    name: string;
    type: string;
    vars: {} = {};
}
