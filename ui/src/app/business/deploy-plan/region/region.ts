import {BaseModel, BaseRequest} from '../../../shared/class/BaseModel';

export class Region extends BaseModel {
    id: string;
    name: string;
    vars: string;
    datacenter: string;
    regionVars: {} = {};
}

export class RegionCreateRequest extends BaseRequest {
    id: string;
    regionVars: {} = {};
    cloudProvider: string;
    datacenter: string;
    vars: string;
}

