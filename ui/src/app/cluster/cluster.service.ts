import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Cluster, ClusterConfigs, ExtraConfig} from './cluster';
import {Observable, throwError} from 'rxjs';
import {catchError} from 'rxjs/operators';
import {HostService} from '../host/host.service';


const baseClusterUrl = '/api/v1/clusters/';
const webKubeCtlUrl = '/api/v1/cluster/{id}/webkubectl/token/';
const checkNameSpaceUrl = '/api/v1/cluster/{project_name}/check/{namespace}/';


@Injectable({
  providedIn: 'root'
})
export class ClusterService {

  constructor(private http: HttpClient, private hostService: HostService) {
  }

  listCluster(): Observable<Cluster[]> {
    return this.http.get<Cluster[]>(baseClusterUrl).pipe(
      catchError(error => throwError(error)));
  }

  listItemClusters(itemName: string): Observable<Cluster[]> {
    return this.http.get<Cluster[]>(baseClusterUrl + '?itemName=' + itemName).pipe(
      catchError(error => throwError(error)));
  }

  getCluster(clusterName): Observable<Cluster> {
    return this.http.get<Cluster>(`${baseClusterUrl}${clusterName}`).pipe(
      catchError(error => throwError(error))
    );
  }

  createCluster(cluster: Cluster): Observable<Cluster> {
    return this.http.post<Cluster>(baseClusterUrl, cluster).pipe(
      catchError(error => throwError(error))
    );
  }

  updateCluster(cluster: Cluster): Observable<Cluster> {
    return this.http.patch<Cluster>(`${baseClusterUrl}${cluster.name}/`, cluster);
  }

  deleteCluster(clusterName): Observable<any> {
    return this.http.delete(`${baseClusterUrl}${clusterName}`).pipe(
      catchError(error => throwError(error))
    );
  }

  getClusterConfigs(): Observable<ClusterConfigs> {
    return this.http.get<ClusterConfigs>('/api/v1/cluster/config');
  }

  getWebkubectlToken(id: string): Observable<any> {
    return this.http.get<any>(webKubeCtlUrl.replace('{id}', id));
  }

  changeStatus(status: string, name: string): Observable<Cluster> {
    return this.http.patch<Cluster>(`${baseClusterUrl}${name}/`, {'status': status});
  }

  checkNameSpace(name: string, namespace: string): Observable<Boolean> {
    return this.http.get<Boolean>(checkNameSpaceUrl.replace('{project_name}', name).replace('{namespace}', namespace));
  }
}
