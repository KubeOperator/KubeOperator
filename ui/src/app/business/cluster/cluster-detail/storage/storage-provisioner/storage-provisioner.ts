import {BaseModel} from '../../../../../shared/class/BaseModel';

export class StorageProvisioner extends BaseModel {
    id: string;
    name: string;
    type: string;
    status: string;
    vars: string;
    message: string;
}

export class CreateStorageProvisionerRequest {
    name: string;
    type: string;
    vars: {} = {};
}

export class ProvisionerSync {
    name: string;
    type: string;
    status: string;
}