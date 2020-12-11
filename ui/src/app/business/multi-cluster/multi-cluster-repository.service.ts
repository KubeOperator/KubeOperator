import {Injectable} from '@angular/core';
import {BaseModelService} from "../../shared/class/BaseModelService";
import {
    FileContent,
    MultiClusterRepository, MultiClusterSyncLog, MultiClusterSyncLogDetail,
    ReadFileRequest, RelateClusterRequest,
    RelatedCluster,
    TreeNode
} from "./multi-cluster-repository";
import {HttpClient} from "@angular/common/http";
import {Observable} from "rxjs";
import {Page} from "../../shared/class/Page";

@Injectable({
    providedIn: 'root'
})
export class MultiClusterRepositoryService extends BaseModelService<MultiClusterRepository> {

    baseUrl = '/api/v1/multicluster/repositories';

    constructor(http: HttpClient) {
        super(http);
    }

    getTree(name: string): Observable<TreeNode> {
        return this.http.get<TreeNode>(`${this.baseUrl}/tree/${name}`);
    }

    readFile(name: string, fileName: string): Observable<FileContent> {
        const content = new FileContent();
        content.path = fileName;
        content.readOnly = true;
        return this.http.post<FileContent>(`${this.baseUrl}/tree/content/${name}`, content);
    }

    saveFile(name: string, content: FileContent): Observable<FileContent> {
        return this.http.post<FileContent>(`${this.baseUrl}/tree/content/${name}`, content);
    }

    createOrDeleteTreeNode(name: string, treeNode: TreeNode): Observable<TreeNode> {
        return this.http.post<TreeNode>(`${this.baseUrl}/tree/${name}`, treeNode);
    }

    updateTreeNode(name: string, treeNode: TreeNode): Observable<TreeNode> {
        return this.http.post<TreeNode>(`${this.baseUrl}/tree/${name}`, treeNode);
    }

    pullRemoteRepository(name: string): Observable<any> {
        return this.http.get(`${this.baseUrl}/remote/${name}`);
    }

    pushRemoteRepository(name: string): Observable<any> {
        return this.http.post(`${this.baseUrl}/remote/${name}`, {});
    }

    listRelations(name: string): Observable<RelatedCluster[]> {
        return this.http.get<RelatedCluster[]>(`${this.baseUrl}/relations/${name}`);
    }

    createRelations(name: string, clusterNames: string[]): Observable<any> {
        const req = new RelateClusterRequest();
        req.delete = false;
        req.clusterNames = clusterNames;
        return this.http.post(`${this.baseUrl}/relations/${name}`, req);
    }

    deleteRelations(name: string, clusterNames: string[]) {
        const req = new RelateClusterRequest();
        req.delete = true;
        req.clusterNames = clusterNames;
        return this.http.post(`${this.baseUrl}/relations/${name}`, req);
    }

    getLog(name: string, page: number, size: number): Observable<Page<MultiClusterSyncLog>> {
        return this.http.get<Page<MultiClusterSyncLog>>(`${this.baseUrl}/logs/${name}?pageNum=${page}&pageSize=${size}`);
    }

    getLogDetail(name: string, logId: string): Observable<MultiClusterSyncLogDetail> {
        return this.http.get<MultiClusterSyncLogDetail>(`${this.baseUrl}/logs/detail/${name}/${logId}`);
    }
}
