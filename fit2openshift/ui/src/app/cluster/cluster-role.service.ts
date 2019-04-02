import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {ClusterRole} from './cluster-role';
import {Observable} from 'rxjs';

export const url = '/api/v1/clusters/{clusterName}/roles/{name}/';

@Injectable({
  providedIn: 'root'
})
export class ClusterRoleService {

  constructor(private http: HttpClient) {
  }

  getClusterRole(clusterName: string, roleName: string): Observable<ClusterRole> {
    return this.http.get<ClusterRole>(url.replace('{clusterName}', clusterName).replace('{name}', roleName));
  }
}
