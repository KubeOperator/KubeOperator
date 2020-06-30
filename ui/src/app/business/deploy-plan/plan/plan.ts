import {BaseModel, BaseRequest} from '../../../shared/class/BaseModel';

export class Plan extends BaseModel {
    name: string;
    zoneId: string;
    deployTemplate: string;
    vars: string;
    regionId: string;
}

export class PlanCreateRequest extends BaseRequest {
    deployTemplate: string;
    vars: string;
    planVars: {} = {};
    regionId: string;
    zone: string;
    zones: string [] = [];
}

export class PlanVmConfig {
    name: string;
    config: VmConfig;
}

export class VmConfig {
    cpu: number;
    memory: number;
    disk: number;
}


