import {BaseModel, BaseRequest} from '../../../shared/class/BaseModel';

export class Zone extends BaseModel {
    name: string;
    vars: string;
    cloudVars: {} = {};
}

export class ZoneCreateRequest extends BaseRequest {
    name: string;
    vars: string;
    region: string;
    regionID: string;
    cloudVars: {} = {};
}

