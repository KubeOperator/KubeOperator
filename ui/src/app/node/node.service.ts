import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable, throwError} from 'rxjs';
import {catchError} from 'rxjs/operators';
import {Node} from './node';
import {Role} from './role';
import {Cluster} from '../cluster/cluster';

const baseNodeUrl = '/api/v1/clusters/{clusterName}/nodes/';
const roleUrl = '/api/v1/clusters/{clusterName}/roles/';

@Injectable({
  providedIn: 'root'
})
export class NodeService {

  checkNodeUrl = '/api/v1/cluster/';

  constructor(private http: HttpClient) {
  }

  listNodes(clusterName): Observable<Node[]> {
    return this.http.get<Node[]>(baseNodeUrl.replace('{clusterName}', clusterName)).pipe(
      catchError(err => throwError(err))
    );
  }

  createNode(clusterName, node: Node): Observable<Node> {
    return this.http.post<Node>(baseNodeUrl.replace('{clusterName}', clusterName), node).pipe(
      catchError(err => throwError(err))
    );
  }

  deleteNode(clusterName, nodeId): Observable<any> {
    return this.http.delete(`${baseNodeUrl.replace('{clusterName}', clusterName)}/${nodeId}`).pipe(
      catchError(err => throwError(err))
    );
  }

  listRoles(clusterName): Observable<Role[]> {
    return this.http.get<Role[]>(`${roleUrl.replace('{clusterName}', clusterName)}`).pipe(
      catchError(err => throwError(err))
    );
  }

  get_grafana_url(nodeIp: string, cluster: Cluster): string {
    const base = cluster.apps['nodes_grafana'];
    if (!base) {
      return null;
    }
    return `${base}?orgId=1&var-server=${nodeIp}:9100`;
  }

  checkNodes(project_name: string): Observable<any> {
    return this.http.get<any>(this.checkNodeUrl + project_name + '/checkNodes/');
  }

  syncHostTime(project_name: string): Observable<any> {
    return this.http.get<any>(this.checkNodeUrl + project_name + '/syncNodeTime/');
  }
}
