import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Role} from './role';

@Injectable({
  providedIn: 'root'
})
export class RoleService {

  constructor(private http: HttpClient) {
  }

  listRoles(clusterName: string): Observable<Role[]> {
    return this.http.get<Role[]>('/api/v1/clusters/{clusterName}/roles/'.replace('{clusterName}', clusterName));
  }
}
