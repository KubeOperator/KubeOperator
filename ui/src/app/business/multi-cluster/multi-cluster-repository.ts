import {BaseModel, BaseRequest} from "../../shared/class/BaseModel";

export class MultiClusterRepository extends BaseModel {
    name: string;
    source: string;
    message: string;
    status: string;
    branch: string;
    gitTimeout: number;
    syncInterval: number;
    syncEnable: boolean;
}

export class MultiClusterRepositoryCreateRequest {
    name: string;
    source: string;
    password: string;
    username: string;
    branch: string;
}

export class MultiClusterRepositoryUpdateRequest extends BaseRequest {
    gitTimeout: number;
    syncInterval: number;
    syncEnable: boolean;
}

export class TreeNode {
    name: string;
    path: string;
    dir: boolean;
    children: TreeNode[] = [];
    content: string;
    originContent: string;
    active: boolean;
    changed: boolean;
    delete = false;
}

export class FileContent {
    content: string;
    path: string;
    readOnly = false;
}

export class ReadFileRequest {
    fileName: string;
}

export class RelatedCluster {
    clusterName: string;
    status: string;
    message: string;
    createdAt: string;
}

export class RelateClusterRequest {
    clusterNames: string[] = [];
    delete = false;
}

export class MultiClusterSyncLog extends BaseModel {
    id: string;
    gitCommitId: string;
    status: string;
}

export class ResourceLog extends BaseModel {
    resourceName: string;
    sourceFile: string;
    status: string;
    message: string;

}

export class MultiClusterSyncClusterLog extends BaseModel {
    clusterName: string;
    multiClusterSyncClusterResourceLogs: ResourceLog[] = [];
}


export class MultiClusterSyncLogDetail extends BaseModel {
    multiClusterSyncClusterLogs: MultiClusterSyncClusterLog[] = [];
}
