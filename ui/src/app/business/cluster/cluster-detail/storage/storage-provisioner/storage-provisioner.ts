import {BaseModel} from '../../../../../shared/class/BaseModel';

export class StorageProvisioner extends BaseModel {
    id: string;
    name: string;
    type: string;
    status: string;
    vars: string;
}

export class CreateStorageProvisionerRequest {
    name: string;
    type: string;
    vars: {} = {};
}
