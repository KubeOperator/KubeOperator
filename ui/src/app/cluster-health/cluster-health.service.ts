import { Injectable } from '@angular/core';
import {Observable} from 'rxjs';
import {ClusterHealth} from './cluster-health';
import {HttpClient} from '@angular/common/http';
import {ClusterHealthHistory} from './cluster-health-history';

@Injectable({
  providedIn: 'root'
})
export class ClusterHealthService {
  baseUrl = '/api/v1/cluster/';
  healthHistoryUrl = '/api/v1/clusterHealthHistory/';

  constructor(private http: HttpClient) { }

  listClusterHealth(project_name: string): Observable<ClusterHealth> {
    return this.http.get<ClusterHealth>(this.baseUrl + project_name + '/health/');
  }

  listClusterHealthHistory(project_id: string): Observable<ClusterHealthHistory[]> {
    return this.http.get<ClusterHealthHistory[]>(this.healthHistoryUrl + project_id + '/');
  }
}
