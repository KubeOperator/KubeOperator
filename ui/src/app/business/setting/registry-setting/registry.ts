import {BaseModel, BaseRequest} from '../../../shared/class/BaseModel';

export class Registry extends BaseModel {
    id: string;
    hostname: string;
    protocol: string;
    architecture: string;
}

export class RegistryCreateRequest extends BaseRequest {
    hostname: string;
    protocol: string;
    architecture: string;
}
