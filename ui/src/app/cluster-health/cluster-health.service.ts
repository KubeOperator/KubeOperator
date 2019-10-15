import { Injectable } from '@angular/core';
import {Observable} from 'rxjs';
import {ClusterHealth} from './cluster-health';
import {HttpClient} from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class ClusterHealthService {
  baseUrl = '/api/v1/cluster/';

  constructor(private http: HttpClient) { }

  listClusterHealth(project_name: string): Observable<ClusterHealth> {
    return this.http.get<ClusterHealth>(this.baseUrl + project_name + '/health/');
  }
}
