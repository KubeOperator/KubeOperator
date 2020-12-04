import {BaseModel, BaseRequest} from '../../../shared/class/BaseModel';

export class System extends BaseModel {
    vars: {} = {};
}

export class SystemCreateRequest extends BaseRequest {
    vars: {} = {};
    tab: string;
}
