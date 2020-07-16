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
    architectures: string;
    networkType: string;
    runtimeType: string;
    dockerStorageDir: string;
    containerdStorageDir: string;
    flannelBackend: string;
    calicoIpv4poolIpip: string;
    kubePodSubnet: string;
    kubeServiceSubnet: string;
    kubeProxyMode: string;
    kubeMaxPod: number;
    certsExpired: number;
    kubernetesAudit: boolean;
    ingressControllerType: string;
    plan: string;
    nodes: CreateNodeRequest[] = [];
    workerAmount: number;
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
