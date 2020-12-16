import {BaseModel} from '../../../shared/class/BaseModel';

export class IpPool extends BaseModel {
    name: string;
    description: string;
    subnet: string;
    status: string;
}

export class IpPoolCreate {
    name: string;
    description: string;
    subnet: string;
    ipStart: string;
    ipEnd: string;
    gateway: string;
    dns1: string;
    dns2: string;
}