import {BaseModel, BaseRequest} from '../../../shared/class/BaseModel';

export class Registry extends BaseModel {
    RegistryHostname: string;
    RegistryProtocol: string;
    Architecture: string;
}

export class RegistryCreateRequest extends BaseRequest {
    RegistryHostname: string;
    RegistryProtocol: string;
    Architecture: string;
}