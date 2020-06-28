import {BaseModel} from '../../shared/class/BaseModel';

export class Cluster extends BaseModel {
    name: string;
    spec: Spec;
    nodeSize: string;
    status: string;
    ingressStatus: string;
}

export class Spec {
    version: string;
    networkType: string;
}

export class ClusterStatus {
    phase: string;
    conditions: Condition[] = [];
}

export class ClusterMonitor {
    enable: boolean;
    domain: string;
    status: string;
    dashboardUrl: string;
}


export class Condition {
    status: string;
    message: string;
    name: string;
}

export class CreateNodeRequest {
    role: string;
    hostName: string;
}


export class ClusterCreateRequest extends BaseModel {
    name: string;
    version: string;
    provider: string;
    networkType: string;
    runtimeType: string;
    dockerStorageDir: string;
    containerdStorageDir: string;
    appDomain: string;
    clusterCIDR: string;
    serviceCIDR: string;
    nodes: CreateNodeRequest[] = [];
}

export class CLusterImportRequest {
    name: string;
    apiServer: string;
    token: string;
    router: string;
}

export class InitClusterResponse {
    message: string;
}
