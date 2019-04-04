import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable, throwError} from 'rxjs';
import {catchError} from 'rxjs/operators';
import {Node} from './node';
import {Role} from './role';

const baseNodeUrl = '/api/v1/clusters/{clusterName}/nodes/';
const roleUrl = '/api/v1/clusters/{clusterName}/roles/';

@Injectable({
  providedIn: 'root'
})
export class NodeService {

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
}
