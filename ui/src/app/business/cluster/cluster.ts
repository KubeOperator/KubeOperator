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
    multiClusterRepository: string;
}

export class Spec {
    dockerStorageDir: string;
    version: string;
    upgradeVersion: string;
    kubeProxyMode: string;
    networkType: string;
    ciliumTunnelMode: string;
    flannelBackend: string;
    calicoIpv4PoolIpip: string;
    architectures: string;
    runtimeType: string;
}

export class ClusterStatus {
    phase: string;
    prePhase: string;
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
    ciliumVersion: string;
    ciliumTunnelMode: string;
    ciliumNativeRoutingCidr: string;
    runtimeType: string;
    dockerStorageDir: string;
    containerdStorageDir: string;
    flannelBackend: string;
    calicoIpv4poolIpip: string;
    kubeProxyMode: string;
    nodeportAddress: string;
    enableDnsCache: string;
    dnsCacheVersion: string;
    kubernetesAudit: string;
    ingressControllerType: string;
    plan: string;
    nodes: CreateNodeRequest[] = [];
    workerAmount: number;
    dockerSubnet: string;
    projectName: string;
    helmVersion: string;
    networkInterface: string;
    supportGpu: string;
    yumOperate: string;
    clusterCidr: string;
    serviceCidr: string;
    maxPodNum: number;
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

export class ClusterHealthCheck {
    hooks: ClusterHealthCheckHook[] = [];
    level: string;
}

export class ClusterHealthCheckHook {
    name: string;
    level: string;
    msg: string;
}

export class ClusterRecoverItem {
    name: string;
    hookName: string;
    result: string;
    msg: string;
}
