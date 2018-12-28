import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Cluster} from './cluster';
import {Observable, throwError} from 'rxjs';
import {catchError} from 'rxjs/operators';


const baseClusterUrl = '/api/v1/clusters/';

@Injectable({
  providedIn: 'root'
})
export class ClusterService {

  constructor(private http: HttpClient) {
  }

  listCluster(): Observable<Cluster[]> {
    return this.http.get<Cluster[]>(baseClusterUrl).pipe(
      catchError(error => throwError(error)));
  }

  getCluster(clusterId): Observable<Cluster> {
    return this.http.get<Cluster>(`${baseClusterUrl}/${clusterId}`).pipe(
      catchError(error => throwError(error))
    );
  }

  createCluster(cluster: Cluster): Observable<Cluster> {
    return this.http.post<Cluster>(baseClusterUrl, cluster).pipe(
      catchError(error => throwError(error))
    );
  }

  deleteCluster(clusterName): Observable<any> {
    return this.http.delete(`${baseClusterUrl}/${clusterName}`).pipe(
      catchError(error => throwError(error))
    );
  }
}
