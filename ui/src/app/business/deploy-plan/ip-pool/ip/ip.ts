import {BaseModel, BaseRequest} from '../../../../shared/class/BaseModel';

export class Ip extends BaseModel {
    address: string;
    gateway: string;
    dns1: string;
    dns2: string;
    status: string;
}

export class IpCreate extends BaseRequest {
    ipStart: string;
    ipEnd: string;
    gateway: string;
    dns1: string;
    dns2: string;
    ipPoolName: string;
    subnet: string;
}

export class IpUpdate extends BaseRequest {
    address: string;
    operation: string;
}