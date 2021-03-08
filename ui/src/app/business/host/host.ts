import {BaseModel, BaseRequest} from '../../shared/class/BaseModel';

export class Host extends BaseModel {
    name: string;
    ip: string;
    port: string;
    os: string;
    osVersion: string;
    memory: string;
    cpuCore: number;
    gpuNum: number;
    gpuInfo: string;
    status: string;
    volumes: Volume[];
    projectName: string;
    clusterName: string;
    clusterId: string;
    zoneName: string;
    message: string;
    hasGpu: boolean;
    architecture: string;
}

export class Volume extends BaseModel {
    size: string;
    name: string;
    hostId: string;
}

export class HostCreateRequest extends BaseRequest {
    name: string;
    ip: string;
    port: string;
    credentialId: string;
}

export class HostSync {
    hostName: string;
    hostStatus: string;
}

export class Project {
    id: string;
    name: string;
}
