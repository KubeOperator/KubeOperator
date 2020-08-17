import {BaseModel, BaseRequest} from '../../../shared/class/BaseModel';
import {Region} from '../region/region';
import {Zone} from '../zone/zone';

export class Plan extends BaseModel {
    name: string;
    zoneId: string;
    deployTemplate: string;
    vars: string;
    regionId: string;
    region: Region = new Region();
    zones: Zone[] = [];
    planVars: {} = {};
    regionName: string;
    zoneNames: string[] = [];
}

export class PlanCreateRequest extends BaseRequest {
    deployTemplate: string;
    vars: string;
    planVars: {} = {};
    // regionId: string;
    zone: string;
    zones: string [] = [];
    region: string;
    projects: string[] = [];
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


