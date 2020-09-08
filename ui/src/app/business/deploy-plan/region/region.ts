import {BaseModel, BaseRequest} from '../../../shared/class/BaseModel';

export class Region extends BaseModel {
    id: string;
    name: string;
    vars: string;
    datacenter: string;
    regionVars: {} = {};
    provider: string;
}

export class RegionCreateRequest extends BaseRequest {
    id: string;
    regionVars: {} = {};
    provider: string;
    datacenter: string;
    vars: string;
}

