import {BaseModel, BaseRequest} from '../../../shared/class/BaseModel';

export class Zone extends BaseModel {
    name: string;
    vars: string;
    cloudVars: {} = {};
}

export class ZoneCreateRequest extends BaseRequest {
    vars: string;
    region: string;
    regionID: string;
    cloudVars: {} = {};
    provider: string;
}

export class CloudZoneRequest extends BaseRequest {
    cloudVars: {} = {};
    datacenter: string;
}



