import {BaseModel} from '../../../../shared/class/BaseModel';

export class Ip extends BaseModel {
    address: string;
    gateway: string;
    dns1: string;
    dns2: string;
    status: string;
}