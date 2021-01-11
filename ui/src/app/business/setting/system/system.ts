import {BaseModel, BaseRequest} from '../../../shared/class/BaseModel';

export class System extends BaseModel {
    vars: {} = {};
    tab: string;
}

export class SystemCreateRequest extends BaseRequest {
    vars: {} = {};
    tab: string;
}
