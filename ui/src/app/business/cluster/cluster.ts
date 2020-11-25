import {BaseModel} from '../../shared/class/BaseModel';

export class Cluster extends BaseModel {
    name: string;
    spec: Spec;
    nodeSize: string;
    status: string;
    preStatus: string;
    provider: string;
    projectName: string;
    source: string;
}

export class Spec {
    dockerStorageDir: string;
    version: string;
    upgradeVersion: string;
    networkType: string;
    architectures: string;
    runtimeType: string;
}

export class ClusterStatus {
    phase: string;
    message: string;
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
    kubeMaxPods: number;
    certsExpired: number;
    kubernetesAudit: string;
    ingressControllerType: string;
    plan: string;
    nodes: CreateNodeRequest[] = [];
    workerAmount: number;
    dockerSubnet: string;
    projectName: string;
    helmVersion: string;
    networkInterface: string;
}

export class CLusterImportRequest {
    name: string;
    apiServer: string;
    token: string;
    router: string;
    projectName: string;
}

export class InitClusterResponse {
    message: string;
}


export class ClusterSecret {
    kubernetesToken: string;
}

export class ClusterUpgradeRequest {
    clusterName: string;
    version: string;
}
