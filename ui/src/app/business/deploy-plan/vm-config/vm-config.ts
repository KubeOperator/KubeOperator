import {BaseModel} from '../../../shared/class/BaseModel';

export class VmConfig extends BaseModel {
    cpu: number;
    memory: number;
    disk: number;
    name: string;
    provider: string;
}
