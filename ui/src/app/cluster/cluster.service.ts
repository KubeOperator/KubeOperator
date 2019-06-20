import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Cluster, ExtraConfig} from './cluster';
import {Observable, throwError} from 'rxjs';
import {catchError} from 'rxjs/operators';
import {HostService} from '../host/host.service';


const baseClusterUrl = '/api/v1/clusters/';

const baseClusterConfigUrl = '/api/v1/clusters/{cluster_name}/configs/';

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

  getCluster(clusterName): Observable<Cluster> {
    return this.http.get<Cluster>(`${baseClusterUrl}${clusterName}`).pipe(
      catchError(error => throwError(error))
    );
  }

  configCluster(clusterName: string, extraConfig: ExtraConfig): Observable<ExtraConfig> {
    return this.http.post<ExtraConfig>(`${baseClusterConfigUrl.replace('{cluster_name}', clusterName)}`, extraConfig).pipe(
      catchError(error => throwError(error))
    );
  }

  configClusterAuth(clusterName: string, auth: string): Observable<Cluster> {
    return this.http.patch<Cluster>(baseClusterUrl + clusterName + '/', {auth_template: auth});
  }

  getClusterConfig(clusterName: string, key: string): Observable<ExtraConfig> {
    return this.http.get<ExtraConfig>(baseClusterConfigUrl.replace('{cluster_name}', clusterName) + key);
  }

  createCluster(cluster: Cluster): Observable<Cluster> {
    return this.http.post<Cluster>(baseClusterUrl, cluster).pipe(
      catchError(error => throwError(error))
    );
  }

  deleteCluster(clusterName): Observable<any> {
    return this.http.delete(`${baseClusterUrl}${clusterName}`).pipe(
      catchError(error => throwError(error))
    );
  }
}
