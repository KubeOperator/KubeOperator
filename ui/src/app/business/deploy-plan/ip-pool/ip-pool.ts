import {BaseModel} from '../../../shared/class/BaseModel';
import {Ip} from './ip/ip';

export class IpPool extends BaseModel {
    name: string;
    description: string;
    subnet: string;
    status: string;
    ipUsed: number;
    ips: Ip[] = [];
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