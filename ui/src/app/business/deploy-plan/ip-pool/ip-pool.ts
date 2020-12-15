import {BaseModel} from '../../../shared/class/BaseModel';

export class IpPool extends BaseModel {
    name: string;
}

export class IpPoolCreate {
    name: string;
}