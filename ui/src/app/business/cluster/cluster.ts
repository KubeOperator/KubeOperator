import {BaseModel} from '../../shared/class/BaseModel';

export class Cluster extends BaseModel {
    name: string;
    spec: Spec;
    status: Status;
}

export class Spec {
    version: string;
}

export class Status {
    phase: string;
}


export class ClusterCreateRequest extends BaseModel {
    name: string;
    version: string;
    provider: string;
    networkType: string;
    runtimeType: string;
    clusterCIDR: string;
    serviceCIDR: string;
}
